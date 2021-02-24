FROM alpine:latest

RUN apk --no-cache add git curl jq \
    && rm -rf /var/lib/apt/lists/*

RUN curl -LO https://github.com/rhysd/changelog-from-release/releases/download/v2.2.0/changelog-from-release_linux_amd64.zip \
    && unzip changelog-from-release_linux_amd64.zip -d /

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
