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
	"fmt"
	"github.com/openconfig/goyang/pkg/yang"
	"github.com/openconfig/ygot/util"
	"strings"
)

// resolveRootName resolves the name of the fakeroot by taking configuration
// and the default values, along with a boolean indicating whether the fake
// root is to be generated. It returns an empty string if the root is not
// to be generated.
func resolveRootName(name, defName string, generateRoot bool) string {
	if !generateRoot {
		return ""
	}

	if name == "" {
		return defName
	}

	return name
}

// resolveTypeArgs is a structure used as an input argument to the yangTypeToGoType
// function which allows extra context to be handed on. This provides the ability
// to use not only the YangType but also the yang.Entry that the type was part of
// to resolve the possible type name.
type resolveTypeArgs struct {
	// yangType is a pointer to the yang.YangType that is to be mapped.
	yangType *yang.YangType
	// contextEntry is an optional yang.Entry which is supplied where a
	// type requires knowledge of the leaf that it is used within to be
	// mapped. For example, where a leaf is defined to have a type of a
	// user-defined type (typedef) that in turn has enumerated values - the
	// context of the yang.Entry is required such that the leaf's context
	// can be established.
	contextEntry *yang.Entry
}

// given yang.Entry of type yang.Yleafref traverse the tree
// to find the entry in the appropriate place
func findLeafRef(entry *yang.Entry) (*yang.Entry, error) {

	if entry.Type.Kind != yang.Yleafref {
		return nil, fmt.Errorf("entry %s is not a leafref", entry.Name)
	}

	path := entry.Type.Path

	if len(strings.Split(path, "/")) < 2 {
		return nil, fmt.Errorf("key %s had an invalid path %s", entry.Name, entry.Path())
	}

	var dp []string // downward path
	var curEntry = entry

	// if the path is absolute (eg: /t1:cont1a/t1:list2a/t1:name) go back to the root and then descend
	if path[0:1] == "/" {
		dp = strings.Split(path[1:], "/")

		// go back till the root
		for {
			if curEntry.Parent != nil {
				curEntry = curEntry.Parent
			} else {
				break
			}
		}
	} else {
		// identify how many levels we have to go up the tree
		lb := strings.Count(path, "../")

		// identify the path to take once we have gone up the tree
		dp = strings.Split(strings.ReplaceAll(path, "../", ""), "/")

		// this is the entry we are moving to

		for i := 0; i < lb; i++ {
			// we're going up the tree till it's needed

			if curEntry.Parent == nil {
				return nil, fmt.Errorf("entry %s does not have a parent", curEntry.Name)
			}

			curEntry = curEntry.Parent
		}
	}

	var downwardPath = []string{}
	// remove the prefix from the pieces
	for _, p := range dp {
		downwardPath = append(downwardPath, util.StripModulePrefix(p))
	}

	for _, k := range downwardPath {
		// and then descending to the leafref path
		_curEntry, ok := curEntry.Dir[k]
		if !ok {
			return nil, fmt.Errorf("entry %s not found in %s", k, curEntry.Name)
		}
		curEntry = _curEntry
	}

	return curEntry, nil
}
