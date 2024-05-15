import jsYaml from 'js-yaml';

import { AppData, TraefikService } from './interfaces';

export const getTraefikDynamicConf = (traefikServices: TraefikService[], data: AppData): string => {
  const routers: any = {};
  const services: any = {};

  for (const service of traefikServices) {
    const routerName = `${service.host.name}_${service.serviceName}`.toLowerCase();
    const serviceName = `${service.host.name}_${service.serviceName}`.toLowerCase();

    const proxy = data.proxies.find((proxy) => proxy.sourceIP === service.host.ip);

    const serviceIp = proxy?.activated ? proxy.destIP : service.host.ip;

    routers[routerName] = {
      entryPoints: ['http'],
      service: serviceName,
      rule: `Host(\`${service.serviceName}.${service.host.name.toLowerCase()}.home\`)`,
    };

    services[serviceName] = {
      loadBalancer: {
        servers: [
          {
            url: `http://${serviceIp}:{{env "${service.portVariable}" }}/`,
          },
        ],
      },
    };
  }

  const dynamicConf = {
    http: {
      routers,
      services,
    },
  };

  return jsYaml.dump(dynamicConf);
};
