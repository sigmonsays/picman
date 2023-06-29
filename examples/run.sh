#!/usr/bin/env sh


incoming=/data/Pictures-Android/AndroidDCIM/Camera

picman autosort \
    --source Phone10 \
    --source-dir "$incoming" \
    --destination-dir /data/Pictures
