/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ResponseMatchResults } from '../models/ResponseMatchResults';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class MatchingService {
    /**
     * Match investors for founder
     * @param requestBody
     * @returns ResponseMatchResults OK
     * @throws ApiError
     */
    public static matchInvestorsForFounder(
        requestBody: {
            founder_id: string;
            top_k?: number;
        },
    ): CancelablePromise<ResponseMatchResults> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/founder/match-investors',
            body: requestBody,
            mediaType: 'application/json',
        });
    }
    /**
     * Match founders for investor
     * @returns ResponseMatchResults OK
     * @throws ApiError
     */
    public static matchFoundersForInvestor(): CancelablePromise<ResponseMatchResults> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/investor/match-founders',
        });
    }
}
