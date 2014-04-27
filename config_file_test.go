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
	. "gopkg.in/check.v1"
	"strings"
	"testing"
)

type ConfigFileSuite struct {
	fields map[string]*concreteField
}

func (s *ConfigFileSuite) SetUpTest(c *C) {
	dest := &testConfig{}
	config, err := New(dest)
	c.Assert(err, IsNil)
	c.Assert(config, NotNil)

	config.Field("BoolField").FileCategory("bool_category").FileKey("bool_key")
	concrete := config.(*concreteConfig)
	s.fields = concrete.fields
}

func TestConfigFile(t *testing.T) {
	Suite(&ConfigFileSuite{})
	TestingT(t)
}

func (s *ConfigFileSuite) TestSuccessfulRead(c *C) {
	file := `
		# Commented out line
		uint_field = 50
		float_32_field = 0.5
		string_field = String! = 

		[ struct_field ]
		bool_field = true

		[bool_category]
		bool_key = false`

	reader := strings.NewReader(file)
	err := readConfigFile(s.fields, reader)
	c.Assert(err, IsNil)

	c.Assert(s.fields["UintField"].found, Equals, true)
	c.Assert(s.fields["UintField"].parsedValue, Equals, "50")
	c.Assert(s.fields["Float32Field"].found, Equals, true)
	c.Assert(s.fields["Float32Field"].parsedValue, Equals, "0.5")
	c.Assert(s.fields["StringField"].found, Equals, true)
	c.Assert(s.fields["StringField"].parsedValue, Equals, "String! =")
	c.Assert(s.fields["StructField.BoolField"].found, Equals, true)
	c.Assert(s.fields["StructField.BoolField"].parsedValue, Equals, "true")
	c.Assert(s.fields["BoolField"].found, Equals, true)
	c.Assert(s.fields["BoolField"].parsedValue, Equals, "false")
	c.Assert(s.fields["IntField"].found, Equals, false)
}

func (s *ConfigFileSuite) TestInvalidConfigLineFails(c *C) {
	file := `
		# Commented out line
		uint_field 50
`

	reader := strings.NewReader(file)
	err := readConfigFile(s.fields, reader)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "Invalid configuration line: uint_field 50")
}

func (s *ConfigFileSuite) TestInvalidKeyFails(c *C) {
	file := `
		# Commented out line
		uint_fied = 50
`

	reader := strings.NewReader(file)
	err := readConfigFile(s.fields, reader)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "Invalid configuration file key: uint_fied")
}
