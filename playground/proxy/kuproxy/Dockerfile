FROM golang:1.14

EXPOSE 8080 8081 

WORKDIR /app
COPY kuproxy .

ENTRYPOINT [ "/app/kuproxy" ]
# sudo docker build -t localhost:32000/kuproxy:v2 .
# sudo docker tag xxx localhost:32000/kuproxy:v2
# sudo docker push localhost:32000/kuproxy
# sudo docker run -P localhost:32000/kuproxy:v0
# microk8s ctr images ls | grep kuproxy