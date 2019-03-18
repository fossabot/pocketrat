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

### Directory indices
PocketRat will generate directory indices automatically, but you can also provide
them: PocketRat supports its own `gopher.index` format as well as the more standard
`Gophermap` format.  If both files exist in a directory, PocketRat will prefer `gopher.index`, though this format is simpler and less capable than `Gophermap`.

#### `gopher.index`
A `gopher.index` file consists of multiple lines with two tab-separated fields: the
file name and the text to associate with that entry.  An entry with no file name
will be sent as an informational message.  PocketRat will attempt to infer the type
of the entry will be inferred from the file name; unrecognized file types will be
assumed to be binary.

##### Example `gopher.index`
```
	Welcome to the PocketRat demonstration server.
README.md	PocketRat README
CHANGELOG.md	PocketRat CHANGELOG
```

#### `Gophermap`
The `Gophermap` format is as described at [Wikipedia]: a selector line beginning
with a type sigil, followed by a display string, a selector, hostname, and port,
or a commnt line containing no tab characters that will be sent to the client as
an informational message.

##### Example `Gophermap`
```
Welcome to the PocketRat demonstration server.
This line and the one before it are comment lines.
0README.md	PocketRat README	/README.md	localhost	70
0CHANGELOG.md	PocketRat CHANGELOG	/CHANGELOG.md	localhost	70
```

[Gopher]: https://en.wikipedia.org/wiki/Gopher_(protocol)
[Wikipedia]: https://en.wikipedia.org/wiki/Gopher_(protocol)#Source_code_of_a_menu
[Go]: http://www.golang.org/
[the RFC]: https://tools.ietf.org/html/rfc1436
