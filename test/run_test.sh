echo "Executing test $1..."
cp -Rf ../bin $1
docker build -t "$1" $1 1> /dev/null
echo "============================="
docker run --cidfile="$1_cid" -i $1
echo "============================="
cid=`cat $1_cid`
exit=`docker wait "$cid"`
docker rm "$cid" 1> /dev/null
rm "$1_cid"
docker rmi -f $1 1> /dev/null

if [ $exit == 0 ]
then
	echo "SUCCESS"
else
	echo "FAILED"
fi

