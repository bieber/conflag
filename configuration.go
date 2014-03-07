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
	"io"
	"os"
	"reflect"
)

// Config defines an interface for representing, manipulating and
// reading program configuration.  You can create one with New(), set
// your desired options on it, and then read your program's
// configuration using Read().
type Config interface {
	// Field retrieves an individual field from the configuration for
	// you to modify.  It expects a field name matching an exported
	// field in the destination configuration struct.  To access a
	// subfield of an anonymous struct field, use
	// "OuterField.InnerField".  Returns nil for an invalid field
	// name.
	Field(field string) Field

	// ConfigReader sets an open io.Reader to read settings directly
	// from.  The caller is responsible for closing the Reader
	// afterwards.  If you intend to simply open a file on disk,
	// consider using the convenience function ConfigFile.
	ConfigReader(file io.Reader) Config

	// ConfigFile sets a file path to read a config file from.  If the
	// file is not present or otherwise unopenable, it will simply be
	// ignored.
	ConfigFile(fileName string) Config

	// ConfigFileRequired sets a file path to read a config file from,
	// and requires that the file is readable.  If the file can't be
	// opened, subsequent calls to Parse will return an error.
	ConfigFileRequired(fileName string) Config

	// Args sets a slice of command-line arguments to parse settings
	// from.  If you intend to use the arguments presented on the
	// command-line by the user, consider using the convenience
	// function OSArgs.
	Args(args []string) Config

	// AllowExtraArgs allows the user to enter command-line arguments
	// after any flags without triggering an error.  usage should
	// specify the usage text to include for the extra arguments in
	// the first line of the program usage text.  These arguments will
	// be returned from Parse.
	AllowExtraArgs(usage string) Config

	// Parse reads configuration from the available sources into the
	// specified fields of the config struct.  It returns a slice of
	// strings with any extra arguments (which will trigger an error
	// if not explicitly allowed via AllowExtraArgs) and an error
	// which will be nil if the configuration was processed
	// successfully.
	Parse() ([]string, error)
}

type concreteConfig struct {
	destination      reflect.Value
	fields           map[string]Field
	fileName         string
	file             io.Reader
	fileRequired     bool
	args             []string
	extraArgsAllowed bool
}

// New creates a new Config based on a destination struct.  The
// destination parameter must be a pointer to a struct containing
// fields of the allowed types (bool, int*, uint*, float*, and
// string).  Top-level fields may also be anonymous structs containing
// fields of the allowed types, but these can only go a single level
// deep.  By default nested structs as fields will represent sections
// of a config file.
//
// For each field, New will set a default file category, file key, and
// long command-line flag.  Both are formed by converting the field
// names from the destination struct into lower-camel-case,
// e.g. "ExampleField" becomes "example_field".  The file key as
// always the field name in lower-camel-case.  The file category will
// be the name of the enclosing anonymous struct field in the same
// style, if it exists.  The long-form command-line flag will be the
// same as the file key for top-level fields, and for nested fields it
// will be of the form category_name.field_name.
func New(destination interface{}) (Config, error) {
	destValue := reflect.ValueOf(destination)

	if destValue.Type().Kind() != reflect.Ptr {
		return nil, errors.New(
			"conflag: The config destination must be a pointer to a struct.",
		)
	}
	destValue = reflect.Indirect(destValue)
	if destValue.Type().Kind() != reflect.Struct {
		return nil, errors.New(
			"conflag: The config destination must be a pointer to a struct.",
		)
	}

	config := &concreteConfig{
		destination:      destValue,
		fields:           map[string]Field{},
		fileName:         "",
		file:             nil,
		fileRequired:     false,
		args:             os.Args,
		extraArgsAllowed: false,
	}
	for i := 0; i < destValue.NumField(); i++ {
		field := destValue.FieldByIndex([]int{i})
		err := processField(
			config.fields,
			field,
			"",
			destValue.Type().Field(i).Name,
		)
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}

func (c *concreteConfig) Field(field string) Field {
	if val, ok := c.fields[field]; ok {
		return val
	}
	panic(
		fmt.Errorf(
			"Field %s isn't present in your configuration struct.",
			field,
		),
	)
}

func (c *concreteConfig) ConfigReader(file io.Reader) Config {
	if c.file != nil || c.fileName != "" {
		panic(
			errors.New("You have already set a config file for this config."),
		)
	}
	c.file = file
	return c
}

func (c *concreteConfig) ConfigFile(fileName string) Config {
	if c.file != nil || c.fileName != "" {
		panic(
			errors.New("You have already set a config file for this config."),
		)
	}
	c.fileName = fileName
	return c
}

func (c *concreteConfig) ConfigFileRequired(fileName string) Config {
	if c.file != nil || c.fileName != "" {
		panic(
			errors.New("You have already set a config file for this config."),
		)
	}
	c.fileName = fileName
	c.fileRequired = true
	return c
}

func (c *concreteConfig) Args(args []string) Config {
	c.args = args
	return c
}

func (c *concreteConfig) AllowExtraArgs(usage string) Config {
	c.extraArgsAllowed = true
	return c
}

func (c *concreteConfig) Parse() ([]string, error) {
	return nil, errors.New("Not implemented yet.")
}
