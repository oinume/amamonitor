#!/bin/bash

set -e

cat<<EOF
============================================================
Usage: ${0} <version> <download_path>
This script installs gcloud command. You can modify gcloud version with arguments.
============================================================
EOF


GOOGLE_CLOUD_SDK_VERSION=${1:-206.0.0}
GOOGLE_CLOUD_SDK_DOWNLOAD_PATH=${2:-/google-cloud-sdk}

echo "GOOGLE_CLOUD_SDK_VERSION = ${GOOGLE_CLOUD_SDK_VERSION}"
echo "GOOGLE_CLOUD_SDK_DOWNLOAD_PATH = ${GOOGLE_CLOUD_SDK_DOWNLOAD_PATH}"

# Some CI service support cache feature.
# If there are cache, just use it and skip downloading.
if [ ! -e ${GOOGLE_CLOUD_SDK_DOWNLOAD_PATH} ]; then
    curl -L -o /tmp/google-cloud-sdk-${GOOGLE_CLOUD_SDK_VERSION}-linux-x86_64.tar.gz \
         https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-$GOOGLE_CLOUD_SDK_VERSION-linux-x86_64.tar.gz
    tar -xz -C /tmp -f /tmp/google-cloud-sdk-${GOOGLE_CLOUD_SDK_VERSION}-linux-x86_64.tar.gz
    mv /tmp/google-cloud-sdk ${GOOGLE_CLOUD_SDK_DOWNLOAD_PATH}

    ${GOOGLE_CLOUD_SDK_DOWNLOAD_PATH}/install.sh --quiet
    ${GOOGLE_CLOUD_SDK_DOWNLOAD_PATH}/bin/gcloud components install gsutil --quiet
fi

# Create symlink to default PATH
ln -s ${GOOGLE_CLOUD_SDK_DOWNLOAD_PATH}/bin/gcloud /usr/local/bin/gcloud
ln -s ${GOOGLE_CLOUD_SDK_DOWNLOAD_PATH}/bin/gsutil /usr/local/bin/gsutil
