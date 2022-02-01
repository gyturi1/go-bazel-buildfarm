load("@io_bazel_rules_docker//container:container.bzl", "container_image", "container_layer")

def test_in_docker_compose(
        name, 
        dataFiles, 
        goTest = "component-test_test",
        composeFiles = ["docker-compose/docker-compose.yaml", "docker-compose/docker-compose-testrunner.yaml"], 
        envFile = "docker-compose/.env", 
        baseImage = "@debian_bullseye_slim//image",
        size = "small"):
    
    container_layer(
        name = name + "_layer_test",
        testonly = True,
        files = [goTest],
        visibility = ["//visibility:private"],
    )

    container_image(
        name = name + "_test_runner_image",
        testonly = True,
        base = "@debian_bullseye_slim//image",
        entrypoint = ["/{name}".format(name = goTest.rpartition(":")[2])],
        layers = [
            ":{x}".format( x = name + "_layer_test"),
        ],
        visibility = ["//visibility:private"],
    )

    native.sh_test(
        name = name + "_with_docker_compose",
        size = "small",
        srcs = ["//tools/test_rules:run-test-in-docker-compose.sh"],
        args = [
            "$(rootpath //tools/test_rules:run-test-in-docker-compose.sh)",
            "$(rootpath :{x})".format(x = name + "_test_runner_image"),
        ] + ["$(location %s)" % d for d in composeFiles],
        data = [
            envFile,
            ":{x}".format(x = name + "_test_runner_image"),
        ] + composeFiles + dataFiles,
        env = {
            "test_runner_image": "bazel/{package}:{x}".format(package =  native.package_name(), x = name + "_test_runner_image"),
        },
        tags = ["no-sandbox"],
        visibility = ["//visibility:public"],
    )


def test_in_docker(name, servieUnderTest, goTest = ":it-test_test", baseImage = "@debian_bullseye_slim//image", size = "small"):
    container_layer(
        name = name + "_layer_sut",
        testonly = True,
        files = [servieUnderTest],
        visibility = ["//visibility:private"],
    )

    container_layer(
        name = name + "_layer_test",
        testonly = True,
        files = [goTest],
        visibility = ["//visibility:private"],
    )

    container_layer(
        name = name + "_layer_test_wrapper_sh",
        testonly = True,
        files = ["//tools/test_rules:wrapper.sh"],
        visibility = ["//visibility:private"],
    )

    container_image(
        name = name + "_test_runner_image",
        testonly = True,
        base = baseImage,
        entrypoint = ["/wrapper.sh"],
        layers = [
            ":{x}".format( x = name + "_layer_sut"),
            ":{x}".format( x = name + "_layer_test"),
            ":{x}".format( x = name + "_layer_test_wrapper_sh"),
        ],
        visibility = ["//visibility:private"],
    )

    native.sh_test(
        name = name + "_with_docker",
        size = size,
        srcs = ["//tools/test_rules:run-test-in-docker.sh"],
        args = [
            "$(rootpath :{x})".format(x = name + "_test_runner_image"),
        ],
        data = [":{x}".format(x = name + "_test_runner_image")],
        visibility = ["//visibility:public"],
        tags = ["no-sandbox"],
    )