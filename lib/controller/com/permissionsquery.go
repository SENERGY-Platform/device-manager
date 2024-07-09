/*
 * Copyright 2022 InfAI (CC SES)
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

package com

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/permission-search/lib/client"
	"log"
	"net/http"
	"runtime/debug"
)

func (this *Com) QueryPermissionsSearch(token string, query QueryMessage, result interface{}) (err error, code int) {
	requestBody := new(bytes.Buffer)
	err = json.NewEncoder(requestBody).Encode(query)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	req, err := http.NewRequest("POST", this.config.PermissionsUrl+"/v3/query", requestBody)
	if err != nil {
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		err = errors.New(buf.String())
		log.Println("ERROR: ", resp.StatusCode, err)
		debug.PrintStack()
		return err, resp.StatusCode
	}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

type QueryMessage = client.QueryMessage

type QueryFind = client.QueryFind

type QueryListIds = client.QueryListIds

type QueryCheckIds = client.QueryCheckIds

type QueryListCommons = client.QueryListCommons

type ListAfter = client.ListAfter

type QueryOperationType = client.QueryOperationType

const (
	QueryEqualOperation             = client.QueryEqualOperation
	QueryUnequalOperation           = client.QueryUnequalOperation
	QueryAnyValueInFeatureOperation = client.QueryAnyValueInFeatureOperation
)

type ConditionConfig = client.ConditionConfig

type Selection = client.Selection
