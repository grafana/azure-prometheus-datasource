// force timezone to UTC to allow tests to work regardless of local timezone
// generally used by snapshots, but can affect specific tests
process.env.TZ = 'UTC';
const { grafanaESModules, nodeModulesToTransform } = require('./.config/jest/utils');

module.exports = {
  // Jest configuration provided by @grafana/create-plugin
  ...require('./.config/jest.config'),
  // d3-color >= 3.1.1 ships an .mjs entry point, which the scaffolded transform regex (js/ts only) skips.
  transform: {
    '^.+\\.(t|j)sx?$|^.+\\.mjs$': require('./.config/jest.config').transform['^.+\\.(t|j)sx?$'],
  },
  // Inform Jest to only transform specific node_module packages.
  transformIgnorePatterns: [nodeModulesToTransform([...grafanaESModules, 'monaco-promql', '@marcbachmann/cel-js'])],
};
