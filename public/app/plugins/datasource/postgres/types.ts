import { MetricFindValue } from '@grafinsight/data';

export interface PostgresQueryForInterpolation {
  alias?: any;
  format?: any;
  rawSql?: any;
  refId?: any;
  hide?: any;
}

export interface PostgresMetricFindValue extends MetricFindValue {
  value?: string;
}
