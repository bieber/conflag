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

/*

Package conflag implements simple program configuration parsing.  It
allows you to create a set of options for your program in a struct
definition, and then load them from any combination of defaults, a
configuration file, and command-line flags.  As an example, consider
the following struct definition:

	type ServerConfig struct {
		Port int
		Path string
	}

This defines some simple information that you may want to use as a
configuration for a web server.  Start by creating an instance of
your configuration struct type, which you may also set default values
on.  For example, setting the default port to 80:

	settings := ServerConfig{Port: 80}

You can then create a new configuration parser by passing a pointer
to your instance to the New() function.

	configParser, err := conflag.New(&settings)

If there's anything wrong with your struct (see the New()
documentation for restrictions on nesting and field types) the
returned error message will notify you.  Also see the New()
documentation for the parser's default settings.  To modify the
default parser settings, use the Config object's Field() method to
access individual fields and modify their settings.  For instance, to
make the path field required and give it a short command-line flag of
-p:

	configParser.Field("Path").Required().ShortFlag('p')
*/
package conflag
