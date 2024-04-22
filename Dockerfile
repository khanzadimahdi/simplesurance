FROM golang:1.22 as build
WORKDIR /dist
COPY . .
ENV GOARCH=amd64 CGO_ENABLED=0
RUN go mod download
RUN go build -v -o app ./main.go && chmod +x app

FROM alpine:latest as deploy
WORKDIR /opt/server
COPY --from=build /dist /opt/server
CMD ["/opt/server/app"]