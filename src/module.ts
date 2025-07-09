import { DataSourcePlugin } from '@grafana/data';
import { initPluginTranslations } from '@grafana/i18n';
import { PromQueryEditorByApp, PrometheusDatasource, PromCheatSheet, loadResources as loadPrometheusResources } from '@grafana/prometheus';

import { ConfigEditor } from './configuration/ConfigEditor';
import pluginJson from './plugin.json';

// top level await is fine in our bundled code, but not in our test environment (yet)
// TODO remove this when our test environment can handle top level await
if (process.env.NODE_ENV === 'test') {
  initPluginTranslations(pluginJson.id, [loadPrometheusResources]);
} else {
  await initPluginTranslations(pluginJson.id, [loadPrometheusResources]);
}

export const plugin = new DataSourcePlugin(PrometheusDatasource)
  .setQueryEditor(PromQueryEditorByApp)
  .setConfigEditor(ConfigEditor)
  .setQueryEditorHelp(PromCheatSheet);
