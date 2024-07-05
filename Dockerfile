FROM golang:latest as build

WORKDIR /go/src/app
COPY ./main.go .
RUN go mod init docker-scraper; go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build

FROM chromedp/headless-shell:latest
RUN apt-get update; apt install dumb-init -y
ENTRYPOINT ["dumb-init", "--"]
COPY --from=build /go/src/app/docker-scraper /tmp
RUN mkdir -p /tmp/images
CMD ["/tmp/docker-scraper"]
