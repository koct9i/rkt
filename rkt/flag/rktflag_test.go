// Copyright 2015 The rkt Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package flag

import (
	"strings"
	"testing"
)

var options = []string{"zero", "one", "two"}

func TestOptionList(t *testing.T) {
	tests := []struct {
		opts string
		ex   string
		err  bool
	}{
		{
			opts: "zero,two",
			ex:   "zero,two",
			err:  false,
		},
		{ // Duplicate test
			opts: "one,two,two",
			ex:   "",
			err:  true,
		},
		{ // Not permissible test
			opts: "one,two,three",
			ex:   "",
			err:  true,
		},
		{ // Empty string
			opts: "",
			ex:   "",
			err:  false,
		},
	}

	for i, tt := range tests {
		// test NewOptionsList
		if _, err := NewOptionList(options, tt.opts); (err != nil) != tt.err {
			t.Errorf("test %d: unexpected error in NewOptionList: %v", i, err)
		}

		// test OptionList.Set()
		ol, err := NewOptionList(options, strings.Join(options, ","))
		if err != nil {
			t.Errorf("test %d: unexpected error preparing test: %v", i, err)
		}

		if err := ol.Set(tt.opts); (err != nil) != tt.err {
			t.Errorf("test %d: could not parse options as expected: %v", i, err)
		}
		if tt.ex != "" && tt.ex != ol.String() {
			t.Errorf("test %d: resulting options not as expected: %s != %s",
				i, tt.ex, ol.String())
		}
	}
}

var bfMap = map[string]int{
	options[0]: 0,
	options[1]: 1,
	options[2]: 1 << 1,
}

func TestBitFlags(t *testing.T) {
	tests := []struct {
		opts     string
		ex       int
		parseErr bool
		logicErr bool
	}{
		{
			opts: "one,two",
			ex:   3,
		},
		{ // Duplicate test
			opts:     "zero,two,two",
			ex:       -1,
			parseErr: true,
		},
		{ // Not included test
			opts:     "zero,two,three",
			ex:       -1,
			parseErr: true,
		},
		{ // Test 10 in 11
			opts: "one,two",
			ex:   1,
		},
		{ // Test 11 not in 01
			opts:     "one",
			ex:       3,
			logicErr: true,
		},
	}

	for i, tt := range tests {
		// test NewBitFlags
		if _, err := NewBitFlags(options, tt.opts, bfMap); (err != nil) != tt.parseErr {
			t.Errorf("test %d: unexpected error in NewBitFlags: %v", i, err)
		}

		bf, err := NewBitFlags(options, strings.Join(options, ","), bfMap)
		if err != nil {
			t.Errorf("test %d: unexpected error preparing test: %v", i, err)
		}

		// test BitFlags.Set()
		if err := bf.Set(tt.opts); (err != nil) != tt.parseErr {
			t.Errorf("test %d: Could not parse options as expected: %v", i, err)
		}
		if tt.ex >= 0 && bf.hasFlag(tt.ex) == tt.logicErr {
			t.Errorf("test %d: Result was unexpected: %d != %d",
				i, tt.ex, bf.flags)
		}
	}
}