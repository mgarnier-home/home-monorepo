import { useContext, useEffect, useMemo, useState } from "react";
import { Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";

import { HwRam } from "@mgarnier11/nodesight-types";
import WidgetInterfaces from "@shared/interfaces/widgetInterfaces";

import { ConfigContext } from "../../../utils/configContext";
import Utils from "../../../utils/utils";
import { WidgetContext } from "../widgetContext";
import CustomTooltip from "./tooltip";
import { useData } from "./useData";

type RamStatsWidgetProps = {
  options: WidgetInterfaces.Stats.Options;
};

function renderRamPayload(
  payload: any,
  renderItem: (name: string, value: number, unit: string, color: string) => JSX.Element
) {
  if (!payload || !payload.length) return null;

  const ram = payload.find((p: any) => p.dataKey === "load");

  return <div className="flex flex-col">{renderItem("Ram", ram?.value, "MB", ram?.stroke)}</div>;
}

function RamStatsWidget(props: RamStatsWidgetProps) {
  const host = useContext(WidgetContext);
  const appConfig = useContext(ConfigContext);
  const { options } = props;

  const [data] = useData<HwRam.History.Value>(appConfig, options);

  return (
    <ResponsiveContainer width="100%" height="100%">
      <LineChart margin={{ left: -40, bottom: -15, top: 5, right: 20 }} data={data} syncId={host.name}>
        <XAxis
          dataKey="timestamp"
          type="number"
          tickSize={4}
          fontSize={8}
          interval="preserveStartEnd"
          domain={["dataMin", "dataMax"]}
          tickFormatter={Utils.dateFormatter}
        />
        <Tooltip content={<CustomTooltip renderPayload={renderRamPayload} />} />
        <YAxis
          dataKey="load"
          domain={[0, options.ramMaxMemory ?? 2048]}
          tickFormatter={(value) => `${Math.round(value / 1024)}G`}
          tickSize={2}
          fontSize={8}
        />
        <Line
          dataKey="load"
          type={"monotone"}
          strokeWidth={1}
          dot={false}
          isAnimationActive={false}
          stroke={Utils.getLineColor("memory")}
        />
      </LineChart>
    </ResponsiveContainer>
  );
}

export default RamStatsWidget;
