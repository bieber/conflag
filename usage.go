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
	"strings"
)

const minUsageWidth = 10
const keyToDescriptionSpacing = 4
const flagIndentDepth = 2

func (c *concreteConfig) Usage(width uint) string {
	if width < minUsageWidth {
		width = minUsageWidth
	}

	sections := []string{}

	if c.name != "" {
		sections = append(sections, restrictWidthByWords(c.name, int(width)))
	}

	if c.description != "" {
		sections = append(
			sections,
			restrictWidthByWords(c.description, int(width)),
		)
	}

	keys, keysWidth := formatFieldKeys(c.fields, c.fieldKeysInOrder)

	descriptions := []string{}
	for _, k := range c.fieldKeysInOrder {
		descriptions = append(descriptions, c.fields[k].description)
	}

	var combinedArgInfo []string
	if int(width)-keysWidth-keyToDescriptionSpacing > minUsageWidth {
		combinedArgInfo = formatArgsSideBySide(
			keys,
			descriptions,
			keysWidth,
			int(width)-keysWidth-keyToDescriptionSpacing,
		)
	} else {
		combinedArgInfo = formatArgsVertical(keys, descriptions, int(width))
	}
	sections = append(sections, combinedArgInfo...)

	return strings.Join(sections, "\n\n")
}

func formatFieldKeys(
	fields map[string]*concreteField,
	fieldKeysInOrder []string,
) (sections []string, maxWidth int) {
	sections = []string{}
	maxWidth = 0
	indentation := strings.Repeat(" ", flagIndentDepth)

	for _, key := range fieldKeysInOrder {
		field := fields[key]
		lines := []string{}

		fieldComponents := []string{}
		if field.shortFlag != 0 {
			fieldComponents = append(
				fieldComponents,
				string([]rune{'-', field.shortFlag}),
			)
		}
		if field.longFlag != "" {
			fieldComponents = append(fieldComponents, "--"+field.longFlag)
		}

		if len(fieldComponents) != 0 {
			lines = append(
				lines,
				indentation+strings.Join(fieldComponents, ", "),
			)
		}

		if field.fileKey != "" {
			fileLine := indentation
			if field.fileCategory != "" {
				fileLine += field.fileCategory + "."
			}
			fileLine += field.fileKey
			lines = append(lines, fileLine)
		}

		for _, line := range lines {
			if strlen(line) > maxWidth {
				maxWidth = strlen(line)
			}
		}
		sections = append(sections, strings.Join(lines, "\n"))
	}

	return
}

func formatArgsSideBySide(
	keys []string,
	descriptions []string,
	keysWidth int,
	descriptionWidth int,
) []string {
	results := []string{}
	for i := range keys {
		results = append(
			results,
			formatArgSideBySide(
				keys[i],
				descriptions[i],
				keysWidth,
				descriptionWidth,
			),
		)
	}

	return results
}

func formatArgSideBySide(
	keys string,
	description string,
	keysWidth int,
	descriptionWidth int,
) string {
	lines := []string{}
	keysLines := strings.Split(keys, "\n")
	descriptionLines := strings.Split(
		restrictWidthByWords(description, descriptionWidth),
		"\n",
	)

	for i := 0; i < len(keysLines) || i < len(descriptionLines); i++ {
		currentLine := ""
		if i < len(keysLines) {
			currentLine += keysLines[i]
		}
		if i < len(descriptionLines) {
			padding := keysWidth - strlen(currentLine) + keyToDescriptionSpacing
			currentLine += strings.Repeat(" ", padding)
			currentLine += descriptionLines[i]
		}

		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}

func formatArgsVertical(
	keys []string,
	descriptions []string,
	maxWidth int,
) []string {
	sections := []string{}
	for i := range keys {
		lines := []string{}
		if keys[i] != "" {
			keysLines := strings.Split(keys[i], "\n")
			for _, l := range keysLines {
				lines = append(lines, restrictWidthByWords(l, maxWidth))
			}
		}
		if descriptions[i] != "" {
			if len(lines) != 0 {
				lines = append(lines, "")
			}
			lines = append(
				lines,
				restrictWidthByWords(descriptions[i], maxWidth),
			)
		}
		sections = append(sections, strings.Join(lines, "\n"))
	}
	return sections
}

func restrictWidthByWords(s string, width int) string {
	words := strings.Split(s, " ")
	rows := []string{}

	currentRow := []string{}
	currentRowWordLength := 0
	for _, word := range words {
		if strlen(word)+currentRowWordLength+len(currentRow) > width {
			rows = append(rows, strings.Join(currentRow, " "))
			currentRow = []string{}
			currentRowWordLength = 0
		}

		currentRow = append(currentRow, word)
		currentRowWordLength += strlen(word)
	}

	rows = append(rows, strings.Join(currentRow, " "))

	return strings.Join(rows, "\n")
}

func strlen(s string) int {
	return len([]rune(s))
}
