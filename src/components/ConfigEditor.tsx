import React, { ChangeEvent } from 'react';
import { InlineField, Input, Switch, SecretInput } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { MyDataSourceOptions, MySecureJsonData } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions, MySecureJsonData> {}

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;
  const { jsonData, secureJsonFields } = options;
  // const { jsonData, secureJsonFields, secureJsonData } = options;

  const onPathChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...jsonData,
        ksqlserver: event.target.value,
      },
    });
  };

  const onUserChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...jsonData,
        username: event.target.value,
      },
    });
  };

  // Function to handle changes for the 'http' toggle switch
  const onHTTPChange = (event: React.SyntheticEvent<HTMLInputElement>) => {
    // The Switch component's onChange event provides the checked state directly
    onOptionsChange({
      ...options,
      jsonData: {
        ...jsonData,
        http: event.currentTarget.checked, // Use event.currentTarget.checked for Switch
      },
    });
  };

  // Secure field (only sent to the backend)
  const onPassChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        password: event.target.value,
      },
    });
  };

  const onResetPass = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        password: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        password: '',
      },
    });
  };

  return (
    <>
      <InlineField label="KsqlDB" labelWidth={14} interactive tooltip={'The URL to your KsqlDB server'}>
        <Input
          id="ksqlserver"
          onChange={onPathChange}
          value={jsonData.ksqlserver}
          placeholder="http://localhost:8008"
          width={40}
        />
      </InlineField>
      <InlineField label="Use HTTP (instead of HTTPS)" tooltip="Toggle to use HTTP instead of HTTPS for connection">
        <Switch
          id="http"
          value={jsonData.http || false} // Provide a default false if http is undefined
          onChange={onHTTPChange}
        />
      </InlineField>
      <InlineField label="Username" labelWidth={14} interactive tooltip={'Optional KsqlDB server username for authentication'}>
        <Input
          id="username"
          onChange={onUserChange}
          value={jsonData.username || ""}
          placeholder="Auth username"
          width={20}
        />
      </InlineField>
      <InlineField label="Password" tooltip="Optional password for KSQLDB authentication">
        <SecretInput
          isConfigured={Boolean(secureJsonFields.password)} // secureJsonFields.password is a boolean indicating if a value is stored
          id="password"
          onChange={onPassChange}
          onReset={onResetPass}
          // The value for SecretInput is not read directly from secureJsonData.password
          // Instead, you check secureJsonFields.password which is a boolean
          // indicating if a secure value is set (true) or cleared (false).
          // If a value is set, the input will show 'configured' or similar.
          // If you *do* want to display the *last entered value* (not recommended for security),
          // you'd access secureJsonData.password, but that's typically not done.
          // For an optional field, you don't need a default value here as it's handled by isConfi
          placeholder="Auth password"
          width={20}
        />
      </InlineField>
    </>
  );
}
