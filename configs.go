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
	"code.google.com/p/goconf/conf"
	"flag"
)

func readFlags(options map[string]*option, configFile string) string {
	flag.StringVar(
		&configFile,
		"config-file",
		configFile,
		"Custom config file path",
	)
	for _, option := range options {
		name := option.name
		if option.section != "default" {
			name = option.section + "." + option.name
		}

		switch option.typeOf {
		case "bool":
			flag.BoolVar(&option.boolVal, name, option.boolVal, option.usage)
		case "float64":
			flag.Float64Var(
				&option.floatVal,
				name,
				option.floatVal,
				option.usage,
			)
		case "int":
			flag.IntVar(&option.intVal, name, option.intVal, option.usage)
		case "string":
			flag.StringVar(
				&option.stringVal,
				name,
				option.stringVal,
				option.usage,
			)
		}
	}

	flag.Parse()
	flag.Visit(
		func(f *flag.Flag) {
			if option, ok := options[f.Name]; ok {
				option.read = true
				option.readFromFlag = true
			}
		},
	)

	return configFile
}

func readFile(options map[string]*option, configFile string) {
	c, err := conf.ReadConfigFile(configFile)
	if err != nil {
		return
	}

	for _, option := range options {
		if !option.readFromFlag && c.HasOption(option.section, option.name) {
			switch option.typeOf {
			case "bool":
				val, err := c.GetBool(option.section, option.name)
				if err == nil {
					option.boolVal = val
					option.read = true
				}
			case "float64":
				val, err := c.GetFloat64(option.section, option.name)
				if err == nil {
					option.floatVal = val
					option.read = true
				}
			case "int":
				val, err := c.GetInt(option.section, option.name)
				if err == nil {
					option.intVal = val
					option.read = true
				}
			case "string":
				val, err := c.GetString(option.section, option.name)
				if err == nil {
					option.stringVal = val
					option.read = true
				}
			}
		}
	}
}
