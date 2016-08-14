# PocketRat
PocketRat is a [Gopher] server written in [Go].

It was written as a learning exercise, and with only infrequent references to
[the RFC] and none to any existing server implementations.  If it worked with
`lynx` I called it good.

As this was written as a learning exercise, and because I am a .NET developer by
day, I am sure the code is horrible, unidiomatic Go.  Please forgive me.

The name "PocketRat" is a calque from the German word for "gopher", *Taschenratte*.

## Usage
Build the server in the usual way (`go build`).

Set `ListenAddr`, `ListenPort`, `ServerName`, and `GopherRoot` in `config.json`.
If `ServerName` is unspecified it defaults to the machine's hostname as returned by
`os.Hostname()`, and there is of course a very good chance that that's not what you
want.

Run `pocketrat`.  It does not daemonize itself because some googling suggested
that there are issues around that when goroutines are involved.  Use whatever
tools your operating system provides to daemonize it instead.

[Gopher]: https://en.wikipedia.org/wiki/Gopher_(protocol)
[Go]: http://www.golang.org/
[the RFC]: https://tools.ietf.org/html/rfc1436
