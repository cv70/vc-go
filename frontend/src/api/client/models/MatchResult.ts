/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { MatchReason } from './MatchReason';
export type MatchResult = {
    target_id: string;
    score: number;
    reasons: Array<MatchReason>;
    risks?: Array<string>;
};

