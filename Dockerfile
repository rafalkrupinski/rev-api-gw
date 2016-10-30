FROM alpine:3.4
MAINTAINER Matte Silver
ADD . /go/src/bitbucket.org/mattesilver/etsygw
ENV GOPATH=/go
WORKDIR /go/src/bitbucket.org/mattesilver/etsygw
RUN apk add --no-cache --update go git &&\
    go get &&\
    go build bitbucket.org/mattesilver/etsygw &&\
    apk del --no-cache git go
ENTRYPOINT ["/go/bin/etsygw"]
CMD ["-v"]
EXPOSE 8080
