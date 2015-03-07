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
)

func readCommandLineFlags(
	dest map[string]*Field,
	src []string,
	extraArgsAllowed bool,
) ([]string, error) {
	fieldsByShortFlag, fieldsByLongFlag := buildFlagIndices(dest)
	extraArgs := make([]string, 0)

	for i := 0; i < len(src); i++ {
		if len(src[i]) > 2 && src[i][0:2] == "--" {
			field, ok := fieldsByLongFlag[src[i][2:]]
			if !ok {
				return nil, fmt.Errorf("Unexpected flag %s", src[i][2:])
			}
			field.found = true
			if field.kind == boolFieldType {
				if src[i][2:] == field.longFlag {
					field.parsedValue = "true"
				} else {
					field.parsedValue = "false"
				}
			} else {
				i++
				if i >= len(src) {
					return nil, errors.New("Expected argument to last flag")
				}
				field.parsedValue = src[i]
			}
		} else if len(src[i]) > 1 && src[i][0:1] == "-" {
			deltaI := 0
			for _, v := range []rune(src[i][1:]) {
				field, ok := fieldsByShortFlag[v]
				if !ok {
					err := fmt.Errorf("Unexpected flag %s", string([]rune{v}))
					return nil, err
				}
				field.found = true
				if field.kind == boolFieldType {
					if field.shortFlag == v {
						field.parsedValue = "true"
					} else {
						field.parsedValue = "false"
					}
				} else {
					if deltaI == 0 {
						deltaI = 1
					}
					if i+1 >= len(src) {
						return nil, errors.New("Expected argument to last flag")
					}
					field.parsedValue = src[i+1]
				}
			}
			i += deltaI
		} else {
			extraArgs = append(extraArgs, src[i])
		}
	}

	if len(extraArgs) > 0 && !extraArgsAllowed {
		return nil, fmt.Errorf("Unexpected argument %s", extraArgs[0])
	}
	return extraArgs, nil
}

func buildFlagIndices(
	fields map[string]*Field,
) (shortIndex map[rune]*Field, longIndex map[string]*Field) {
	shortIndex = make(map[rune]*Field, len(fields))
	longIndex = make(map[string]*Field, len(fields))
	for _, v := range fields {
		for _, shortFlag := range []rune{v.shortFlag, v.inverseShortFlag} {
			if shortFlag != 0 {
				if _, ok := shortIndex[shortFlag]; ok {
					panic(
						fmt.Errorf(
							"conflag: Short flag %s used twice",
							string([]rune{shortFlag}),
						),
					)
				}
				shortIndex[shortFlag] = v
			}
		}
		for _, longFlag := range []string{v.longFlag, v.inverseLongFlag} {
			if longFlag != "" {
				if _, ok := longIndex[longFlag]; ok {
					panic(
						fmt.Errorf(
							"conflag: Long flag %s used twice",
							longFlag,
						),
					)
				}
				longIndex[longFlag] = v
			}
		}
	}
	return
}
