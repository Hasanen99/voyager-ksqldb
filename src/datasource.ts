import { DataSourceInstanceSettings, CoreApp, ScopedVars } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv, getGrafanaLiveSrv} from '@grafana/runtime';

import { MyQuery, MyDataSourceOptions, DEFAULT_QUERY } from './types';

import {
  DataQueryRequest,
  DataQueryResponse,
  LiveChannelScope, // For LiveChannelScope.DataSource
} from '@grafana/data'; // Core Grafana data types

import {
  Observable, // For the Observable return type
  merge,      // For merging multiple observables
} from 'rxjs'; // Reactive Extensions for JavaScript (RxJS)


export class DataSource extends DataSourceWithBackend<MyQuery, MyDataSourceOptions> {

  ksqlserver: string;
  username: string;
  // pass: string;
  http: boolean;

  constructor(instanceSettings: DataSourceInstanceSettings<MyDataSourceOptions>) {
    super(instanceSettings);

    this.ksqlserver = instanceSettings.jsonData.ksqlserver;
    this.http = instanceSettings.jsonData.http;
    this.username = instanceSettings.jsonData.username ?? "";
    // this.user = instanceSettings.jsonData.password;

  }

  getDefaultQuery(_: CoreApp): Partial<MyQuery> {
    return DEFAULT_QUERY;
  }

  applyTemplateVariables(query: MyQuery, scopedVars: ScopedVars) {
    return {
      ...query,
      queryText: getTemplateSrv().replace(query.queryText, scopedVars),
    };
  }

  filterQuery(query: MyQuery): boolean {
    // if no query has been provided, prevent the query from being executed
    return !!query.queryText;
  }
  
  query(request: DataQueryRequest<MyQuery>): Observable<DataQueryResponse> {
    const observables = request.targets.map((query, index) => {

      return getGrafanaLiveSrv().getDataStream({
        addr: {
          scope: LiveChannelScope.DataSource,
          namespace: this.uid,
          path: `ksql/custom-${query.refId}-${query.timeout}`, // this will allow each new query to create a new connection
          data: {
            ...query,
          },
        },
      });
    });

    return merge(...observables);
  }

  // Mandatory: Implement the testDatasource method for validation
  async testDatasource() {
    // Basic validation: Check if URL is provided
    if (!this.ksqlserver) {
      return {
        status: 'error',
        message: 'KSQLDB URL is a mandatory field. Please provide a valid URL.',
      };
    }

    // You can add more complex URL validation here if needed
    // For example, checking if it starts with http(s)://
    if (!this.ksqlserver.startsWith('http://') && !this.ksqlserver.startsWith('https://')) {
        return {
            status: 'error',
            message: 'KSQLDB URL must start with http:// or https://',
        };
    }

    // Now, attempt to connect to your KSQLDB backend using the URL
    // You'd typically make a simple HTTP request to your backend plugin
    // to perform a health check or a basic version query.
    // Example using this.getResource, assuming your backend has a /health endpoint:
    try {
      // You can pass the URL to your backend as part of the test payload if needed,
      // or your backend can just use the URL from its own instance settings.
      const response = await this.getResource('/health'); // Assuming your backend has a /health endpoint
      // You might check response.status or a specific field in the response.
      if (response && response.status === 'ok') { // Assuming your backend /health returns { status: 'ok' }
          return {
              status: 'success',
              message: 'Successfully connected to KSQLDB via backend plugin.',
          };
      } else {
          // Handle specific errors from your backend's health check
          return {
              status: 'error',
              message: response?.message || 'Failed to connect to KSQLDB. Please check the URL.',
          };
      }
    } catch (err) {
      // Catch network errors or errors from the backend HTTP request
      return {
        status: 'error',
        message: `${this.ksqlserver} ${this.http} ${this.username} Failed to connect to KSQLDB: Unknown error: ${err}`,
      };
    }
  }
}
