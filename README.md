# efp

> An embedded forward proxy.

## Features

- SOCKS5 proxy
- Only supports TCP

## Basic Usage

### Remote Proxy Server

```shell
efp -s :8488 -cipher dummy -v
```

```shell
efp -s :8488 -cipher aes-128-gcm -key 1234567890abcdef1234567890abcdef -v
```

### Local Proxy Client

```shell
efp -c localhost:8488 -cipher dummy -socks :1080 -v
```

```shell
efp -c localhost:8488 -cipher aes-128-gcm -key 1234567890abcdef1234567890abcdef -socks :1080 -tcptunnel :80=180.101.49.11:80 -v
```
