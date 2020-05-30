#!/bin/sh

start=$SECONDS
sh ../setup/setup_rad_wordload_cluster.sh
duration=$(( SECONDS - start ))
echo $(duration)
echo "^seconds"

kubectl delete -f ../../../../capi-quickstart.yaml 
