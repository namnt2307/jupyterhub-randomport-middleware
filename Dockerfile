FROM repo.bigdata.local/alpine:3.14
LABEL maintainer="namnt96@fpt.com.vn"

RUN apk update && apk add --no-cache tzdata && apk add curl jq gettext && \
    cp /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime && \
    rm -rf /var/lib/{apt,dpkg,cache,log}/

WORKDIR /app 

COPY bin/main /app/main

EXPOSE 2307

CMD ["/app/main"]