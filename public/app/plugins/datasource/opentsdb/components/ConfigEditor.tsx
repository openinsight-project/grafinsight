import React from 'react';
import { DataSourceHttpSettings } from '@grafinsight/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafinsight/data';
import { OpenTsdbDetails } from './OpenTsdbDetails';
import { OpenTsdbOptions } from '../types';

export const ConfigEditor = (props: DataSourcePluginOptionsEditorProps<OpenTsdbOptions>) => {
  const { options, onOptionsChange } = props;

  return (
    <>
      <DataSourceHttpSettings
        defaultUrl="http://localhost:4242"
        dataSourceConfig={options}
        onChange={onOptionsChange}
      />
      <OpenTsdbDetails value={options} onChange={onOptionsChange} />
    </>
  );
};
