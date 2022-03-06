# efp

> An embedded forward proxy.

## Features

- SOCKS5 proxy
- Only supports TCP

## Basic Usage

### Key generation

A random key is almost always better than a password. You can generate a base64url-encoded 16-byte random key with

```shell
efp -keygen 16
```

```shell
qARC56jnsb3bjOLXfmItKQ==
```

### Remote Proxy Server

#### No Encryption

```shell
efp -s :8488 -cipher dummy -v
```

#### Hexadecimal Key

```shell
efp -s :8488 -cipher aes-128-gcm -key 1234567890abcdef1234567890abcdef -v
```

#### Base64url-Encoded Key

```shell
efp -s :8488 -cipher aes-128-gcm -key64 qARC56jnsb3bjOLXfmItKQ== -v
```

### Local Proxy Client

#### No Encryption

```shell
efp -c localhost:8488 -cipher dummy -socks :1080 -v
```

#### Hexadecimal Key

```shell
efp -c localhost:8488 -cipher aes-128-gcm -key 1234567890abcdef1234567890abcdef -socks :1080 -tcptunnel :80=180.101.49.11:80 -v
```

#### Base64url-Encoded Key

```shell
efp -c localhost:8488 -cipher aes-128-gcm -key64 qARC56jnsb3bjOLXfmItKQ== -socks :1080 -tcptunnel :80=180.101.49.11:80 -v
```
