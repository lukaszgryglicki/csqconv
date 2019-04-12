#!/bin/bash
for t in "$@"
do
  u161=`echo "((${t}+40)/160)*65535" | bc -l`
  u162=`echo "((${t}+20)/140)*65535" | bc -l`
  u163=`echo "(${t}/650)*65535" | bc -l`
  echo "${t} --> $u161 or $u162 or $u163 (*)"
done
