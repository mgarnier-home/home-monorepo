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

    const entryPoints = service.entryPoints ? service.entryPoints.split(',') : ['http'];
    const middlewares = service.middlewares ? service.middlewares.split(',') : undefined;
    const tlsResolver = service.tlsResolver;
    const host = service.rule || `${service.serviceName}.${service.host.name.toLowerCase()}.home`;

    routers[routerName] = {
      entryPoints,
      service: serviceName,
      middlewares,
      rule: `Host(\`${host}\`)`,
      tls: tlsResolver ? { certResolver: tlsResolver } : undefined,
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
