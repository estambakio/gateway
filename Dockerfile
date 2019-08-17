FROM golang as binary

WORKDIR /app
COPY / ./
ENV GOPROXY https://proxy.golang.org
RUN make build

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=binary /app/dist/proxy /app/
ENTRYPOINT ["/app/proxy"]
