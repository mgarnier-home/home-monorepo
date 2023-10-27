import { NameType, ValueType } from "recharts/types/component/DefaultTooltipContent";
import { ContentType } from "recharts/types/component/Tooltip";

import Utils from "../../../utils/utils";

const CustomTooltip = ({ active, payload, label, renderPayload }: any) => {
  const renderItem = (name: string, value: number, unit: string, color: string) => (
    <div className="flex justify-between" style={{ color }}>
      <div>{name}</div>
      &emsp;
      <div>
        {new Intl.NumberFormat().format(value)} {unit}
      </div>
    </div>
  );
  if (active && payload && payload.length) {
    return (
      <div className="bg-background-darker w-32 text-primary text-xs border-2 border-primary leading-tight	rounded-lg px-2 py-1">
        {/* {Utils.dateFormatter(label)} */}
        <div>{renderPayload(payload, renderItem)}</div>
      </div>
    );
  }

  return null;
};

export default CustomTooltip;
