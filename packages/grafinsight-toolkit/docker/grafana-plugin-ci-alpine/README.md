# Using this docker image

Uploaded to dockerhub as grafinsight/grafinsight-plugin-ci:latest-alpine

Based off of `circleci/node:12-browsers` 

## User
The user will be `circleci`
The home directory will be `/home/circleci`

## Node
- node 12 is installed
- yarn is installed globally
- npm is installed globally

## Go
- Go is installed in `/usr/local/bin/go`
- golangci-lint is installed in `/usr/local/bin/golangci-lint`
- mage is installed in `/home/circleci/go/bin/mage`

All of the above directories are in the path, so there is no need to specify fully qualified paths.

## GrafInsight
- Installed in `/home/circleci/src/grafinsight`
- `yarn install` has been run

## Integration/Release Testing
There are 4 previous versions pre-downloaded to /usr/local/grafinsight. These versions are:
1. 6.6.2
2. 6.5.3
3. 6.4.5
4. 6.3.7

To test, your CircleCI config will need a run section with something similar to the following
```
- run:
        name: Setup GrafInsight (local install)
        command: |
          sudo dpkg -i /usr/local/grafinsight/deb/grafinsight_6.6.2_amd64.deb
          sudo cp ci/grafinsight-test-env/custom.ini /usr/share/grafinsight/conf/custom.ini
          sudo cp ci/grafinsight-test-env/custom.ini /etc/grafinsight/grafinsight.ini
          sudo service grafinsight-server start
          grafinsight-cli --version
```


# Building
To build, cd to `<srcroot>/packages/grafinsight-toolkit/docker/grafinsight-plugin-ci-alpine`
```
./build.sh
```

# Developing/Testing
To test, you should have docker-compose installed.
```
cd test
./start.sh
```

You will be in /home/circleci/test with the buildscripts installed to the local directory.
Do your edits/run tests. When saving, your edits will be available in the container immediately.
