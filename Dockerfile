############################################
# STEP 1 build docker image for detour-proxy
############################################
FROM golang:1.19.1-bullseye AS builder
WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .
RUN go mod tidy
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/detour-proxy
############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /bin/detour-proxy /app/detour-proxy
COPY --from=builder /go/src/mypackage/myapp/config.yaml /app/config.yaml
WORKDIR /app
# Run the hello binary.
ENTRYPOINT ["/app/detour-proxy"]
