import * as si from 'systeminformation';

import { HwGpu } from '@libs/nodesight-types';

const normalizeGpuBrand = (brand: string) => {
  return brand ? brand.replace(/(corporation)/gi, '').trim() : undefined;
};

const normalizeGpuName = (name: string) => {
  return name ? name.replace(/(nvidia|amd|intel)/gi, '').trim() : undefined;
};

const normalizeGpuModel = (model: string) => {
  return model ? model.replace(/\[.*\]/gi, '').trim() : undefined;
};

const sortControllers = (
  list: si.Systeminformation.GraphicsControllerData[]
): si.Systeminformation.GraphicsControllerData[] => {
  return list.sort((a, b) => (a.name ?? '').localeCompare(b.name ?? '') ?? 0);
};

const getStaticInfo = async (): Promise<HwGpu.Static> => {
  const gpuInfo = await si.graphics();

  console.log(gpuInfo);

  return {
    gpus: sortControllers(gpuInfo.controllers).map((controller, index) => ({
      index,
      brand: normalizeGpuBrand(controller.vendor) ?? '',
      model: normalizeGpuName(controller.name ?? '') ?? normalizeGpuModel(controller.model) ?? '',
      memoryTotal: controller.memoryTotal ?? controller.vram ?? 0,
    })),
  };
};

export class Gpu {
  public static staticInfo = getStaticInfo();

  public static current = async (): Promise<HwGpu.Load> => {
    const gpuInfo = await si.graphics();

    console.log(gpuInfo);

    return {
      gpus: sortControllers(gpuInfo.controllers).map((controller, index) => ({
        index,
        model: normalizeGpuName(controller.name ?? '') ?? normalizeGpuModel(controller.model) ?? '',
        load: controller.utilizationGpu ?? 0,
        memoryUsed: controller.memoryUsed ?? 0,
        powerDraw: controller.powerDraw ?? 0,
        temp: controller.temperatureGpu ?? 0,
      })),
    };
  };
}
