#!/bin/bash
bold="\\033[1m"
green="\\033[1;32m"
red="\\033[1;31m"
normal="\\033[0m"

testDir=`dirname $0`
dockerImage=`basename $1`

echo -e "Executing test ${bold}${dockerImage}${normal}..."

#Copy bin sources into the docker build context.
#see https://docs.docker.com/reference/builder/#add
cp -Rf $testDir/../bin $1

#Build a docker image for the test
docker build -t $dockerImage $1 &> /dev/null
echo

#Run the test within a container
docker run --cidfile="cid" -i $dockerImage
cid=`cat cid`
#Exit code of the container represent the test result
exit=`docker wait "$cid"`

#Cleanup cid file, docker container and docker image
rm cid
docker rm "$cid" 1> /dev/null
docker rmi -f $dockerImage 1> /dev/null

echo

#Print test result
if [ $exit == 0 ]
then
	echo -e "${green}$dockerImage SUCCESS${normal}"
else
	echo -e "${red}$dockerImage FAILED${normal}"
	exit 1
fi

