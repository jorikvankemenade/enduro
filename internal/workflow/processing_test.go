package workflow

import (
	"errors"
	"testing"
	"time"

	"github.com/artefactual-labs/enduro/internal/collection"
	collectionfake "github.com/artefactual-labs/enduro/internal/collection/fake"
	nha_activities "github.com/artefactual-labs/enduro/internal/nha/activities"
	"github.com/artefactual-labs/enduro/internal/pipeline"
	"github.com/artefactual-labs/enduro/internal/watcher"
	watcherfake "github.com/artefactual-labs/enduro/internal/watcher/fake"
	"github.com/artefactual-labs/enduro/internal/workflow/manager"

	logrtesting "github.com/go-logr/logr/testing"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/cadence"
	cadencetestsuite "go.uber.org/cadence/testsuite"
	cadenceworkflow "go.uber.org/cadence/workflow"
)

type ProcessingWorkflowTestSuite struct {
	suite.Suite
	cadencetestsuite.WorkflowTestSuite

	env *cadencetestsuite.TestWorkflowEnvironment

	manager *manager.Manager

	// Each test registers the workflow with a different name to avoid dups.
	workflow string
}

func (s *ProcessingWorkflowTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
	s.manager = buildManager(s.T(), gomock.NewController(s.T()))

	s.workflow = uuid.New().String()
	cadenceworkflow.RegisterWithOptions(NewProcessingWorkflow(s.manager).Execute, cadenceworkflow.RegisterOptions{Name: s.workflow})
}

func (s *ProcessingWorkflowTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

// Workflow ignores an error in parseName when NHA hooks are disabled.
func (s *ProcessingWorkflowTestSuite) TestParseErrorIsIgnored() {
	s.manager.Hooks["hari"]["disabled"] = true
	s.manager.Hooks["prod"]["disabled"] = true

	// Collection is persisted.
	s.env.OnActivity(createPackageLocalActivity, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(uint(12345), nil).Once()

	// parseName is executed with errors.
	s.env.OnActivity(nha_activities.ParseNameLocalActivity, mock.Anything, "key").Return(nil, errors.New("parse error")).Once()

	// loadConfig is executed (workflow continued), returning an error.
	s.env.OnActivity(loadConfigLocalActivity, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("pipeline is unavailable")).Once()

	// Defer updates the package with the error status before returning.
	s.env.OnActivity(updatePackageLocalActivity, mock.Anything, mock.Anything, mock.Anything, &updatePackageLocalActivityParams{
		CollectionID: uint(12345),
		Key:          "key",
		PipelineID:   "",
		TransferID:   "",
		SIPID:        "",
		StoredAt:     time.Time{},
		Status:       collection.StatusError,
	}).Return(nil).Once()

	retentionPeriod := time.Second
	s.env.ExecuteWorkflow(s.workflow, &collection.ProcessingWorkflowRequest{
		CollectionID: 0,
		Event: &watcher.BlobEvent{
			WatcherName:      "watcher",
			PipelineName:     "pipeline",
			RetentionPeriod:  &retentionPeriod,
			StripTopLevelDir: true,
			Key:              "key",
			Bucket:           "bucket",
		},
	})

	s.True(s.env.IsWorkflowCompleted())
	s.EqualError(s.env.GetWorkflowError(), "error loading configuration: pipeline is unavailable")
}

// Workflow does not ignore an error in parseName when NHA hooks are enabled.
func (s *ProcessingWorkflowTestSuite) TestParseError() {
	s.manager.Hooks["hari"]["disabled"] = false
	s.manager.Hooks["prod"]["disabled"] = false

	// Collection is persisted.
	s.env.OnActivity(createPackageLocalActivity, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(uint(12345), nil).Once()

	// parseName is executed, inject error.
	s.env.OnActivity(nha_activities.ParseNameLocalActivity, mock.Anything, "key").Return(nil, errors.New("parse error")).Once()

	// Defer updates the package with the error status before returning.
	s.env.OnActivity(updatePackageLocalActivity, mock.Anything, mock.Anything, mock.Anything, &updatePackageLocalActivityParams{
		CollectionID: uint(12345),
		Key:          "key",
		PipelineID:   "",
		TransferID:   "",
		SIPID:        "",
		StoredAt:     time.Time{},
		Status:       collection.StatusError,
	}).Return(nil).Once()

	retentionPeriod := time.Second
	s.env.ExecuteWorkflow(s.workflow, &collection.ProcessingWorkflowRequest{
		CollectionID: 0,
		Event: &watcher.BlobEvent{
			WatcherName:      "watcher",
			PipelineName:     "pipeline",
			RetentionPeriod:  &retentionPeriod,
			StripTopLevelDir: true,
			Key:              "key",
			Bucket:           "bucket",
		},
	})

	s.True(s.env.IsWorkflowCompleted())
	s.EqualError(s.env.GetWorkflowError(), "error parsing transfer name: parse error")
}

func TestProcessingWorkflow(t *testing.T) {
	suite.Run(t, new(ProcessingWorkflowTestSuite))
}

func buildManager(t *testing.T, ctrl *gomock.Controller) *manager.Manager {
	t.Helper()

	return manager.NewManager(
		logrtesting.NullLogger{},
		collectionfake.NewMockService(ctrl),
		watcherfake.NewMockService(ctrl),
		&pipeline.Registry{},
		map[string]map[string]interface{}{
			"prod": {"disabled": "false"},
			"hari": {"disabled": "false"},
		},
	)
}
func assertNilWorkflowError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		return
	}

	if perr, ok := err.(*cadence.CustomError); ok {
		var details string
		perr.Details(&details)
		t.Fatal(details)
	} else {
		t.Fatal(err.Error())
	}

}
