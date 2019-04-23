FROM golang:alpine AS build

ARG project="go-freeproxy"

WORKDIR /go/src/github.com/NoUseFreak/${project}/
RUN apk add --no-cache dep git

COPY ./Gopkg.* /go/src/github.com/NoUseFreak/${project}/
COPY ./cmd /go/src/github.com/NoUseFreak/${project}/cmd
COPY ./internal /go/src/github.com/NoUseFreak/${project}/internal

RUN cd /go/src/github.com/NoUseFreak/${project}/ \
    && dep ensure \
    && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${project} ./cmd/${project}

FROM alpine
ARG project="go-freeproxy"

RUN apk add --no-cache ca-certificates

WORKDIR /root/
COPY --from=build /go/src/github.com/NoUseFreak/${project}/${project} .

ENTRYPOINT ["./go-freeproxy"]
