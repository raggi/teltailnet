# teltailnet

Teltailnet is a Telnet proxy that may be useful if you want to expose a machine
that supports telnet to computers on your Tailscale network, but don't want to
expose it to the internet.

## Usage

First start the server:

```
go run ./cmd/teltailnet -connect bat.org:23
```

Authorize the server using the URL in the log, or set `$TS_AUTHKEY` beforehand.

The server will print out an IP and port that it is proxying on, for example:

```
teltailnet: listening on: 100.100.100.100:23
```

Then connect to the server using any telnet client. For convenience, an
extremely basic telnet client is included in this repository:

```
go run ./cmd/telnet -connect 100.100.100:23
```

## Building

To build binaries, use `go build ./cmd/teltailnet`.

## How it works

This proxy is implemented using [tsnet](https://tailscale.com/kb/1244/tsnet/)
the embedded Tailscale library. It starts a Tailscale node inside the teltailnet
process. The tsnet node listens for incoming connections on port 23, and on
receiving a connection it opens a new upstream connection to the IP:port
specified by the `-connect` flag.

The connections are established using the
[ziutek/telnet](https://github.com/ziutek/telnet) library.

## License

MIT
