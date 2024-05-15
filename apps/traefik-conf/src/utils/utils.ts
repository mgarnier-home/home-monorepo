import jsYaml from 'js-yaml';
import { Utils } from 'utils';

export const mergeYamls = (yamlFiles: string[]): string => {
  const yamls: any[] = yamlFiles.map((yaml) => jsYaml.load(yaml));

  const mergedYaml = yamls.reduce((acc, curr) => {
    return Utils.deepMerge(acc, curr);
  }, {});

  return jsYaml.dump(mergedYaml);
};
