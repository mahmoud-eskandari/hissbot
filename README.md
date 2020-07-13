# hissbot (ربات ناشناس تلگرام)
A simple golang project, just for fun

## Build
`go build -o ./docker/bot/hissbot`

## Run
```
cd docker
docker-compose build bot
docker-compose up -d
```

## Configuration
open `./docker/.env` file and set your config on it