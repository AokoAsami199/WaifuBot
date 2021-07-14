FROM golang:1.17rc1
WORKDIR /work

ext:
    RUN go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
    RUN go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
    RUN go install golang.org/x/tools/cmd/goimports@latest
    RUN go install github.com/derision-test/go-mockgen/...@latest

deps:
    COPY go.mod go.sum .
    RUN go mod download

    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum


test:
    FROM +deps
    COPY . .
    RUN /bin/bash -c "set -a; source .env; go test ./... -v -cover"

lint:
    FROM golangci/golangci-lint:v1.41
    WORKDIR /work
    COPY . .
    RUN golangci-lint run -v

build:
    FROM +deps
    COPY . .
    RUN go build -ldflags="-w -s" -o ./bot
    SAVE ARTIFACT ./bot

docker:
    ARG BOT_IMAGE_TAG=latest
    FROM alpine:3.14

    WORKDIR /bin
    
    COPY +build/bot /bin/bot
    ENTRYPOINT ["/bin/bot"]

    SAVE IMAGE --push ghcr.io/karitham/waifubot:$BOT_IMAGE_TAG

docker-otc:
    ARG OTC_IMAGE_TAG=otc_latest
    FROM ghcr.io/karitham/otc
    RUN apk add --no-cache postgresql-client

    SAVE IMAGE --push ghcr.io/karitham/waifubot:$OTC_IMAGE_TAG

mock-search:
    FROM +ext
    COPY . .

    RUN go-mockgen -f github.com/Karitham/WaifuBot/discord -i SearchProvider -o discord/SearchMock_test.go
    SAVE ARTIFACT discord/SearchMock_test.go AS LOCAL discord/search_mock.go

mock-roll:
    FROM +ext
    COPY . .

    RUN go-mockgen -f github.com/Karitham/WaifuBot/discord -i Randomer -o discord/randomer_mock_test.go
    RUN go-mockgen -f github.com/Karitham/WaifuBot/discord -i Storager -o discord/storager_mock_test.go

    SAVE ARTIFACT discord/randomer_mock_test.go AS LOCAL discord/randomer_mock_test.go
    SAVE ARTIFACT discord/storager_mock_test.go AS LOCAL discord/storager_mock_test.go

run:
    FROM +deps
    COPY . .
    RUN go run . || :

sqlc:
    FROM +ext
    COPY . .
    RUN sqlc generate

    RUN fieldalignment -fix ./service/store/ || :
    RUN fieldalignment -fix ./service/store/ || :
    RUN fieldalignment -fix ./service/store/ || :

    SAVE ARTIFACT ./service/store/* AS LOCAL ./service/store/

docker-all:
    BUILD +docker
    BUILD +docker-otc