FROM alpine:3.3

RUN mkdir /gobin/

COPY ../sample_file /gobin

WORKDIR /gobin

EXPOSE 8080

ENTRYPOINT ["/gobin/sample_file"]