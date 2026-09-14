ARG GO_VERSION=1.27.0

FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /app

COPY lib/ ./lib/
COPY services/cell_a/ ./services/cell_a/

WORKDIR /app/services/cell_a

RUN go mod edit -replace zygarde/lib=../../lib -require zygarde/lib@v0.0.0 \
    && go mod tidy && go mod verify

RUN CGO_ENABLED=0 go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /app/bin \
    ./cmd/api


FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/bin /bin

ENTRYPOINT ["/bin"]
