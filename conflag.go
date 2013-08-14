/*
 * Copyright (c) 2013, Robert Bieber
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 *
 * Redistributions in binary form must reproduce the above copyright
 * notice, this list of conditions and the following disclaimer in the
 * documentation and/or other materials provided with the distribution.
 *
 * Neither the name of the project's author nor the names of its
 * contributors may be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
 * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
 * TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
 * PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

/*
Package conflag simplifies program configuration by accumulating
options from both a config file and command line flags.  You specify
the options you require by defining a struct type to express your
program's configuration.

Your configuration struct can contain bool, int, float64, and string
values.  Values will be set first to their default, then to an option
from the config file if it exists, and then to an option from the
command line if it exists.  If you don't override the default
behavior, options will be read from config file values and command
line flags of the field name converted to all lower-case.

You can override this behavior by using field tags in your
configuration struct.  Consider the following example struct:

	type TestConfig struct {
		Name        string
		Port        int     `conflag:"net.port|Port to listen on."`
		UseDB       bool    `conflag:"db.use|Flag, set to true to use DB|true"`
		Probability float64 `conflag:"user.probability|Odds of success|0.5"`
	}

In this case, Name will be read from the field "name" in the config
file, or specified on the command line with --name Value or
--name=Value.  The other fields in the struct have custom
configuration, which comes in the form

	"section.name|Usage text|default_value"

Any of the three options may be left out.  The name, if specified,
comes in the form section.name, where the section is the section of
the config file the option should be located in, and name is the name
of the option.  If no section is included, then the option should be
found in the default section of the config file.  To set an option
with a section on the command line, use the full section.name
specifier as the flag, for instance --net.port=5005.

If usage text is specified for an option, it will appear in the usage
message that displays when a user enters --help on the command line,
or attempts to set a disallowed option.

If a default value is specified for an option, it will be used if the
option doesn't appear in the config file or the command line flags.
Setting a default value for the option makes it optional, and no
error will be thrown if the user doesn't provide a value for it.

conflag reserves --config-file as a special command line option.  If a
user specifies it, conflag will read from that config file.
Otherwise, it will attempt to read from the config file specified as
the default.  If no config file is found at the specified location,
conflag will use only defaults and command line flags.

Example of a simple program that reads configuration with conflag:

	package main

	import (
		"fmt"
		"github.com/bieber/conflag"
		"os/user"
		"path/filepath"
	)

	type TestConfig struct {
		Name        string
		Port        int     `conflag:"net.port|Port to listen on."`
		UseDB       bool    `conflag:"db.use|Flag, set to true to use DB|true"`
		Probability float64 `conflag:"user.probability|Odds of success|0.5"`
	}

	func main() {
		t := new(TestConfig)
		u, _ := user.Current()
		err := conflag.ReadConfig(
			t,
			filepath.Join(u.HomeDir, "test.conf"),
		)
		fmt.Println(t, err)
	}
*/
package conflag

import (
	"fmt"
	"strings"
)

type option struct {
	// Identification
	name        string
	section     string
	structField string
	usage       string

	// Type of the field
	typeOf string

	// Has it been read yet?
	read         bool
	readFromFlag bool

	// Possible value fields
	boolVal   bool
	floatVal  float64
	intVal    int
	stringVal string
}

// ReadConfig reads configuration from the configuration file and the
// command line, filling in options in the provided destination
// struct.  If your destination struct is improperly configured, it
// will panic.  If required options aren't provided by the user, it
// will return an error.
//
// dest should be a pointer to your result struct, and configFile the
// default location to load the configuration file from.
func ReadConfig(dest interface{}, configFile string) error {
	options := getOptions(dest)
	configFile = readFlags(options, configFile)
	readFile(options, configFile)
	setOptions(dest, options)

	missingOptions := make([]string, 0, 10)
	for _, option := range options {
		if !option.read {
			missingOptions = append(
				missingOptions,
				option.section+"."+option.name,
			)
		}
	}
	if len(missingOptions) > 0 {
		return fmt.Errorf(
			"Missing options: %s",
			strings.Join(missingOptions, ","),
		)
	}

	return nil
}
