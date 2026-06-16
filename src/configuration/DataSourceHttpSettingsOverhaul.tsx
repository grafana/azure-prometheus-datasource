import { t, Trans } from '@grafana/i18n';
import { Auth, AuthMethod, ConnectionSettings, convertLegacyAuthProps } from '@grafana/plugin-ui';
import { overhaulStyles } from '@grafana/prometheus';
import { Alert, SecureSocksProxySettings, TextLink, useTheme2 } from '@grafana/ui';
import React, { ReactElement, useState } from 'react';
import { useEffectOnce } from 'react-use';

import { AzurePromDataSourceSettings } from './AzureCredentialsConfig';

type Props = {
  options: AzurePromDataSourceSettings;
  onOptionsChange: (options: AzurePromDataSourceSettings) => void;
  azureAuthEditor: React.ReactNode;
  secureSocksDSProxyEnabled: boolean;
};

export const DataSourceHttpSettingsOverhaul = (props: Props) => {
  const { options, onOptionsChange, azureAuthEditor, secureSocksDSProxyEnabled } = props;

  const newAuthProps = convertLegacyAuthProps({
    config: options,
    onChange: onOptionsChange,
  });

  useEffectOnce(() => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
      },
    });
  });

  const theme = useTheme2();
  const styles = overhaulStyles(theme);

  const [hasPromTypeMig] = useState<boolean>(options.jsonData['prometheus-type-migration'] || false);

  let customMethods: CustomMethod[] = [];

  const azureAuthId = 'custom-azureAuthId';

  const azureAuthOption: CustomMethod = {
    id: azureAuthId,
    label: t('configuration.data-source-http-settings-overhaul.azure-auth-option.label.azure-auth', 'Azure auth'),
    description: t(
      'configuration.data-source-http-settings-overhaul.azure-auth-option.description.authenticate-with-azure',
      'Authenticate with Azure'
    ),
    component: <>{azureAuthEditor}</>,
  };

  customMethods.push(azureAuthOption);

  function returnSelectedMethod(): `custom-${string}` | AuthMethod {
    return azureAuthId;
  }

  // Do we need this switch anymore? Update the language.
  let urlTooltip;
  switch (options.access) {
    case 'direct':
      urlTooltip = (
        <>
          <Trans i18nKey="configuration.data-source-http-settings-overhaul.tooltip-direct">
            Your access method is <em>Browser</em>, this means the URL needs to be accessible from the browser.
          </Trans>
        </>
      );
      break;
    case 'proxy':
      urlTooltip = (
        <>
          <Trans i18nKey="configuration.data-source-http-settings-overhaul.tooltip-proxy">
            Your access method is <em>Server</em>, this means the URL needs to be accessible from the grafana
            backend/server.
          </Trans>
        </>
      );
      break;
    default:
      urlTooltip = (
        <>
          <Trans i18nKey="configuration.data-source-http-settings-overhaul.specify-complete-example-httpyourserver">
            Specify a complete HTTP URL (for example http://your_server:8080){' '}
          </Trans>
        </>
      );
  }

  return (
    <>
      <ConnectionSettings
        urlPlaceholder="http://localhost:9090"
        config={options}
        onChange={onOptionsChange}
        urlLabel="Prometheus server URL"
        urlTooltip={urlTooltip}
      />
      <hr className={`${styles.hrTopSpace} ${styles.hrBottomSpace}`} />
      {hasPromTypeMig && (
        <Alert
          severity="warning"
          title={t('components.logs-query-builder.title-prometheus-migration-occurred', 'Data source migrated')}
        >
          <Trans i18nKey="components.logs-query-builder.prometheus-migration-occurred">
            This data source has been migrated from Prometheus to Azure Monitor Managed Service for Prometheus. Refer to{' '}
            <TextLink
              href="https://grafana.com/docs/grafana-cloud/connect-externally-hosted/data-sources/prometheus/configure/azure-authentication/"
              external
            >
              Connect to Azure Monitor Managed Service for Prometheus
            </TextLink>{' '}
            for more information.
          </Trans>
        </Alert>
      )}
      <Auth
        {...newAuthProps}
        customMethods={customMethods}
        onAuthMethodSelect={(method) => {
          onOptionsChange({
            ...options,
            basicAuth: method === AuthMethod.BasicAuth,
            withCredentials: method === AuthMethod.CrossSiteCredentials,
            jsonData: {
              ...options.jsonData,
              azureCredentials: method === azureAuthId ? options.jsonData.azureCredentials : undefined,
              oauthPassThru: method === AuthMethod.OAuthForward,
            },
          });
        }}
        selectedMethod={returnSelectedMethod()}
        visibleMethods={[azureAuthId]}
      />
      <div className={styles.sectionBottomPadding} />
      {secureSocksDSProxyEnabled && (
        <>
          <SecureSocksProxySettings options={options} onOptionsChange={onOptionsChange} />
          <div className={styles.sectionBottomPadding} />
        </>
      )}
    </>
  );
};

export type CustomMethodId = `custom-${string}`;

export type CustomMethod = {
  id: CustomMethodId;
  label: string;
  description: string;
  component: ReactElement;
};
