# Initial stage: download modules
FROM golang:1.15 as modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

# Intermediate stage: Build the binary
FROM golang:1.15 as builder

COPY --from=modules /go/pkg /go/pkg

# add a non-privileged user
RUN useradd -u 10001 wifi_auth


RUN mkdir -p /wifi_auth
ADD . /wifi_auth



WORKDIR /wifi_auth



# Build the binary with go build
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o ./bin/wifi_auth ./cmd/wifi_auth


# Final stage: Run the binary
FROM scratch

# don't forget /etc/passwd from previous stage

COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Europe/Moscow

USER wifi_auth

# and finally the binary
COPY --from=builder /wifi_auth/bin/wifi_auth /wifi_auth

CMD ["/wifi_auth"]