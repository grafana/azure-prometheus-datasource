import { AzureCredentials } from '@grafana/azure-sdk';
import { config } from '@grafana/runtime';

import { AzurePromDataSourceSettings, getCredentials, updateCredentials } from './AzureCredentialsConfig';

const workloadIdentityOptions = (azureCredentials: object) =>
  ({
    jsonData: { azureCredentials: { authType: 'workloadidentity', ...azureCredentials } },
    secureJsonData: {},
    secureJsonFields: {},
  }) as unknown as AzurePromDataSourceSettings;

const setToggle = (enabled: boolean) => {
  config.featureToggles = {
    ...config.featureToggles,
    workloadIdentityInMemory: enabled,
  } as typeof config.featureToggles;
};

describe('workload identity tenant and client IDs', () => {
  beforeEach(() => {
    config.azure.workloadIdentityEnabled = true;
    setToggle(true);
  });

  it('should ignore tenantId and clientId when the feature toggle is off', () => {
    setToggle(false);
    const options = workloadIdentityOptions({ tenantId: 'tenant-1', clientId: 'client-1' });

    expect(getCredentials(options)).toEqual({ authType: 'workloadidentity' });

    const updated = updateCredentials(options, {
      authType: 'workloadidentity',
      tenantId: 'tenant-1',
      clientId: 'client-1',
    } as AzureCredentials);
    expect(updated.jsonData.azureCredentials).toEqual({ authType: 'workloadidentity' });
  });

  it('should restore saved tenantId and clientId', () => {
    const options = workloadIdentityOptions({ tenantId: 'tenant-1', clientId: 'client-1' });

    expect(getCredentials(options)).toEqual({
      authType: 'workloadidentity',
      tenantId: 'tenant-1',
      clientId: 'client-1',
    });
  });

  it('should persist tenantId and clientId in azureCredentials', () => {
    const credentials = { authType: 'workloadidentity', tenantId: 'tenant-1', clientId: 'client-1' };

    const updated = updateCredentials(workloadIdentityOptions({}), credentials as AzureCredentials);

    expect(updated.jsonData.azureCredentials).toEqual(credentials);
  });

  it('should omit tenantId and clientId when they are cleared', () => {
    const options = workloadIdentityOptions({ tenantId: 'tenant-1', clientId: 'client-1' });
    const credentials = { authType: 'workloadidentity', tenantId: '', clientId: '' };

    const updated = updateCredentials(options, credentials as AzureCredentials);

    expect(updated.jsonData.azureCredentials).toEqual({ authType: 'workloadidentity' });
  });
});
