import { HwCpu } from 'nodesight-types';
import { useContext, useMemo } from 'react';
import { Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts';

import { ConfigContext } from '../../../utils/configContext';
import Utils from '../../../utils/utils';
import { WidgetContext } from '../widgetContext';
import CustomTooltip from './tooltip';
import { useData } from './useData';

import type { WidgetInterfaces } from '@shared/interfaces/widgetInterfaces';
type CpuStatsWidgetProps = {
  options: WidgetInterfaces.Stats.Options;
};

function renderCpuPayload(
  payload: any,
  renderItem: (name: string, value: number, unit: string, color: string) => JSX.Element
) {
  if (!payload || !payload.length) return null;

  const load = payload.find((p: any) => p.dataKey === 'load');
  const temp = payload.find((p: any) => p.dataKey === 'temp');

  return (
    <div className='flex flex-col'>
      {renderItem('CPU Load', load?.value, '%', load?.stroke)}
      {renderItem('Temp', temp?.value, '°C', temp?.stroke)}
    </div>
  );
}

function CpuStatsWidget(props: CpuStatsWidgetProps) {
  const host = useContext(WidgetContext);
  const appConfig = useContext(ConfigContext);
  const { options } = props;

  const [data] = useData<HwCpu.History.Value>(appConfig, options);

  const sanitizedData = useMemo(() => {
    console.log(data);
    return data.map((d) => ({
      temp: parseFloat(d.temp.toFixed(1)),
      load: parseFloat(d.load.toFixed(1)),
      timestamp: d.timestamp,
    }));
  }, [data]);

  return (
    <ResponsiveContainer width='100%' height='100%'>
      <LineChart margin={{ left: -40, bottom: -15, top: 5, right: 20 }} data={sanitizedData} syncId={host.name}>
        <XAxis
          dataKey='timestamp'
          type='number'
          tickSize={4}
          fontSize={8}
          interval='preserveStartEnd'
          domain={['dataMin', 'dataMax']}
          tickFormatter={Utils.dateFormatter}
        />
        <Tooltip content={<CustomTooltip renderPayload={renderCpuPayload} />} />
        <YAxis dataKey='load' domain={[0, 100]} tickSize={2} fontSize={8} />
        <Line
          dataKey='load'
          type={'monotone'}
          strokeWidth={1}
          dot={false}
          isAnimationActive={false}
          stroke={Utils.getLineColor('load')}
        />
        <Line
          dataKey='temp'
          type={'monotone'}
          strokeWidth={1}
          dot={false}
          isAnimationActive={false}
          stroke={Utils.getLineColor('temp')}
        />
      </LineChart>
    </ResponsiveContainer>
  );
}

export default CpuStatsWidget;
