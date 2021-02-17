FROM golang:1.15-alpine as build
COPY *.go /app/
COPY go.mod go.sum /app/
WORKDIR /app
RUN go mod download
RUN go build

FROM alpine:3.13 as run
COPY --from=build /app/doo /
ENTRYPOINT ["/doo"]
