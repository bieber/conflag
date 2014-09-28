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
	"io"
	"os"
)

func (c *concreteConfig) findConfigFile() (
	f io.Reader,
	args []string,
	err error,
) {
	err = nil
	f = c.file
	args = c.args

	fileName := c.fileName
	foundFlag := false
	if c.fileShortFlag != 0 || c.fileLongFlag != "" {
		shortFlag := ""
		longFlag := ""
		if c.fileShortFlag != 0 {
			shortFlag = "-" + string([]rune{c.fileShortFlag})
		}
		if c.fileLongFlag != "" {
			longFlag = "--" + c.fileLongFlag
		}

		for i := len(args) - 1; i >= 0; i-- {
			if args[i] == longFlag || args[i] == shortFlag {
				if foundFlag {
					if closer, ok := f.(io.Closer); ok {
						closer.Close()
					}
					f = nil
					err = errors.New("conflag: Duplicate config file flags")
					return
				} else {
					foundFlag = true
				}

				if i+1 > len(args)-1 {
					err = errors.New("conflag: Missing config file name")
					return
				}

				fileName = args[i+1]
				// Removing file names and flags from the args so they
				// don't cause problems later
				args = append(args[:i], args[i+2:]...)
			}
		}
	}

	if (f == nil && fileName != "") || foundFlag {
		if closer, ok := f.(io.Closer); ok {
			closer.Close()
		}
		f, err = os.Open(fileName)
	}
	return
}
