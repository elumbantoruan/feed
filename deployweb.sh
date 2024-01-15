echo docker build
echo
docker build -f DockerfileWeb . -t web-feed
sleep 1
echo docker tag 
echo 
docker tag web-feed edisonlt/web-feed
sleep 1
echo docker push
echo
docker push edisonlt/web-feed
sleep 1
echo delete newsfeed-web deployment
kubectl delete -f localdeploy/newsfeed-web.yaml
sleep 1
echo apply newsfeed-web deployment
kubectl apply -f localdeploy/newsfeed-web.yaml