import { DashboardDatasource } from './datasource';
import { DataSourcePlugin } from '@grafinsight/data';

export const plugin = new DataSourcePlugin(DashboardDatasource);
