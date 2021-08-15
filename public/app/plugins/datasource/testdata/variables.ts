import { StandardVariableQuery, StandardVariableSupport } from '@grafinsight/data';

import { TestDataDataSource } from './datasource';
import { TestDataQuery } from './types';

export class TestDataVariableSupport extends StandardVariableSupport<TestDataDataSource> {
  toDataQuery(query: StandardVariableQuery): TestDataQuery {
    return {
      refId: 'TestDataDataSource-QueryVariable',
      stringInput: query.query,
      scenarioId: 'variables-query',
      csvWave: null,
      points: [],
    };
  }
}
