echo docker build
echo
docker build -f DockerfileGRPC . -t grpc-feed
sleep 1
echo docker tag
echo
docker tag grpc-feed edisonlt/grpc-feed
sleep 1
echo docker push
echo
docker push edisonlt/grpc-feed
sleep 1
echo delete newsfeed-grpc deployment
kubectl delete -f localdeploy/newsfeed-grpc.yaml
sleep 1
echo apply newsfeed-grpc deployment
kubectl apply -f localdeploy/newsfeed-grpc.yaml