FROM golang:1.21 as build

RUN apt-get update && apt-get install -y zsh
RUN wget https://github.com/casey/just/releases/download/1.23.0/just-1.23.0-x86_64-unknown-linux-musl.tar.gz && \
  tar xf just-1.23.0-x86_64-unknown-linux-musl.tar.gz just && \
  mv just /bin/

WORKDIR /go/src/app
COPY . .

RUN just build

# google only supports 'latest' tag for their distroless images
FROM gcr.io/distroless/static-debian11:latest
COPY --from=build /go/src/app/bin/ianpod /
CMD ["/ianpod"]
