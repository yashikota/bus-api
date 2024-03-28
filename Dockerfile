# ======
# Dev
# ======
FROM golang:1.22-bookworm AS dev

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum .air.toml ./
RUN go mod download

COPY src/ ./src/

CMD ["air", "-c", ".air.toml"]

# ======
# Build
# ======
FROM golang:1.22-bookworm AS build

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY src/ ./src/

RUN go build -o /bin/main -ldflags="-s -w" ./src

# ======
# Deploy
# ======
FROM gcr.io/distroless/static-debian12 AS deploy

COPY --from=build /bin/main /main

EXPOSE 8080
USER nonroot:nonroot

CMD ["/main"]
