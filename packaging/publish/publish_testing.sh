#! /usr/bin/env bash
deb_ver=5.1.0-beta1
rpm_ver=5.1.0-beta1

wget https://s3-us-west-2.amazonaws.com/grafinsight-releases/release/grafinsight_${deb_ver}_amd64.deb

package_cloud push grafinsight/testing/debian/jessie grafinsight_${deb_ver}_amd64.deb
package_cloud push grafinsight/testing/debian/wheezy grafinsight_${deb_ver}_amd64.deb
package_cloud push grafinsight/testing/debian/stretch grafinsight_${deb_ver}_amd64.deb

wget https://s3-us-west-2.amazonaws.com/grafinsight-releases/release/grafinsight-${rpm_ver}.x86_64.rpm

package_cloud push grafinsight/testing/el/6 grafinsight-${rpm_ver}.x86_64.rpm
package_cloud push grafinsight/testing/el/7 grafinsight-${rpm_ver}.x86_64.rpm

rm grafinsight*.{deb,rpm}
