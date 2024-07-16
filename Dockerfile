                                                                                                                                                                                 gchrome.dockerfile                                                                                                                                                                                             FROM debian:bookworm-slim as builder
RUN apt-get update && apt-get install -y \
        apt-transport-https \
        ca-certificates \
        golang

COPY ./main.go /project/main.go
WORKDIR /project
RUN go mod init docker-scraper; go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build

FROM debian:stable-slim
# FROM ubuntu:xenial
# FROM google/debian:jessie

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
        apt-transport-https \
        ca-certificates \
        curl \
        gnupg \
        --no-install-recommends \
        && curl -sSL https://dl.google.com/linux/linux_signing_key.pub | apt-key add - \
        && echo "deb https://dl.google.com/linux/chrome/deb/ stable main" > /etc/apt/sources.list.d/google-chrome.list \
        && apt-get update && apt-get install -y \
        google-chrome-stable \
        fontconfig \
        fonts-ipafont-gothic \
        fonts-wqy-zenhei \
        fonts-thai-tlwg \
        fonts-kacst \
        fonts-symbola \
        fonts-noto \
        fonts-freefont-ttf \
        --no-install-recommends \
        && apt-get purge --auto-remove -y curl gnupg \
        && rm -rf /var/lib/apt/lists/*

RUN useradd headless --shell /bin/bash --create-home \
  && usermod -a -G sudo headless \
  && echo 'ALL ALL = (ALL) NOPASSWD: ALL' >> /etc/sudoers \
  && echo 'headless:nopassword' | chpasswd

RUN mkdir /data && chown -R headless:headless /data

COPY --from=builder /project/docker-scraper  /project/docker-scraper
WORKDIR /
RUN mkdir -p /project/images
COPY run.sh /project/run.sh

ENTRYPOINT ["/project/run.sh"]

