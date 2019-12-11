load(
    "@io_bazel_rules_docker//container:container.bzl",
    "container_pull",
)

def docker_repositories():
    container_pull(
        name = "etcd_docker",
        registry = "index.docker.io",
        repository = "xieyanze/etcd3",
        tag = "latest",
    )

    container_pull(
        name = "openzipkin_docker",
        registry = "index.docker.io",
        repository = "openzipkin/zipkin",
        tag = "latest",
    )
