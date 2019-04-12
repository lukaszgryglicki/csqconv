#!/bin/bash
for u in "$@"
do
  t1=`echo "${u}/65535*160-40" | bc -l`
  t2=`echo "${u}/65535*120-20" | bc -l`
  t3=`echo "${u}/65535*650" | bc -l`
  echo "${u} --> $t1 or $t2 or $t3 (*)"
done
