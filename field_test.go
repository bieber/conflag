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
	"testing"
)

type FieldSuite struct {
	dest testConfig
}

func (s *FieldSuite) SetUpTest(c *C) {
	s.dest = testConfig{}
}

func TestFlag(t *testing.T) {
	Suite(&FieldSuite{})
	TestingT(t)
}

func (s *FieldSuite) TestValidConfig(c *C) {
	config, err := New(&s.dest)
	c.Assert(err, IsNil)
	c.Assert(config, NotNil)

	intField := config.Field("IntField")
	c.Assert(intField.kind, Equals, intFieldType)
	c.Assert(intField.required, Equals, false)
	c.Assert(intField.found, Equals, false)
	c.Assert(intField.longFlag, Equals, "int-field")
	c.Assert(intField.shortFlag, Equals, int32(0))
	c.Assert(intField.fileCategory, Equals, "")
	c.Assert(intField.fileKey, Equals, "int_field")

	nestedIntField := config.Field("StructField.IntField")
	c.Assert(nestedIntField.kind, Equals, intFieldType)
	c.Assert(nestedIntField.required, Equals, false)
	c.Assert(nestedIntField.found, Equals, false)
	c.Assert(nestedIntField.longFlag, Equals, "struct-field.int-field")
	c.Assert(nestedIntField.shortFlag, Equals, int32(0))
	c.Assert(nestedIntField.fileCategory, Equals, "struct_field")
	c.Assert(nestedIntField.fileKey, Equals, "int_field")
}

func (s *FieldSuite) TestFieldModifiers(c *C) {
	config, err := New(&s.dest)
	c.Assert(err, IsNil)
	c.Assert(config, NotNil)

	stringField := config.Field("StringField")
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

	boolField := config.Field("BoolField")
	boolField.
		InverseShortFlag('i').
		InverseLongFlag("inverse-bool")
	c.Assert(boolField.inverseShortFlag, Equals, 'i')
	c.Assert(boolField.inverseLongFlag, Equals, "inverse-bool")
}

func (s *FieldSuite) TestInverseShortFlagFailure(c *C) {
	config, err := New(&s.dest)
	c.Assert(err, IsNil)
	c.Assert(config, NotNil)

	defer func() {
		c.Assert(recover(), NotNil)
	}()
	config.Field("StringField").InverseShortFlag('i')
}

func (s *FieldSuite) TestInverseLongFlagFailure(c *C) {
	config, err := New(&s.dest)
	c.Assert(err, IsNil)
	c.Assert(config, NotNil)

	defer func() {
		c.Assert(recover(), NotNil)
	}()
	config.Field("StringField").InverseLongFlag("inverse")
}

func (s *FieldSuite) TestNonPointerFails(c *C) {
	config, err := New(s.dest)
	c.Assert(err, NotNil)
	c.Assert(config, IsNil)
}

func (s *FieldSuite) TestNonStructFails(c *C) {
	x := 5
	config, err := New(&x)
	c.Assert(err, NotNil)
	c.Assert(config, IsNil)
}

func (s *FieldSuite) TestDeepNestingFails(c *C) {
	dest := struct{ A struct{ B struct{ C int } } }{}
	config, err := New(&dest)
	c.Assert(err, NotNil)
	c.Assert(config, IsNil)
}

func (s *FieldSuite) TestWrongFieldTypeFails(c *C) {
	dest := struct{ A *int }{}
	config, err := New(&dest)
	c.Assert(err, NotNil)
	c.Assert(config, IsNil)
}

func (s *FieldSuite) TestCategoryWithDotFails(c *C) {
	defer func() {
		c.Assert(recover(), NotNil)
	}()

	config, err := New(&s.dest)
	c.Assert(err, IsNil)
	c.Assert(config, NotNil)

	config.Field("UintField").FileCategory("test.category")
}

func (s *FieldSuite) TestKeyWithDotFails(c *C) {
	defer func() {
		c.Assert(recover(), NotNil)
	}()

	config, err := New(&s.dest)
	c.Assert(err, IsNil)
	c.Assert(config, NotNil)

	config.Field("UintField").FileKey("test.key")
}
