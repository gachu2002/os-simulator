#!/usr/bin/env bash
set -euo pipefail

SIM_COVER=$(go test ./internal/sim -cover | awk '/coverage:/ {print $5}' | tr -d '%')
LESSON_COVER=$(go test ./internal/lessons -cover | awk '/coverage:/ {print $5}' | tr -d '%')

echo "internal/sim coverage: ${SIM_COVER}%"
echo "internal/lessons coverage: ${LESSON_COVER}%"

if awk "BEGIN {exit !(${SIM_COVER} >= 85)}"; then
  :
else
  echo "coverage check failed: internal/sim must be >= 85%"
  exit 1
fi

if awk "BEGIN {exit !(${LESSON_COVER} >= 85)}"; then
  :
else
  echo "coverage check failed: internal/lessons must be >= 85%"
  exit 1
fi
