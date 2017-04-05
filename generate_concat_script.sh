#!/bin/bash

# Example
#   ./generate_concat_script.sh > temp-concat.sh
#   bash temp-concat.sh
#
# Requirement
#   ffmpeg <http://ffmpeg.org/>

echo -e "ffmpeg -i concat:'\c"

isFirst=1
for ts in $(ls *.ts | sort); do
  if [ $isFirst -eq 1 ]; then
    isFirst=0
  else
    echo -e '|\c'
  fi
  echo -e "${ts}\c"
done

echo -e "' -c copy output.m4a\c"
echo
