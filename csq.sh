#!/bin/bash
# RF=... - specify red func to apply
# GF=... - specify green func to apply
# BF=... - specify blue func to apply
if [ -z "$RF" ]
then
  RF='if(x1<.95,if(x1<.33333,x1*3,1),1-2*(x1-.95))'
fi
if [ -z "$GF" ]
then
  GF='if(x1<.95,if(x1<.33333,0,if(x1>.66667,1,(x1-.33333)*3)),1-2*(x1-.95))'
fi
if [ -z "$BF" ]
then
  BF='if(x1<.66667,0,(x1-.66667)*3)'
fi
RC=1 GC=1 BC=1 RLO=1 RHI=1 GLO=1 GHI=1 BLO=1 BHI=1 NA=1 INF=20 PQ=1 HINT=1 MF=30 MODE=veryslow CRF=14 RF="$RF" GF="$GF" BF="$BF" csqconv "$@"
