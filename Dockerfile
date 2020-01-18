FROM golang:alpine AS builder
ENV GO111MODULE=on
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w' -o home-server ./cmd/home-server/main.go
#RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/home-server ./cmd/home-server

FROM scratch
COPY --from=builder /app/home-server /app/
ENTRYPOINT ["/app/home-server"]
