import { Point } from "@influxdata/influxdb-client";
import { HwGpu } from "@mgarnier11/nodesight-types";

import { config } from "../utils/config.js";
import { Database } from "./database.js";
import { DbUtils } from "./utils.js";

type GpuQueryResult = {
  _time: string;
  _measurement: string;
  host: string;
  index: string;
  model: string;
  load: number;
  temp: number;
  memory: number;
  powerDraw: number;
};

class DatabaseGpu extends Database<HwGpu.History.Value, HwGpu.Load> {
  constructor() {
    super();
  }

  public async logLoad(host: string, data: HwGpu.Load, timestamp: Date = new Date()) {
    try {
      const writeApi = this.getWriteApi(host);

      for (const gpu of data.gpus) {
        writeApi.writePoint(
          new Point("gpu")
            .tag("index", gpu.index.toString())
            .tag("model", gpu.model)
            .floatField("load", gpu.load)
            .floatField("memory", gpu.memoryUsed)
            .floatField("temp", gpu.temp ?? 0)
            .floatField("powerDraw", gpu.powerDraw ?? 0)
            .timestamp(timestamp)
        );
      }

      await writeApi.close();

      console.log(`Logged Gpu load for ${host} at ${timestamp.toISOString()}`);
    } catch (e) {
      console.error(e);
    }
  }

  override getQuery(host: string, rangeStart: string, every: string): string {
    return `
    from(bucket: "${config.dbBucket}")
      |> range(start: ${rangeStart})
      |> filter(fn: (r) => r._measurement == "gpu" and r.host == "${host}")
      |> aggregateWindow(every: ${every}, fn: mean, createEmpty: true)
      |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
      |> yield(name: "mean")
    `;
  }

  override async getDatas(query: string): Promise<HwGpu.History.Value[]> {
    const data = (await this.executeQuery<GpuQueryResult>(query)).filter(
      (d) => !(d.temp === null && d.load === null && d.memory === null && d.powerDraw === null)
    );
    const result: HwGpu.History.Value[] = [];

    const uniqueTimestamps = new Set(data.map((d) => d._time));

    for (const timestamp of uniqueTimestamps) {
      const gpus = data
        .filter((d) => d._time === timestamp)
        .map((point) => ({
          index: parseInt(point.index),
          load: point.load,
          memoryUsed: point.memory,
          temp: point.temp,
          powerDraw: point.powerDraw,
          model: point.model,
        }));

      result.push({
        gpus,
        timestamp: new Date(timestamp).getTime(),
      });
    }

    return result;
  }
}

export const databaseGpu = new DatabaseGpu();
