import path from 'path';

export namespace Utils {
  export async function fetchWithTimeout(url: RequestInfo | URL, timeout = 5000, init?: RequestInit | undefined) {
    const controller = new AbortController();
    const id = setTimeout(() => controller.abort(), timeout);
    const response = await fetch(url, {
      ...init,
      signal: controller.signal,
    });
    clearTimeout(id);
    return response;
  }

  export async function timeout(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }

  export function padStart(value: number, length: number, padChar: string = '0') {
    return value.toString().padStart(length, padChar);
  }

  export enum Unit {
    B = 'B',
    KB = 'KB',
    MB = 'MB',
    GB = 'GB',
    TB = 'TB',
  }

  export function convert(input: number, inputUnit: keyof typeof Unit, outputUnit: keyof typeof Unit) {
    const units = {
      B: 1,
      KB: 1024,
      MB: 1024 * 1024,
      GB: 1024 * 1024 * 1024,
      TB: 1024 * 1024 * 1024 * 1024,
    };

    return (input * units[inputUnit]) / units[outputUnit];
  }

  export const capFirst = (str: string) => {
    return str[0]?.toUpperCase() + str.slice(1);
  };

  export const platformIsWindows = (platform: string) => platform === 'Windows' || platform === 'win32';

  export const deepMerge = (obj1: any, obj2: any) => {
    const clone1 = structuredClone(obj1);

    const clone2 = structuredClone(obj2);

    for (let key in clone2) {
      if (clone2[key] instanceof Object && clone1[key] instanceof Object) {
        clone1[key] = deepMerge(clone1[key], clone2[key]);
      } else {
        clone1[key] = clone2[key];
      }
    }

    return clone1;
  };

  export const hash = (str: string) => {
    return str.split('').reduce((prevHash, currVal) => ((prevHash << 5) - prevHash + currVal.charCodeAt(0)) | 0, 0);
  };

  export const hashCode = (str: string) => {
    return hash(str).toString(16);
  };

  export const resolveConfigPath = (pathToResolve: string, callerDirname?: string): string => {
    return pathToResolve.startsWith('/')
      ? pathToResolve
      : path.resolve(callerDirname || __dirname, '../', pathToResolve);
  };
}
