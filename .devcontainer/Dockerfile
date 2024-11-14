FROM ubuntu:jammy 

# Install baseline packages
RUN apt-get update && \
  apt-get install --no-install-recommends --yes \
  apt-transport-https \
  bash \
  bash-completion \
  build-essential \
  ca-certificates \
  cifs-utils \
  curl \
  dmidecode \
  gnupg-agent \
  gpg \
  htop \
  iputils-ping \
  jq \
  lm-sensors \
  locales \
  man \
  nano \
  openssh-server \
  p7zip-full \
  pciutils \
  python3 \
  python3-pip \
  software-properties-common \
  sudo \
  unzip \
  util-linux \
  vim \
  wget \
  rsync \
  && \
  # Install latest Git using their official PPA
  add-apt-repository ppa:git-core/ppa && \
  apt-get install --yes git \
  && rm -rf /var/lib/apt/lists/*


RUN \
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -; \
  apt-key fingerprint 0EBFCD88; \
  add-apt-repository \
  "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) \
  stable"; \
  apt-get update && apt-get install --no-install-recommends -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin


ARG SSHD_PORT=2222
RUN \
  sed -i -e "s/^#\(Port\) .*/\1 ${SSHD_PORT}/" /etc/ssh/sshd_config; \
  mkdir /run/sshd; \
  systemctl enable ssh


# Make typing unicode characters in the terminal work.
ENV LANG=en_US.UTF-8
ARG USER=mgarnier


# Create a non-root user
RUN \
  useradd --groups sudo --no-create-home --shell /bin/bash ${USER}; \
  echo "${USER} ALL=(ALL) NOPASSWD:ALL" >/etc/sudoers.d/${USER}; \
  chmod 0440 /etc/sudoers.d/${USER}; \
  sudo usermod -aG docker ${USER}


# Switch to the non-root user
USER ${USER}
WORKDIR /home/${USER}

# Setup user home directory and SSH keys
COPY home ./
RUN \
  sudo chmod 644 ./.p10k.zsh; \
  sudo chown -R ${USER}:${USER} .ssh; \
  sudo chmod 700 .ssh; \
  sudo chmod 600 .ssh/*

# Zsh Installation
RUN sh -c "$(wget -O- https://github.com/deluan/zsh-in-docker/releases/download/v1.2.0/zsh-in-docker.sh)" && \
  sudo chsh -s /usr/bin/zsh $USER
SHELL ["/usr/bin/zsh", "-c"]

RUN echo "[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh" >> ~/.zshrc

# Tasks Installation
ARG INSTALL_TASKS=false
RUN \
  if [ "${INSTALL_TASKS}" = "true" ]; then \
  sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b .local; \
  echo "export PATH=\$PATH:~/.local" >> ~/.zshrc; \
  fi

# GO Installation
ARG INSTALL_GO=false
RUN \
  if [ "${INSTALL_GO}" = "true" ]; then \
  curl -o go.tar.gz https://dl.google.com/go/go1.23.0.linux-amd64.tar.gz; \
  sudo tar -C /usr/local -xzf go.tar.gz; \
  rm go.tar.gz; \
  echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.zshrc; \
  echo "Go Installed"; \
  fi

# NVM Installation
ARG INSTALL_NVM=false
RUN \
  if [ "${INSTALL_NVM}" = "true" ]; then \
  curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash; \
  echo "NVM Installed"; \
  source ~/.zshrc; \
  nvm install 20; \
  fi


# Ansible Installation
ARG INSTALL_ANSIBLE=false
RUN \
  if [ "${INSTALL_ANSIBLE}" = "true" ]; then \
  sudo add-apt-repository --yes --update ppa:ansible/ansible; \
  sudo DEBIAN_FRONTEND=noninteractive TZ=Europe/Paris apt-get install -y ansible; \
  echo "Ansible Installed"; \
  fi

# Copy entrypoint before setting up vscode, because vscode install script will modify the entrypoint
COPY --chmod=755 --chown=${USER}:${USER} ./entrypoint.sh /entrypoint.sh

ARG INSTALL_VSCODE=false
COPY --chmod=755 --chown=${USER}:${USER} ./setup-vscode /setup-vscode
RUN \
  if [ "${INSTALL_VSCODE}" = "true" ]; then \
  bash /setup-vscode/setup.sh "$(tr '\n' ' ' </setup-vscode/vscode-params.json)"; \
  fi


VOLUME [ "/mnt/dev" ]

RUN sudo mkdir -p /mnt/dev; \
  sudo chown -R ${USER}:${USER} /mnt/dev

COPY --chmod=755 --chown=${USER}:${USER} ./setup /setup


ENTRYPOINT [ "/entrypoint.sh" ]


HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
  CMD [ -f ~/.setup_complete ] && service ssh status > /dev/null 2>&1 || exit 1

