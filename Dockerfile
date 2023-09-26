FROM node:18.18.0-bullseye AS build-frontend-stage

WORKDIR /app
COPY ./web/ ./

RUN npm ci
RUN npm run build

FROM golang:1.21.1-bullseye AS build-backend-stage

WORKDIR /app
COPY ./ ./
COPY --from=build-frontend-stage /app/dist ./web/dist

ARG RELEASE_VERSION
ARG NOW

RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$RELEASE_VERSION -X main.buildTime=$NOW -X main.isRelease=true" -o  N_m3u8DL-RE-Web ./cmd
RUN chmod a+x ./N_m3u8DL-RE-Web

FROM alpine AS download-stage

WORKDIR /

ADD https://github.com/nilaoda/N_m3u8DL-RE/releases/download/v0.2.0-beta/N_m3u8DL-RE_Beta_linux-x64_20230628.tar.gz /
RUN tar -xf /N_m3u8DL-RE_Beta_linux-x64_20230628.tar.gz
RUN chmod a+x /N_m3u8DL-RE_Beta_linux-x64/N_m3u8DL-RE

FROM jrottenberg/ffmpeg:4.1.10-ubuntu2004

WORKDIR /

RUN apt update && DEBIAN_FRONTEND=noninteractive TZ=Asia/Taipai apt install -y libicu-dev && rm -rf /var/lib/apt/lists/*
COPY --from=download-stage /N_m3u8DL-RE_Beta_linux-x64/N_m3u8DL-RE /usr/local/bin
COPY --from=build-backend-stage /app/N_m3u8DL-RE-Web /

ENTRYPOINT ["/N_m3u8DL-RE-Web"]