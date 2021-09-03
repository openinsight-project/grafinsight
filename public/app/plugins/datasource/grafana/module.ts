import { DataSourcePlugin } from '@grafinsight/data';
import { GrafInsightDatasource } from './datasource';
import { QueryEditor } from './components/QueryEditor';
import { GrafInsightQuery } from './types';
import { GrafInsightAnnotationsQueryCtrl } from './annotation_ctrl';

export const plugin = new DataSourcePlugin<GrafInsightDatasource, GrafInsightQuery>(GrafInsightDatasource)
  .setQueryEditor(QueryEditor)
  .setAnnotationQueryCtrl(GrafInsightAnnotationsQueryCtrl);
