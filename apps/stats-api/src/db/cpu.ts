import { HwCpu } from 'nodesight-types';

import { Point } from '@influxdata/influxdb-client';

import { config } from '../utils/config.js';
import { Database } from './database.js';

type CpuQueryResult = {
  _time: string;
  _measurement: string;
  host: string;
  core: string;
  load: number;
  temp: number;
};

class DatabaseCpu extends Database<HwCpu.History.Value, HwCpu.Load> {
  constructor() {
    super();
  }

  public async logLoad(host: string, data: HwCpu.Load, timestamp: Date = new Date()) {
    try {
      const writeApi = this.getWriteApi(host);

      const avgLoad = data.cores.reduce((acc, curr) => acc + curr.load, 0) / data.cores.length;
      const avgTemp = data.cores.reduce((acc, curr) => acc + curr.temp, 0) / data.cores.length;

      for (const dataPoint of data.cores) {
        writeApi.writePoint(
          new Point('cpu')
            .tag('core', dataPoint.core.toString())
            .floatField('load', dataPoint.load)
            .floatField('temp', dataPoint.temp ?? 0)
            .timestamp(timestamp)
        );
      }

      writeApi.writePoint(
        new Point('cpu') //
          .tag('core', 'avg')
          .floatField('load', avgLoad)
          .floatField('temp', avgTemp)
          .timestamp(timestamp)
      );

      await writeApi.close();

      console.log(`Logged CPU load for ${host} at ${timestamp.toISOString()}`);
    } catch (e) {
      console.error(e);
    }
  }

  override getQuery(host: string, rangeStart: string, every: string): string {
    return `
    from(bucket: "${config.dbBucket}")
      |> range(start: ${rangeStart})
      |> filter(fn: (r) => r._measurement == "cpu" and r.host == "${host}" and r.core == "avg")
      |> aggregateWindow(every: ${every}, fn: mean, createEmpty: true)
      |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
      |> yield(name: "mean")
    `;
  }

  override async getDatas(query: string): Promise<HwCpu.History.Value[]> {
    const data = (await this.executeQuery<CpuQueryResult>(query)).filter((d) => !(d.temp === null && d.load === null));

    const result: HwCpu.History.Value[] = [];

    for (const point of data) {
      result.push({
        core: point.core === 'avg' ? -1 : parseInt(point.core),
        timestamp: new Date(point._time).getTime(),
        load: point.load,
        temp: point.temp,
      });
    }

    return result;
  }
}

export const databaseCpu = new DatabaseCpu();
