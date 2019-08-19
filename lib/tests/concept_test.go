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

package tests

import (
	"encoding/json"
	"github.com/SENERGY-Platform/device-manager/lib/config"
	"github.com/SENERGY-Platform/device-manager/lib/model"
	"github.com/SENERGY-Platform/device-manager/lib/tests/helper"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

func testConcepts(t *testing.T, conf config.Config) {
	createConcept := model.Concept{
		Name: "c1",
		Characteristics: []model.Characteristic{{
			Name: "ch1",
			SubCharacteristics: []model.Characteristic{{
				Name: "ch1b",
			}},
		}},
	}
	resp, err := helper.Jwtpost(adminjwt, "http://localhost:"+conf.ServerPort+"/concepts", createConcept)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(resp.Status, resp.StatusCode, string(b))
	}

	concept := model.Concept{}
	err = json.NewDecoder(resp.Body).Decode(&concept)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("ids are set", func(t *testing.T) {
		conceptWithIds(t, concept)
	})

	t.Run("concept preserved structure", func(t *testing.T) {
		conceptHasStructure(t, concept, createConcept)
	})

	resp, err = helper.Jwtget(userjwt, conf.SemanticRepoUrl+"/concepts/"+url.PathEscape(concept.Id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(resp.Status, resp.StatusCode, string(b), conf.SemanticRepoUrl+"/concepts/"+url.PathEscape(concept.Id))
	}

	conceptGet := model.Concept{}
	err = json.NewDecoder(resp.Body).Decode(&conceptGet)
	if err != nil {
		t.Fatal(err)
	}

	expected := createConcept

	t.Run("concept preserved structure", func(t *testing.T) {
		conceptHasStructure(t, conceptGet, expected)
	})

}

func conceptHasStructure(t *testing.T, concept model.Concept, expected model.Concept) {
	concept = removeIdsFromConcept(concept)
	if !reflect.DeepEqual(concept, expected) {
		t.Fatal(concept, expected)
	}
}

func conceptWithIds(t *testing.T, concept model.Concept) {
	if concept.Id == "" {
		t.Fatal(concept)
	}
	for i, characteristic := range concept.Characteristics {
		t.Run("concept characteristics "+strconv.Itoa(i), func(t *testing.T) {
			characteristicWithId(t, characteristic)
		})
	}
}

func characteristicWithId(t *testing.T, characteristic model.Characteristic) {
	if characteristic.Id == "" {
		t.Fatal(characteristic)
	}
	for _, sub := range characteristic.SubCharacteristics {
		characteristicWithId(t, sub)
	}
}

func removeIdsFromConcept(concept model.Concept) model.Concept {
	concept.Id = ""
	for i, ch := range concept.Characteristics {
		concept.Characteristics[i] = removeIdsCharacteristic(ch)
	}
	return concept
}

func removeIdsCharacteristic(characteristic model.Characteristic) model.Characteristic {
	characteristic.Id = ""
	for i, ch := range characteristic.SubCharacteristics {
		characteristic.SubCharacteristics[i] = removeIdsCharacteristic(ch)
	}
	return characteristic
}
