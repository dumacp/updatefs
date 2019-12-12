FROM golang:1.13

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

VOLUME /data/all

EXPOSE 8000/tcp

CMD ["app"]
