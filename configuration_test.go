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

type ConfigSuite struct {
	dest testConfig
}

func TestConfig(t *testing.T) {
	Suite(&ConfigSuite{})
	TestingT(t)
}

func (s *ConfigSuite) SetUpTest(c *C) {
	s.dest = testConfig{}
}

func (s *ConfigSuite) TestConfigFileReader(c *C) {
	config, err := New(&s.dest)
	c.Assert(err, IsNil)
	c.Assert(config, NotNil)

	reader := strings.NewReader("config string")
	config.ConfigReader(reader)

	concrete := config.(*concreteConfig)
	c.Assert(concrete.fileName, Equals, "")
	c.Assert(concrete.fileRequired, Equals, false)
	c.Assert(concrete.file, NotNil)
}

func (s *ConfigSuite) TestConfigFileName(c *C) {
	config, err := New(&s.dest)
	c.Assert(err, IsNil)
	c.Assert(config, NotNil)

	config.ConfigFile("/file/name")

	concrete := config.(*concreteConfig)
	c.Assert(concrete.fileName, Equals, "/file/name")
	c.Assert(concrete.fileRequired, Equals, false)
	c.Assert(concrete.file, IsNil)
}

func (s *ConfigSuite) TestConfigFileRequired(c *C) {
	config, err := New(&s.dest)
	c.Assert(err, IsNil)
	c.Assert(config, NotNil)

	config.ConfigFileRequired()

	concrete := config.(*concreteConfig)
	c.Assert(concrete.fileName, Equals, "")
	c.Assert(concrete.fileRequired, Equals, true)
	c.Assert(concrete.file, IsNil)
}

func (s *ConfigSuite) TestMultipleConfigFilesFailure(c *C) {
	defer func() {
		c.Assert(recover(), NotNil)
	}()

	config, err := New(&s.dest)
	c.Assert(err, IsNil)
	c.Assert(config, NotNil)

	reader := strings.NewReader("config file")
	config.ConfigReader(reader)
	config.ConfigFile("/file/name")
}

func (s *ConfigSuite) TestArgs(c *C) {
	config, err := New(&s.dest)
	c.Assert(err, IsNil)
	c.Assert(config, NotNil)

	args := []string{"slice", "of", "args"}
	config.Args(args)

	concrete := config.(*concreteConfig)
	c.Assert(len(concrete.args), Equals, 3)
	c.Assert(concrete.args, DeepEquals, []string{"slice", "of", "args"})
}
