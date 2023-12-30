import jsYaml from 'js-yaml';
import { Utils } from 'utils';

import { Host } from './interfaces';

export const mergeYamls = (yamlFiles: string[]): string => {
  const yamls: any[] = yamlFiles.map((yaml) => jsYaml.load(yaml));

  const mergedYaml = yamls.reduce((acc, curr) => {
    return Utils.deepMerge(acc, curr);
  }, {});

  return jsYaml.dump(mergedYaml);
};

export const getHost = (hostName: string, hosts: Host[]) => {
  return hosts.find((host) => host.name === hostName);
};
