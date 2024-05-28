import jsYaml from 'js-yaml';
import { Utils } from 'utils';

import { AppData, Host } from './interfaces';

export const parseTraefikLabels = (
  host: Host,
  labels: { [key: string]: string },
  data: AppData
): { services: any; routers: any } => {
  const proxy = data.proxies.find((proxy) => proxy.sourceIP === host.ip);
  const serverIp = proxy?.activated ? proxy.destIP : host.ip;

  const services: any = {};

  const personnalisedRouters: any = {};

  const defaultRouters: any = {};

  for (const [key, value] of Object.entries(labels)) {
    const trimmedKey = key.trim();
    const trimmedValue = value.trim();

    try {
      // if (trimmedKey.startsWith('my-traefik.http.services')) {
      //   const serviceName = trimmedKey.split('.')[3];

      //   if (serviceName) {
      //     services[serviceName] = {};
      //     if (trimmedKey.endsWith('port')) {
      //       services[serviceName].loadBalancer = {
      //         servers: [{ url: `http://${serverIp}:{{env "${trimmedValue}" }}/` }],
      //       };
      //     }
      //   }
      // }

      if (trimmedKey.startsWith('my-traefik.http.routers')) {
        // console.log('key', trimmedKey, 'value', trimmedValue);
        const splitKey = trimmedKey.split('.');
        const routerName = splitKey[3];

        if (routerName) {
          if (!personnalisedRouters[routerName]) personnalisedRouters[routerName] = {};

          if (trimmedKey.endsWith('rule')) {
            personnalisedRouters[routerName].rule = trimmedValue;
          }

          if (trimmedKey.endsWith('entrypoints')) {
            personnalisedRouters[routerName].entryPoints = trimmedValue.split(',');
          }

          if (trimmedKey.endsWith('service')) {
            personnalisedRouters[routerName].service = trimmedValue;
          }

          if (trimmedKey.endsWith('middlewares')) {
            personnalisedRouters[routerName].middlewares = trimmedValue.split(',');
          }

          if (splitKey[4] === 'tls') {
            if (!personnalisedRouters[routerName].tls) {
              personnalisedRouters[routerName].tls = {};
            }

            personnalisedRouters[routerName].tls[splitKey[5]!] = trimmedValue;
          }
        }
      }

      if (trimmedKey === 'traefik-conf.port') {
        const serviceName = labels['traefik-conf.name'] || labels['com.docker.compose.service'] || '';
        const rule = `Host(\`${serviceName}.${data.domainName}\`)`;
        const entrypoints = data.defaultEntrypoints;
        const middlewares = data.defaultMiddlewares;
        const tls = data.defaultTls
          ? {
              certResolver: data.defaultTls.defaultCertResolver ?? undefined,
              options: data.defaultTls.defaultOptions ?? undefined,
            }
          : undefined;

        const port = isNaN(parseInt(trimmedValue, 10)) ? `{{env "${trimmedValue}" }}` : parseInt(trimmedValue, 10);

        const routerName = `${host.name}_${serviceName}`.toLowerCase();

        defaultRouters[routerName] = {
          entryPoints: entrypoints.split(',') || [],
          middlewares: middlewares.split(',') || [],
          service: routerName,
          tls,
          rule,
        };

        services[routerName] = {
          loadBalancer: {
            servers: [
              {
                url: `http://${serverIp}:${port}`,
              },
            ],
          },
        };
      }
    } catch (error) {
      console.error('Error while parsing traefik labels : ', error);
    }
  }

  return { services, routers: Utils.deepMerge(defaultRouters, personnalisedRouters) };
};
