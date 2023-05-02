FROM golang:1.19

#RUN apt-get update && apt-get install -y git-buildpackage debhelper zlib1g-dev libssl-dev libsasl2-dev liblz4-dev

#RUN cd /tmp && git clone https://github.com/edenhill/librdkafka && cd librdkafka && ./configure && make && make install

RUN go install github.com/makiuchi-d/arelo@latest
RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

WORKDIR /app

ENV LD_LIBRARY_PATH=/usr/local/lib

ENTRYPOINT [ "/app/scripts/docker-entrypoint.sh"]
