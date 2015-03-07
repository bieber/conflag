/*
 * Copyright (c) 2015, Robert Bieber
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

type UsageSuite struct {
	config Config
}

func (s *UsageSuite) SetUpTest(c *C) {
	configStruct := &struct {
		Verbose              bool
		Source               string
		MaxConcurrentThreads int
		Net                  struct {
			HostName string
			Port     int
		}
	}{}

	programDescription :=
		"This is a program with a pretty long description.  In fact, it " +
			"should extend pretty significantly beyond 80 columns, such that " +
			"it will need to be broken down onto several lines."

	config, err := New(configStruct)
	config.
		ProgramName("Test Program Name Line").
		ProgramDescription(programDescription)

	config.Field("Verbose").
		ShortFlag('v').
		LongFlag("verbose").
		FileKey("verbose").
		Description(
		"Verbosity flag.  Set to display debug information of a particularly " +
			"verbose nature.  Like this description, for instance.",
	)

	config.Field("Source").
		Description("The directory from which to read the things.")

	config.Field("MaxConcurrentThreads").
		ShortFlag('m').
		LongFlag("").
		FileKey("").
		Description("The maximum number of threads.")

	config.Field("Net.HostName").
		ShortFlag('h').
		LongFlag("host-name").
		FileCategory("").
		Description("Hostname to serve on.")

	config.Field("Net.Port").
		LongFlag("").
		Description("Port to serve on.")

	s.config = config
	c.Assert(err, IsNil)
}

func TestUsage(t *testing.T) {
	Suite(&UsageSuite{})
	TestingT(t)
}

func (s *UsageSuite) TestFormatting(c *C) {
	sixtyColsOutput := "" +
		"Test Program Name Line\n" +
		"\n" +
		"This is a program with a pretty long description.  In fact,\n" +
		"it should extend pretty significantly beyond 80 columns,\n" +
		"such that it will need to be broken down onto several lines.\n" +
		"\n" +
		"  -v, --verbose      Verbosity flag.  Set to display debug\n" +
		"  verbose            information of a particularly verbose\n" +
		"                     nature.  Like this description, for\n" +
		"                     instance.\n" +
		"\n" +
		"  --source           The directory from which to read the\n" +
		"  source             things.\n" +
		"\n" +
		"  -m                 The maximum number of threads.\n" +
		"\n" +
		"  -h, --host-name    Hostname to serve on.\n" +
		"  host_name\n" +
		"\n" +
		"  net.port           Port to serve on."
	c.Assert(s.config.Usage(60), Equals, sixtyColsOutput)

	fortyColsOutput := "" +
		"Test Program Name Line\n" +
		"\n" +
		"This is a program with a pretty long\n" +
		"description.  In fact, it should extend\n" +
		"pretty significantly beyond 80 columns,\n" +
		"such that it will need to be broken down\n" +
		"onto several lines.\n" +
		"\n" +
		"  -v, --verbose      Verbosity flag. \n" +
		"  verbose            Set to display\n" +
		"                     debug information\n" +
		"                     of a particularly\n" +
		"                     verbose nature. \n" +
		"                     Like this\n" +
		"                     description, for\n" +
		"                     instance.\n" +
		"\n" +
		"  --source           The directory from\n" +
		"  source             which to read the\n" +
		"                     things.\n" +
		"\n" +
		"  -m                 The maximum number\n" +
		"                     of threads.\n" +
		"\n" +
		"  -h, --host-name    Hostname to serve\n" +
		"  host_name          on.\n" +
		"\n" +
		"  net.port           Port to serve on."

	c.Assert(s.config.Usage(40), Equals, fortyColsOutput)

	twentyColsOutput := "" +
		"Test Program Name\n" +
		"Line\n" +
		"\n" +
		"This is a program\n" +
		"with a pretty long\n" +
		"description.  In\n" +
		"fact, it should\n" +
		"extend pretty\n" +
		"significantly beyond\n" +
		"80 columns, such\n" +
		"that it will need to\n" +
		"be broken down onto\n" +
		"several lines.\n" +
		"\n" +
		"  -v, --verbose\n" +
		"  verbose\n" +
		"\n" +
		"Verbosity flag.  Set\n" +
		"to display debug\n" +
		"information of a\n" +
		"particularly verbose\n" +
		"nature.  Like this\n" +
		"description, for\n" +
		"instance.\n" +
		"\n" +
		"  --source\n" +
		"  source\n" +
		"\n" +
		"The directory from\n" +
		"which to read the\n" +
		"things.\n" +
		"\n" +
		"  -m\n" +
		"\n" +
		"The maximum number\n" +
		"of threads.\n" +
		"\n" +
		"  -h, --host-name\n" +
		"  host_name\n" +
		"\n" +
		"Hostname to serve\n" +
		"on.\n" +
		"\n" +
		"  net.port\n" +
		"\n" +
		"Port to serve on."
	c.Assert(s.config.Usage(20), Equals, twentyColsOutput)

	c.Assert(s.config.Usage(1), Equals, s.config.Usage(minUsageWidth))
}
