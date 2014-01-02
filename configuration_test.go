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
	. "launchpad.net/gocheck"
	"testing"
)

type testConfig struct {
	BoolField    bool
	UintField    uint
	IntField     int
	Float32Field float32
	Float64Field float64
	StringField  string
	StructField  struct {
		BoolField    bool
		UintField    uint
		IntField     int
		Float32Field float32
		Float64Field float64
		StringField  string
	}
}

type ConfigSuite struct {
	dest testConfig
}

func (s *ConfigSuite) SetUpTest(c *C) {
	s.dest = testConfig{}
}

func Test(t *testing.T) {
	Suite(&ConfigSuite{})
	TestingT(t)
}

func (s *ConfigSuite) TestValidConfig(c *C) {
	config, err := New(&s.dest)
	c.Assert(err, IsNil)
	c.Assert(config, NotNil)

	intField := config.Field("IntField").(*concreteField)
	c.Assert(intField.required, Equals, false)
	c.Assert(intField.found, Equals, false)
	c.Assert(intField.longFlag, Equals, "int_field")
	c.Assert(intField.shortFlag, Equals, int32(0))
	c.Assert(intField.fileCategory, Equals, "")
	c.Assert(intField.fileKey, Equals, "int_field")

	nestedIntField := config.Field("StructField.IntField").(*concreteField)
	c.Assert(nestedIntField.required, Equals, false)
	c.Assert(nestedIntField.found, Equals, false)
	c.Assert(nestedIntField.longFlag, Equals, "struct_field.int_field")
	c.Assert(nestedIntField.shortFlag, Equals, int32(0))
	c.Assert(nestedIntField.fileCategory, Equals, "struct_field")
	c.Assert(nestedIntField.fileKey, Equals, "int_field")
}

func (s *ConfigSuite) TestFieldModifiers(c *C) {
	config, err := New(&s.dest)
	c.Assert(err, IsNil)
	c.Assert(config, NotNil)

	stringField := config.Field("StringField").(*concreteField)
	stringField.
		Required().
		ShortFlag('s').
		LongFlag("string").
		FileCategory("category").
		FileKey("key")
	c.Assert(stringField.required, Equals, true)
	c.Assert(stringField.shortFlag, Equals, 's')
	c.Assert(stringField.longFlag, Equals, "string")
	c.Assert(stringField.fileCategory, Equals, "category")
	c.Assert(stringField.fileKey, Equals, "key")
}

func (s *ConfigSuite) TestNonPointerFails(c *C) {
	config, err := New(s.dest)
	c.Assert(err, NotNil)
	c.Assert(config, IsNil)
}

func (s *ConfigSuite) TestNonStructFails(c *C) {
	x := 5
	config, err := New(&x)
	c.Assert(err, NotNil)
	c.Assert(config, IsNil)
}

func (s *ConfigSuite) TestDeepNestingFails(c *C) {
	dest := struct{ A struct{ B struct{ C int } } }{}
	config, err := New(&dest)
	c.Assert(err, NotNil)
	c.Assert(config, IsNil)
}

func (s *ConfigSuite) TestWrongFieldTypeFails(c *C) {
	dest := struct{ A *int }{}
	config, err := New(&dest)
	c.Assert(err, NotNil)
	c.Assert(config, IsNil)
}
