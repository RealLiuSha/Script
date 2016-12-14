#!/bin/sh
#
# Usage: ./redis-movekeys.sh [-h host] [-p port] [-n src] [-m dest] pattern
#
# Move keys matching pattern from the src Redis database to the
# dest Redis database.

set -e

HOST="localhost"
PORT="6379"
SRCDB="0"
DESTDB="0"
while getopts "h:p:n:m:" opt; do
    case $opt in
        h)  HOST=$OPTARG;;
        p)  PORT=$OPTARG;;
        n)  SRCDB=$OPTARG;;
        m)  DESTDB=$OPTARG;;
        \?) echo "invalid option: -$OPTARG" >&2; exit 1;;
    esac
done
shift $(( $OPTIND -1 )) 
PATTERN="$@"
if [ -z "$PATTERN" ]; then
    echo "pattern required" >&2
    exit 2
fi

redis-cli -h "$HOST" -p "$PORT" -n "$SRCDB" --raw keys "$PATTERN" |
    xargs -I{} redis-cli -h "$HOST" -p "$PORT" -n "$SRCDB" move {} "$DESTDB"
