# Change Log

## 0.1.1 - 2019-03-17
### Added
- `Gophermap` files are now supported in addition to `gopher.index` files.
However, if both files exist, `gopher.index` will be used.

### Fixed
- The "fix" in 0.1.0 that ceased terminating directory listings and
text files with `.` has been reverted and corrected.

## 0.1.0 - 2016-03-01
### Added
- Custom directory indices.  In any directory in your gopherspace
you may place a `gopher.index` file.  Each line consists of two
tab-separated fields: the file, and the description that should
appear in the Gopher index.  Blank file names will add entries
of type `i` (informational message).  Files not listed in `gopher.index`
will not appear in directory indices sent to clients.

### Fixed
- The [Gopher RFC](https://tools.ietf.org/html/rfc1436) specifies that
when sending menus or text files, the transmission should be terminated
by `.` on a single line.  However, at least one prominent Gopher client
does not respect this, instead displaying the terminating period as part
of the content.  Consequently, all transfers are now ended simply by
closing the connection after all data is sent.

## 0.0.1 - not released
### Fixed
- Fixed directory traversal vulnerability.

## 0.0.0 - 2016-02-28
### Initial release
### Known issues
- PocketRat v0.0.0 is vulnerable to directory traversal attacks.
Do not use PocketRat in a production environment.
