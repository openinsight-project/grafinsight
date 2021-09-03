#!/bin/bash

# shellcheck source=./scripts/helpers/exit-if-fail.sh
source "$(dirname "$0")/helpers/exit-if-fail.sh"

start=$(date +%s)

export TEST_MAX_WORKERS=2

/tmp/grabpl test-frontend --github-token "${GITHUB_GRAFINSIGHTBOT_TOKEN}" "$@"

end=$(date +%s)
seconds=$((end - start))

exit_if_fail ./scripts/ci-frontend-metrics.sh

if [ "${CIRCLE_BRANCH}" == "master" ]; then
	exit_if_fail ./scripts/ci-metrics-publisher.sh grafinsight.ci-performance.frontend-tests=$seconds
fi

