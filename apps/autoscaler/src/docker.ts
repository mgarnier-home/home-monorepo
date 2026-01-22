import Dockerode, { MountSettings } from 'dockerode';

import { Utils } from '@libs/utils';

import { config } from './utils/config';
import { DockerHost } from './utils/interfaces';
import { logger } from '@libs/logger';

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

  await pullImage(dockerApi, config.runnerImage);

  await dockerApi.createVolume({ Name: 'github-runner-docker-cache' });
  await dockerApi.createVolume({ Name: 'github-actions-runner-files' });
  const mounts: MountSettings[] = [];

  if (config.runtime !== '') {
    mounts.push({
      Target: '/var/lib/docker',
      Source: 'github-runner-docker-cache',
      Type: 'volume',
    });
  }

  mounts.push({
    Target: '/actions-runner',
    Source: 'github-actions-runner-files',
    Type: 'volume',
  });

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
      `DOCKER_REGISTRY_URL=${config.dockerRegistryUrl}`,
      `DOCKER_REGISTRY_USERNAME=${config.dockerRegistryUsername}`,
      `DOCKER_REGISTRY_PASSWORD=${config.dockerRegistryPassword}`,
      'START_DOCKER_SERVICE=true',
      'CONFIGURED_ACTIONS_RUNNER_FILES_DIR=/actions-runner',
      'DISABLE_AUTOMATIC_DEREGISTRATION=true',
    ],
    HostConfig: {
      // Quand on est en mode production (sur le serveur), on utilise sysbox-runc pour gérer les conteneurs imbriqués
      // afin qu'on ne puisse pas accéder aux autres conteneurs du host
      Runtime: config.runtime !== '' ? config.runtime : undefined,
      // Quand on est en mode développement, on bind le socket docker de l'hôte pour permettre au runner d'utiliser docker sans avoir besoin d'installer sysbox-runc
      Binds: config.runtime === '' ? ['/var/run/docker.sock:/var/run/docker.sock'] : [],
      Mounts: mounts,
    },
    Labels: {
      'autoscaler.runner': jobId.toString(),
    },
  });

  await newContainer.start();
  logger.info(`Started new runner container for job ${jobId} on host ${host.label}`);
};

export const stopRunner = async (host: DockerHost, jobId: number): Promise<void> => {
  await Utils.sleep(2000);

  const container = await getRunnerContainer(host, jobId);

  if (container) {
    await removeContainer(container);
  }
};

const pullImage = async (dockerApi: Dockerode, imageName: string): Promise<void> => {
  try {
    logger.info(`Pulling image: ${imageName}`);

    const stream = await dockerApi.pull(imageName, {
      authconfig: {
        username: config.dockerRegistryUsername,
        password: config.dockerRegistryPassword,
        serveraddress: config.dockerRegistryUrl,
      },
    });

    await new Promise<void>((resolve, reject) => {
      dockerApi.modem.followProgress(stream, (err, _) => {
        if (err) {
          reject(err);
        } else {
          logger.info(`Successfully pulled image: ${imageName}`);
          resolve();
        }
      });
    });
  } catch (error) {
    logger.error(`Failed to pull image ${imageName}:`, error);
    throw error;
  }
};
