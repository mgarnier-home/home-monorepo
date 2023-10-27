import { HwCpu, HwGpu, HwNetwork, HwRam, HwStorage } from "@mgarnier11/nodesight-types";
import WidgetInterfaces from "@shared/interfaces/widgetInterfaces";
import { SimpleCache } from "@shared/simpleCache";

import { Api } from "./api";

export namespace StatsApi {
  function getApiH(history?: WidgetInterfaces.Stats.OptionsTimeHistory) {
    if (history === WidgetInterfaces.Stats.OptionsTimeHistory.LastHour) return "hour";
    else if (history === WidgetInterfaces.Stats.OptionsTimeHistory.LastDay) return "day";
    else if (history === WidgetInterfaces.Stats.OptionsTimeHistory.LastWeek) return "week";
    else return "";
  }

  function getUri(options: WidgetInterfaces.Stats.Options, hostname: string) {
    return `${hostname}/${options.type}/${getApiH(options.history)}`;
  }

  const currentCache: SimpleCache<{
    cpu: HwCpu.Load;
    gpu: HwGpu.Load;
    ram: HwRam.Load;
    network: HwNetwork.Load;
    storage: HwStorage.Load;
  }> = new SimpleCache();

  export async function getHistory<T>(
    statsApiUrl: string,
    options: WidgetInterfaces.Stats.Options,
    hostname: string
  ): Promise<T> {
    const uri = getUri(options, hostname);
    const url = `${statsApiUrl}/api/${uri}`;

    const response = await Api.makeServerRequest<T>(url, "GET");

    return response.data!;
  }

  type Current = {
    cpu: HwCpu.Load;
    gpu: HwGpu.Load;
    ram: HwRam.Load;
    network: HwNetwork.Load;
    storage: HwStorage.Load;
  };

  export async function getCurrent(nodesightUrl: string): Promise<Current> {
    const cacheHit = currentCache.get(nodesightUrl);

    if (cacheHit) {
      return cacheHit;
    }

    const response = await Api.makeServerRequest<Current>(`${nodesightUrl}/current`, "GET");

    if (response.data) {
      currentCache.set(nodesightUrl, response.data!);
    }

    return response.data as Current;
  }

  export async function getCurrentType<T>(nodesightUrl: string, options: WidgetInterfaces.Stats.Options): Promise<T> {
    const current = await getCurrent(nodesightUrl);

    let data;

    if (current) {
      if (options.type === WidgetInterfaces.Stats.OptionsType.Cpu) {
        data = {
          load: current.cpu.cores.reduce((acc, core) => acc + core.load, 0) / current.cpu.cores.length,
          temp: current.cpu.cores.reduce((acc, core) => acc + core.temp, 0) / current.cpu.cores.length,
          core: -1,
        };
      } else if (options.type === WidgetInterfaces.Stats.OptionsType.Gpu) {
        data = current.gpu;
      } else if (options.type === WidgetInterfaces.Stats.OptionsType.Ram) {
        data = current.ram;
      } else if (options.type === WidgetInterfaces.Stats.OptionsType.Network) {
        data = current.network;
      }
    }

    return { ...data, timestamp: Date.now() } as T;
  }
}
