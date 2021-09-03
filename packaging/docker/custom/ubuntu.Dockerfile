ARG GRAFINSIGHT_VERSION="latest"

FROM grafinsight/grafinsight:${GRAFINSIGHT_VERSION}-ubuntu

USER root

# Set DEBIAN_FRONTEND=noninteractive in environment at build-time
ARG DEBIAN_FRONTEND=noninteractive

ARG GF_INSTALL_IMAGE_RENDERER_PLUGIN="false"

ARG GF_GID="0"
ENV GF_PATHS_PLUGINS="/var/lib/grafinsight-plugins"

RUN mkdir -p "$GF_PATHS_PLUGINS" && \
    chown -R grafinsight:${GF_GID} "$GF_PATHS_PLUGINS"

RUN if [ $GF_INSTALL_IMAGE_RENDERER_PLUGIN = "true" ]; then \
    apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y gdebi-core && \
    cd /tmp && \
    curl -LO https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb && \
    gdebi --n google-chrome-stable_current_amd64.deb && \
    apt-get autoremove -y && \
    rm -rf /var/lib/apt/lists/*; \
fi

USER grafinsight

ENV GF_RENDERER_PLUGIN_CHROME_BIN="/usr/bin/google-chrome"

RUN if [ $GF_INSTALL_IMAGE_RENDERER_PLUGIN = "true" ]; then \
    grafinsight-cli \
        --pluginsDir "$GF_PATHS_PLUGINS" \
        --pluginUrl https://github.com/grafinsight/grafinsight-image-renderer/releases/latest/download/plugin-linux-x64-glibc-no-chromium.zip \
        plugins install grafinsight-image-renderer; \
fi

ARG GF_INSTALL_PLUGINS=""

RUN if [ ! -z "${GF_INSTALL_PLUGINS}" ]; then \
    OLDIFS=$IFS; \
        IFS=','; \
    for plugin in ${GF_INSTALL_PLUGINS}; do \
        IFS=$OLDIFS; \
        grafinsight-cli --pluginsDir "$GF_PATHS_PLUGINS" plugins install ${plugin}; \
    done; \
fi
