#!/bin/bash

COUNTER=1
while [  $COUNTER -lt 4 ]; do
 sleep 1
 #echo COUNTER $COUNTER
 let COUNTER+=1
done