import { useContext, useEffect, useState } from 'react';

import { ApiInterfaces } from '@shared/interfaces/apiInterfaces';

import { Api } from '../../utils/api';
import { ConfigContext, StatusChecksContext, WidgetContext } from '../../utils/contexts';
import Utils from '../../utils/utils';
import { loaderSm } from '../dotLoader/dotLoader';
import Icon from '../icon/icon';
import Service from '../service/service';
import StatusIndicator from '../statusIndicator/statusIndicator';
import { getWidget } from '../widgets/widget';

import type { AppInterfaces } from '@shared/interfaces/appInterfaces';
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
    <div className='flex text-black bg-primary h-8'>
      {host.icon && <Icon icon={host.icon} size='sm' />}
      <h3 className='text-xl pl-1 pt-0.5 flex-1'>{host.name}</h3>
      <div className='flex flex-col items-end justify-center'>
        {pingInfos.loading ? (
          loaderSm('yellow')
        ) : pingInfos.ping ? (
          pingInfos.ms < 50 ? (
            <StatusIndicator color='bg-success' />
          ) : (
            <div className='text grow flex' style={{ color: Utils.getPingColor(pingInfos.ms) }}>
              {pingInfos.ms}
            </div>
          )
        ) : (
          <StatusIndicator color='bg-danger' />
        )}
      </div>
    </div>
  );
}

function Host(props: HostProps) {
  const { host } = props;
  const appConfig = useContext(ConfigContext);

  const { services, widgets } = host;

  const [statusChecks, setStatusChecks] = useState<Record<string, ApiInterfaces.StatusChecks.ResponseData>>({});

  const getServiceId = (service: AppInterfaces.HostService) => `${host.id}_${service.name}`;

  const refreshStatusChecks = async () => {
    const statusChecksRequest = {
      statusChecks: services
        .map((service) =>
          service.statusChecks.map((statusCheck) => ({
            id: `${getServiceId(service)}_${statusCheck.name ?? 'Service'}`,
            url: statusCheck.url || service.url,
          }))
        )
        .flat(),
    };

    const response = await Api.getStatusChecks(statusChecksRequest);

    const statusChecksMap: Record<string, ApiInterfaces.StatusChecks.ResponseData> = {};

    response.statusChecks.forEach((statusCheck) => {
      statusChecksMap[statusCheck.id] = statusCheck;
    });

    setStatusChecks(statusChecksMap);
  };

  useEffect(() => {
    console.log('useEffect');

    refreshStatusChecks();

    const pingInterval = setInterval(refreshStatusChecks, appConfig.globalConfig.statusCheckInterval);

    return () => clearInterval(pingInterval);
  }, []);

  return (
    <div className='bg-background-darker m-4 border-4 border-primary rounded-lg' style={{ order: host.order ?? 0 }}>
      <Title host={host} />

      <div className='flex flex-wrap p-2'>
        <StatusChecksContext.Provider value={statusChecks}>
          {(services || []).map((service, index) => (
            <Service key={`${host.id}${service.name}${index}`} service={service} serviceId={getServiceId(service)} />
          ))}
        </StatusChecksContext.Provider>
      </div>
      <div className='p-2'>
        <WidgetContext.Provider value={host}>
          {(widgets || []).map((widget, index) => (
            <div key={`${host.id}${widget.name}${index}`}>{getWidget(widget)}</div>
          ))}
        </WidgetContext.Provider>
      </div>
    </div>
  );
}

export default Host;