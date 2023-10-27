import { useContext, useEffect, useRef, useState } from "react";

import AppInterfaces from "@shared/interfaces/appInterfaces";
import WidgetInterfaces from "@shared/interfaces/widgetInterfaces";

import { StatsApi } from "../../../utils/statsApi";
import { WidgetContext } from "../widgetContext";

export function useData<T>(
  config: AppInterfaces.AppConfig,
  options: WidgetInterfaces.Stats.Options
): [T[], boolean, Error | null] {
  const [data, setData] = useState<T[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<Error | null>(null);
  const host = useContext(WidgetContext);
  const isMounted = useRef(true);

  useEffect(() => {
    let liveInterval: NodeJS.Timeout | null = null;

    const fetchData = async () => {
      setLoading(true);

      try {
        if (options.time === WidgetInterfaces.Stats.OptionsTime.Live) {
          const history = await StatsApi.getHistory<T[]>(config.globalConfig.statsApiUrl, options, host.name);

          if (isMounted.current) {
            setData(history);
            setLoading(false);
          }

          liveInterval = setInterval(async () => {
            try {
              const current = await StatsApi.getCurrentType<T>(host.nodesightUrl!, options);

              if (isMounted.current) {
                setData((oldData) => [...oldData.slice(1), current]);
              }
            } catch (err: any) {
              if (isMounted.current) setError(err);
            }
          }, 15000);
        } else if (options.time === WidgetInterfaces.Stats.OptionsTime.History) {
          const historyData = await StatsApi.getHistory<T[]>(config.globalConfig.statsApiUrl, options, host.name);
          if (isMounted.current) {
            setData(historyData);
            setLoading(false);
          }
        }
      } catch (err: any) {
        if (isMounted.current) setError(err);
      }
    };

    fetchData();

    return () => {
      if (liveInterval) {
        clearInterval(liveInterval);
      }
      isMounted.current = false;
    };
  }, [config, options]);

  return [data, loading, error];
}
