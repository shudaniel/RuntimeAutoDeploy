FROM alpine:3.3

RUN mkdir /gobin/

COPY /Users/aartij17/go/src/RuntimeAutoDeploy/buildRAD/sample_file /gobin

WORKDIR /gobin

EXPOSE 8080

ENTRYPOINT ["/gobin/sample_file"]