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

export enum Unit {
  B = "B",
  KB = "KB",
  MB = "MB",
  GB = "GB",
  TB = "TB",
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

export const platformIsWindows = (platform: string) => platform === "Windows" || platform === "win32";

export const capFirst = (str: string) => {
  return str[0].toUpperCase() + str.slice(1);
};
