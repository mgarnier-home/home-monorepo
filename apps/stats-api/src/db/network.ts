import { logger } from 'logger';
import { HwNetwork } from 'nodesight-types';

import { Point } from '@influxdata/influxdb-client';

import { config } from '../utils/config.js';
import { Database } from './database.js';

type NetworkQueryResult = {
  _time: string;
  _measurement: string;
  host: string;
  up: number;
  down: number;
};

class DatabaseNetwork extends Database<HwNetwork.History.Value, HwNetwork.Load> {
  constructor() {
    super();
  }

  public async logLoad(host: string, data: HwNetwork.Load, timestamp: Date = new Date()) {
    try {
      const writeApi = this.getWriteApi(host);

      writeApi.writePoint(
        new Point('network') //
          .floatField('up', data.up)
          .floatField('down', data.down)
          .timestamp(timestamp)
      );

      await writeApi.close();
    } catch (e) {
      logger.error(`Error Network load for ${host}`, e);
    }
  }

  override getQuery(host: string, rangeStart: string, every: string): string {
    return `
    from(bucket: "${config.dbBucket}")
      |> range(start: ${rangeStart})
      |> filter(fn: (r) => r._measurement == "network" and r.host == "${host}")
      |> aggregateWindow(every: ${every}, fn: mean, createEmpty: true)
      |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
      |> yield(name: "mean")
    `;
  }

  override async getDatas(query: string): Promise<HwNetwork.History.Value[]> {
    const data = (await this.executeQuery<NetworkQueryResult>(query)).filter(
      (d) => !(d.down === null && d.up === null)
    );

    return data.map((d) => ({
      timestamp: new Date(d._time).getTime(),
      up: d.up,
      down: d.down,
    }));
  }
}

export const databaseNetwork = new DatabaseNetwork();
