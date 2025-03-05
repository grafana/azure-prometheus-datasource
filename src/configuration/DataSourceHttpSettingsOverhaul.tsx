import { Auth, AuthMethod, ConnectionSettings, convertLegacyAuthProps } from '@grafana/plugin-ui';
import { overhaulStyles } from '@grafana/prometheus';
import { SecureSocksProxySettings, useTheme2 } from '@grafana/ui';
import React, { ReactElement, useState } from 'react';
import { useEffectOnce } from 'react-use';

import { AzurePromDataSourceSettings } from './AzureCredentialsConfig';
import { AzureAuthSettings } from './types';


type Props = {
  options: AzurePromDataSourceSettings;
  onOptionsChange: (options: AzurePromDataSourceSettings) => void;
  azureAuthSettings: AzureAuthSettings;
  secureSocksDSProxyEnabled: boolean;
};

export const DataSourceHttpSettingsOverhaul = (props: Props) => {
  const { options,
    onOptionsChange, 
    azureAuthSettings,
    secureSocksDSProxyEnabled
  } = props;

  const newAuthProps = convertLegacyAuthProps({
    config: options,
    onChange: onOptionsChange,
  });

  useEffectOnce(() => {
    // Since we are not allowing users to select another auth,
    // need to update sigV4Auth field to true for auth to work.
    onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
      },
    });
  });

  const theme = useTheme2();
  const styles = overhaulStyles(theme);

  // for custom auth methods sigV4
  let customMethods: CustomMethod[] = [];


  const azureAuthEnabled: boolean =
    (azureAuthSettings?.azureAuthSupported && azureAuthSettings.getAzureAuthEnabled(options)) || false;

  
  const [azureAuthSelected, setAzureAuthSelected] = useState<boolean>(azureAuthEnabled);

  const azureAuthId = 'custom-azureAuthId';


  const azureAuthOption: CustomMethod = {
    id: azureAuthId,
    label: 'Azure auth',
    description: 'Authenticate with Azure',
    component: (
      <>
        {azureAuthSettings.azureSettingsUI && (
          <azureAuthSettings.azureSettingsUI dataSourceConfig={options} onChange={onOptionsChange} />
        )}
      </>
    ),
  };

  // allow the option to show in the dropdown
  if (azureAuthSettings?.azureAuthSupported) {
    customMethods.push(azureAuthOption);
  }

  function returnSelectedMethod() {

    if (azureAuthSelected) {
      return azureAuthId;
    }

    return newAuthProps.selectedMethod;
  }

  // Do we need this switch anymore? Update the language.
  let urlTooltip;
  switch (options.access) {
    case 'direct':
      urlTooltip = (
        <>
          Your access method is <em>Browser</em>, this means the URL needs to be accessible from the browser.

        </>
      );
      break;
    case 'proxy':
      urlTooltip = (
        <>
          Your access method is <em>Server</em>, this means the URL needs to be accessible from the grafana
          backend/server.
        </>
      );
      break;
    default:
      urlTooltip = <>Specify a complete HTTP URL (for example http://your_server:8080) </>;
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
      <Auth
        {...newAuthProps}
        customMethods={customMethods}
        onAuthMethodSelect={(method) => {

          // Azure
          if (azureAuthSettings?.azureAuthSupported) {
            setAzureAuthSelected(method === azureAuthId);
            azureAuthSettings.setAzureAuthEnabled(options, method === azureAuthId);
          }

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
        //DataSourceSettings<AzurePromDataSourceOptions, AzureDataSourceSecureJsonData>.jsonData: AzurePromDataSourceOptions
        //DataSourceSettings<AzurePromDataSourceOptions, AzureDataSourceSecureJsonData>.jsonData: AzurePromDataSourceOptions
        selectedMethod={returnSelectedMethod()}
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
