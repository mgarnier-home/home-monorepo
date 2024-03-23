import { WidgetInterfaces } from '@shared/interfaces/widgetInterfaces';

import StatsWidget from './stats/stats.widget';

export const getWidget = (widget: WidgetInterfaces.Widget) => {
  switch (widget.type) {
    case WidgetInterfaces.Type.Stats:
      return <StatsWidget options={widget.options} />;
    default:
      return <div>Unknown widget type</div>;
  }
};
