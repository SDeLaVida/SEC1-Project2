#!/bin/bash

for i in {0..3}; do
    x=$((i))
    bash "go run . $x | tee -a ./logs/text$x.txt"
done