# Contributing

See [README.md](README.md) for how to build and run the plugin locally (backend via Mage,
frontend via `yarn`).

## Data Source Configuration Schema

`pkg/schema/dsconfig.json` is the **single source of truth** for the data source's
configuration surface — every field a user can set, where it is stored (`root`, `jsonData`,
`secureJsonData`), its type, validation rules and UI hints. It is consumed by provisioning
tooling, documentation and automation.

The schema format is defined and documented by [`grafana/dsconfig`](https://github.com/grafana/dsconfig/tree/main/dsconfig):

- [README](https://github.com/grafana/dsconfig/tree/main/dsconfig#readme) — concepts and a worked example for each field shape (root / jsonData / secret / array / virtual), plus current gaps and limitations.
- [`schema.md`](https://github.com/grafana/dsconfig/blob/main/dsconfig/schema.md) — full property reference.
- [`schema.json`](https://github.com/grafana/dsconfig/blob/main/dsconfig/schema.json) — the JSON Schema `dsconfig.json` validates against. It is pinned via the `$schema` key at the top of our file, so editors autocomplete from it; bump that URL when you bump `github.com/grafana/dsconfig/schema` in `go.mod`.

The rest of this section covers only what is specific to this repository.

### Layout

| File in `pkg/schema/`        | Description                                                                          |
| ----------------------------- | ------------------------------------------------------------------------------------ |
| `dsconfig.json`               | Source of truth — **edit this**                                                      |
| `models.go`                   | `AzurePromOptions` — the typed settings model the schema is checked against          |
| `dsconfig_test.go`            | Wires the schema into the shared conformance suite; also holds `SecureKeys` and the provisioning examples shipped with the plugin |
| `*.gen.json`                  | Generated artifacts — **never hand-edit**; `npm run build` copies them into `dist/schema/` via `webpack.config.ts` |

### This plugin's settings model is split across two repos

Most of this plugin's `jsonData` surface (`httpMethod`, TLS settings, `keepCookies`,
`seriesLimit`, exemplars, alerting, …) isn't rendered by this repo at all — it comes from the
shared `PromSettings` / `AlertingSettingsOverhaul` components in `@grafana/prometheus` and the
`Auth` / `AdvancedHttpSettings` components in `@grafana/plugin-ui`, and is parsed on the backend
by `PromOptions` in `github.com/grafana/grafana-prometheus-datasource/pkg/promlib/models`
(an external dependency pinned in `go.mod`, not vendored here). Only the Azure-specific fields —
`azureCredentials`, `azureEndpointResourceId`, `prometheus-type-migration` — are genuinely local
to this plugin, modeled by `AzurePromOptions` in `pkg/schema/models.go`, which embeds `PromOptions`
and adds those three fields.

This means **where you add a field depends on what it is**:

- **A field shared with vanilla Prometheus** (anything the `@grafana/prometheus` / `@grafana/plugin-ui` components already render) must be added upstream first: open a PR against
  [`grafana/grafana-prometheus-datasource`](https://github.com/grafana/grafana-prometheus-datasource) adding the field to `PromOptions`
  (see [its own `CONTRIBUTING.md`](https://github.com/grafana/grafana-prometheus-datasource/blob/main/CONTRIBUTING.md#data-source-configuration-schema) and [PR #213](https://github.com/grafana/grafana-prometheus-datasource/pull/213) for the pattern). Once that's released, bump the dependency here:

  ```bash
  go get github.com/grafana/grafana-prometheus-datasource/pkg/promlib@latest
  go mod tidy
  ```

  Then declare the field in `dsconfig.json` as described below — `AzurePromOptions` picks it
  up automatically through the embedded `PromOptions` struct, no Go change needed in this repo.

- **A field specific to this plugin** (Azure auth, the migration banner, anything only
  `src/configuration/` renders) is added entirely in this repo: add it to `AzurePromOptions` in
  `pkg/schema/models.go` as well as `dsconfig.json`.

### Adding a new settings option

1. **Declare the field** in `pkg/schema/dsconfig.json` under `fields`, and add its `id` to
   the appropriate `groups[].fieldRefs` entry. Field ids follow the `<target>_<key>`
   convention, e.g. `jsonData_azureEndpointResourceId`.
2. **Add the matching Go field**, per the split above: either bump `promlib` (shared field) or
   add it to `AzurePromOptions` in `pkg/schema/models.go` (Azure-specific field), with a json tag
   equal to the schema `key`. This parity is enforced in both directions — a field in the
   schema but not the struct (or vice versa) fails the test suite. Secrets
   (`target: secureJsonData`) are the exception: they get no struct field, but their key
   must be added to `SecureKeys` in `pkg/schema/dsconfig_test.go`.
3. **Regenerate the artifacts** and commit them with your change:

   ```bash
   go generate ./pkg/schema/...
   ```

4. **Verify**:

   ```bash
   go test ./pkg/schema/...
   ```

If you add a setting that changes what a typical configuration looks like, update the
`SettingsExamples` map in `pkg/schema/dsconfig_test.go` too — those are the provisioning
payloads shipped with the plugin. Use placeholders like `REPLACE_WITH_CLIENT_SECRET`, never
real credentials.

### When the conformance suite fails

Most failures are self-explanatory from the assertion message. The three you are most
likely to hit:

- `SchemaArtifactInSync` — a `.gen.json` file has drifted. Run `go generate ./pkg/schema/...` and commit the result.
- `JSONDataMatchesStruct` / `JSONDataTypesMatchStruct` — the schema and `AzurePromOptions` disagree on keys or types. Update whichever side is behind — remembering that most keys come from the embedded `PromOptions`, not this repo.
- `SecureValuesMatchLoadSettings` — the schema's `secureJsonData` fields and `SecureKeys` disagree.
