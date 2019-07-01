/*
 * Copyright 2019 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"net/http"
	"net/url"
	"testing"
)

func TestQueryParams(t *testing.T) {
	urlObj, err := url.Parse("http://foo.bar?str=foo&int=42&float=4.2&bool=true")
	if err != nil {
		t.Fatal(err)
	}
	request := &http.Request{URL: urlObj}

	params, err := QueryParams(request).
		Define("str", "default").
		Define("int", 13).
		Define("float", 1.3).
		Define("bool", true).
		Strict()

	if err != nil {
		t.Fatal(err)
	}

	var str string
	var b bool
	var f float64
	var i int
	params.
		GetAs("str", &str).
		GetAs("bool", &b).
		GetAs("float", &f).
		GetAs("int", &i)

	if str != "foo" {
		t.Fatal(str)
	}
	if i != 42 {
		t.Fatal(i)
	}
	if f != 4.2 {
		t.Fatal()
	}
	if b != true {
		t.Fatal(b)
	}
}

func TestQueryParamsGet(t *testing.T) {
	urlObj, err := url.Parse("http://foo.bar?str=foo&int=42&float=4.2&bool=true")
	if err != nil {
		t.Fatal(err)
	}
	request := &http.Request{URL: urlObj}

	params, err := QueryParams(request).
		Define("str", "default").
		Define("int", 13).
		Define("float", 1.3).
		Define("bool", true).
		Strict()

	if err != nil {
		t.Fatal(err)
	}

	if params.Get("str").(string) != "foo" {
		t.Fatal()
	}
	if params.Get("int").(int) != 42 {
		t.Fatal()
	}
	if params.Get("float").(float64) != 4.2 {
		t.Fatal()
	}
	if params.Get("bool").(bool) != true {
		t.Fatal()
	}
}

func TestQueryParamsDefault(t *testing.T) {
	urlObj, err := url.Parse("http://foo.bar")
	if err != nil {
		t.Fatal(err)
	}
	request := &http.Request{URL: urlObj}

	params, err := QueryParams(request).
		Define("str", "foo").
		Define("int", 42).
		Define("float", 4.2).
		Define("bool", true).
		Strict()

	if err != nil {
		t.Fatal(err)
	}

	var str string
	var b bool
	var f float64
	var i int
	params.
		GetAs("str", &str).
		GetAs("bool", &b).
		GetAs("float", &f).
		GetAs("int", &i)

	if str != "foo" {
		t.Fatal(str)
	}
	if i != 42 {
		t.Fatal(i)
	}
	if f != 4.2 {
		t.Fatal()
	}
	if b != true {
		t.Fatal(b)
	}
}

func TestQueryParamsStrict(t *testing.T) {
	urlObj, err := url.Parse("http://foo.bar?string=foo")
	if err != nil {
		t.Fatal(err)
	}
	request := &http.Request{URL: urlObj}

	_, err = QueryParams(request).
		Define("str", "default").
		Strict()

	if err == nil {
		t.Fatal("missing error")
	}
}

func TestQueryParamsErrors(t *testing.T) {
	urlObj, err := url.Parse("http://foo.bar")
	if err != nil {
		t.Fatal(err)
	}
	request := &http.Request{URL: urlObj}

	params, err := QueryParams(request).
		Define("str", 13).
		Strict()

	if err != nil {
		t.Fatal(err)
	}

	var str string
	defer func() {
		if recover() != nil {
			return
		} else {
			t.Fatal(str)
		}
	}()
	params.GetAs("str", &str)
}

func TestQueryParamsUserErr(t *testing.T) {
	urlObj, err := url.Parse("http://foo.bar?int=bar")
	if err != nil {
		t.Fatal(err)
	}
	request := &http.Request{URL: urlObj}

	_, err = QueryParams(request).
		Define("int", 13).
		Strict()

	if err == nil {
		t.Fatal("missing error")
	}
}
