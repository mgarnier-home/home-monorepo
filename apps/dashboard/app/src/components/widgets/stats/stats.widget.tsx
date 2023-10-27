import WidgetInterfaces from "@shared/interfaces/widgetInterfaces";

import CpuStatsWidget from "./cpu.stats.widget";
import GpuStatsWidget from "./gpu.stats.widget";
import NetworkStatsWidget from "./network.stats.widget";
import RamStatsWidget from "./ram.stats.widget";

type StatsWidgetProps = {
  options: WidgetInterfaces.Stats.Options;
};

function StatsWidget(props: StatsWidgetProps) {
  const { options } = props;

  function getWidget() {
    switch (options.type) {
      case WidgetInterfaces.Stats.OptionsType.Cpu:
        return <CpuStatsWidget options={options} />;
      case WidgetInterfaces.Stats.OptionsType.Gpu:
        return <GpuStatsWidget options={options} />;
      case WidgetInterfaces.Stats.OptionsType.Ram:
        return <RamStatsWidget options={options} />;
      case WidgetInterfaces.Stats.OptionsType.Network:
        return <NetworkStatsWidget options={options} />;
      default:
        return <div>Unknown widget type</div>;
    }
  }

  return <div className="h-20 m-2 bg-background rounded-md flex">{getWidget()}</div>;
}

export default StatsWidget;
