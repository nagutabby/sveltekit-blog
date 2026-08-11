FROM golang:1.26.5-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates && adduser -D -H app
WORKDIR /app

COPY --from=build /out/server ./server
COPY --from=build /src/content ./content

ENV CONTENT_DIR=/app/content
USER app

EXPOSE 8080
ENTRYPOINT ["./server"]
