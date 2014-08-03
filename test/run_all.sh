#!/bin/bash
tests=`ls -1d *_test`
for test in $tests; do
	./run_test.sh "$test"
	echo "============================="
done