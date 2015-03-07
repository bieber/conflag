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

// Config stores metadata about your configuration.  You can create
// one with New(), set your desired options on it, and then read your
// program's configuration using Read().
type Config struct {
	name             string
	description      string
	destination      reflect.Value
	fields           map[string]*Field
	fieldKeysInOrder []string
	fileName         string
	file             io.Reader
	fileShortFlag    rune
	fileLongFlag     string
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
func New(destination interface{}) (*Config, error) {
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

	config := &Config{
		name:             "",
		description:      "",
		destination:      destValue,
		fields:           map[string]*Field{},
		fieldKeysInOrder: []string{},
		fileName:         "",
		file:             nil,
		fileShortFlag:    0,
		fileLongFlag:     "",
		fileRequired:     false,
		args:             os.Args[1:],
		extraArgsAllowed: false,
	}
	for i := 0; i < destValue.NumField(); i++ {
		field := destValue.FieldByIndex([]int{i})
		err := processField(
			config.fields,
			&config.fieldKeysInOrder,
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

// ProgramName sets a name for the program to display in the usage
// text.
func (c *Config) ProgramName(name string) *Config {
	c.name = name
	return c
}

// ProgramDescription sets a description of the program to display in
// the usage text.
func (c *Config) ProgramDescription(description string) *Config {
	c.description = description
	return c
}

// Field retrieves an individual field from the configuration for you
// to modify.  It expects a field name matching an exported field in
// the destination configuration struct.  To access a subfield of an
// anonymous struct field, use "OuterField.InnerField".  Returns nil
// for an invalid field name.
func (c *Config) Field(field string) *Field {
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

// ConfigReader sets an open io.Reader to read settings directly from.
// If the value also implements io.Closer, it will be closed after
// reading.  If you intend to simply open a file on disk, consider
// using the convenience function ConfigFile.
func (c *Config) ConfigReader(file io.Reader) *Config {
	if c.file != nil || c.fileName != "" {
		panic(
			errors.New("You have already set a config file for this config."),
		)
	}
	c.file = file
	return c
}

// ConfigFile sets a file path to read a config file from.  If the
// file is not present or otherwise unopenable, it will simply be
// ignored.
func (c *Config) ConfigFile(fileName string) *Config {
	if c.file != nil || c.fileName != "" {
		panic(
			errors.New("You have already set a config file for this config."),
		)
	}
	c.fileName = fileName
	return c
}

// ConfigFileShortFlag sets a short command-line flag with which the
// user can specify a config file.  If this option is set and the user
// sets a config file, it will take precedence over a file specified
// with the ConfigReader or ConfigFile options.
func (c *Config) ConfigFileShortFlag(flag rune) *Config {
	c.fileShortFlag = flag
	return c
}

// ConfigFileLongFlag sets a long command-line flag with which the
// user can specify a config file.  If this option is set and the user
// sets a config file, it will take precedence over a file specified
// with the ConfigReader or ConfigFile options.
func (c *Config) ConfigFileLongFlag(flag string) *Config {
	c.fileLongFlag = flag
	return c
}

// ConfigFileRequired requires that a config file path is set (by
// calling ConfigFile or ConfigReader, or by command line argument
// with ConfigFileShortFlag or ConfigFileLongFlag) and the file
// exists.  If the file can't be opened, subsequent calls to Parse
// will return an error.
func (c *Config) ConfigFileRequired() *Config {
	c.fileRequired = true
	return c
}

// Args sets a slice of command-line arguments to parse settings from.
// If you don't explicitly set the command-line arguments, os.Args
// will be used as the default.
func (c *Config) Args(args []string) *Config {
	c.args = args
	return c
}

// AllowExtraArgs allows the user to enter command-line arguments
// after any flags without triggering an error.  usage should specify
// the usage text to include for the extra arguments in the first line
// of the program usage text.  These arguments will be returned from
// Parse.
func (c *Config) AllowExtraArgs(usage string) *Config {
	c.extraArgsAllowed = true
	return c
}
