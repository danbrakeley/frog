#!/bin/bash -e

if [ $# -eq 0 ]; then
  echo "No arguments supplied: please include a unique name for this benchmark."
  exit 1
fi

DATE=$(date +'%Y-%m-%d')
go test -bench=. -benchmem -cpuprofile "testdata/benchmarks/$DATE-$1.cpu" -cpu 8 -count 10 -timeout 2h > "testdata/benchmarks/$DATE-$1.txt"
