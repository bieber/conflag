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
 options from both a config file and command-line flags.  You specify
 the options you require by defining a struct type to express your
 program's configuration.  For example, consider this struct:

	type ExampleConfig struct {
		 Name     string
		 Port     int     `conflag:"dest_port"`
		 UseDB    bool    `conflag:"db.use,true"`
		 FloatArg float64 `conflag:"example.float,0.5"`
	}

 conflag accepts string, int, bool, and float64 as types.  By default,
 defining a field with one of these types will result in conflag
 attempting to find an option in the config file's default section
 named with a lower-cased version of the field name (for instance,
 "name" instead of "Name").

 If you specify a field tag with a "conflag" section of a single
 string, that string will be used instead of the default name for the
 option.  For instance, in the example configuration the Port field
 would be read from an option or command-line flag named "dest_port".

 If your special field name includes a dot, then the portion of the
 name before the dot will be used as the section in the config file,
 and the portion after used as the name of the option.  On the
 command-line, users should use the option's full section.option name.
 In the example configuration, users would specify the UseDB option by
 either setting the "use" option in the "db" section of the config
 file, or providing a -db.use flag on the command-line.

 Your conflag field tag may also include a default value, separated
 from the option name by a comma.  Default values are used if the user
 doesn't provide one, and setting one renders the field optional.  Any
 field without a default is considered required, and must be provided
 either in the config file or at the command-line.

 conflag reserves the command-line option "config-file" for its own
 use.  Providing a --config-file option on the command-line will
 override the default config file location, and prompt conflag to
 attempt to load a config file at the provided path.
*/
package conflag

import (
	"fmt"
)

type option struct {
	// Identification
	name    string
	section string

	// Type of the field
	typeOf string

	// Has it been read yet?
	read bool

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
func ReadConfig(dest interface{}) error {
	options := getOptions(dest)
	fmt.Println(options)
	return nil
}
