FROM golang:1.21 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /server

FROM alpine:latest
COPY ./certificates/localhost.crt /opt/ssl/localhost.crt
COPY ./certificates/localhost.unencrypted.key /opt/ssl/localhost.key

COPY ./certificates/ca.cert.pem /usr/local/share/ca-certificates/ca.cert.crt

RUN cat /usr/local/share/ca-certificates/ca.cert.crt >> /etc/ssl/certs/ca-certificates.crt
RUN apk --no-cache add ca-certificates && rm -rf /var/cache/apk/*
RUN update-ca-certificates

COPY --from=builder /server /server
EXPOSE 8443

CMD ["/server"]