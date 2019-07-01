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
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

type QueryParamsStruct struct {
	query     url.Values
	values    map[string]interface{}
	strictErr error
}

func QueryParams(request *http.Request) *QueryParamsStruct {
	return &QueryParamsStruct{query: request.URL.Query(), values: map[string]interface{}{}}
}

func (this *QueryParamsStruct) Define(param string, defaultValue interface{}) *QueryParamsStruct {
	this.values[param] = defaultValue
	if this.strictErr != nil {
		return this
	}
	if _, ok := this.query[param]; ok {
		value := this.query.Get(param)
		var err error
		switch defaultValue.(type) {
		case int:
			this.values[param], err = strconv.Atoi(value)
		case int64:
			this.values[param], err = strconv.ParseInt(value, 10, 64)
		case float64:
			this.values[param], err = strconv.ParseFloat(value, 64)
		case bool:
			this.values[param], err = strconv.ParseBool(value)
		case string:
			this.values[param] = value
		default:
			err = errors.New("")
		}
		if err != nil {
			t := reflect.Indirect(reflect.ValueOf(defaultValue)).Type()
			this.strictErr = errors.New("unable to interpret " + param + " as " + t.String())
		}
		return this
	}
	return this
}

func (this *QueryParamsStruct) Strict() (*QueryParamsStruct, error) {
	if this.strictErr != nil {
		return this, this.strictErr
	}
	for key, _ := range this.query {
		if _, ok := this.values[key]; !ok {
			return this, errors.New("unknown query-parameter: " + key)
		}
	}
	return this, nil
}

func (this *QueryParamsStruct) Get(param string) (result interface{}) {
	if value, ok := this.values[param]; !ok {
		panic("unknown parameter " + param)
	} else {
		return value
	}
}

func (this *QueryParamsStruct) GetAs(param string, target interface{}) *QueryParamsStruct {
	value := this.Get(param)
	targetvalue := reflect.Indirect(reflect.ValueOf(target))
	targetvalue.Set(reflect.ValueOf(value))
	return this
}
