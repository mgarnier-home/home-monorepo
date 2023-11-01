export namespace Utils {
  export const bytesToMb = (bytes: number): number => {
    return Math.round(bytes / 1024 / 1024);
  };

  export const bytesToKb = (bytes: number): number => {
    return Math.round(bytes / 1024);
  };

  export const capFirst = (str: string) => {
    return str[0].toUpperCase() + str.slice(1);
  };
  export const platformIsWindows = (platform: string) => platform === "Windows" || platform === "win32";
}
