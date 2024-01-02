#!/usr/bin/env bash

# build
go build -o . ./...

# 1: echo
maelstrom test -w echo --bin 01-echo --node-count 1 --time-limit 10 || exit 1

# 2: unique-ids
maelstrom test -w unique-ids --bin 02-unique-ids --time-limit 30 --rate 1000 --node-count 3 --availability total \
  --nemesis partition || exit 1

# 3a: single-node broadcast
maelstrom test -w broadcast --bin 3a-single-node-broadcast --node-count 1 --time-limit 20 --rate 10 || exit 1

# 3b: multi-node broadcast
maelstrom test -w broadcast --bin 3b-multi-node-broadcast --node-count 5 --time-limit 20 --rate 10 || exit 1

# 3c: fault tolerant broadcast
maelstrom test -w broadcast --bin 3c-fault-tolerant-broadcast --node-count 5 --time-limit 20 --rate 10 \
  --nemesis partition || exit 1
