/* tslint:disable */
/* eslint-disable */
/**
 * Enduro API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { exists, mapValues } from '../runtime';
/**
 * Collection not found
 * @export
 * @interface CollectionCancelNotFoundResponseBody
 */
export interface CollectionCancelNotFoundResponseBody {
    /**
     * Identifier of missing collection
     * @type {number}
     * @memberof CollectionCancelNotFoundResponseBody
     */
    id: number;
    /**
     * Message of error
     * @type {string}
     * @memberof CollectionCancelNotFoundResponseBody
     */
    message: string;
}

export function CollectionCancelNotFoundResponseBodyFromJSON(json: any): CollectionCancelNotFoundResponseBody {
    return CollectionCancelNotFoundResponseBodyFromJSONTyped(json, false);
}

export function CollectionCancelNotFoundResponseBodyFromJSONTyped(json: any, ignoreDiscriminator: boolean): CollectionCancelNotFoundResponseBody {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'message': json['message'],
    };
}

export function CollectionCancelNotFoundResponseBodyToJSON(value?: CollectionCancelNotFoundResponseBody | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'message': value.message,
    };
}


