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

type fieldType int

const (
	invalidFieldType fieldType = iota
	boolFieldType
	intFieldType
	uintFieldType
	floatFieldType
	stringFieldType
)

var allowedTypes = map[fieldType]map[reflect.Kind]bool{
	boolFieldType: map[reflect.Kind]bool{
		reflect.Bool: true,
	},
	intFieldType: map[reflect.Kind]bool{
		reflect.Int:   true,
		reflect.Int8:  true,
		reflect.Int16: true,
		reflect.Int32: true,
		reflect.Int64: true,
	},
	uintFieldType: map[reflect.Kind]bool{
		reflect.Uint:   true,
		reflect.Uint8:  true,
		reflect.Uint16: true,
		reflect.Uint32: true,
		reflect.Uint64: true,
	},
	floatFieldType: map[reflect.Kind]bool{
		reflect.Float32: true,
		reflect.Float64: true,
	},
	stringFieldType: map[reflect.Kind]bool{
		reflect.String: true,
	},
}

// Field represents a single field in a configuration.  You can get it
// from the Config struct using its Field() method, and then set
// command-line and config-file properties of the field with it.
type Field interface {
	// Usage sets the usage text to display for the given field.
	Usage(usage string) Field

	// Required indicates that the field must be specified in either
	// the config file or a command-line parameter.
	Required() Field

	// LongFlag sets the long command-line flag for the option, to be
	// found on the command line in the form --long-flag.
	LongFlag(flag string) Field

	// InverseLongFlag sets the command-line flag to set the option to
	// false, to be found on the command line in the form
	// --inverse-long-flag.  Only usable on boolean fields.
	InverseLongFlag(flag string) Field

	// ShortFlag sets the short command-line flag for the option, to
	// be found on the command line in the form -f.
	ShortFlag(flag rune) Field

	// InverseShortFlag sets the short command-line flag to set the
	// option to false, to be found on the command line in the form
	// -i.  Only usable on boolean fields.
	InverseShortFlag(flag rune) Field

	// FileCategory sets the config file category the option will be
	// found under.  An empty string indicates none.
	FileCategory(category string) Field

	// FileKey sets the key in the config file for the option.
	FileKey(key string) Field
}

type concreteField struct {
	destination      reflect.Value
	kind             fieldType
	usage            string
	required         bool
	found            bool
	parsedValue      string
	longFlag         string
	shortFlag        rune
	inverseLongFlag  string
	inverseShortFlag rune
	fileCategory     string
	fileKey          string
}

func processField(
	fields map[string]*concreteField,
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

	kind := invalidFieldType
	for currentKind, subMap := range allowedTypes {
		if _, ok := subMap[field.Type().Kind()]; ok {
			kind = currentKind
			break
		}
	}
	if kind == invalidFieldType {
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
	longFlag := strings.Replace(fileKey, "_", "-", -1)
	if prefix != "" {
		fileCategory = strings.ToLower(
			caseTransition.ReplaceAllString(prefix, "${1}${3}_${2}${4}"),
		)
		longFlag = strings.Replace(fileCategory, "_", "-", -1) + "." + longFlag
	}

	key := name
	if prefix != "" {
		key = prefix + "." + key
	}

	fields[key] = &concreteField{
		usage:        "",
		destination:  field,
		kind:         kind,
		required:     false,
		found:        false,
		parsedValue:  "",
		longFlag:     longFlag,
		shortFlag:    0,
		fileCategory: fileCategory,
		fileKey:      fileKey,
	}

	return nil
}

func (f *concreteField) Usage(usage string) Field {
	f.usage = usage
	return f
}

func (f *concreteField) Required() Field {
	f.required = true
	return f
}

func (f *concreteField) LongFlag(flag string) Field {
	f.longFlag = flag
	return f
}

func (f *concreteField) InverseLongFlag(flag string) Field {
	f.inverseLongFlag = flag
	if f.kind != boolFieldType {
		panic(errors.New("Only boolean fields may have inverse flags."))
	}
	return f
}

func (f *concreteField) ShortFlag(flag rune) Field {
	f.shortFlag = flag
	return f
}

func (f *concreteField) InverseShortFlag(flag rune) Field {
	f.inverseShortFlag = flag
	if f.kind != boolFieldType {
		panic(
			errors.New("conflag: Only boolean fields may have inverse flags."),
		)
	}
	return f
}

func (f *concreteField) FileCategory(category string) Field {
	if strings.Contains(category, ".") {
		panic(
			errors.New("conflag: File category names cannot include '.'"),
		)
	}
	f.fileCategory = category
	return f
}

func (f *concreteField) FileKey(key string) Field {
	if strings.Contains(key, ".") {
		panic(errors.New("File key names cannot include '.'"))
	}
	f.fileKey = key
	return f
}
