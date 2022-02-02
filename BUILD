load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/gyturi1/go-bazel-buildfarm
# gazelle:proto disable
# gazelle:go_naming_convention import_alias
# gazelle:exclude tools/m10s_template
gazelle(
    name = "gen_buildfiles",
    command = "update",
)

gazelle(
    name = "update_repos",
    args = [
        "-from_file=go.mod",
        "-to_macro=deps.bzl%go_dependencies",
        "-prune=true",
        "-build_file_proto_mode=disable",
    ],
    command = "update-repos",
)

config_setting(
    name = "debug_build",
    values = {
        "compilation_mode": "dbg",
    },
)

#filegroup(
#    name = "gomodfile",
#    srcs = [
#        "go.mod",
#    ],
#    visibility = ["//visibility:public"],
#)
#
#filegroup(
#    name = "gosumfile",
#    srcs = [
#        "go.sum",
#    ],
#    visibility = ["//visibility:public"],
#)
