/**
 * TypeScript distribution of the Mediary cloud API contract — re-exports
 * `cloud.gen.ts`, generated from `openapi/cloud.openapi.yaml`. Nothing here
 * is hand-written; see docs/contracts-and-codegen.md.
 */

export type { components, operations, paths, webhooks } from './cloud.gen.js'

export type { CloudSchemas, CloudOperation, CloudRequestBody, CloudResponse } from './helpers.js'

export { CLOUD_CONTRACT_VERSION } from './version.gen.js'
