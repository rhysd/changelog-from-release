FROM alpine:latest

RUN apk --no-cache add git curl jq \
    && rm -rf /var/lib/apt/lists/*

RUN curl -LO https://github.com/rhysd/changelog-from-release/releases/download/v2.2.5/changelog-from-release_2.2.5_linux_amd64.tar.gz \
    && tar xf changelog-from-release_2.2.5_linux_amd64.tar.gz -C / \
    && rm changelog-from-release_2.2.5_linux_amd64.tar.gz

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
