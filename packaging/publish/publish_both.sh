#! /usr/bin/env bash
version=5.4.2

# wget https://dl.grafinsight.com/oss/release/grafinsight_${version}_amd64.deb
#
# package_cloud push grafinsight/stable/debian/jessie grafinsight_${version}_amd64.deb
# package_cloud push grafinsight/stable/debian/wheezy grafinsight_${version}_amd64.deb
# package_cloud push grafinsight/stable/debian/stretch grafinsight_${version}_amd64.deb
#
# package_cloud push grafinsight/testing/debian/jessie grafinsight_${version}_amd64.deb
# package_cloud push grafinsight/testing/debian/wheezy grafinsight_${version}_amd64.deb --verbose
# package_cloud push grafinsight/testing/debian/stretch grafinsight_${version}_amd64.deb --verbose

wget https://dl.grafinsight.com/oss/release/grafinsight-${version}-1.x86_64.rpm

package_cloud push grafinsight/testing/el/6 grafinsight-${version}-1.x86_64.rpm --verbose
package_cloud push grafinsight/testing/el/7 grafinsight-${version}-1.x86_64.rpm --verbose

package_cloud push grafinsight/stable/el/7 grafinsight-${version}-1.x86_64.rpm --verbose
package_cloud push grafinsight/stable/el/6 grafinsight-${version}-1.x86_64.rpm --verbose

rm grafinsight*.{deb,rpm}
