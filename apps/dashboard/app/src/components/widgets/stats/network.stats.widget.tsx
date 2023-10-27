import { useContext, useEffect, useMemo, useState } from "react";
import { Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";

import { HwNetwork } from "@mgarnier11/nodesight-types";
import WidgetInterfaces from "@shared/interfaces/widgetInterfaces";
import { convert } from "@shared/utils";

import { ConfigContext } from "../../../utils/configContext";
import { StatsApi } from "../../../utils/statsApi";
import Utils from "../../../utils/utils";
import { WidgetContext } from "../widgetContext";
import CustomTooltip from "./tooltip";
import { useData } from "./useData";

type NetworkStatsWidgetProps = {
  options: WidgetInterfaces.Stats.Options;
};

function renderNetworkPayload(
  payload: any,
  renderItem: (name: string, value: number, unit: string, color: string) => JSX.Element
) {
  if (!payload || !payload.length) return null;

  const up = payload.find((p: any) => p.dataKey === "up");
  const down = payload.find((p: any) => p.dataKey === "down");

  return (
    <div className="flex flex-col">
      {renderItem("Up", up?.value, "KB/s", up?.stroke)}
      {renderItem("Down", down?.value, "KB/s", down?.stroke)}
    </div>
  );
}

function NetworkStatsWidget(props: NetworkStatsWidgetProps) {
  const host = useContext(WidgetContext);
  const appConfig = useContext(ConfigContext);
  const { options } = props;

  const [data] = useData<HwNetwork.History.Value>(appConfig, options);

  const sanitizedData = data.map((d) => ({
    up: Math.round(convert(d.up, "B", "KB")),
    down: Math.round(convert(d.down, "B", "KB")),
    timestamp: d.timestamp,
  }));

  return (
    <ResponsiveContainer width="100%" height="100%">
      <LineChart margin={{ left: -40, bottom: -15, top: 5, right: 20 }} data={sanitizedData} syncId={host.name}>
        <XAxis
          dataKey="timestamp"
          type="number"
          tickSize={4}
          fontSize={8}
          interval="preserveStartEnd"
          domain={["dataMin", "dataMax"]}
          tickFormatter={Utils.dateFormatter}
        />
        <Tooltip content={<CustomTooltip renderPayload={renderNetworkPayload} />} />
        <YAxis tickSize={2} fontSize={8} />
        <Line
          dataKey="up"
          type={"monotone"}
          strokeWidth={1}
          dot={false}
          isAnimationActive={false}
          stroke={Utils.getLineColor("netUp")}
        />
        <Line
          dataKey="down"
          type={"monotone"}
          strokeWidth={1}
          dot={false}
          isAnimationActive={false}
          stroke={Utils.getLineColor("netDown")}
        />
      </LineChart>
    </ResponsiveContainer>
  );
}

export default NetworkStatsWidget;
