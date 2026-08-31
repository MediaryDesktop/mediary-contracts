import type { components, operations } from './cloud.gen.js'

/**
 * Named shortcuts into the generated tree — reaching into openapi-typescript's
 * deep `components` object at every call site is noisy. The only hand-written
 * types in the package; purely structural, adding no new information.
 */

/** Every schema declared in the cloud contract, by name. */
export type CloudSchemas = components['schemas']

/** One operation, by its OpenAPI operationId (`getHealth`, `pingCatalog`, …). */
export type CloudOperation<Id extends keyof operations> = operations[Id]

/** The JSON request body of an operation, if it has one. */
export type CloudRequestBody<Id extends keyof operations> =
	operations[Id] extends { requestBody?: { content: { 'application/json': infer Body } } }
		? Body
		: never

/** The 200-response JSON body of an operation. */
export type CloudResponse<Id extends keyof operations> =
	operations[Id]['responses'] extends { 200: { content: { 'application/json': infer Body } } }
		? Body
		: never
