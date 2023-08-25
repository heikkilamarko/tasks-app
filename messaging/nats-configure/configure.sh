#!/bin/sh
set -e

nats stream add --config "/streams/tasks.json"
nats consumer add tasks --config "/consumers/tasks.json"
