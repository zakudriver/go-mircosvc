workspace(name = "go_mircosvc")

load("//:load.bzl", "repositories")

repositories()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("//:repos.bzl", "go_repositories")

go_repositories()

# docker
load("@io_bazel_rules_docker//repositories:repositories.bzl", docker_rules_repositories = "repositories")

docker_rules_repositories()

#load("@io_bazel_rules_docker//repositories:deps.bzl", docker_rules_deps = "deps")
#
#docker_rules_deps()

load(
    "@io_bazel_rules_docker//go:image.bzl",
    go_image_repos = "repositories",
)

go_image_repos()
