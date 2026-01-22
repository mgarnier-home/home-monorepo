import { Utils } from '@libs/utils';

import { AppData, Host } from './interfaces';

export const parseTraefikLabels = (
  host: Host,
  labels: { [key: string]: string },
  data: AppData,
  subDomain?: string,
): { services: any; routers: any; middlewares: any } => {
  const proxy = data.proxies.find((proxy) => proxy.sourceIP === host.ip);
  const serverIp = proxy?.activated ? proxy.destIP : host.ip;

  const services: any = {};
  const middlewares: any = {};

  const personnalisedRouters: any = {};

  const defaultRouters: any = {};

  let sablierMiddleware = '';

  for (const [key, value] of Object.entries(labels)) {
    const trimmedKey = key.trim();
    const trimmedValue = value.trim();

    try {
      if (trimmedKey.startsWith('traefik-conf.http.services')) {
        const serviceName = trimmedKey.split('.')[3];

        if (serviceName) {
          services[serviceName] = {};
          if (trimmedKey.endsWith('port')) {
            services[serviceName].loadBalancer = {
              servers: [{ url: `http://${serverIp}:{{env "${trimmedValue}" }}/` }],
            };
          }
        }
      }

      if (trimmedKey.startsWith('traefik-conf.http.routers')) {
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

      if (trimmedKey === 'sablier.group') {
        // Create a new middleware for this sablier group
        sablierMiddleware = `sablier-${trimmedValue}`;

        middlewares[sablierMiddleware] = {
          plugin: {
            sablier: {
              group: trimmedValue,
              sablierUrl: data.sablierUrl,
              blocking: {
                timeout: '1m',
              },
            },
          },
        };
      }

      if (trimmedKey === 'sablier.group.strategy') {
        if (sablierMiddleware && middlewares[sablierMiddleware]) {
          if (trimmedValue === 'dynamic') {
            // default blocking, switch to dynamic
            middlewares[sablierMiddleware].plugin.sablier.blocking = undefined;
            middlewares[sablierMiddleware].plugin.sablier.dynamic = {
              showDetails: true,
              theme: data.sablierTheme || 'default',
              refreshFrequency: '2s',
              sessionDuration: data.sablierSessionDuration || '5m',
            };
          }
        }
      }

      if (trimmedKey === 'traefik-conf.port') {
        const serviceName = labels['traefik-conf.name'] || labels['com.docker.compose.service'] || '';
        const serviceNameWithSubDomain = subDomain ? `${serviceName}.${subDomain}` : serviceName;
        const rule = `Host(\`${serviceNameWithSubDomain}.${data.domainName}\`)`;
        const serviceEntrypoints = data.defaultEntrypoints;
        const serviceMiddlewares = labels['traefik-conf.middlewares'] || data.defaultMiddlewares;

        const tls = data.defaultTls
          ? {
              certResolver: data.defaultTls.defaultCertResolver ?? undefined,
              options: data.defaultTls.defaultOptions ?? undefined,
            }
          : undefined;

        const port = isNaN(parseInt(trimmedValue, 10)) ? `{{env "${trimmedValue}" }}` : parseInt(trimmedValue, 10);

        const routerName = `${host.name}_${serviceName}`.toLowerCase();

        defaultRouters[routerName] = {
          entryPoints: serviceEntrypoints.split(',') || [],
          middlewares: serviceMiddlewares.split(',') || [],
          service: routerName,
          tls,
          rule,
        };
        if (sablierMiddleware) {
          defaultRouters[routerName].middlewares.push(sablierMiddleware + '@http');
        }

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

  return { services, routers: Utils.deepMerge(defaultRouters, personnalisedRouters), middlewares };
};
