import {
  AzureCredentials,
  AzureDataSourceJsonData,
  AzureDataSourceSecureJsonData,
  AzureDataSourceSettings,
  AzureWorkloadIdentityCredentials,
  getAzureClouds,
  getDatasourceCredentials,
  getDefaultAzureCloud,
  updateDatasourceCredentials,
} from '@grafana/azure-sdk';
import { DataSourceSettings, SelectableValue } from '@grafana/data';
import { PromOptions } from '@grafana/prometheus';
import { config } from '@grafana/runtime';

// The SDK's workload identity credentials only carry the auth type and drop any tenant/client ID on
// load and save. The backend accepts them in jsonData.azureCredentials, so they are preserved here.
export type AzureWorkloadIdentityCredentialsWithIds = AzureWorkloadIdentityCredentials & {
  tenantId?: string;
  clientId?: string;
};

export const getWorkloadIdentityInMemoryEnabled = (): boolean =>
  !!(config.featureToggles as Record<string, boolean | undefined>)['workloadIdentityInMemory'];

export function getAzureCloudOptions(): Array<SelectableValue<string>> {
  const cloudInfo = getAzureClouds();

  return cloudInfo.map((cloud) => ({
    value: cloud.name,
    label: cloud.displayName,
  }));
}

export function getDefaultCredentials(): AzureCredentials {
  if (config.azure?.userIdentityEnabled) {
    return { authType: 'currentuser' };
  }
  return { authType: 'clientsecret', azureCloud: getDefaultAzureCloud() };
}

export function getCredentials(options: AzureDataSourceSettings): AzureCredentials {
  const credentials = getDatasourceCredentials(options);
  if (credentials) {
    return credentials.authType === 'workloadidentity' && getWorkloadIdentityInMemoryEnabled()
      ? withWorkloadIdentityIds(credentials, options.jsonData.azureCredentials)
      : credentials;
  }

  return getDefaultCredentials();
}

export function updateCredentials(
  options: AzurePromDataSourceSettings,
  credentials: AzureCredentials
): AzurePromDataSourceSettings {
  const updated = updateDatasourceCredentials(options, credentials);
  if (credentials.authType !== 'workloadidentity' || !getWorkloadIdentityInMemoryEnabled()) {
    return updated;
  }

  const { tenantId, clientId } = credentials as AzureWorkloadIdentityCredentialsWithIds;
  return {
    ...updated,
    jsonData: {
      ...updated.jsonData,
      azureCredentials: {
        ...updated.jsonData.azureCredentials,
        ...(tenantId && { tenantId }),
        ...(clientId && { clientId }),
      } as AzureCredentials,
    },
  };
}

function withWorkloadIdentityIds(credentials: AzureCredentials, saved: unknown): AzureCredentials {
  const { tenantId, clientId } = (saved ?? {}) as AzureWorkloadIdentityCredentialsWithIds;
  return { ...credentials, ...(tenantId && { tenantId }), ...(clientId && { clientId }) } as AzureCredentials;
}

export function setDefaultCredentials(options: AzurePromDataSourceSettings): AzurePromDataSourceSettings {
  return {
    ...options,
    jsonData: {
      ...options.jsonData,
      azureCredentials: getDefaultCredentials(),
    },
  };
}

export function resetCredentials(options: AzurePromDataSourceSettings): Partial<AzurePromDataSourceSettings> {
  return {
    jsonData: {
      ...options.jsonData,
      azureCredentials: undefined,
      azureEndpointResourceId: undefined,
    },
  };
}

export interface AzurePromDataSourceOptions extends PromOptions, AzureDataSourceJsonData {
  azureEndpointResourceId?: string;
  'prometheus-type-migration'?: boolean;
}

export type AzurePromDataSourceSettings = DataSourceSettings<AzurePromDataSourceOptions, AzureDataSourceSecureJsonData>;
