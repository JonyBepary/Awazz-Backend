FROM golang:1.21 AS build

WORKDIR /app
COPY --chown=nonroot:nonroot  .. .
RUN go build -o app

FROM alpine:3.14
WORKDIR /app
COPY --from=build /app/app .
CMD ["./app"]
