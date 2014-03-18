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

type CommandLineSuite struct {
	fields map[string]*concreteField
}

func (s *CommandLineSuite) SetUpTest(c *C) {
	dest := &testConfig{}
	config, err := New(dest)
	c.Assert(err, IsNil)
	c.Assert(config, NotNil)

	config.Field("BoolField").
		ShortFlag('b').
		LongFlag("bool").
		InverseShortFlag('i').
		InverseLongFlag("inverse-bool")
	config.Field("UintField").
		ShortFlag('u').
		LongFlag("unsigned-int")
	config.Field("StructField.BoolField").
		ShortFlag('c')
	concrete := config.(*concreteConfig)
	s.fields = concrete.fields
}

func TestCommandLine(t *testing.T) {
	Suite(&CommandLineSuite{})
	TestingT(t)
}

func (s *CommandLineSuite) TestSimpleUsage(c *C) {
	extras, err := readCommandLineFlags(
		s.fields,
		[]string{"--bool", "--unsigned-int", "5"},
		false,
	)
	c.Assert(err, IsNil)
	c.Assert(len(extras), Equals, 0)
	c.Assert(s.fields["BoolField"].found, Equals, true)
	c.Assert(s.fields["BoolField"].parsedValue, Equals, "true")
	c.Assert(s.fields["UintField"].found, Equals, true)
	c.Assert(s.fields["UintField"].parsedValue, Equals, "5")

	extras, err = readCommandLineFlags(
		s.fields,
		[]string{"-u", "5", "-b"},
		false,
	)
	c.Assert(err, IsNil)
	c.Assert(len(extras), Equals, 0)
	c.Assert(s.fields["BoolField"].found, Equals, true)
	c.Assert(s.fields["BoolField"].parsedValue, Equals, "true")
	c.Assert(s.fields["UintField"].found, Equals, true)
	c.Assert(s.fields["UintField"].parsedValue, Equals, "5")
}

func (s *CommandLineSuite) TestShortFlagCombination(c *C) {
	extras, err := readCommandLineFlags(
		s.fields,
		[]string{"-bc"},
		false,
	)
	c.Assert(err, IsNil)
	c.Assert(len(extras), Equals, 0)
	c.Assert(s.fields["BoolField"].found, Equals, true)
	c.Assert(s.fields["BoolField"].parsedValue, Equals, "true")
	c.Assert(s.fields["StructField.BoolField"].found, Equals, true)
	c.Assert(s.fields["StructField.BoolField"].parsedValue, Equals, "true")

	extras, err = readCommandLineFlags(
		s.fields,
		[]string{"-bu", "5"},
		false,
	)
	c.Assert(err, IsNil)
	c.Assert(len(extras), Equals, 0)
	c.Assert(s.fields["BoolField"].found, Equals, true)
	c.Assert(s.fields["BoolField"].parsedValue, Equals, "true")
	c.Assert(s.fields["UintField"].found, Equals, true)
	c.Assert(s.fields["UintField"].parsedValue, Equals, "5")
}

func (s *CommandLineSuite) TestExtraArgs(c *C) {
	extras, err := readCommandLineFlags(
		s.fields,
		[]string{"extra", "-b", "flags"},
		true,
	)
	c.Assert(err, IsNil)
	c.Assert(extras, DeepEquals, []string{"extra", "flags"})
}

func (s *CommandLineSuite) TestExtraArgFailure(c *C) {
	extras, err := readCommandLineFlags(
		s.fields,
		[]string{"-b", "extra"},
		false,
	)
	c.Assert(err, NotNil)
	c.Assert(extras, IsNil)
}

func (s *CommandLineSuite) TestExpectedArgFailure(c *C) {
	extras, err := readCommandLineFlags(
		s.fields,
		[]string{"-u"},
		false,
	)
	c.Assert(err, NotNil)
	c.Assert(extras, IsNil)
}

func (s *CommandLineSuite) TestShortFlagCollisionError(c *C) {
	dest := &testConfig{}
	config, err := New(dest)
	concrete := config.(*concreteConfig)
	c.Assert(err, IsNil)
	c.Assert(config, NotNil)

	config.Field("BoolField").
		ShortFlag('b')
	config.Field("UintField").
		ShortFlag('b')

	defer func() {
		c.Assert(recover(), NotNil)
	}()
	readCommandLineFlags(
		concrete.fields,
		[]string{},
		false,
	)
}

func (s *CommandLineSuite) TestLongFlagCollisionError(c *C) {
	dest := &testConfig{}
	config, err := New(dest)
	concrete := config.(*concreteConfig)
	c.Assert(err, IsNil)
	c.Assert(config, NotNil)

	config.Field("BoolField").
		LongFlag("flag")
	config.Field("UintField").
		LongFlag("flag")

	defer func() {
		c.Assert(recover(), NotNil)
	}()
	readCommandLineFlags(
		concrete.fields,
		[]string{},
		false,
	)
}
