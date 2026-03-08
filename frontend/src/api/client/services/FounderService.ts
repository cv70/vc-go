/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ResponseFounder } from '../models/ResponseFounder';
import type { ResponseFounderSearch } from '../models/ResponseFounderSearch';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class FounderService {
    /**
     * Register founder profile
     * @param requestBody
     * @returns ResponseFounder OK
     * @throws ApiError
     */
    public static registerFounder(
        requestBody: {
            name: string;
            stage: string;
            sectors?: Array<string>;
        },
    ): CancelablePromise<ResponseFounder> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/founder/register',
            body: requestBody,
            mediaType: 'application/json',
        });
    }
    /**
     * Search founders
     * @param requestBody
     * @returns ResponseFounderSearch OK
     * @throws ApiError
     */
    public static searchFounders(
        requestBody: {
            keyword?: string | null;
            stages?: Array<string>;
            sectors?: Array<string>;
            cursor?: string | null;
            limit?: number;
        },
    ): CancelablePromise<ResponseFounderSearch> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/founder/search',
            body: requestBody,
            mediaType: 'application/json',
        });
    }
}
