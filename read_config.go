/*
 * Copyright (c) 2014, Robert Bieber
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
 */

package conflag

import (
	"errors"
	"fmt"
	"strconv"
)

// Read reads configuration from the available sources into the
// specified fields of the config struct.  It returns a slice of
// strings with any extra arguments (which will trigger an error if
// not explicitly allowed via AllowExtraArgs) and an error which will
// be nil if the configuration was processed successfully.
func (c *Config) Read() ([]string, error) {
	fin, args, err := c.findConfigFile()
	if err != nil {
		return nil, err
	}
	if fin == nil && c.fileRequired {
		return nil, errors.New("conflag: Required config file not found.")
	}

	if fin != nil {
		err = readConfigFile(c.fields, fin)
		if err != nil {
			return nil, err
		}
	}

	extraArgs, err := readCommandLineFlags(c.fields, args, c.extraArgsAllowed)
	if err != nil {
		return nil, err
	}

	for _, field := range c.fields {
		err := field.readValue()
		if err != nil {
			return nil, err
		}
	}

	return extraArgs, nil
}

// Attempts to read the raw string value from the struct and fill in
// the corresponding field in the destination
func (f *Field) readValue() error {
	if !f.found {
		if f.required {
			return fmt.Errorf(
				"conflag: Required configuration value %s not found.",
				f.description,
			)
		}
		return nil
	}

	switch f.kind {
	case boolFieldType:
		val := false
		if f.parsedValue == "true" {
			val = true
		}
		f.destination.SetBool(val)
	case intFieldType:
		val, err := strconv.ParseInt(f.parsedValue, 10, 64)
		if err != nil {
			return fmt.Errorf(
				"conflag: Couldn't parse %s as integer.",
				f.parsedValue,
			)
		}
		f.destination.SetInt(val)
	case uintFieldType:
		val, err := strconv.ParseUint(f.parsedValue, 10, 64)
		if err != nil {
			return fmt.Errorf(
				"conflag: Couldn't parse %s as unsigned integer.",
				f.parsedValue,
			)
		}
		f.destination.SetUint(val)
	case floatFieldType:
		val, err := strconv.ParseFloat(f.parsedValue, 64)
		if err != nil {
			return fmt.Errorf(
				"conflag: Couldn't parse %s as floating point number.",
				f.parsedValue,
			)
		}
		f.destination.SetFloat(val)
	case stringFieldType:
		f.destination.SetString(f.parsedValue)
	}

	return nil
}
