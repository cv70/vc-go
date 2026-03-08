/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type Founder = {
    id: string;
    name: string;
    stage: Founder.stage;
    sectors: Array<string>;
    location?: string | null;
    funding_target_cny?: number | null;
    created_at: string;
    updated_at: string;
};
export namespace Founder {
    export enum stage {
        IDEA = 'idea',
        PRE_SEED = 'pre_seed',
        SEED = 'seed',
        PRE_A = 'pre_a',
        A = 'a',
    }
}

