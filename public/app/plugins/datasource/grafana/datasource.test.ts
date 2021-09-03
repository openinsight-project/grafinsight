import { DataSourceInstanceSettings, dateTime, AnnotationQueryRequest } from '@grafinsight/data';

import { backendSrv } from 'app/core/services/backend_srv'; // will use the version in __mocks__
import { GrafInsightDatasource } from './datasource';
import { GrafInsightQuery, GrafInsightAnnotationQuery, GrafInsightAnnotationType } from './types';

jest.mock('@grafinsight/runtime', () => ({
  ...((jest.requireActual('@grafinsight/runtime') as unknown) as object),
  getBackendSrv: () => backendSrv,
  getTemplateSrv: () => ({
    replace: (val: string) => {
      return val.replace('$var2', 'replaced__delimiter__replaced2').replace('$var', 'replaced');
    },
  }),
}));

describe('grafinsight data source', () => {
  const getMock = jest.spyOn(backendSrv, 'get');

  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('when executing an annotations query', () => {
    let calledBackendSrvParams: any;
    let ds: GrafInsightDatasource;
    beforeEach(() => {
      getMock.mockImplementation((url: string, options: any) => {
        calledBackendSrvParams = options;
        return Promise.resolve([]);
      });

      ds = new GrafInsightDatasource({} as DataSourceInstanceSettings);
    });

    describe('with tags that have template variables', () => {
      const options = setupAnnotationQueryOptions({ tags: ['tag1:$var'] });

      beforeEach(() => {
        return ds.annotationQuery(options);
      });

      it('should interpolate template variables in tags in query options', () => {
        expect(calledBackendSrvParams.tags[0]).toBe('tag1:replaced');
      });
    });

    describe('with tags that have multi value template variables', () => {
      const options = setupAnnotationQueryOptions({ tags: ['$var2'] });

      beforeEach(() => {
        return ds.annotationQuery(options);
      });

      it('should interpolate template variables in tags in query options', () => {
        expect(calledBackendSrvParams.tags[0]).toBe('replaced');
        expect(calledBackendSrvParams.tags[1]).toBe('replaced2');
      });
    });

    describe('with type dashboard', () => {
      const options = setupAnnotationQueryOptions(
        {
          type: GrafInsightAnnotationType.Dashboard,
          tags: ['tag1'],
        },
        { id: 1 }
      );

      beforeEach(() => {
        return ds.annotationQuery(options);
      });

      it('should remove tags from query options', () => {
        expect(calledBackendSrvParams.tags).toBe(undefined);
      });
    });
  });
});

function setupAnnotationQueryOptions(annotation: Partial<GrafInsightAnnotationQuery>, dashboard?: { id: number }) {
  return ({
    annotation,
    dashboard,
    range: {
      from: dateTime(1432288354),
      to: dateTime(1432288401),
    },
    rangeRaw: { from: 'now-24h', to: 'now' },
  } as unknown) as AnnotationQueryRequest<GrafInsightQuery>;
}
