import { useContext, useEffect, useState } from "react";

import AppInterfaces from "@shared/interfaces/appInterfaces";

import { Api } from "../../utils/api";
import { ConfigContext } from "../../utils/configContext";
import Utils from "../../utils/utils";
import { loaderSm } from "../dotLoader/dotLoader";
import Icon from "../icon/icon";
import Service from "../service/service";
import StatusIndicator from "../statusIndicator/statusIndicator";
import { getWidget } from "../widgets/widget";
import { WidgetContext } from "../widgets/widgetContext";

type HostProps = {
  host: AppInterfaces.Host;
};

function Title(props: { host: AppInterfaces.Host }) {
  const { host } = props;
  const appConfig = useContext(ConfigContext);

  const [pingInfos, setPingInfos] = useState<{ ping: boolean; duration: number; loading: boolean; ms: number }>({
    ping: false,
    duration: 0,
    loading: true,
    ms: NaN,
  });

  const pingHost = async () => {
    const res = await Api.pingHost(host);

    setPingInfos({
      ping: res.ping,
      duration: res.duration,
      loading: false,
      ms: res.ms,
    });
  };

  useEffect(() => {
    pingHost();

    const pingInterval = setInterval(pingHost, appConfig.globalConfig.pingInterval);

    return () => clearInterval(pingInterval);
  }, []);

  return (
    <div className="flex text-black bg-primary h-8">
      {host.icon && <Icon icon={host.icon} size="sm" />}
      <h3 className="text-xl pl-1 pt-0.5 flex-1">{host.name}</h3>
      <div className="flex flex-col items-end justify-center">
        {pingInfos.loading ? (
          loaderSm("yellow")
        ) : pingInfos.ping ? (
          pingInfos.ms < 50 ? (
            <StatusIndicator color="bg-success" />
          ) : (
            <div className="text grow flex" style={{ color: Utils.getPingColor(pingInfos.ms) }}>
              {pingInfos.ms}
            </div>
          )
        ) : (
          <StatusIndicator color="bg-danger" />
        )}
      </div>
    </div>
  );
}

function Host(props: HostProps) {
  const { host } = props;

  const { services, widgets } = host;

  return (
    <div className="bg-background-darker m-4 border-4 border-primary rounded-lg" style={{ order: host.order ?? 0 }}>
      <Title host={host} />

      <div className="flex flex-wrap p-2">
        {(services || []).map((service, index) => (
          <Service key={`${host.name}${service.name}${index}`} service={service} />
        ))}
      </div>
      <div className="p-2">
        <WidgetContext.Provider value={host}>
          {(widgets || []).map((widget, index) => (
            <div key={`${host.name}${widget.name}${index}`}>{getWidget(widget)}</div>
          ))}
        </WidgetContext.Provider>
      </div>
    </div>
  );
}

export default Host;
