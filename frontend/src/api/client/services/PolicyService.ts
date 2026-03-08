/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class PolicyService {
    /**
     * List policies
     * @returns any OK
     * @throws ApiError
     */
    public static listPolicies(): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/policy/list',
        });
    }
    /**
     * Get one policy
     * @returns any OK
     * @throws ApiError
     */
    public static getPolicy(): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/policy/get',
        });
    }
    /**
     * Search policies
     * @returns any OK
     * @throws ApiError
     */
    public static searchPolicies(): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/policy/search',
        });
    }
    /**
     * Match policies for founder profile
     * @returns any OK
     * @throws ApiError
     */
    public static matchPolicies(): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/policy/match',
        });
    }
}
