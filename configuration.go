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
	"reflect"
	"regexp"
	"strings"
)

var allowedTypes = map[reflect.Kind]bool{
	reflect.Bool:    true,
	reflect.Int:     true,
	reflect.Int8:    true,
	reflect.Int16:   true,
	reflect.Int32:   true,
	reflect.Int64:   true,
	reflect.Uint:    true,
	reflect.Uint8:   true,
	reflect.Uint16:  true,
	reflect.Uint32:  true,
	reflect.Uint64:  true,
	reflect.Float32: true,
	reflect.Float64: true,
	reflect.String:  true,
}

// Field represents a single field in a configuration.  You can get it
// from the Config struct using its Field() method, and then set
// command-line and config-file properties of the field with it.
type Field interface {
	// Required indicates that the field must be specified in either
	// the config file or a command-line parameter.
	Required() Field

	// LongFlag sets the long command-line flag for the option, to be
	// found on the command line in the form --long-flag.
	LongFlag(flag string) Field

	// ShortFlag sets the short command-line flag for the option, to
	// be found on the command line in the form -f.
	ShortFlag(flag rune) Field

	// FileCategory sets the config file category the option will be
	// found under.  An empty string indicates none.
	FileCategory(category string) Field

	// FileKey indicates the key in the config file for the option.
	FileKey(key string) Field
}

type concreteField struct {
	destination  reflect.Value
	required     bool
	found        bool
	longFlag     string
	shortFlag    rune
	fileCategory string
	fileKey      string
}

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
}

type concreteConfig struct {
	destination reflect.Value
	fields      map[string]Field
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
		destination: destValue,
		fields:      map[string]Field{},
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

func processField(
	fields map[string]Field,
	field reflect.Value,
	prefix string,
	name string,
) error {
	if field.Type().Kind() == reflect.Struct {
		if prefix != "" {
			return errors.New(
				"conflag: Configuration structs may only be nested one level.",
			)
		}
		for i := 0; i < field.NumField(); i++ {
			err := processField(
				fields,
				field.FieldByIndex([]int{i}),
				name,
				field.Type().Field(i).Name,
			)
			if err != nil {
				return err
			}
		}
		return nil
	}

	if _, ok := allowedTypes[field.Type().Kind()]; !ok {
		return fmt.Errorf(
			"conflag: Type %s is not allowed in configuration structs.",
			field.Type().Kind().String(),
		)
	}
	caseTransition, err := regexp.Compile("([a-z0-9])([A-Z])|([a-z])([A-Z0-9])")
	if err != nil {
		return err
	}

	fileCategory := ""
	fileKey := strings.ToLower(
		caseTransition.ReplaceAllString(name, "${1}${3}_${2}${4}"),
	)
	longFlag := fileKey
	if prefix != "" {
		fileCategory = strings.ToLower(
			caseTransition.ReplaceAllString(prefix, "${1}${3}_${2}${4}"),
		)
		longFlag = fileCategory + "." + longFlag
	}

	key := name
	if prefix != "" {
		key = prefix + "." + key
	}

	fields[key] = &concreteField{
		destination:  field,
		required:     false,
		found:        false,
		longFlag:     longFlag,
		shortFlag:    0,
		fileCategory: fileCategory,
		fileKey:      fileKey,
	}

	return nil
}

func (c *concreteConfig) Field(field string) Field {
	if val, ok := c.fields[field]; ok {
		return val
	}
	return nil
}

func (f *concreteField) Required() Field {
	f.required = true
	return f
}

func (f *concreteField) LongFlag(flag string) Field {
	f.longFlag = flag
	return f
}

func (f *concreteField) ShortFlag(flag rune) Field {
	f.shortFlag = flag
	return f
}

func (f *concreteField) FileCategory(category string) Field {
	f.fileCategory = category
	return f
}

func (f *concreteField) FileKey(key string) Field {
	f.fileKey = key
	return f
}
