#!/usr/bin/env bash
set -o pipefail

_basedir=$(dirname "$0")
cd "$_basedir" || exit

docker_build () {
	dockerfile=$1
	repo_file=$2
	package=$3
	tag=$4

  docker build -f "$dockerfile"          \
    --build-arg "REPO_CONFIG=$repo_file" \
    --build-arg "PACKAGE=$package"       \
    --tag "$tag"                         \
    --no-cache                           \
    .

    retval=$(docker run --rm "$tag" cat /usr/share/grafinsight/VERSION)
}

CHECK_BETA=$1
if [[ $CHECK_BETA == "beta" ]]; then
  # Testing deb repos
  docker_build "Dockerfile.deb" "deb-oss-beta.list" "grafinsight" "gf-oss-deb-repo-test"
  _oss_deb_v="$retval"

  docker_build "Dockerfile.deb" "deb-ee-beta.list" "grafinsight-enterprise" "gf-ee-deb-repo-test"
  _ee_deb_v="$retval"

  # Testing rpm repos
  docker_build "Dockerfile.rpm" "rpm-oss-beta.list" "grafinsight" "gf-oss-rpm-repo-test"
  _oss_rpm_v="$retval"

  docker_build "Dockerfile.rpm" "rpm-ee-beta.list" "grafinsight-enterprise" "gf-ee-rpm-repo-test"
  _ee_rpm_v="$retval"
else
  # Testing deb repos
  docker_build "Dockerfile.deb" "deb-oss-stable.list" "grafinsight" "gf-oss-deb-repo-test"
  _oss_deb_v="$retval"

  docker_build "Dockerfile.deb" "deb-ee-stable.list" "grafinsight-enterprise" "gf-ee-deb-repo-test"
  _ee_deb_v="$retval"

  # Testing rpm repos
  docker_build "Dockerfile.rpm" "rpm-oss-stable.list" "grafinsight" "gf-oss-rpm-repo-test"
  _oss_rpm_v="$retval"

  docker_build "Dockerfile.rpm" "rpm-ee-stable.list" "grafinsight-enterprise" "gf-ee-rpm-repo-test"
  _ee_rpm_v="$retval"
fi

echo Versions:
echo OSS deb = "${_oss_deb_v}"
echo OSS rpm = "${_oss_rpm_v}"
echo EE  deb = "${_ee_deb_v}"
echo EE  rpm = "${_ee_rpm_v}"
