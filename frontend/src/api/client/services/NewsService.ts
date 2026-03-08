/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class NewsService {
    /**
     * Get news
     * @returns any OK
     * @throws ApiError
     */
    public static getNews(): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/news/get',
        });
    }
    /**
     * Search news
     * @returns any OK
     * @throws ApiError
     */
    public static searchNews(): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/news/search',
        });
    }
    /**
     * Add news item
     * @returns any OK
     * @throws ApiError
     */
    public static addNews(): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/news/add',
        });
    }
}
