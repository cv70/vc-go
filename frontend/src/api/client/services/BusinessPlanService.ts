/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class BusinessPlanService {
    /**
     * Create business plan
     * @returns any OK
     * @throws ApiError
     */
    public static createBusinessPlan(): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/finance/business/create',
        });
    }
    /**
     * Get one business plan
     * @returns any OK
     * @throws ApiError
     */
    public static getBusinessPlan(): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/finance/business/plan',
        });
    }
    /**
     * List business plans
     * @returns any OK
     * @throws ApiError
     */
    public static listBusinessPlans(): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/finance/business/plans',
        });
    }
    /**
     * Update business plan
     * @returns any OK
     * @throws ApiError
     */
    public static updateBusinessPlan(): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/finance/business/update',
        });
    }
    /**
     * Delete business plan
     * @returns any OK
     * @throws ApiError
     */
    public static deleteBusinessPlan(): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/finance/business/delete',
        });
    }
}
