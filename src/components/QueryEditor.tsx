import React, { ChangeEvent } from 'react';
import { InlineField, Input, Stack } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../datasource';
import { MyDataSourceOptions, MyQuery } from '../types';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export function QueryEditor({ query, onChange, onRunQuery }: Props) {
  
  const onTimeoutChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, timeout: event.target.valueAsNumber });
  };

  const onQueryTextChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, queryText: event.target.value });
  };

  const { queryText, timeout } = query;

  return (
    <Stack gap={0}>
      <InlineField label="Timeout" labelWidth={16} tooltip="Timeout of each run of this query (hint: the lower this value, the lowwer computational time will consumed when it's unused)">
        <Input onChange={onTimeoutChange} onBlur={onRunQuery} value={timeout || ''} type="number" />
      </InlineField>
      <InlineField label="Query Text" labelWidth={16} tooltip="The KsqlDB query to run">
        <Input
          id="query-text"
          onChange={onQueryTextChange}
          value={queryText || ''}
          required
          placeholder="Enter a query"
        />
      </InlineField>
    </Stack>
  );
}
