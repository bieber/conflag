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

package conflag

import (
	"fmt"
	"reflect"
	"strings"
)

func getOptions(dest interface{}) map[string]*option {
	t := reflect.Indirect(reflect.ValueOf(dest)).Type()
	options := make(map[string]*option, 10)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		opt := option{section: "default", structField: f.Name}

		opt.name = strings.ToLower(f.Name)
		if f.Type.PkgPath() != "" {
			panic(
				fmt.Errorf(
					"conflag: Field %s's type has non-default package path %s.",
					f.Name,
					f.Type.PkgPath(),
				),
			)
		}
		switch f.Type.Name() {
		case "bool":
			fallthrough
		case "float64":
			fallthrough
		case "int":
			fallthrough
		case "string":
			opt.typeOf = f.Type.Name()
		default:
			panic(
				fmt.Errorf(
					"conflag: Field %s has non-supported type %s",
					f.Name,
					f.Type.Name(),
				),
			)
		}

		tag := f.Tag.Get("conflag")
		if tag != "" {
			parts := strings.Split(tag, "|")

			name := parts[0]
			nameParts := strings.Split(name, ".")
			if len(nameParts) == 1 {
				opt.name = nameParts[0]
			} else if len(nameParts) == 2 {
				opt.section = nameParts[0]
				opt.name = nameParts[1]
			} else {
				panic(
					fmt.Errorf(
						"conflag: Malformed conflag field name %s",
						name,
					),
				)
			}

			if len(parts) > 1 {
				opt.usage = parts[1]
			}

			if len(parts) > 2 {
				opt.read = true
				def := parts[2]
				var format string
				var dest interface{}
				switch opt.typeOf {
				case "bool":
					format = "%t"
					dest = &opt.boolVal
				case "float64":
					format = "%f"
					dest = &opt.floatVal
				case "int":
					format = "%d"
					dest = &opt.intVal
				case "string":
					format = "%s"
					dest = &opt.stringVal
				}
				n, err := fmt.Sscanf(def, format, dest)
				if n != 1 || err != nil {
					fmt.Println(n, err, opt.typeOf, format)
					panic(
						fmt.Errorf(
							"conflag: Malformed conflag default value %s.",
							def,
						),
					)
				}
			}

			if len(parts) > 3 {
				panic(
					fmt.Errorf(
						"conflag: Too many parts to field tag for %s.",
						f.Name,
					),
				)
			}
		}

		fullName := opt.name
		if opt.section != "default" {
			fullName = opt.section + "." + opt.name
		}
		options[fullName] = &opt
	}

	return options
}

func setOptions(dest interface{}, options map[string]*option) {
	val := reflect.Indirect(reflect.ValueOf(dest))
	for _, option := range options {
		field := val.FieldByName(option.structField)
		switch option.typeOf {
		case "bool":
			field.SetBool(option.boolVal)
		case "float64":
			field.SetFloat(option.floatVal)
		case "int":
			field.SetInt(int64(option.intVal))
		case "string":
			field.SetString(option.stringVal)
		}
	}
}
