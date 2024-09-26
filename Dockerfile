# Step 1: Modules caching
FROM golang:1.22-alpine AS modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.22-alpine AS builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags ip2country -o /bin/app ./cmd/app

# Step 3: Final
FROM scratch
COPY --from=builder /app/config /config
COPY --from=builder /app/internal/repositories/disk/data.json /data.json
COPY --from=builder /bin/app /app
ENV DISK_REPOSITORY_RELATIVE_PATH=/data.json

CMD ["/app"]