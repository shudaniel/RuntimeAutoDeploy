#!/bin/sh

start=$SECONDS
sh ../setup_rad_wordload_cluster.sh
duration=$(( SECONDS - start ))
echo duration
echo "^seconds"


