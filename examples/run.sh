#!/usr/bin/env sh


incoming=/data/Pictures-Android/AndroidDCIM/Camera
incoming=/data/Pictures-tmp

set -x

picman autosort \
    --source Phone10 \
    --source-dir "$incoming" \
    --destination-dir /data/Pictures
