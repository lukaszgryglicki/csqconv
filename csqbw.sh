#!/bin/bash
# F=... - specify GS func to apply
if [ -z "$F" ]
then
  F="(1-x1)"
fi
MODE=veryslow HINT=1 INF=20 PQ=1 OGS=1 NA=1 GSR=1 GSG=0 GSB=0 RLO=1 RHI=1 RF="$F" RC=1 RGA=1.41 MF=30 CRF=14 csqconv "$@"
