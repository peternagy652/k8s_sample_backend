FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

COPY ./certificates/localhost.crt /opt/ssl/localhost.crt
COPY ./certificates/localhost.unencrypted.key /opt/ssl/localhost.key

RUN CGO_ENABLED=0 GOOS=linux go build -o /server

EXPOSE 8443

CMD ["/server"]