import { config } from './config';

export const log = (...args: any[]) => {
  if (config.devMode) {
    return;
  }

  console.log(...args);
};

export const time = (label: string) => {
  if (config.devMode) {
    return;
  }

  console.time(label);
};

export const timeEnd = (label: string) => {
  if (config.devMode) {
    return;
  }

  console.timeEnd(label);
};
