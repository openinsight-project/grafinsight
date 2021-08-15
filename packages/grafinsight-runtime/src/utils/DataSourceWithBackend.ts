import {
  DataSourceApi,
  DataQueryRequest,
  DataQueryResponse,
  DataSourceInstanceSettings,
  DataQuery,
  DataSourceJsonData,
  ScopedVars,
} from '@grafinsight/data';
import { Observable, of } from 'rxjs';
import { map, catchError } from 'rxjs/operators';
import { getBackendSrv, getDataSourceSrv } from '../services';
import { toDataQueryResponse } from './queryResponse';

const ExpressionDatasourceID = '__expr__';

/**
 * Describes the current health status of a data source plugin.
 *
 * @public
 */
export enum HealthStatus {
  Unknown = 'UNKNOWN',
  OK = 'OK',
  Error = 'ERROR',
}

/**
 * Describes the payload returned when checking the health of a data source
 * plugin.
 *
 * @public
 */
export interface HealthCheckResult {
  status: HealthStatus;
  message: string;
  details?: Record<string, any>;
}

/**
 * Extend this class to implement a data source plugin that is depending on the GrafInsight
 * backend API.
 *
 * @public
 */
export class DataSourceWithBackend<
  TQuery extends DataQuery = DataQuery,
  TOptions extends DataSourceJsonData = DataSourceJsonData
> extends DataSourceApi<TQuery, TOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<TOptions>) {
    super(instanceSettings);
  }

  /**
   * Ideally final -- any other implementation may not work as expected
   */
  query(request: DataQueryRequest<TQuery>): Observable<DataQueryResponse> {
    const { intervalMs, maxDataPoints, range, requestId } = request;
    let targets = request.targets;

    if (this.filterQuery) {
      targets = targets.filter((q) => this.filterQuery!(q));
    }

    const queries = targets.map((q) => {
      let datasourceId = this.id;

      if (q.datasource === ExpressionDatasourceID) {
        return {
          ...q,
          datasourceId,
        };
      }

      if (q.datasource) {
        const ds = getDataSourceSrv().getInstanceSettings(q.datasource);

        if (!ds) {
          throw new Error('Unknown Datasource: ' + q.datasource);
        }

        datasourceId = ds.id;
      }

      return {
        ...this.applyTemplateVariables(q, request.scopedVars),
        datasourceId,
        intervalMs,
        maxDataPoints,
      };
    });

    // Return early if no queries exist
    if (!queries.length) {
      return of({ data: [] });
    }

    const body: any = { queries };

    if (range) {
      body.range = range;
      body.from = range.from.valueOf().toString();
      body.to = range.to.valueOf().toString();
    }

    return getBackendSrv()
      .fetch({
        url: '/api/ds/query',
        method: 'POST',
        data: body,
        requestId,
      })
      .pipe(
        map((rsp: any) => {
          return toDataQueryResponse(rsp, queries as DataQuery[]);
        }),
        catchError((err) => {
          return of(toDataQueryResponse(err));
        })
      );
  }

  /**
   * Override to skip executing a query
   *
   * @returns false if the query should be skipped
   *
   * @virtual
   */
  filterQuery?(query: TQuery): boolean;

  /**
   * Override to apply template variables.  The result is usually also `TQuery`, but sometimes this can
   * be used to modify the query structure before sending to the backend.
   *
   * NOTE: if you do modify the structure or use template variables, alerting queries may not work
   * as expected
   *
   * @virtual
   */
  applyTemplateVariables(query: TQuery, scopedVars: ScopedVars): Record<string, any> {
    return query;
  }

  /**
   * Make a GET request to the datasource resource path
   */
  async getResource(path: string, params?: any): Promise<any> {
    return getBackendSrv().get(`/api/datasources/${this.id}/resources/${path}`, params);
  }

  /**
   * Send a POST request to the datasource resource path
   */
  async postResource(path: string, body?: any): Promise<any> {
    return getBackendSrv().post(`/api/datasources/${this.id}/resources/${path}`, { ...body });
  }

  /**
   * Run the datasource healthcheck
   */
  async callHealthCheck(): Promise<HealthCheckResult> {
    return getBackendSrv()
      .request({ method: 'GET', url: `/api/datasources/${this.id}/health`, showErrorAlert: false })
      .then((v) => {
        return v as HealthCheckResult;
      })
      .catch((err) => {
        return err.data as HealthCheckResult;
      });
  }

  /**
   * Checks the plugin health
   * see public/app/features/datasources/state/actions.ts for what needs to be returned here
   */
  async testDatasource(): Promise<any> {
    return this.callHealthCheck().then((res) => {
      if (res.status === HealthStatus.OK) {
        return {
          status: 'success',
          message: res.message,
        };
      }
      throw new Error(res.message);
    });
  }
}
