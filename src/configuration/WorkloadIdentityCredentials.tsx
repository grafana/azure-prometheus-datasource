import { cx } from '@emotion/css';
import { AzureCredentials } from '@grafana/azure-sdk';
import { Trans } from '@grafana/i18n';
import { InlineFormLabel, Input } from '@grafana/ui';
import React, { ChangeEvent } from 'react';

import { AzureWorkloadIdentityCredentialsWithIds } from './AzureCredentialsConfig';

export interface Props {
  credentials: AzureWorkloadIdentityCredentialsWithIds;
  onCredentialsChange: (updatedCredentials: AzureCredentials) => void;
  disabled?: boolean;
}

export const WorkloadIdentityCredentials = ({ credentials, onCredentialsChange, disabled }: Props) => {
  const onTenantIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    onCredentialsChange({ ...credentials, tenantId: event.target.value } as AzureCredentials);
  };

  const onClientIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    onCredentialsChange({ ...credentials, clientId: event.target.value } as AzureCredentials);
  };

  return (
    <>
      <div className="gf-form-inline">
        <div className="gf-form">
          <InlineFormLabel className="width-12">
            <Trans i18nKey="configuration.azure-credentials-form.directory-tenant-id">Directory (tenant) ID</Trans>
          </InlineFormLabel>
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
          <InlineFormLabel className="width-12">
            <Trans i18nKey="configuration.azure-credentials-form.application-client-id">Application (client) ID</Trans>
          </InlineFormLabel>
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
    </>
  );
};
