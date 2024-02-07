import fs from 'fs';
import jsYaml from 'js-yaml';
import path from 'path';

import { AppData, Host, StackInfos, TraefikService } from './interfaces';

export const getComposeStacksPaths = async (
  hosts: Host[],
  composeDir: string,
  stacksToIgnore: string[]
): Promise<StackInfos[]> => {
  const stacks = (await listStacks(composeDir)).filter((stack) => !stacksToIgnore.includes(stack));
  const stacksInfos: StackInfos[] = [];

  for (const host of hosts) {
    for (const stack of stacks) {
      const stackFilePath = path.join(composeDir, stack, `${host.name.toLowerCase()}.${stack}.yml`);
      if (fs.existsSync(stackFilePath)) {
        stacksInfos.push({ path: stackFilePath, host, stack });
      }
    }
  }

  return stacksInfos;
};

export const listStacks = async (dir: string): Promise<string[]> => {
  const filesInFolder = (await fs.promises.readdir(dir)).map((file) => file.toLowerCase());

  return await Promise.all(
    filesInFolder.filter(async (file) => {
      const stat = await fs.promises.lstat(path.join(dir, file));

      return stat.isDirectory();
    })
  );
};

export const getTraefikServices = (stackInfos: StackInfos, yaml: string): TraefikService[] => {
  const yamlObj: any = jsYaml.load(yaml);

  const services: TraefikService[] = [];

  for (const serviceName in yamlObj.services) {
    const service = yamlObj.services[serviceName];
    const traefikPorts: string[] = (service.annotations?.['traefik-ports'] || '').split(',');
    const serviceNameOverride = service.annotations?.['traefik-name'];

    if (traefikPorts.length > 0 && traefikPorts[0] !== '') {
      for (const port of traefikPorts) {
        services.push({
          host: stackInfos.host,
          serviceName: serviceNameOverride || serviceName,
          stackName: stackInfos.stack,
          portVariable: port.trim(),
        });
      }
    }
  }

  return services;
};

export const getTraefikDynamicConf = (traefikServices: TraefikService[], data: AppData): string => {
  const routers: any = {};
  const services: any = {};

  for (const service of traefikServices) {
    const routerName = `${service.host.name}_${service.stackName}_${service.serviceName}`;
    const serviceName = `${service.host.name}_${service.stackName}_${service.serviceName}`;

    const proxy = data.proxies.find((proxy) => proxy.sourceIP === service.host.ip);

    const serviceIp = proxy?.activated ? proxy.destIP : service.host.ip;

    routers[routerName] = {
      entryPoints: ['http'],
      service: serviceName,
      rule: `Host(\`${service.serviceName}.${service.host.name}.home\`)`,
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
