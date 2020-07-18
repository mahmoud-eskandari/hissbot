# hissbot (ربات ناشناس تلگرام)
A simple golang project, just for fun

## 1.Install Dependencies

`go mod tidy`

## 2.Build

Mac:
`GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./docker/bot/hissbot`

Linux:
`go build -o ./docker/bot/hissbot`

Windows:
`SET GOOS=linux && SET GOARCH=amd64 && SET CGO_ENABLED=0 && go build -o ./docker/bot/hissbot`


## 3.Configuration
open `./docker/.env` file and set your config on it

## 4.Run
```
cd docker
docker-compose build bot
docker-compose up -d
```


# Run With docker image
### From docker hub (no need to build or etc...)  [Docker Image](https://hub.docker.com/repository/docker/mahmoudetc/hissbot)

```
version: '3'
services:
  bot:
    image: mahmoudetc/hissbot:0.1.0
    container_name: 'bot'
    links:
      - mysql:mysql
      - redis:redis
    restart: always
    env_file:
      - ".env"
    environment:
      - TELEGRAM_API=${TELEGRAM_API}
      - DEBUG=${DEBUG}
      - DB=${MYSQL_DATABASE}
      - DB_HOST=tcp(${MYSQL_HOST}:3306)
      - DB_USER=${MYSQL_USER}
      - DB_PASSWORD=${MYSQL_PASSWORD}
      - DB_LOCATION=Local
      - REDIS_HOST=${REDIS_HOST}
      - REDIS_PORT="6379"
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - RANDOM_INT=${RANDOM_HASH_INT}
      - HASH_TABLE=${RANDOM_HASH_SRT}
    depends_on:
      - mysql
      - redis

  redis:
    image: bitnami/redis:${REDIS_VERSION}
    container_name: ${REDIS_HOST}
    env_file:
      - ".env"
    environment:
      - REDIS_PASSWORD=${REDIS_PASSWORD}

  myadmin:
    image: phpmyadmin/phpmyadmin:4.7.6-1
    container_name: mybotadmin
    links:
      - mysql:mysql
    ports:
      - "8099:80"
    env_file:
      - ".env"
    environment:
      - PMA_ARBITRARY=1
      - PMA_HOST=${MYSQL_HOST}
    restart: always
    depends_on:
      - mysql

  mysql:
    image: mysql:${MYSQL_VERSION}
    container_name: ${MYSQL_HOST}
    restart: always
    command: mysqld --sql_mode="STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION" --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
    env_file:
      - ".env"
    environment:
      - MYSQL_DATABASE=${MYSQL_DATABASE}
      - MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
      - MYSQL_USER=${MYSQL_USER}
      - MYSQL_PASSWORD=${MYSQL_PASSWORD}
      - TZ=Asia/Tehran
    volumes:
      - "./mysql:/var/lib/mysql"

```

## .env

```

# Telegram bot api form bot father
TELEGRAM_API=

DEBUG=true
# MySQL
MYSQL_VERSION=5.7.22
MYSQL_HOST=mysqlbot
MYSQL_DATABASE=bot
MYSQL_ROOT_USER=root
MYSQL_ROOT_PASSWORD=AstrengthPasswordHere
MYSQL_USER=botttt
MYSQL_PASSWORD=AstrengthPasswordHereToo

# Redis config
REDIS_VERSION=5.0.8
REDIS_HOST=redisbot
REDIS_PASSWORD=AstrengthPasswordHere

# Hash salt
RANDOM_HASH_INT=1342
# Random salt 10 unique char
RANDOM_HASH_SRT=abcdefghij

```
