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
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

type FindConfigFileSuite struct {
	defaultTestFileName   string
	defaultTestFileReader io.Reader
	extraTestFileName     string
	config                *concreteConfig
}

func (s *FindConfigFileSuite) SetUpTest(c *C) {
	dest := &testConfig{}
	config, err := New(dest)
	c.Assert(err, IsNil)
	c.Assert(config, NotNil)
	s.config = config.(*concreteConfig)

	tempDir := c.MkDir()

	s.defaultTestFileName = path.Join(tempDir, "default")
	s.extraTestFileName = path.Join(tempDir, "extra")
	defaultFout, err := os.OpenFile(
		s.defaultTestFileName,
		os.O_WRONLY|os.O_CREATE,
		0666,
	)
	c.Assert(err, IsNil)
	extraFout, err := os.OpenFile(
		s.extraTestFileName,
		os.O_WRONLY|os.O_CREATE,
		0666,
	)
	c.Assert(err, IsNil)
	defer defaultFout.Close()
	defer extraFout.Close()

	defaultFout.Write([]byte("DEFAULT_FILE"))
	extraFout.Write([]byte("EXTRA_FILE"))
	s.defaultTestFileReader = strings.NewReader("DEFAULT_READER")
}

func TestFindConfigFile(t *testing.T) {
	Suite(&FindConfigFileSuite{})
	TestingT(t)
}

func (s *FindConfigFileSuite) TestNothing(c *C) {
	s.config.Args([]string{"some", "random", "flags"})
	fin, _, err := s.config.findConfigFile()
	c.Assert(fin, IsNil)
	c.Assert(err, IsNil)
}

func (s *FindConfigFileSuite) TestDefaultReader(c *C) {
	s.config.ConfigReader(s.defaultTestFileReader)
	s.config.Args([]string{"some", "random", "flags"})
	fin, _, err := s.config.findConfigFile()
	c.Assert(err, IsNil)
	assertFileContents(c, fin, "DEFAULT_READER")
}

func (s *FindConfigFileSuite) TestDefaultFile(c *C) {
	s.config.ConfigFile(s.defaultTestFileName)
	s.config.Args([]string{"some", "random", "flags"})
	fin, _, err := s.config.findConfigFile()
	c.Assert(err, IsNil)
	assertFileContents(c, fin, "DEFAULT_FILE")
}

func (s *FindConfigFileSuite) TestShortFlag(c *C) {
	s.config.ConfigFileShortFlag('c')
	s.config.Args([]string{"-c", s.extraTestFileName})
	fin, _, err := s.config.findConfigFile()
	c.Assert(err, IsNil)
	assertFileContents(c, fin, "EXTRA_FILE")
}

func (s *FindConfigFileSuite) TestLongFlag(c *C) {
	s.config.ConfigFileLongFlag("config-file")
	s.config.Args([]string{"--config-file", s.extraTestFileName})
	fin, _, err := s.config.findConfigFile()
	c.Assert(err, IsNil)
	assertFileContents(c, fin, "EXTRA_FILE")
}

func (s *FindConfigFileSuite) TestShortFlagOverrideReader(c *C) {
	reader := &closerStringReader{
		Reader: s.defaultTestFileReader,
		closed: false,
	}
	s.config.ConfigReader(reader)
	s.config.ConfigFileShortFlag('c')
	s.config.Args([]string{"-c", s.extraTestFileName})
	fin, _, err := s.config.findConfigFile()
	c.Assert(err, IsNil)
	c.Assert(reader.closed, Equals, true)
	assertFileContents(c, fin, "EXTRA_FILE")
}

func (s *FindConfigFileSuite) TestShortFlagOverrideFileName(c *C) {
	s.config.ConfigFile(s.defaultTestFileName)
	s.config.ConfigFileShortFlag('c')
	s.config.Args([]string{"-c", s.extraTestFileName})
	fin, _, err := s.config.findConfigFile()
	c.Assert(err, IsNil)
	assertFileContents(c, fin, "EXTRA_FILE")
}

func (s *FindConfigFileSuite) TestLongFlagOverrideFileReader(c *C) {
	reader := &closerStringReader{
		Reader: s.defaultTestFileReader,
		closed: false,
	}
	s.config.ConfigReader(reader)
	s.config.ConfigFileLongFlag("config-file")
	s.config.Args([]string{"--config-file", s.extraTestFileName})
	fin, _, err := s.config.findConfigFile()
	c.Assert(err, IsNil)
	c.Assert(reader.closed, Equals, true)
	assertFileContents(c, fin, "EXTRA_FILE")
}

func (s *FindConfigFileSuite) TestLongFlagOverrideFileName(c *C) {
	s.config.ConfigFile(s.defaultTestFileName)
	s.config.ConfigFileLongFlag("config-file")
	s.config.Args([]string{"--config-file", s.extraTestFileName})
	fin, _, err := s.config.findConfigFile()
	c.Assert(err, IsNil)
	assertFileContents(c, fin, "EXTRA_FILE")
}

func (s *FindConfigFileSuite) TestMultipleFlags(c *C) {
	reader := &closerStringReader{
		Reader: strings.NewReader("CLOSABLE_FILE_SIMULATION"),
		closed: false,
	}
	s.config.ConfigReader(reader)
	s.config.ConfigFileShortFlag('c')
	s.config.ConfigFileLongFlag("config-file")
	s.config.Args([]string{"-c", "file", "--config-file", "otherfile"})
	fin, _, err := s.config.findConfigFile()
	c.Assert(fin, IsNil)
	c.Assert(err, NotNil)
	c.Assert(reader.closed, Equals, true)
}

func (s *FindConfigFileSuite) TestFileNameRemoval(c *C) {
	s.config.ConfigFileLongFlag("config-file")
	s.config.Args([]string{"a", "b", "--config-file", s.defaultTestFileName})
	fin, args, err := s.config.findConfigFile()
	c.Assert(fin, NotNil)
	c.Assert(err, IsNil)
	c.Assert(args, DeepEquals, []string{"a", "b"})
}

func (s *FindConfigFileSuite) TestMissingFile(c *C) {
	s.config.ConfigFileShortFlag('c')
	s.config.Args([]string{"-c", "/missing/file/"})
	fin, _, err := s.config.findConfigFile()
	c.Assert(fin, IsNil)
	c.Assert(err, NotNil)

	s.config.ConfigFile("/missing/file/")
	s.config.Args([]string{"-c", "/missing/file"})
	fin, _, err = s.config.findConfigFile()
	c.Assert(fin, IsNil)
	c.Assert(err, NotNil)
}

func assertFileContents(c *C, fin io.Reader, expected string) {
	bytes, err := ioutil.ReadAll(fin)
	c.Assert(err, IsNil)
	c.Assert(string(bytes), Equals, expected)
	if closer, ok := fin.(io.Closer); ok {
		closer.Close()
	}
}
