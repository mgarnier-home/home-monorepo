import { FluxTableMetaData, InfluxDB } from '@influxdata/influxdb-client';

import { config } from '../utils/config';

export abstract class Database<H, L> {
  private _client = new InfluxDB({ url: config.dbHost, token: config.dbToken });

  protected getWriteApi(host: string) {
    return this._client.getWriteApi(config.dbOrg, config.dbBucket).useDefaultTags({ host });
  }

  protected getQueryApi() {
    return this._client.getQueryApi(config.dbOrg);
  }

  protected async executeQuery<T>(query: string): Promise<T[]> {
    const queryApi = this.getQueryApi();

    return new Promise((resolve, reject) => {
      const result: T[] = [];
      queryApi.queryRows(query, {
        next(row: any, tableMeta: FluxTableMetaData) {
          const o = tableMeta.toObject(row);
          result.push(o as any);
        },
        error(error: any) {
          reject(error);
        },
        complete() {
          resolve(result);
        },
      });
    });
  }

  protected abstract getQuery(host: string, rangeStart: string, every: string): string;
  protected abstract getDatas(query: string): Promise<H[]>;

  public abstract logLoad(host: string, data: L): Promise<void>;

  public getDataLast15Min = (host: string) => this.getDatas(this.getQuery(host, '-15m', '15s'));

  public getDataLast1Hour = (host: string) => this.getDatas(this.getQuery(host, '-1h', '1m'));

  public getDataLast6Hour = (host: string) => this.getDatas(this.getQuery(host, '-6h', '6m'));

  public getDataLast1Day = (host: string) => this.getDatas(this.getQuery(host, '-1d', '15m'));

  public getDataLast1Week = (host: string) => this.getDatas(this.getQuery(host, '-7d', '2h'));
}
