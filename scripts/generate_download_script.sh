#!/bin/bash

# Example
#   ./generate_download_script.sh v-low-tokyo1 2865 11 20170402092959 360 > temp-download.sh
#   bash temp-download.sh
#
#   where v-low-tokyo1   : area
#         2865           : dir
#         11             : sequence (00-99)
#         20170402092959 : start
#         360            : max (e.g. 30min => 30*60/5 => 360)
#
# To list 100-latest ts files, curl either of:
#   https://smartcast.hs.llnwd.net/v-low-tokyo1/2865/2865.txt
#   https://smartcast.hs.llnwd.net/v-low-nagoya1/2C65/2C65.txt
#   https://smartcast.hs.llnwd.net/v-low-osaka1/3065/3065.txt
#   https://smartcast.hs.llnwd.net/v-low-fukuoka1/3865/3865.txt

area=$1
dir=$2
seq=$3
ts=$4
max=$5

echo 'set -o errexit'
for ((i=0; i < $max; i++)); do
  url=$(printf 'curl https://smartcast.hs.llnwd.net/%s/%d/%02s_%s.ts -o %s.ts' $area $dir $seq $ts $ts)
  echo $url

  seq=$(($seq + 1))
  if [ $seq -eq 100 ]; then
    seq="0"
  fi
  ts=$(($ts + 5))
  if [ $(($ts % 100)) -ge 60 ]; then
    ts=$(($ts - 60 + 100))
    if [ $(($ts % 10000)) -ge 6000 ]; then
      ts=$(($ts - 6000 + 10000))
      if [ $(($ts % 1000000)) -ge 240000 ]; then
        ts=$(($ts - 240000 + 1000000))
        # FIXME: Increment month, and year
      fi
    fi
  fi
done
