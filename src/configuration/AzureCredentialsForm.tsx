import { cx } from '@emotion/css';
import { AzureAuthType, AzureCredentials } from '@grafana/azure-sdk';
import { SelectableValue } from '@grafana/data';
import { InlineFormLabel, Button, Select, Input } from '@grafana/ui';
import React, { ChangeEvent, useMemo } from 'react';

import { getAzureCloudOptions } from './AzureCredentialsConfig';
import CurrentUserFallbackCredentials from './CurrentUserFallbackCredentials';

export interface Props {
  managedIdentityEnabled: boolean;
  workloadIdentityEnabled: boolean;
  userIdentityEnabled: boolean;
  credentials: AzureCredentials;
  azureCloudOptions?: SelectableValue[];
  onCredentialsChange: (updatedCredentials: AzureCredentials) => void;
  getSubscriptions?: () => Promise<SelectableValue[]>;
  disabled?: boolean;
}

export const AzureCredentialsForm = (props: Props) => {
  const {
    credentials,
    azureCloudOptions,
    onCredentialsChange,
    disabled,
    managedIdentityEnabled,
    workloadIdentityEnabled,
    userIdentityEnabled,
  } = props;

  const authTypeOptions = useMemo(() => {
    let opts: Array<SelectableValue<AzureAuthType>> = [
      {
        value: 'clientsecret',
        label: 'App Registration',
      },
    ];

    if (managedIdentityEnabled) {
      opts.push({
        value: 'msi',
        label: 'Managed Identity',
      });
    }

    if (workloadIdentityEnabled) {
      opts.push({
        value: 'workloadidentity',
        label: 'Workload Identity',
      });
    }

    if (userIdentityEnabled) {
      opts.unshift({
        value: 'currentuser',
        label: 'Current User',
      });
    }

    return opts;
  }, [managedIdentityEnabled, workloadIdentityEnabled, userIdentityEnabled]);

  const onAuthTypeChange = (selected: SelectableValue<AzureAuthType>) => {
    const defaultAuthType = managedIdentityEnabled
      ? 'msi'
      : workloadIdentityEnabled
        ? 'workloadidentity'
        : 'clientsecret';
    const updated: AzureCredentials = {
      ...credentials,
      authType: selected.value || defaultAuthType,
    };
    onCredentialsChange(updated);
  };

  const onAzureCloudChange = (selected: SelectableValue<string>) => {
    if (credentials.authType === 'clientsecret') {
      const updated: AzureCredentials = {
        ...credentials,
        azureCloud: selected.value,
      };
      onCredentialsChange(updated);
    }
  };

  const onTenantIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    if (credentials.authType === 'clientsecret') {
      const updated: AzureCredentials = {
        ...credentials,
        tenantId: event.target.value,
      };
      onCredentialsChange(updated);
    }
  };

  const onClientIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    if (credentials.authType === 'clientsecret') {
      const updated: AzureCredentials = {
        ...credentials,
        clientId: event.target.value,
      };
      onCredentialsChange(updated);
    }
  };

  const onClientSecretChange = (event: ChangeEvent<HTMLInputElement>) => {
    if (credentials.authType === 'clientsecret') {
      const updated: AzureCredentials = {
        ...credentials,
        clientSecret: event.target.value,
      };
      onCredentialsChange(updated);
    }
  };

  const onClientSecretReset = () => {
    if (credentials.authType === 'clientsecret') {
      const updated: AzureCredentials = {
        ...credentials,
        clientSecret: '',
      };
      onCredentialsChange(updated);
    }
  };

  return (
    <div className="gf-form-group">
      {authTypeOptions.length > 1 && (
        <div className="gf-form-inline">
          <div className="gf-form">
            <InlineFormLabel className="width-12" tooltip="Choose the type of authentication to Azure services">
              Authentication
            </InlineFormLabel>
            <Select
              className="width-15"
              value={authTypeOptions.find((opt) => opt.value === credentials.authType)}
              options={authTypeOptions}
              onChange={onAuthTypeChange}
              isDisabled={disabled}
            />
          </div>
        </div>
      )}
      {credentials.authType === 'clientsecret' && (
        <>
          {azureCloudOptions && (
            <div className="gf-form-inline">
              <div className="gf-form">
                <InlineFormLabel className="width-12" tooltip="Choose an Azure Cloud">
                  Azure Cloud
                </InlineFormLabel>
                <Select
                  className="width-15"
                  value={azureCloudOptions.find((opt) => opt.value === credentials.azureCloud)}
                  options={azureCloudOptions}
                  onChange={onAzureCloudChange}
                  isDisabled={disabled}
                />
              </div>
            </div>
          )}
          <div className="gf-form-inline">
            <div className="gf-form">
              <InlineFormLabel className="width-12">Directory (tenant) ID</InlineFormLabel>
              <div className="width-15">
                <Input
                  className={cx('width-20')}
                  placeholder="XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
                  value={credentials.tenantId || ''}
                  onChange={onTenantIdChange}
                  disabled={disabled}
                />
              </div>
            </div>
          </div>
          <div className="gf-form-inline">
            <div className="gf-form">
              <InlineFormLabel className="width-12">Application (client) ID</InlineFormLabel>
              <div className="width-15">
                <Input
                  className={cx('width-20')}
                  placeholder="XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
                  value={credentials.clientId || ''}
                  onChange={onClientIdChange}
                  disabled={disabled}
                />
              </div>
            </div>
          </div>
          {typeof credentials.clientSecret === 'symbol' ? (
            <div className="gf-form-inline">
              <div className="gf-form">
                <InlineFormLabel htmlFor="azure-client-secret" className="width-12">
                  Client Secret
                </InlineFormLabel>
                <Input id="azure-client-secret" className={cx('width-20')} placeholder="configured" disabled />
              </div>
              {!disabled && (
                <div className="gf-form">
                  <div className={cx('max-width-20 gf-form-inline')}>
                    <Button variant="secondary" type="button" onClick={onClientSecretReset}>
                      reset
                    </Button>
                  </div>
                </div>
              )}
            </div>
          ) : (
            <div className="gf-form-inline">
              <div className="gf-form">
                <InlineFormLabel className="width-12">Client Secret</InlineFormLabel>
                <div className="width-15">
                  <Input
                    className={cx('width-20')}
                    placeholder="XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
                    value={credentials.clientSecret || ''}
                    onChange={onClientSecretChange}
                    disabled={disabled}
                  />
                </div>
              </div>
            </div>
          )}
        </>
      )}
      {credentials.authType === 'currentuser' && (
        <CurrentUserFallbackCredentials
          credentials={credentials}
          azureCloudOptions={getAzureCloudOptions()}
          onCredentialsChange={onCredentialsChange}
          disabled={disabled}
          managedIdentityEnabled={managedIdentityEnabled}
          workloadIdentityEnabled={workloadIdentityEnabled}
        />
      )}
    </div>
  );
};

export default AzureCredentialsForm;
