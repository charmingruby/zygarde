ARG GO_VERSION=1.27.0

FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /app

COPY lib/ ./lib/
COPY services/cell_c/go.mod services/cell_c/go.sum ./services/cell_c/

WORKDIR /app/services/cell_c

RUN go mod edit -replace zygarde/lib=../lib \
    && go mod download && go mod verify

COPY services/cell_c/ ./

RUN CGO_ENABLED=0 go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /app/bin \
    ./cmd/api


FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/bin /bin

ENTRYPOINT ["/bin"]
