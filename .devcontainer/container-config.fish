if status is-interactive
    set -x ZEPHYR_HOST ssh://mgarnier@100.64.98.97
    set -x ATLAS_HOST ssh://mgarnier@100.64.98.98
    set -x APOLLON_HOST ssh://mgarnier@100.64.98.99
    set -x ATHENA_HOST ssh://mgarnier@100.64.98.100

    set -x COMPOSE_DIR /mnt/docker-data/zephyr/orchestrator
    set -x ENV_DIR /mnt/docker-data/zephyr/orchestrator

    home-cli completion fish | source
    task --completion fish | source
end
