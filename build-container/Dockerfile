FROM debian:bullseye-slim

ARG BAZEL_VERSION=4.2.1

RUN export DEBIAN_FRONTEND=noninteractive && apt-get update \
    && apt-get install -y apt-transport-https curl software-properties-common git gcc gnupg2 g++ openjdk-11-jdk-headless python-dev zip wget xz-utils \
    && apt-get install -y make bash-completion parallel make \
    && curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add - \
    && add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/debian $(lsb_release -cs) stable" \
    && apt-get update && apt-get install -y docker-ce \
    && rm -rf /var/lib/apt/lists/*

RUN wget https://github.com/bazelbuild/bazel/releases/download/${BAZEL_VERSION}/bazel-${BAZEL_VERSION}-installer-linux-x86_64.sh -O ./bazel-installer.sh && chmod 755 ./bazel-installer.sh \
    && ./bazel-installer.sh \
    && rm bazel-installer.sh