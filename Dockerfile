FROM golang:1.21.4 AS file-api-builder

RUN go version

COPY . /github.com/Woodfyn/file-api/
WORKDIR /github.com/Woodfyn/file-api/

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=file-api-builder /github.com/Woodfyn/file-api/.bin/app .

COPY --from=file-api-builder /github.com/Woodfyn/file-api/configs/prod.yml ./configs/prod.yml

CMD ["./app"]