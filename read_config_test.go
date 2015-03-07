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

// Most every component being used in the Read() method already has
// its own unit test, so this is largely an integration test of them
// all. The only really unique functionality that it's checking is
// that of the field parsing.

type ReadConfigSuite struct {
	destination *testConfig
	config      *Config
}

func (s *ReadConfigSuite) SetUpTest(c *C) {
	s.destination = &testConfig{}
	config, err := New(s.destination)
	s.config = config
	c.Assert(err, IsNil)
}

func TestReadConfig(t *testing.T) {
	Suite(&ReadConfigSuite{})
	TestingT(t)
}

func (s *ReadConfigSuite) TestFullConfig(c *C) {
	s.config.AllowExtraArgs("files")
	s.config.Field("BoolField").
		ShortFlag('b').
		Required()
	s.config.Field("UintField").
		LongFlag("uint").
		Required()
	s.config.Field("IntField").
		FileKey("integer").
		Required()
	s.config.Field("Float32Field").
		FileCategory("floats").
		FileKey("floating_point_32").
		Required()
	s.config.Field("Float64Field").
		FileCategory("floats").
		FileKey("floating_point_64").
		Required()
	s.config.Field("StringField").
		ShortFlag('s').
		LongFlag("string").
		Required()
	s.config.Field("StructField.BoolField").
		ShortFlag('p').
		InverseShortFlag('q').
		Required()
	s.config.Field("StructField.UintField").
		ShortFlag('t').
		Required()
	s.config.Field("StructField.Float32Field").
		FileCategory("").
		FileKey("unboxed_32").
		LongFlag("float_32").
		Required()
	s.config.Field("StructField.Float64Field").
		Required()

	s.destination.StructField.IntField = 897

	file := `
      # Section-less values
      integer=127
      unboxed_32 =   2.4

      [floats]
      floating_point_32 = 54.29
      floating_point_64 = 23.23

      [struct_field]
      float_64_field = 56.49`

	reader := &closerStringReader{Reader: strings.NewReader(file)}
	s.config.ConfigReader(reader)
	s.config.Args(
		[]string{
			"-b",
			"--uint", "50",
			"-sqt", "200",
			"some", "extra", "args",
			"--float_32", "5.4",
		},
	)

	correctResult := &testConfig{
		BoolField:    true,
		UintField:    50,
		IntField:     127,
		Float32Field: 54.29,
		Float64Field: 23.23,
		StringField:  "200",
	}
	correctResult.StructField.BoolField = false
	correctResult.StructField.UintField = 200
	correctResult.StructField.IntField = 897
	correctResult.StructField.Float32Field = 5.4
	correctResult.StructField.Float64Field = 56.49

	extraArgs, err := s.config.Read()
	c.Assert(err, IsNil)
	c.Assert(extraArgs, DeepEquals, []string{"some", "extra", "args"})
	c.Assert(s.destination, DeepEquals, correctResult)
}

func (s *ReadConfigSuite) TestExtraArgsFailure(c *C) {
	extraArgs, err := s.config.
		Args([]string{"--unexpected-arg"}).
		Read()
	c.Assert(extraArgs, IsNil)
	c.Assert(err, NotNil)
}

func (s *ReadConfigSuite) TestRequiredFileFailure(c *C) {
	s.config.Args([]string{}).ConfigFileRequired()
	extraArgs, err := s.config.Read()
	c.Assert(extraArgs, IsNil)
	c.Assert(err, NotNil)
}

func (s *ReadConfigSuite) TestRequiredFieldFailure(c *C) {
	s.config.Args([]string{}).
		Field("IntField").
		Required()
	extraArgs, err := s.config.Read()
	c.Assert(extraArgs, IsNil)
	c.Assert(err, NotNil)
}

func (s *ReadConfigSuite) TestIntParseFailure(c *C) {
	extraArgs, err := s.config.
		Args([]string{"--int_field", "non-int"}).
		Read()
	c.Assert(extraArgs, IsNil)
	c.Assert(err, NotNil)
}

func (s *ReadConfigSuite) TestUintParseFailure(c *C) {
	extraArgs, err := s.config.
		Args([]string{"--uint_field", "-5"}).
		Read()
	c.Assert(extraArgs, IsNil)
	c.Assert(err, NotNil)
}

func (s *ReadConfigSuite) TestFloatParseFailure(c *C) {
	extraArgs, err := s.config.
		Args([]string{"--float_32_field", "non-float"}).
		Read()
	c.Assert(extraArgs, IsNil)
	c.Assert(err, NotNil)
}
