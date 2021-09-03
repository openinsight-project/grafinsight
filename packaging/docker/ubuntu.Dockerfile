ARG BASE_IMAGE=ubuntu:20.04
FROM ${BASE_IMAGE} AS grafinsight-builder

ARG GRAFINSIGHT_TGZ="grafinsight-latest.linux-x64.tar.gz"

COPY ${GRAFINSIGHT_TGZ} /tmp/grafinsight.tar.gz

RUN mkdir /tmp/grafinsight && tar xzf /tmp/grafinsight.tar.gz --strip-components=1 -C /tmp/grafinsight

FROM ${BASE_IMAGE}

EXPOSE 3000

# Set DEBIAN_FRONTEND=noninteractive in environment at build-time
ARG DEBIAN_FRONTEND=noninteractive
ARG GF_UID="472"
ARG GF_GID="0"

ENV PATH=/usr/share/grafinsight/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin \
    GF_PATHS_CONFIG="/etc/grafinsight/grafinsight.ini" \
    GF_PATHS_DATA="/var/lib/grafinsight" \
    GF_PATHS_HOME="/usr/share/grafinsight" \
    GF_PATHS_LOGS="/var/log/grafinsight" \
    GF_PATHS_PLUGINS="/var/lib/grafinsight/plugins" \
    GF_PATHS_PROVISIONING="/etc/grafinsight/provisioning"

WORKDIR $GF_PATHS_HOME

# Install dependencies
# We need curl in the image
RUN apt-get update && apt-get install -y ca-certificates curl tzdata && \
    apt-get autoremove -y && rm -rf /var/lib/apt/lists/*;

COPY --from=grafinsight-builder /tmp/grafinsight "$GF_PATHS_HOME"

RUN if [ ! $(getent group "$GF_GID") ]; then \
      addgroup --system --gid $GF_GID grafinsight; \
    fi

RUN export GF_GID_NAME=$(getent group $GF_GID | cut -d':' -f1) && \
    mkdir -p "$GF_PATHS_HOME/.aws" && \
    adduser --system --uid $GF_UID --ingroup "$GF_GID_NAME" grafinsight && \
    mkdir -p "$GF_PATHS_PROVISIONING/datasources" \
             "$GF_PATHS_PROVISIONING/dashboards" \
             "$GF_PATHS_PROVISIONING/notifiers" \
             "$GF_PATHS_PROVISIONING/plugins" \
             "$GF_PATHS_LOGS" \
             "$GF_PATHS_PLUGINS" \
             "$GF_PATHS_DATA" && \
    cp "$GF_PATHS_HOME/conf/sample.ini" "$GF_PATHS_CONFIG" && \
    cp "$GF_PATHS_HOME/conf/ldap.toml" /etc/grafinsight/ldap.toml && \
    chown -R "grafinsight:$GF_GID_NAME" "$GF_PATHS_DATA" "$GF_PATHS_HOME/.aws" "$GF_PATHS_LOGS" "$GF_PATHS_PLUGINS" "$GF_PATHS_PROVISIONING" && \
    chmod -R 777 "$GF_PATHS_DATA" "$GF_PATHS_HOME/.aws" "$GF_PATHS_LOGS" "$GF_PATHS_PLUGINS" "$GF_PATHS_PROVISIONING"

COPY ./run.sh /run.sh

USER "$GF_UID"
ENTRYPOINT [ "/run.sh" ]
