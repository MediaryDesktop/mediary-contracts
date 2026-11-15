/**
 * Regenerate every derived file in this repository.
 *
 *   openapi/cloud.openapi.yaml ──oapi-codegen──▶ go/cloudclient/client.gen.go
 *
 * The spec is authored by huma in mediary-cloud and copied here by
 * ../scripts/sync-contracts.mjs — nothing here is hand-edited. The desktop
 * generates its own TypeScript types from the same spec file.
 *
 *   node scripts/generate.mjs            regenerate in place
 *   node scripts/generate.mjs --check    fail if the checked-in output is stale
 */
import { spawnSync } from 'node:child_process'
import { existsSync, readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

const ROOT = resolve(dirname(fileURLToPath(import.meta.url)), '..')

const SPEC = resolve(ROOT, 'openapi/cloud.openapi.yaml')
const GO_OUT = resolve(ROOT, 'go/cloudclient/client.gen.go')

const OUTPUTS = [GO_OUT]

const checkOnly = process.argv.includes('--check')

if (!existsSync(SPEC)) {
	console.error(
		[
			`[generate] missing ${SPEC}`,
			'',
			'The spec is produced by mediary-cloud, not written here. Run:',
			'  node ../scripts/sync-contracts.mjs'
		].join('\n')
	)
	process.exit(1)
}

// Snapshot outputs so --check can tell whether generation changed anything.
const before = new Map(
	OUTPUTS.map(file => [
		file,
		existsSync(file) ? readFileSync(file, 'utf8') : null
	])
)

function run(command, args) {
	console.log(`$ ${command} ${args.join(' ')}`)
	const result = spawnSync(command, args, {
		cwd: ROOT,
		stdio: 'inherit',
		shell: true
	})

	if (result.status !== 0) {
		console.error(
			`\n[generate] "${command}" failed with code ${result.status}`
		)
		if (command === 'oapi-codegen') {
			console.error(
				'Install it with: go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest'
			)
		}
		process.exit(result.status ?? 1)
	}
}

// Go client for the self-hosted node.
run('oapi-codegen', [
	'-config',
	'oapi-codegen.cloud.yaml',
	'openapi/cloud.openapi.yaml'
])

if (!checkOnly) {
	console.log('\n[generate] done')
	process.exit(0)
}

const stale = OUTPUTS.filter(file => {
	const now = existsSync(file) ? readFileSync(file, 'utf8') : null
	return now !== before.get(file)
})

if (stale.length > 0) {
	console.error(
		[
			'',
			'[generate] checked-in output is stale:',
			...stale.map(file => `  ${file.replace(ROOT, '.')}`),
			'',
			'Run `npm run generate` and commit the result.'
		].join('\n')
	)
	process.exit(1)
}

console.log('\n[generate] generated output is current')
