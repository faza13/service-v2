docker install

* docker run -p 3000:3000 -p 4317:4317 -p 4318:4318 --rm -ti grafana/otel-lgtm
* docker run -d -p 9092:9092 --name broker apache/kafka\:latest 

checklist

* ~~gin router~~
* ~~mariadb~~
* ~~kafka (watermill)~~
* sqs
* config(via aws param)
* mongodb
* cache (redis)
* elastic