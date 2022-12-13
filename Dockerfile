FROM golang:1.19-alpine

WORKDIR /go/src/app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go .
COPY internal ./internal

RUN go build -o /registry

CMD [ "/registry" ]