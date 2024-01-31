import { logger } from 'logger';
import { useContext, useEffect, useState } from 'react';

import { WidgetInterfaces } from '@shared/interfaces/widgetInterfaces';

import { StatsApi } from '../../../utils/statsApi';
import { WidgetContext } from '../widgetContext';

import type { AppInterfaces } from '@shared/interfaces/appInterfaces';
export function useData<T>(
  config: AppInterfaces.AppConfig,
  options: WidgetInterfaces.Stats.Options
): [T[], boolean, Error | null] {
  const [data, setData] = useState<T[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<Error | null>(null);
  const host = useContext(WidgetContext);

  useEffect(() => {
    let liveInterval: NodeJS.Timeout | null = null;

    const fetchData = async () => {
      setLoading(true);

      try {
        if (options.time === WidgetInterfaces.Stats.OptionsTime.Live) {
          const history = await StatsApi.getHistory<T[]>(config.globalConfig.statsApiUrl, options, host.id);
          setData(history);

          liveInterval = setInterval(async () => {
            try {
              const current = await StatsApi.getCurrentType<T>(host.nodesightUrl!, options);

              setData((oldData) => [...oldData.slice(1), current]);
            } catch (err: any) {
              logger.error(err);
            }
          }, 15000);
        } else if (options.time === WidgetInterfaces.Stats.OptionsTime.History) {
          const historyData = await StatsApi.getHistory<T[]>(config.globalConfig.statsApiUrl, options, host.id);
          setData(historyData);
        }
      } catch (err: any) {
        logger.error(err);
      }
    };

    fetchData();

    return () => {
      if (liveInterval) {
        clearInterval(liveInterval);
      }
    };
  }, [config, options]);

  return [data, loading, error];
}
