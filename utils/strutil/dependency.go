// Copyright (C) 2024 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 依赖分析

package strutil

// DS_M DependencySorter 依赖排序入参
type DS_M struct {
	Name         string
	Dependencies []string
}

// DependencySorter 依赖排序(计算依赖排序)
func DependencySorter(modules ...DS_M) []string {
	dps := new(dependencySorter)
	dps.modules = modules
	return dps.dependencyOrder()
}

// dependencySorter 依赖排序器
type dependencySorter struct {
	modules     []DS_M
	moduleIndex map[string]DS_M
	visited     map[string]bool
}

func (sd *dependencySorter) init() {
	sd.visited = make(map[string]bool)
	sd.moduleIndex = make(map[string]DS_M)

	if lm := len(sd.modules); lm > 0 {
		for i := 0; i < lm; i++ {
			sd.moduleIndex[sd.modules[i].Name] = sd.modules[i]
		}
	}
}

func (sd *dependencySorter) topologicalSort(module DS_M, stack *[]string) {
	sd.visited[module.Name] = true
	for _, dependency := range module.Dependencies {
		if !sd.visited[dependency] {
			if val, ok := sd.moduleIndex[dependency]; ok {
				sd.topologicalSort(val, stack)
			}
		}
	}

	*stack = append(*stack, module.Name)
}

func (sd *dependencySorter) dependencyOrder() (stack []string) {
	if len(sd.modules) == 0 {
		return
	} else {
		sd.init()
	}

	for _, module := range sd.modules {
		if !sd.visited[module.Name] {
			sd.topologicalSort(module, &stack)
		}
	}
	return stack
}
