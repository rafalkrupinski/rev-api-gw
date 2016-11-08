FROM alpine:3.4
MAINTAINER Rafał Krupiński
ADD . /go/src/github.com/rafalkrupinski/revapigw
ENV GOPATH=/go
WORKDIR /go/src/bitbucket.org/rafalkrupinski/revapigw
RUN apk add --no-cache --update go git &&\
    go get &&\
    go build github/rafalkrupinski/revapigw &&\
    apk del --no-cache git go
ENTRYPOINT ["/go/bin/etsygw"]
CMD ["-v"]
EXPOSE 8080
