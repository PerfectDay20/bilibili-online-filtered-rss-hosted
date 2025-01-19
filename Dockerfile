FROM golang:1.23 AS build-stage

WORKDIR /app
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /bilibili-rss

FROM gcr.io/distroless/static-debian12 AS build-release-stage

WORKDIR /

COPY --from=build-stage /bilibili-rss /bilibili-rss

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/bilibili-rss"]
