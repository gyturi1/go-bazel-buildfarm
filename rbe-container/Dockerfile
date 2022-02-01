FROM debian:bullseye-slim

RUN export DEBIAN_FRONTEND=noninteractive && apt-get -y update \
    && apt-get -y --no-install-recommends install \
        #base bazel deps
        openjdk-11-jdk-headless clang python3-dev\
        #from https://github.com/GoogleCloudPlatform/layer-definitions/tree/master/layers/ubuntu1804/rbe_tools
        binutils ca-certificates curl ed file git less locales locales-all netcat openssh-client patch unzip wget zip \
    #clean up
    && rm -rf /var/lib/apt/lists/* \
    && ln -s /usr/bin/python3 /usr/bin/python

RUN export DEBIAN_FRONTEND=noninteractive && apt-get -y update \
    && apt-get -y --no-install-recommends install \
        apt-transport-https ca-certificates gnupg lsb-release \
    #docker cli
    && curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg \
    && echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian \
        $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null \
    && apt-get -q -y update \
    && apt-get -y --no-install-recommends install docker-ce-cli \
    #clean up
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p ~/.docker/cli-plugins/ \
    && curl -SL https://github.com/docker/compose/releases/download/v2.2.3/docker-compose-linux-x86_64 -o ~/.docker/cli-plugins/docker-compose \
    && chmod +x ~/.docker/cli-plugins/docker-compose

ENV JAVA_HOME=/usr/lib/jvm/java-11-openjdk-amd64