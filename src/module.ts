import { DataSourcePlugin } from '@grafana/data';
import { initPluginTranslations } from '@grafana/i18n';
import { PromQueryEditorByApp, PrometheusDatasource, PromCheatSheet } from '@grafana/prometheus';

import { ConfigEditor } from './configuration/ConfigEditor';
import pluginJson from './plugin.json';

initPluginTranslations(pluginJson.id);

export const plugin = new DataSourcePlugin(PrometheusDatasource)
  .setQueryEditor(PromQueryEditorByApp)
  .setConfigEditor(ConfigEditor)
  .setQueryEditorHelp(PromCheatSheet);
