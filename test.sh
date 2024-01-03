#!/usr/bin/env bash

set -euo pipefail

# build
go build -o . ./...

# 1: echo
maelstrom test -w echo --bin 01-echo --node-count 1 --time-limit 10

# 2: unique-ids
maelstrom test -w unique-ids --bin 02-unique-ids --time-limit 30 --rate 1000 --node-count 3 --availability total \
  --nemesis partition

# 3a: single-node broadcast
maelstrom test -w broadcast --bin 3a-single-node-broadcast --node-count 1 --time-limit 20 --rate 10

# 3b: multi-node broadcast
maelstrom test -w broadcast --bin 3b-multi-node-broadcast --node-count 5 --time-limit 20 --rate 10

# 3c: fault tolerant broadcast
maelstrom test -w broadcast --bin 3c-fault-tolerant-broadcast --node-count 5 --time-limit 20 --rate 10 \
  --nemesis partition

# 3d & 3e: efficient broadcast
# test if implementation is correct
maelstrom test -w broadcast --bin 3d-3e-efficient-broadcast --node-count 25 --time-limit 20 --rate 100 \
  --latency 100 --nemesis partition

# test performance
output="$(maelstrom test -w broadcast --bin 3d-3e-efficient-broadcast --node-count 25 --time-limit 20 \
  --rate 100 --latency 100 2> /dev/null)"

# check msgs-per-op
echo "$output" | grep -oP ':msgs-per-op\s*\K((\d+\.?\d*)|(\.\d+))' \
  | while read -r line ; do
      echo "msgs-per-op: $line"
      if [[ $(echo "if (${line} > 20) 1 else 0" | bc) -eq 1 ]]; then
        echo "msgs-per-op is greater than 20"
        exit 1
      fi
    done

# check latency
latency_median=$(echo "$output" | grep -A4 ':stable-latencies' | grep -oP -m 1 '(?<=0.5).*?(?=,)' | head -1)
latency_max=$(echo "$output" | grep -A4 ':stable-latencies' | grep -oP -m 1 '(?<=1).*?(?=\},)' | head -1)

echo "latency_median: $latency_median ms"
echo "latency_max: $latency_max ms"

if [[ $(echo "if (${latency_median} > 1000) 1 else 0" | bc) -eq 1 ]]; then
  echo "latency_median is greater than 1 second"
  exit 1
fi

if [[ $(echo "if (${latency_max} > 2000) 1 else 0" | bc) -eq 1 ]]; then
  echo "latency_max is greater than 2 seconds"
  exit 1
fi
