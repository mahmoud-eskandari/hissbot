# hissbot (ربات ناشناس تلگرام)
A simple golang project, just for fun

## Build

Mac:
`GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./docker/bot/hissbot`

Linux:
`go build -o ./docker/bot/hissbot`

Windows:
`SET GOOS=linux && SET GOARCH=amd64 && SET CGO_ENABLED=0 && go build -o ./docker/bot/hissbot`

## Run
```
cd docker
docker-compose build bot
docker-compose up -d
```

## Configuration
open `./docker/.env` file and set your config on it