SHELL = /bin/bash

##################################################################
#
# Daily targets
#
##################################################################

.PHONY: all
all: update_deps gen_buildfiles build test

.PHONY: gen_buildfiles
gen_buildfiles:
	@echo "Generating build files from go source/test files using Gazelle"
	@echo "See: https://github.com/bazelbuild/bazel-gazelle#fix-and-update"
	@bazel $(BAZEL_ROOT_OPTS) run //:gen_buildfiles $(BAZEL_OPTS)

.PHONY: update_deps
update_deps:
	@echo "Generating deps.bzl from go.mod using Gazelle"
	@echo "See: https://github.com/bazelbuild/bazel-gazelle#update-repos"
	@bazel $(BAZEL_ROOT_OPTS) run //:update_repos $(BAZEL_OPTS)

.PHONY: build
build:
	@bazel $(BAZEL_ROOT_OPTS) build ... $(BAZEL_OPTS)

.PHONY: test
test:
	@bazel $(BAZEL_ROOT_OPTS) test ... $(BAZEL_OPTS)

##################################################################
#
# Remote execution
#
##################################################################

.PHONY: all_remote 
all_remote:
	@$(MAKE) all BAZEL_OPTS=--config="dev-remote $(BAZEL_OPTS)"

rbe_image_tag := v2
rbe_image := localhost:5000/test-rbe-container
.PHONY: rbe_container
rbe_container: run_local_docker_registry
	@echo "This is the container in which the remote worker will be running"
	@echo "Google has official rbe container: https://console.cloud.google.com/marketplace/details/google/rbe-ubuntu16-04"
	@echo "But you can build your own of course"

	@cd rbe-container && docker build -t $(rbe_image):$(rbe_image_tag) .
	@docker push $(rbe_image):$(rbe_image_tag)
	@echo "DONT FORGET to copy THE DIGEST into"
	@echo "\t- $$(pwd)/.bazelrc (--experimental_docker_image flag)"
	@echo "\t- $$(pwd)/.rbe-config/config/BUILD (exec_properties section)"
	@echo "\t- $(bazel_buildfarm_repo_root)/examples/worker.config.example (container-image properties section)"
	@echo "\t- $(bazel_buildfarm_repo_root)/images.bzl (rbe_image_base digest)"
	@echo "If you does not have an rbe-config/config/BUILD generate it with rbconfig"


bazel_toolchain_repo_root := ~/git/bazel-toolchains
.PHONY: rbconfig
rbconfig:
	@echo "This is the toolchain config generated with: https://github.com/bazelbuild/bazel-toolchains"
	@echo "In order to smoothly use remote execution you should configure bazel toolchains"
	@cd $(bazel_toolchain_repo_root) && go build -o rbe_configs_gen ./cmd/rbe_configs_gen/rbe_configs_gen.go
	@$(bazel_toolchain_repo_root)/rbe_configs_gen \
    	--bazel_version=4.2.1 \
    	--toolchain_container=$$(docker images --filter "reference=$(rbe_image)" --format="{{.Digest}}") \
    	--output_src_root=$$(pwd) \
    	--output_config_path=rbe-config \
    	--exec_os=linux \
    	--target_os=linux

bazel_buildfarm_repo_root := ~/git/bazel-buildfarm
.PHONY: run_buildfarm_worker
run_buildfarm_worker:
	@echo "Start the worker inside the rbe_container and use the appropriate worker config"
	@echo "See: https://bazelbuild.github.io/bazel-buildfarm/docs/architecture/worker-execution-environment/"
	@stat $(bazel_buildfarm_repo_root) &> /dev/null || (echo please check out bazel-buildfarm repo && exit 1)
	@cd $(bazel_buildfarm_repo_root) && (grep buildfarm-worker-test BUILD || (echo "apply patch from buidfarm_patch to bazel buildfarm repo and replace the digest!" && exit 1))
	@cd $(bazel_buildfarm_repo_root) && bazel run //:buildfarm-worker-test -- --norun
	@cd $(bazel_buildfarm_repo_root) \
		&& docker run --net=host --rm --name=buildfarm-worker \
			-v $$PWD/examples:/config \
			-v /var/run/docker.sock:/var/run/docker.sock \
			-e JAVA_OPTS="-Djava.util.logging.config.file=/config/debug.logging.properties" \
			bazel:buildfarm-worker-test /config/worker.config.example

.PHONY: run_buildfarm_server
run_buildfarm_server:
	@echo "Buildfarm implements remote execution api, so can be used by bazel (as a client) to send remote excution request"
	@stat $(bazel_buildfarm_repo_root) &> /dev/null || (echo please check out bazel-buildfarm repo && exit 1)
	@cd $(bazel_buildfarm_repo_root) \
		&& bazel run //src/main/java/build/buildfarm:buildfarm-server -- --jvm_flag=-Djava.util.logging.config.file=$$PWD/examples/debug.logging.properties $$PWD/examples/server.config.example

.PHONY: run_local_docker_registry
run_local_docker_registry:
	@echo "To build the buildfarm worker with rbe_container, the image needs to be pushed into a registry"
	@echo "See: https://bazelbuild.github.io/bazel-buildfarm/docs/architecture/worker-execution-environment/"
	@docker volume create local_registry_volume
	@docker start local_regsitry || docker run --net=host --name local_registry -v local_registry_volume:/var/lib/registry registry:2

##################################################################
#
# Running bazel inside a build container
#
##################################################################

.PHONY: build_isolated_runner_bullseye
build_isolated_runner_bullseye:
	@echo "This is a container for using locally to prepare for remote execution"
	@cd build-container && docker build -t $(build_image):bullseye .

.PHONY: build_isolated_runner_alpine
build_isolated_runner_alpine:
	@echo "This is a container for using locally to prepare for remote execution, using alpine to isolate more"
	@echo "NOTE libcompat6 is installed inside, because https://go.dev/dl/ provides only libc dinamically linked version"
	@echo "rules_go downloads the Go SDK from https://go.dev/dl/, but in alpine it won't work"
	@echo "TODO: Better download a Musl based go toolchain with https://github.com/bazelbuild/rules_go/blob/master/go/toolchains.rst#go_download_sdk"
	@cd build-container && docker build -t $(build_image):alpine -f Dockerfile.alpine .

build_image := localhost:5000/isolated-runner
TAG := bullseye
ISOLATED_MAKE_TARGET := all
.PHONY: all_isolated
all_isolated: build_container_$(TAG)
	@echo "This is a test preparing for remote execution and/or debugging failures in remote execution"
	@echo "See: https://docs.bazel.build/versions/main/remote-execution-sandbox.html"
	@docker run --rm -it --name=remote-test-build-container \
			--net=host \
			-v /var/tmp:/var/tmp \
			-v /var/run/docker.sock:/var/run/docker.sock \
			-v $$(pwd):/src \
			-w /src \
			$(build_image):$(TAG) \
			bash -c 'make $(ISOLATED_MAKE_TARGET) BAZEL_ROOT_OPTS="--output_user_root=/var/tmp/bazel_docker_root"'

##################################################################
#
#Running bazel in docker-sandbox mode
#
##################################################################

.PHONY: all_in_docker_sandbox
all_in_docker_sandbox:
	@echo "This is a test preparing for remote execution and/or debugging failures in remote execution"
	@echo "See: https://docs.bazel.build/versions/main/remote-execution-sandbox.html"
	@$(MAKE) all BAZEL_OPTS=--config="dev-docker-sandbox --verbose_failures $(BAZEL_OPTS)" BAZEL_ROOT_OPTS=$(BAZEL_ROOT_OPTS)
