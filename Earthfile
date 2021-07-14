FROM golang:1.17rc1-alpine
WORKDIR /work

deps:
    COPY go.mod go.sum .
    RUN CGO_ENABLED=0 go mod download
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

tidy:
    FROM +deps
    COPY . .

    RUN CGO_ENABLED=0 go mod tidy

    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

test:
    FROM +tidy
    RUN CGO_ENABLED=0 go test ./... -v -cover

lint:
    FROM golangci/golangci-lint:v1.41
    WORKDIR /work
    COPY . .
    RUN golangci-lint run -v

build:
    FROM +tidy

    RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o ./bot
    SAVE ARTIFACT ./bot

docker:
    ARG BOT_IMAGE_TAG=latest
    FROM alpine:3.14

    WORKDIR /bin
    
    COPY +build/bot /bin/bot
    ENTRYPOINT ["/bin/bot"]

    SAVE IMAGE --push ghcr.io/karitham/waifubot:$BOT_IMAGE_TAG

docker-otc:
    ARG OTC_IMAGE_TAG=latest

    FROM ghcr.io/karitham/otc
    RUN apk add --no-cache postgresql-client

    SAVE IMAGE --push ghcr.io/karitham/waifubot:$OTC_IMAGE_TAG

mock-search:
    FROM +deps
    RUN CGO_ENABLED=0 go install golang.org/x/tools/cmd/goimports@latest
    RUN CGO_ENABLED=0 go install github.com/derision-test/go-mockgen/...@latest

    COPY . .

    RUN go-mockgen -f github.com/Karitham/WaifuBot/discord -i SearchProvider -o discord/SearchMock_test.go
    SAVE ARTIFACT discord/SearchMock_test.go AS LOCAL discord/search_mock.go

run:
    FROM +tidy

    RUN CGO_ENABLED=0 go run . || :

docker-all:
    BUILD +docker
    BUILD +docker-otc