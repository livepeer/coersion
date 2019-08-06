FROM golang:1.12-alpine as builder

RUN apk add --no-cache make gcc musl-dev git musl
#ENV CGO_ENABLED=1

COPY . /coersion
RUN cd /coersion && go build . && cd / && \
    mv /coersion/coersion /usr/local/bin/ && \
    mv /coersion/swagger.json / && \
    rm -rf /coersion

FROM alpine:latest

RUN apk add --no-cache ca-certificates ffmpeg
COPY --from=builder /usr/local/bin/coersion /usr/local/bin/
COPY --from=builder /swagger.json /

EXPOSE 8088
CMD ["coersion"]
