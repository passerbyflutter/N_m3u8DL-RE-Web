# N_m3u8DL-RE-Web

[![GitHub Release](https://img.shields.io/github/release/passerbyflutter/N_m3u8DL-RE-Web.svg?style=flat)](https://github.com/passerbyflutter/N_m3u8DL-RE-Web/releases/latest)  

Web Server base on [N_m3u8DL-RE](https://github.com/nilaoda/N_m3u8DL-RE) with Golang backend and Vue frontend.

## Environment Variable

- DOWNLOAD_POOL_SIZE: Default size of download pool, default value: `3`.
- SAVE_PATH: Path to save downloaded videos, default value: `./download`.

## Run Backend

```sh
go run cmd/main.go
```

## Run Frontend

In `./web` directory

```sh
npm run dev
```
