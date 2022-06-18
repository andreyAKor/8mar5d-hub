# Environment
FROM golang:1.17 as build-env

RUN mkdir -p /opt/8mar5d-hub
WORKDIR /opt/8mar5d-hub
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /opt/bin/8mar5d-hub ./cmd/8mar5d-hub/main.go

# Release
FROM alpine:latest

COPY --from=build-env /opt/8mar5d-hub/configs/8mar5d-hub.yml /etc/8mar5d-hub.yml
COPY --from=build-env /opt/bin/8mar5d-hub /bin/8mar5d-hub

RUN mkdir -p /var/log/8mar5d-hub
ENV LOGGING_FILE=/var/log/8mar5d-hub/8mar5d-hub.log

ENTRYPOINT ["/bin/8mar5d-hub", "--config=/etc/8mar5d-hub.yml"]
