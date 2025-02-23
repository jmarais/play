// Copyright 2015 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package main

import (
	encjson "encoding/json"
	encxml "encoding/xml"
	"fmt"

	"github.com/gopherjs/gopherjs/js"
	"github.com/jmarais/relapseviz"
	"github.com/katydid/katydid/parser"
	"github.com/katydid/katydid/parser/json"
	"github.com/katydid/katydid/parser/xml"
	"github.com/katydid/katydid/relapse/ast"
	"github.com/katydid/katydid/relapse/mem"
	relapseparser "github.com/katydid/katydid/relapse/parser"
)

func main() {
	js.Global.Set("gofunctions", map[string]interface{}{
		"RelapsePlayground": RelapsePlayground,
		"RelapseASTDraw":    RelapseASTDraw,
	})
}

func RelapsePlayground(mode string, katydidStr, input string) string {
	v, err := relapsePlayground(mode, katydidStr, input)
	if err != nil {
		return "Error: " + err.Error()
	}
	return fmt.Sprintf("%v", v)
}

func newParser(mode string, inputStr string) (parser.Interface, error) {
	switch mode {
	case "json":
		m := make(map[string]interface{})
		if err := encjson.Unmarshal([]byte(inputStr), &m); err != nil {
			return nil, err
		}
		p := json.NewJsonParser()
		err := p.Init([]byte(inputStr))
		if err != nil {
			return nil, err
		}
		return p, nil
	case "xml":
		var m interface{}
		if err := encxml.Unmarshal([]byte(inputStr), &m); err != nil {
			return nil, err
		}
		p := xml.NewXMLParser()
		err := p.Init([]byte(inputStr))
		if err != nil {
			return nil, err
		}
		return p, nil
	}
	return nil, fmt.Errorf("unknown mode %s", mode)
}

func relapsePlayground(mode, katydidStr, inputStr string) (match bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	var g *ast.Grammar
	g, err = relapseparser.ParseGrammar(katydidStr)
	if err != nil {
		return
	}
	var p parser.Interface
	p, err = newParser(mode, inputStr)
	if err != nil {
		return
	}
	var m *mem.Mem
	m, err = mem.New(g)
	if err != nil {
		return
	}
	match, err = m.Validate(p)
	return
}

func RelapseASTDraw(katydidStr string, full bool) string {
	v, err := relapseASTDraw(katydidStr, full)
	if err != nil {
		return "Error: " + err.Error()
	}
	return v
}

func relapseASTDraw(katydidStr string, full bool) (res string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	var g *ast.Grammar
	g, err = relapseparser.ParseGrammar(katydidStr)
	if err != nil {
		return
	}
	graph := relapseviz.TranslateGrammar(g, full)
	res = graph.String()
	return
}
