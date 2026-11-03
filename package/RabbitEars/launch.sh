#!/bin/sh
# Rabbit Ears — a 90s TV set for OnionOS
appdir=/mnt/SDCARD/App/RabbitEars
sysdir=/mnt/SDCARD/.tmp_update
cd "$appdir" || exit 1
mkdir -p data
export LD_LIBRARY_PATH="/lib:/config/lib:/mnt/SDCARD/miyoo/lib:$sysdir/lib:$sysdir/lib/parasyte:$LD_LIBRARY_PATH"
"$appdir/rabbitears" >> "$appdir/data/stdout.log" 2>&1
rm -f /tmp/stay_awake
sync
