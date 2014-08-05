#!/bin/bash
testDir=`dirname $0`
tests=`ls -1d $testDir/*_test`
for test in $tests; do
	$testDir/run_test.sh "$test"
	echo "==========================================================="
done