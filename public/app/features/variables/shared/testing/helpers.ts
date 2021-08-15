import { DataSourceInstanceSettings, DataSourceJsonData, DataSourcePluginMeta } from '@grafinsight/data';

export function getDataSourceInstanceSetting(name: string, meta: DataSourcePluginMeta): DataSourceInstanceSettings {
  return {
    id: 1,
    uid: '',
    type: '',
    name,
    meta,
    jsonData: ({} as unknown) as DataSourceJsonData,
  };
}
