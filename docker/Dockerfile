FROM golang:1.19-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh
WORKDIR /app
COPY ../go.mod go.sum ./
RUN go mod download
COPY .. .
RUN go build -ldflags '-w -s' -o chat cmd/chat/main.go
RUN cp cmd/chat/.env .env
CMD ["./chat"]