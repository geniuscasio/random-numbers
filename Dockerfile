FROM golang:1.16 AS builder

WORKDIR /workspace

COPY ./ ./

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN CGO_ENABLED=0 go build -o ./bin/server ./cmd/main.go

FROM alpine:latest
COPY --from=builder /workspace/bin/server /bin/server

EXPOSE 8080

ENTRYPOINT ["/bin/server"]