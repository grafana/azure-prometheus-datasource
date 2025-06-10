import { useTranslate, Trans } from '@grafana/i18n';
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
  const { t } = useTranslate();
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
        label: t("configuration.azure-credentials-form.auth-type-options.opts.label.app-registration", "App Registration"),
      },
    ];

    if (managedIdentityEnabled) {
      opts.push({
        value: 'msi',
        label: t("configuration.azure-credentials-form.auth-type-options.label.managed-identity", "Managed Identity"),
      });
    }

    if (workloadIdentityEnabled) {
      opts.push({
        value: 'workloadidentity',
        label: t("configuration.azure-credentials-form.auth-type-options.label.workload-identity", "Workload Identity"),
      });
    }

    if (userIdentityEnabled) {
      opts.unshift({
        value: 'currentuser',
        label: t("configuration.azure-credentials-form.auth-type-options.label.current-user", "Current User"),
      });
    }

    return opts;
  }, [managedIdentityEnabled, workloadIdentityEnabled, userIdentityEnabled, t]);

  const onAuthTypeChange = (selected: SelectableValue<AzureAuthType>) => {
    const defaultAuthType = userIdentityEnabled
      ? 'currentuser'
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
            <InlineFormLabel className="width-12" tooltip={t("configuration.azure-credentials-form.tooltip-choose-authentication-azure-services", "Choose the type of authentication to Azure services")}><Trans i18nKey="configuration.azure-credentials-form.authentication">
              Authentication
            </Trans></InlineFormLabel>
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
                <InlineFormLabel className="width-12" tooltip={t("configuration.azure-credentials-form.tooltip-choose-an-azure-cloud", "Choose an Azure Cloud")}><Trans i18nKey="configuration.azure-credentials-form.azure-cloud">
                  Azure Cloud
                </Trans></InlineFormLabel>
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
              <InlineFormLabel className="width-12"><Trans i18nKey="configuration.azure-credentials-form.directory-tenant-id">Directory (tenant) ID</Trans></InlineFormLabel>
              <div className="width-15">
                <Input
                  className={cx('width-20')}
                  // eslint-disable-next-line @grafana/i18n/no-untranslated-strings
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
              <InlineFormLabel className="width-12"><Trans i18nKey="configuration.azure-credentials-form.application-client-id">Application (client) ID</Trans></InlineFormLabel>
              <div className="width-15">
                <Input
                  className={cx('width-20')}
                  // eslint-disable-next-line @grafana/i18n/no-untranslated-strings
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
                <InlineFormLabel htmlFor="azure-client-secret" className="width-12"><Trans i18nKey="configuration.azure-credentials-form.client-secret">
                  Client Secret
                </Trans></InlineFormLabel>
                <Input id="azure-client-secret" className={cx('width-20')} placeholder={t("configuration.azure-credentials-form.azure-client-secret-placeholder-configured", "configured")} disabled />
              </div>
              {!disabled && (
                <div className="gf-form">
                  <div className={cx('max-width-20 gf-form-inline')}>
                    <Button variant="secondary" type="button" onClick={onClientSecretReset}><Trans i18nKey="configuration.azure-credentials-form.reset">
                      reset
                    </Trans></Button>
                  </div>
                </div>
              )}
            </div>
          ) : (
            <div className="gf-form-inline">
              <div className="gf-form">
                <InlineFormLabel className="width-12"><Trans i18nKey="configuration.azure-credentials-form.client-secret">Client Secret</Trans></InlineFormLabel>
                <div className="width-15">
                  <Input
                    className={cx('width-20')}
                    // eslint-disable-next-line @grafana/i18n/no-untranslated-strings
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
