// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ygen

import (
	"github.com/openconfig/goyang/pkg/yang"
	"reflect"
	"testing"
)

func TestResolveRootName(t *testing.T) {
	tests := []struct {
		name           string
		inName         string
		inDefName      string
		inGenerateRoot bool
		want           string
	}{{
		name:           "generate root false",
		inGenerateRoot: false,
	}, {
		name:           "name specified",
		inName:         "value",
		inDefName:      "invalid",
		inGenerateRoot: true,
		want:           "value",
	}, {
		name:           "name not specified",
		inDefName:      "default",
		inGenerateRoot: true,
		want:           "default",
	}}

	for _, tt := range tests {
		if got := resolveRootName(tt.inName, tt.inDefName, tt.inGenerateRoot); got != tt.want {
			t.Errorf("%s: resolveRootName(%s, %s, %v): did not get expected result, got: %s, want: %s", tt.name, tt.inName, tt.inDefName, tt.inGenerateRoot, got, tt.want)
		}
	}
}

func Test_findLeafRef(t *testing.T) {
	type args struct {
		entry *yang.Entry
	}
	tests := []struct {
		name    string
		in      *yang.Entry
		want    *yang.Entry
		wantErr bool
	}{
		{
			name: "not-leafref",
			in: &yang.Entry{
				Name: "number",
				Type: &yang.YangType{
					Kind: yang.Yuint16,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "basic-list-key",
			in: &yang.Entry{
				Name: "keyleafref",
				Type: &yang.YangType{
					Kind: yang.Yleafref,
					Path: "../pfx:config/pfx:keyleafref",
				},
				Parent: &yang.Entry{
					Name:     "list",
					ListAttr: &yang.ListAttr{},
					Key:      "keyleafref",
					Dir: map[string]*yang.Entry{
						"config": {
							Name: "config",
							Dir: map[string]*yang.Entry{
								"keyleafref": {
									Name: "keyleafref",
									Type: &yang.YangType{Kind: yang.Ystring},
								},
							},
						},
					},
				},
			},
			want: &yang.Entry{
				Name: "keyleafref",
				Type: &yang.YangType{Kind: yang.Ystring},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findLeafRef(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("findLeafRef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findLeafRef() got = %v, want %v", got, tt.want)
			}
		})
	}
}
