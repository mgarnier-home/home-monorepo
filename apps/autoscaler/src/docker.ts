import Dockerode from 'dockerode';

import { Utils } from '@libs/utils';

import { config } from './utils/config';
import { DockerHost } from './utils/interfaces';

interface DockerContainer {
  container: Dockerode.Container;
  infos: Dockerode.ContainerInfo;
}
const getDockerApi = (host: DockerHost): Dockerode => {
  return new Dockerode({
    protocol: 'ssh',
    host: host.ip,
    port: host.sshPort,
    username: config.sshUser,
    sshOptions: {
      privateKey: config.sshPrivateKey,
    },
  });
};

const getRunnerContainer = async (host: DockerHost, jobId: number): Promise<DockerContainer | null> => {
  const dockerApi = getDockerApi(host);

  const containers = await dockerApi.listContainers({
    all: true,
    filters: {
      label: [`autoscaler.runner=${jobId}`],
    },
  });

  if (containers.length === 0) {
    return null;
  }

  return {
    container: dockerApi.getContainer(containers[0].Id),
    infos: containers[0],
  };
};

const removeContainer = async (container: DockerContainer): Promise<void> => {
  await container.container.remove({ force: true });
};

export const startRunner = async (host: DockerHost, jobId: number): Promise<void> => {
  const dockerApi = getDockerApi(host);

  const oldContainer = await getRunnerContainer(host, jobId);

  if (oldContainer) {
    await removeContainer(oldContainer);
  }

  const newContainer = await dockerApi.createContainer({
    Image: config.runnerImage,
    name: `runner-${jobId}`,
    Env: [
      `ACCESS_TOKEN=${config.runnerAccessToken}`,
      `RUNNER_SCOPE=org`,
      `ORG_NAME=${config.runnerOrgName}`,
      `REPO_URL=${config.runnerRepoUrl}`,
      `RUNNER_NAME=runner-${host.label.replace('/', '-')}-${jobId}`,
      `LABELS=${host.label}`,
      `EPHEMERAL=true`,
      `RUNNER_WORKDIR=/tmp`,
    ],
    HostConfig: {
      Binds: ['/var/run/docker.sock:/var/run/docker.sock'],
    },
    Labels: {
      'autoscaler.runner': jobId.toString(),
    },
  });

  await newContainer.start();
};

export const stopRunner = async (host: DockerHost, jobId: number): Promise<void> => {
  await Utils.sleep(2000);

  const container = await getRunnerContainer(host, jobId);

  if (container) {
    await removeContainer(container);
  }
};
