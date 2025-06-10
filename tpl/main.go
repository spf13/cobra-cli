// Copyright © 2021 Steve Francia <spf@spf13.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tpl

import (
	_ "embed"
)

func MainTemplate() []byte {
	return []byte(`/*
{{ .Copyright }}
{{ if .Legal.Header }}{{ .Legal.Header }}{{ end }}
*/
package main

import "{{ .PkgName }}/cmd"

func main() {
	cmd.Execute()
}
`)
}

//go:embed root.local.tmpl
var rootLocalTemplate []byte

func RootTemplateLocal() []byte {
	return rootLocalTemplate
}

//go:embed root.global.tmpl
var rootGlobalTemplate []byte

func RootTemplateGlobal() []byte {
	return rootGlobalTemplate
}

//go:embed add-command.local.tmpl
var localAddCommandTemplate []byte

func AddCommandTemplateLocal() []byte {
	return localAddCommandTemplate
}

//go:embed add-command.global.tmpl
var globalAddCommandTemplate []byte

func AddCommandTemplateGlobal() []byte {
	return globalAddCommandTemplate
}
