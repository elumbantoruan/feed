echo docker build
echo
docker build -f DockerfileCronJob . -t cron-feed
sleep 1
echo docker tag
echo
docker tag cron-feed edisonlt/cron-feed
sleep 1
echo docker push
echo
docker push edisonlt/cron-feed
sleep 1
echo delete newsfeed-cronjob deployment
kubectl delete -f localdeploy/newsfeed-cronjob.yaml
sleep 1
echo apply newsfeed-cronjob deployment
kubectl apply -f localdeploy/newsfeed-cronjob.yaml