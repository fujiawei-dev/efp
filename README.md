# efp

> An embedded forward proxy.

## Features

- SOCKS5 proxy

## Basic Usage

### Remote Proxy Server

```shell
./efp-windows-amd64.exe -s :8488 -cipher dummy -v
```

```shell
./efp-windows-amd64.exe -s :8488 -cipher aes-128-gcm -key 1234567890abcdef1234567890abcdef -v
```

### Local Proxy Client

```shell
./efp-windows-amd64.exe -c localhost:8488 -cipher dummy -socks :1080 -v
```

```shell
./efp-windows-amd64.exe -c localhost:8488 -cipher aes-128-gcm -key 1234567890abcdef1234567890abcdef -socks :1080 -tcptunnel :80=180.101.49.11:80 -v
```
