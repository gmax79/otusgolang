#!/bin/bash

proc=$(ps -a | grep $1)
if [[ ! -z $proc ]]; then
  pid=$(echo $proc | awk '{print $1}')
  echo $pid
  kill $pid
fi
