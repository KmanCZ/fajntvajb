FROM golang:alpine AS builder
WORKDIR /src

RUN apk add --no-cache nodejs npm

RUN npm install -g tailwindcss@3.4.16

COPY . /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

RUN tailwindcss -i ./internal/files/static/css/input.css -o ./internal/files/static/css/style.css

RUN go build -o /bin/server ./cmd/fajtvajb

RUN --mount=type=cache,target=/go/pkg/mod/ \
    CGO_ENABLED=0 go build -o ./bin/server ./cmd/fajtvajb

FROM alpine:latest
EXPOSE 8080/tcp

COPY --from=builder /bin/server /bin/

CMD ["/bin/server"]
