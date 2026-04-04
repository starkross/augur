# golang:1.26-alpine
FROM golang@sha256:2389ebfa5b7f43eeafbd6be0c3700cc46690ef842ad962f6c5bd6be49ed82039 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN rm -rf internal/rules/policy && mkdir -p internal/rules/policy && cp -R policy/main policy/lib internal/rules/policy/
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /augur ./cmd/augur

# alpine:3.20
FROM alpine@sha256:a4f4213abb84c497377b8544c81b3564f313746700372ec4fe84653e4fb03805
RUN addgroup -S augur && adduser -S augur -G augur
COPY --from=builder /augur /usr/local/bin/augur
USER augur
WORKDIR /work
ENTRYPOINT ["augur"]
