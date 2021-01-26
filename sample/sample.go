/*
 * Options
 *
 * The MIT License (MIT)
 *
 * Copyright (c) 2020 Sebastian Borchers <sebastian[at]desertbit.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"fmt"

	"github.com/desertbit/options"
	yaml "gopkg.in/yaml.v3"
)

type NestedOptions struct {
	F float32 `yaml:"f,omitempty"`
	I int     `yaml:"i,omitempty"`
}

type NestedInlineOptions struct {
	SS string `yaml:"inline-ss,omitempty"`
}

type Options struct {
	S   string         `yaml:"s,omitempty"`
	I   int            `yaml:"i,omitempty"`
	IS  []int          `yaml:"is,omitempty"`
	MIS map[int]string `yaml:"mis,omitempty"`

	NO NestedOptions `yaml:"no,omitempty"`

	NestedInlineOptions `yaml:",inline,omitempty"`

	PointerTest *NestedOptions `yaml:"pointer-test,omitempty"`
}

func DefaultOptions() Options {
	return Options{
		S:                   "hello",
		I:                   8,
		IS:                  []int{1, 2, 3},
		MIS:                 map[int]string{1: "first", 2: "second", 3: "third"},
		NO:                  NestedOptions{F: 0.8, I: 1},
		NestedInlineOptions: NestedInlineOptions{SS: "test"},
	}
}

func main() {
	o := DefaultOptions()
	// Every field that is changed here will be represented in the YAML.
	o.S = "something-different"
	o.NO.F = 1.5
	o.PointerTest = &NestedOptions{}

	err := options.StripDefaults(&o, DefaultOptions())
	if err != nil {
		panic(err)
	}

	data, err := yaml.Marshal(o)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", string(data))

	// Now set every field with a zero value to its default counterpart.
	// Every field that is set here will not have its default value.
	o = Options{
		I:           5,
		IS:          []int{-1},
		NO:          NestedOptions{F: 0.666},
		PointerTest: &NestedOptions{F: 3.66},
	}

	err = options.SetDefaults(&o, DefaultOptions())
	if err != nil {
		panic(err)
	}

	data, err = yaml.Marshal(o)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", string(data))
}
