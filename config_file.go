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
	"bufio"
	"fmt"
	"io"
	"strings"
)

func readConfigFile(dest map[string]*Field, src io.Reader) error {
	fields := buildConfigFileIndex(dest)
	scanner := bufio.NewScanner(src)

	category := ""
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		if line[0] == '[' && line[len(line)-1] == ']' {
			category = strings.TrimSpace(line[1 : len(line)-1])
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("Invalid configuration line: %s", line)
		}
		key := strings.TrimSpace(parts[0])
		if category != "" {
			key = category + "." + key
		}
		value := strings.TrimSpace(parts[1])

		if _, ok := fields[key]; !ok {
			return fmt.Errorf("Invalid configuration file key: %s", key)
		}
		field := fields[key]
		field.parsedValue = value
		field.found = true
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	if closer, ok := src.(io.Closer); ok {
		closer.Close()
	}
	return nil
}

// Get fields indexed by their file category and key instead of config struct
func buildConfigFileIndex(
	fields map[string]*Field,
) map[string]*Field {
	index := make(map[string]*Field, len(fields))
	for _, v := range fields {
		if key := v.fileKey; key != "" {
			if v.fileCategory != "" {
				key = v.fileCategory + "." + key
			}
			index[key] = v
		}
	}
	return index
}
