docker build -t "$1" $1
docker run -i --rm $1
docker rmi -f $1
