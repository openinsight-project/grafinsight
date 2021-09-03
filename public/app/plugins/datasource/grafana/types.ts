import { AnnotationQuery, DataQuery } from '@grafinsight/data';
import { MeasurementsQuery } from '@grafinsight/runtime/src';

//----------------------------------------------
// Query
//----------------------------------------------

export enum GrafInsightQueryType {
  RandomWalk = 'randomWalk',
  LiveMeasurements = 'measurements',
}

export interface GrafInsightQuery extends DataQuery {
  queryType: GrafInsightQueryType; // RandomWalk by default
  channel?: string;
  measurements?: MeasurementsQuery;
}

export const defaultQuery: GrafInsightQuery = {
  refId: 'A',
  queryType: GrafInsightQueryType.RandomWalk,
};

//----------------------------------------------
// Annotations
//----------------------------------------------

export enum GrafInsightAnnotationType {
  Dashboard = 'dashboard',
  Tags = 'tags',
}

export interface GrafInsightAnnotationQuery extends AnnotationQuery<GrafInsightQuery> {
  type: GrafInsightAnnotationType; // tags
  limit: number; // 100
  tags?: string[];
  matchAny?: boolean; // By default GrafInsight only shows annotations that match all tags in the query. Enabling this returns annotations that match any of the tags in the query.
}
