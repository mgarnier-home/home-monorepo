import { HwRam } from 'nodesight-types';

import { Point } from '@influxdata/influxdb-client';

import { config } from '../utils/config.js';
import { Database } from './database.js';

type RamQueryResult = {
  _time: string;
  _measurement: string;
  host: string;
  load: number;
};

class DatabaseRam extends Database<HwRam.History.Value, HwRam.Load> {
  constructor() {
    super();
  }

  public async logLoad(host: string, data: HwRam.Load, timestamp: Date = new Date()) {
    try {
      const writeApi = this.getWriteApi(host);

      writeApi.writePoint(
        new Point('ram') //
          .intField('load', data.load)
          .timestamp(timestamp)
      );

      await writeApi.close();

      console.log(`Logged RAM load for ${host} at ${timestamp.toISOString()}`);
    } catch (e) {
      console.error(e);
    }
  }

  override getQuery(host: string, rangeStart: string, every: string): string {
    return `
    from(bucket: "${config.dbBucket}")
      |> range(start: ${rangeStart})
      |> filter(fn: (r) => r._measurement == "ram" and r.host == "${host}")
      |> aggregateWindow(every: ${every}, fn: mean, createEmpty: true)
      |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
      |> yield(name: "mean")
    `;
  }

  override async getDatas(query: string): Promise<HwRam.History.Value[]> {
    const data = (await this.executeQuery<RamQueryResult>(query)).filter((d) => !(d.load === null));

    return data.map((d) => ({
      timestamp: new Date(d._time).getTime(),
      load: d.load,
    }));
  }
}

export const databaseRam = new DatabaseRam();
