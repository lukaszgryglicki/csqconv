#!/bin/bash
for f in `ls ${1}_*.jpegls`
do
  f2="${f%.*}.png"
  ffmpeg -f image2 -vcodec jpegls -i "$f" -y -f image2 -vcodec png "$f2"
  # convert -auto-level -auto-gamma "$f2" "$f2"
  convert -auto-level "$f2" "$f2"
done
# ffmpeg -f image2 -vcodec png -r 30 -i "${1}%08d.png" -pix_fmt gray16be -vcodec png "$1.mp4"
ffmpeg -f image2 -vcodec png -r 30 -i "${1}_%08d.png" -vcodec png "$1.mp4"
