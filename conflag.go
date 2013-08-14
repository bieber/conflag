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
 program's configuration.
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
