import { useContext, useEffect, useState } from 'react';

import { Api } from '../../utils/api';
import { ConfigContext } from '../../utils/configContext';
import Icon from '../icon/icon';
import StatusIndicator from '../statusIndicator/statusIndicator';

import type { AppInterfaces } from '@shared/interfaces/appInterfaces';
//#region useStatusChecks hook
type StatusCheckInfo =
  | {
      loading: true;
      name: string;
      clickAction?: AppInterfaces.ClickAction;
    }
  | {
      loading: false;
      name: string;

      color: string;
      clickAction?: AppInterfaces.ClickAction;
    };

const useStatusChecks = (service: AppInterfaces.HostService, statusCheckInterval: number) => {
  const [statusChecks, setStatusChecks] = useState<StatusCheckInfo[]>(
    service.statusChecks.map((statusCheck) => ({
      loading: true,
      name: statusCheck.name,
      clickAction: statusCheck.clickAction as AppInterfaces.ClickAction,
    }))
  );

  const [serviceStatusCheck, setServiceStatusCheck] = useState<StatusCheckInfo | null>(null);

  const updateStatusChecks = (lst: StatusCheckInfo[]) => {
    if (lst.length === 1) {
      lst[0].name = 'Service';
    }

    const serviceStatusCheckIndex = lst.findIndex((statusCheck) => statusCheck.name === 'Service');
    setServiceStatusCheck(lst[serviceStatusCheckIndex]);

    if (serviceStatusCheckIndex !== -1) {
      lst.splice(serviceStatusCheckIndex, 1);
    }

    setStatusChecks(lst);
  };

  const refreshStatusChecks = async () => {
    try {
      const promises = service.statusChecks.map((statusCheck) =>
        Api.makeServerRequest(statusCheck.url || service.url, 'GET').then((response) => ({ statusCheck, response }))
      );

      const results = await Promise.all(promises);

      const updatedStatusChecks: StatusCheckInfo[] = results.map((result) => {
        const { statusCheck, response } = result;

        let color = 'bg-danger';

        if (statusCheck.type === 'singleCode') {
          color = statusCheck.success === response.code ? statusCheck.color || 'bg-success' : 'bg-danger';
        } else if (statusCheck.type === 'multipleCodes') {
          const code = statusCheck.codes.find((c) => c.code === response.code);
          if (code) {
            color = code.color || 'bg-success';
          }
        }

        return {
          loading: false,
          name: statusCheck.name,
          color,
          clickAction: statusCheck.clickAction as AppInterfaces.ClickAction,
        };
      });

      updateStatusChecks(updatedStatusChecks);
    } catch (error) {
      console.error('Error refreshing status checks:', error);
      // Handle any errors - maybe set all statusChecks to an error state, or show a notification.
    }
  };

  useEffect(() => {
    updateStatusChecks(statusChecks);

    refreshStatusChecks();

    const interval = setInterval(refreshStatusChecks, statusCheckInterval);
    return () => clearInterval(interval);
  }, []);

  return { statusChecks, serviceStatusCheck };
};

//#endregion

//#region StatusIndicators component

const renderStatusIndicator = (statusCheck: StatusCheckInfo) => {
  return <StatusIndicator color={statusCheck.loading ? 'bg-warning' : statusCheck.color} />;
};

const StatusIndicators = ({
  statusChecks,
  handleStatusIndicatorClick,
}: {
  statusChecks: StatusCheckInfo[];
  handleStatusIndicatorClick: (e: React.MouseEvent<HTMLElement, MouseEvent>, statusCheck: StatusCheckInfo) => void;
}) => {
  return (
    <div className='grid grid-rows-3 grid-cols-2 grid-flow-col' style={{ direction: 'rtl' }}>
      {statusChecks.map((statusCheck, index) => {
        const hoverClassNames = statusCheck.clickAction ? 'hover:underline cursor-pointer' : '';

        return (
          <div className='flex text-xs leading-3 mb-0.5 order-1' style={{ direction: 'ltr' }} key={index}>
            <div className='flex-grow text-center'>
              <span className={hoverClassNames} onClick={(e) => handleStatusIndicatorClick(e, statusCheck)}>
                {statusCheck.name}
              </span>
            </div>

            {renderStatusIndicator(statusCheck)}
          </div>
        );
      })}
    </div>
  );
};

//#endregion

//#region Service component

type ServiceProps = {
  service: AppInterfaces.HostService;
};

const Service = (props: ServiceProps) => {
  const { service } = props;

  const appConfig = useContext(ConfigContext);

  const { statusChecks, serviceStatusCheck } = useStatusChecks(service, appConfig.globalConfig.statusCheckInterval);

  const handleServiceClick = (e: React.MouseEvent<HTMLDivElement, MouseEvent>) => {
    if (service.clickAction) Api.handleClickAction(service.clickAction as AppInterfaces.ClickAction);
  };

  const handleStatusIndicatorClick = (e: React.MouseEvent<HTMLElement, MouseEvent>, statusCheck: StatusCheckInfo) => {
    e.stopPropagation();

    if (statusCheck.clickAction) Api.handleClickAction(statusCheck.clickAction as AppInterfaces.ClickAction);
  };

  const serviceClickableClassNames =
    'hover:bg-background-hover hover:shadow hover:shadow-primary transition-all cursor-pointer';

  return (
    <div className='text-primary grow basis-1/2 flex-1' style={{ order: service.order ?? 0 }}>
      <div
        className={`p-0.5 pl-1.5 m-2 bg-background h-15 rounded-md flex ${
          service.clickAction ? serviceClickableClassNames : ''
        }`}
        onClick={handleServiceClick}
      >
        <Icon icon={service.icon || 'fa fa-globe'} size='lg' />
        <div className='flex flex-grow'>
          <div className='flex-grow pl-1'>
            <div className='text-base text-center leading-4 flex'>
              <span className='flex-grow'>{service.name}</span>
              {serviceStatusCheck && <div className='mt-0.5'>{renderStatusIndicator(serviceStatusCheck)}</div>}
            </div>

            <StatusIndicators statusChecks={statusChecks} handleStatusIndicatorClick={handleStatusIndicatorClick} />
          </div>
        </div>
      </div>
    </div>
  );
};

export default Service;

//#endregion
