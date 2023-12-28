import { HwCpu, HwGpu, HwNetwork, HwRam } from 'nodesight-types';

import { config } from './utils/config.js';
import { Current } from './utils/interfaces.js';

const getUrl = (type: string) => `${config.statsApiUrl}/api/${config.hostname}/${type}`;

const postToStatsApi = async (type: string, body: string) => {
  try {
    const response = await fetch(getUrl(type), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body,
    });

    const text = await response.text();

    return { code: response.status, text, type };
  } catch (e: any) {
    console.error(e);
    return { code: 0, text: e.message, type };
  }
};

const sendGpuToStatsApi = (gpu: HwGpu.Load) => {
  if (gpu.gpus.length > 0) {
    return postToStatsApi('gpu', JSON.stringify(gpu));
  } else {
    return Promise.resolve({ code: 200, text: 'No GPU', type: 'gpu' });
  }
};

const sendCpuToStatsApi = (cpu: HwCpu.Load) => {
  if (cpu.cores.length > 0) {
    return postToStatsApi('cpu', JSON.stringify(cpu));
  } else {
    return Promise.resolve({ code: 200, text: 'No CPU', type: 'cpu' });
  }
};

const sendRamToStatsApi = (ram: HwRam.Load) => {
  return postToStatsApi('ram', JSON.stringify(ram));
};

const sendNetworkToStatsApi = (network: HwNetwork.Load) => {
  return postToStatsApi('network', JSON.stringify(network));
};

export const sendToStatsApi = async (data: Current): Promise<boolean> => {
  try {
    console.log(`Sending data Stats api`);

    const responses = await Promise.all([
      sendCpuToStatsApi(data.cpu),
      sendGpuToStatsApi(data.gpu),
      sendRamToStatsApi(data.ram),
      sendNetworkToStatsApi(data.network),
    ]);

    console.log(responses);

    return responses.every((response) => response.code === 200);
  } catch (e) {
    console.error(e);

    return false;
  }
};
