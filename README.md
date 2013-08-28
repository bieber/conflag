Conflag
=======

conflag aims to make program configuration as simple as possible.  It
allows you to specify your configuration schema as a struct type with
field tags for more fine-grained control.  Just pass a pointer to an
instance of your config struct to the ReadConfig function along with a
file path to a config file, and conflag will automatically read
options first from the command-line, then from the config file, and
finally from a default if necessary.

Installation
------------

`go get github.com/bieber/conflag`

Usage
-----

See the godoc documentation.