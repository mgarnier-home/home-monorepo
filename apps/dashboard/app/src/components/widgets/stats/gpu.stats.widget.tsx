import { HwGpu } from 'nodesight-types';
import { useContext, useMemo } from 'react';
import { Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts';

import { ConfigContext, WidgetContext } from '../../../utils/contexts';
import Utils from '../../../utils/utils';
import CustomTooltip from './tooltip';
import { useData } from './useData';

import type { WidgetInterfaces } from '@shared/interfaces/widgetInterfaces';
type GpuStatsWidgetProps = {
  options: WidgetInterfaces.Stats.Options;
};

function renderGpuPayload(
  payload: any,
  renderItem: (name: string, value: number, unit: string, color: string) => JSX.Element
) {
  if (!payload || !payload.length) return null;

  const load = payload.find((p: any) => p.dataKey === 'load');
  // const power = payload.find((p: any) => p.dataKey === "powerDraw");
  const memory = payload.find((p: any) => p.dataKey === 'memoryUsed');
  const temp = payload.find((p: any) => p.dataKey === 'temp');

  return (
    <div className='flex flex-col'>
      {renderItem('GPU Load', load?.value, '%', load?.stroke)}
      {/* {renderItem("Power", power?.value, "W", power?.stroke)} */}
      {renderItem('VRam', memory?.value, 'MB', memory?.stroke)}
      {renderItem('Temp', temp?.value, '°C', temp?.stroke)}
    </div>
  );
}

function GpuStatsWidget(props: GpuStatsWidgetProps) {
  const host = useContext(WidgetContext);
  const appConfig = useContext(ConfigContext);
  const { options } = props;

  const [data] = useData<HwGpu.History.Value>(appConfig, options);

  const sanitizedData = useMemo(() => {
    return (data ?? []).map((d) => {
      const gpu = d.gpus.find((gpu) => gpu.index === options.gpuId);

      if (!gpu) return null;

      return {
        load: gpu.load,
        powerDraw: gpu.powerDraw,
        memoryUsed: gpu.memoryUsed,
        temp: gpu.temp,
        timestamp: d.timestamp,
      };
    });
  }, [data]);

  return (
    <ResponsiveContainer width='100%' height='100%'>
      <LineChart margin={{ left: -40, bottom: -15, top: 5, right: -40 }} data={sanitizedData} syncId={host.id}>
        <XAxis
          dataKey='timestamp'
          type='number'
          tickSize={4}
          fontSize={8}
          interval='preserveStartEnd'
          domain={['dataMin', 'dataMax']}
          tickFormatter={Utils.dateFormatter}
        />
        <Tooltip content={<CustomTooltip renderPayload={renderGpuPayload} />} />
        <YAxis dataKey='load' domain={[0, 100]} tickSize={2} fontSize={8} />
        <YAxis
          dataKey='memory'
          domain={[0, options.gpuMaxMemory ?? 2048]}
          tickSize={2}
          fontSize={8}
          tickFormatter={(value) => `${Math.round(value / 1024)}G`}
          orientation='right'
          yAxisId='gpumemory'
        />
        <Line
          dataKey='load'
          type='monotone'
          strokeWidth={1}
          dot={false}
          isAnimationActive={false}
          stroke={Utils.getLineColor('load')}
        />
        {/* <Line
          dataKey="powerDraw"
          type="monotone"
          strokeWidth={1}
          dot={false}
          isAnimationActive={false}
          stroke={Utils.getLineColor("power")}
        /> */}
        <Line
          dataKey='memoryUsed'
          yAxisId='gpumemory'
          type='monotone'
          strokeWidth={1}
          dot={false}
          isAnimationActive={false}
          stroke={Utils.getLineColor('memory')}
        />
        <Line
          dataKey='temp'
          type='monotone'
          strokeWidth={1}
          dot={false}
          isAnimationActive={false}
          stroke={Utils.getLineColor('temp')}
        />
      </LineChart>
    </ResponsiveContainer>
  );
}

export default GpuStatsWidget;
