#!/bin/bash
set -eo pipefail

source ./common.sh

#
# No longer required, but useful to keep just in case we want to deploy
# changes in toolkit directly to the docker image
#
if [ -n "$INCLUDE_TOOLKIT" ]; then
	/bin/rm -rfv install/grafinsight-toolkit
	mkdir -pv install/grafinsight-toolkit
	cp -rv ../../bin install/grafinsight-toolkit
	cp -rv ../../src install/grafinsight-toolkit
	cp -v ../../package.json install/grafinsight-toolkit
	cp -v ../../tsconfig.json install/grafinsight-toolkit
fi

docker build -t ${DOCKER_IMAGE_NAME} .
docker push $DOCKER_IMAGE_NAME
docker tag ${DOCKER_IMAGE_NAME} ${DOCKER_IMAGE_BASE_NAME}:latest
docker push ${DOCKER_IMAGE_BASE_NAME}:latest

[ -n "$INCLUDE_TOOLKIT" ] && /bin/rm -rfv install/grafinsight-toolkit
