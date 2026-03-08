/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { ResponseInvestor } from '../models/ResponseInvestor';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class InvestorService {
    /**
     * Register investor profile
     * @returns ResponseInvestor OK
     * @throws ApiError
     */
    public static registerInvestor(): CancelablePromise<ResponseInvestor> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/investor/register',
        });
    }
    /**
     * Search investors
     * @returns any OK
     * @throws ApiError
     */
    public static searchInvestors(): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/investor/search',
        });
    }
}
