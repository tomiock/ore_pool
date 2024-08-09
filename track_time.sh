##!/usr/bin/env bash

# Usage: ./track_time.sh <name> <command>
# Example: ./track_time.sh Alice "sleep 5"

NAME=$1
shift
COMMAND=$@

# Record the start time
START_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# POST start time
curl -X POST -d "${NAME},${START_TIME}" localhost:8080/track_start

# Execute the command and capture its exit status
eval "$COMMAND"
STATUS=$?

# Record the end time
END_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# POST end time
curl -X POST -d "${NAME},${END_TIME}" localhost:8080/track_end

# Exit with the status of the command
exit $STATUS
