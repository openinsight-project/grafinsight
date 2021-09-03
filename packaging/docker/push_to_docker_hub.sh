#!/bin/sh
set -e

UBUNTU_BASE=0
TAG_SUFFIX=""

while [ "$1" != "" ]; do
  case "$1" in
    "--ubuntu")
      UBUNTU_BASE=1
      TAG_SUFFIX="-ubuntu"
      echo "Ubuntu base image enabled"
      shift
      ;;
    * )
      # unknown param causes args to be passed through to $@
      break
      ;;
  esac
done

_grafinsight_tag=${1:-}
_docker_repo=${2:-grafinsight/grafinsight}

# If the tag starts with v, treat this as an official release
if echo "$_grafinsight_tag" | grep -q "^v"; then
	_grafinsight_version=$(echo "${_grafinsight_tag}" | cut -d "v" -f 2)
else
	_grafinsight_version=$_grafinsight_tag
fi

export DOCKER_CLI_EXPERIMENTAL=enabled

echo "pushing ${_docker_repo}:${_grafinsight_version}${TAG_SUFFIX}"

docker_push_all () {
	repo=$1
	tag=$2

  # Push each image individually
  docker push "${repo}:${tag}${TAG_SUFFIX}"
  docker push "${repo}-arm32v7-linux:${tag}${TAG_SUFFIX}"
  docker push "${repo}-arm64v8-linux:${tag}${TAG_SUFFIX}"

  # Create and push a multi-arch manifest
  docker manifest create "${repo}:${tag}${TAG_SUFFIX}" \
    "${repo}:${tag}${TAG_SUFFIX}" \
    "${repo}-arm32v7-linux:${tag}${TAG_SUFFIX}" \
    "${repo}-arm64v8-linux:${tag}${TAG_SUFFIX}"

  docker manifest push "${repo}:${tag}${TAG_SUFFIX}"
}

if echo "$_grafinsight_tag" | grep -q "^v" && echo "$_grafinsight_tag" | grep -vq "beta"; then
	echo "pushing ${_docker_repo}:latest${TAG_SUFFIX}"
	docker_push_all "${_docker_repo}" "latest"
	docker_push_all "${_docker_repo}" "${_grafinsight_version}"
	# Push to the grafinsight-dev repository with the expected tag
	# for running the end to end tests successfully
  docker push "grafinsight/grafinsight-dev:${_grafinsight_tag}${TAG_SUFFIX}"
elif echo "$_grafinsight_tag" | grep -q "^v" && echo "$_grafinsight_tag" | grep -q "beta"; then
	docker_push_all "${_docker_repo}" "${_grafinsight_version}"
	# Push to the grafinsight-dev repository with the expected tag
	# for running the end to end tests successfully
  docker push "grafinsight/grafinsight-dev:${_grafinsight_tag}${TAG_SUFFIX}"
elif echo "$_grafinsight_tag" | grep -q "master"; then
	docker_push_all "${_docker_repo}" "master"
  docker push "grafinsight/grafinsight-dev:${_grafinsight_version}${TAG_SUFFIX}"
fi
