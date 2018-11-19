FROM alpine

# for https requests
RUN apk update && apk add ca-certificates && update-ca-certificates && apk add openssl

COPY crawler.sh crawler.sh
RUN chmod +x crawler.sh

ENTRYPOINT ["./crawler.sh"]