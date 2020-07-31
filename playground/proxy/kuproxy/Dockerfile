FROM alpine

EXPOSE 8080 8081 

WORKDIR /app
# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
COPY kuproxy .

ENTRYPOINT [ "./kuproxy" ]
# sudo docker build -t localhost:32000/kuproxy:v0 .
# sudo docker tag xxx localhost:32000/kuproxy:v0
# sudo docker push localhost:32000/kuproxy
# sudo docker run -P localhost:32000/kuproxy:v0
# microk8s ctr images ls | grep kuproxy