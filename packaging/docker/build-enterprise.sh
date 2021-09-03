#!/bin/sh
set -e

UBUNTU_BASE=0

while [ "$1" != "" ]; do
  case "$1" in
    "--ubuntu")
      UBUNTU_BASE=1
      echo "Ubuntu base image enabled"
      shift
      ;;
    * )
      # unknown param causes args to be passed through to $@
      break
      ;;
  esac
done

_raw_grafinsight_tag=$1
_docker_repo=${2:-grafinsight/grafinsight-enterprise}

if echo "$_raw_grafinsight_tag" | grep -q "^v"; then
  _grafinsight_tag=$(echo "${_raw_grafinsight_tag}" | cut -d "v" -f 2)
elif echo "$_raw_grafinsight_tag" | grep -q "^master-"; then
  _grafinsight_tag="master"
else
  _grafinsight_tag="${_raw_grafinsight_tag}"
fi

if [ ${UBUNTU_BASE} = "0" ]; then
  TAG_SUFFIX=""
  DOCKERFILE="Dockerfile"
else
  TAG_SUFFIX="-ubuntu"
  DOCKERFILE="ubuntu.Dockerfile"
fi

echo "Building and deploying ${_docker_repo}:${_grafinsight_tag}${TAG_SUFFIX}"

docker build \
  --tag "${_docker_repo}:${_grafinsight_tag}${TAG_SUFFIX}" \
  --no-cache=true \
  -f ${DOCKERFILE} \
  .

docker push "${_docker_repo}:${_grafinsight_tag}${TAG_SUFFIX}"

if echo "$_raw_grafinsight_tag" | grep -q "^v" && echo "$_raw_grafinsight_tag" | grep -qv "beta"; then
  docker tag "${_docker_repo}:${_grafinsight_tag}${TAG_SUFFIX}" "${_docker_repo}:latest${TAG_SUFFIX}"
  docker push "${_docker_repo}:latest${TAG_SUFFIX}"
fi


if echo "${_raw_grafinsight_tag}" | grep -q "^master-" && [ ${UBUNTU_BASE} = "1" ]; then
  docker tag "${_docker_repo}:${_grafinsight_tag}${TAG_SUFFIX}" "grafinsight/grafinsight-enterprise-dev:${_raw_grafinsight_tag}"
  docker push "grafinsight/grafinsight-enterprise-dev:${_raw_grafinsight_tag}"
fi

