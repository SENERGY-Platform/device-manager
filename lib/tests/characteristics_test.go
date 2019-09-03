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

func testCharacteristics(t *testing.T, conf config.Config) {
	createCharacteristic := model.Characteristic{
		Name: "char1",
		Type: model.Structure,
		SubCharacteristics: []model.Characteristic{{
			Name: "char2",
			Type: model.Float,
			SubCharacteristics: nil,
		}},
	}
	resp, err := helper.Jwtpost(adminjwt, "http://localhost:"+conf.ServerPort+"/concepts/4711/characteristics", createCharacteristic)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(resp.Status, resp.StatusCode, string(b))
	}

	characteristic := model.Characteristic{}
	err = json.NewDecoder(resp.Body).Decode(&characteristic)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("create: ids are set", func(t *testing.T) {
		characteristicWithIds(t, characteristic)
		t.Log(characteristic)
	})

	t.Run("create: Characteristic preserved structure", func(t *testing.T) {
		characteristicHasStructure(t, characteristic, createCharacteristic)
	})

	t.Run("create: Characteristic exists at semantic repo", func(t *testing.T) {
		checkCharacteristic(t, conf, characteristic.Id, createCharacteristic)
	})
	//
	//updateConcept := model.Concept{
	//	Id:   concept.Id,
	//	Name: "c2",
	//	CharacteristicIds: []string{"urn:infai:ses:characteristic:4712322a"},
	//}
	//resp, err = helper.Jwtput(adminjwt, "http://localhost:"+conf.ServerPort+"/concepts/"+url.PathEscape(concept.Id), updateConcept)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//defer resp.Body.Close()
	//
	//if resp.StatusCode != http.StatusOK {
	//	b, _ := ioutil.ReadAll(resp.Body)
	//	t.Fatal(resp.Status, resp.StatusCode, string(b))
	//}
	//
	//concept2 := model.Concept{}
	//err = json.NewDecoder(resp.Body).Decode(&concept2)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//t.Run("update: ids are set", func(t *testing.T) {
	//	conceptWithIds(t, concept2)
	//})
	//
	//t.Run("update: concept preserved structure", func(t *testing.T) {
	//	conceptHasStructure(t, concept2, updateConcept)
	//})
	//
	//t.Run("update: concept exists at semantic repo", func(t *testing.T) {
	//	checkConcept(t, conf, concept2.Id, updateConcept)
	//})
	//
	//resp, err = helper.Jwtdelete(adminjwt, "http://localhost:"+conf.ServerPort+"/concepts/"+url.PathEscape(concept.Id))
	//if err != nil {
	//	t.Fatal(err)
	//}
	//defer resp.Body.Close()
	//
	//t.Run("delete: concept removed at semantic repo", func(t *testing.T) {
	//	checkConceptDelete(t, conf, concept2.Id)
	//})
}

//func checkConceptDelete(t *testing.T, conf config.Config, id string) {
//	resp, err := helper.Jwtget(userjwt, conf.SemanticRepoUrl+"/concepts/"+url.PathEscape(id))
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusNotFound {
//		b, _ := ioutil.ReadAll(resp.Body)
//		t.Fatal(resp.Status, resp.StatusCode, string(b))
//	}
//}

func checkCharacteristic(t *testing.T, conf config.Config, id string, expected model.Characteristic) {
	resp, err := helper.Jwtget(userjwt, conf.SemanticRepoUrl+"/characteristics/"+url.PathEscape(id))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(resp.Status, resp.StatusCode, string(b))
	}

	characteristic := model.Characteristic{}
	err = json.NewDecoder(resp.Body).Decode(&characteristic)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("characteristic preserved structure", func(t *testing.T) {
		characteristicHasStructure(t, characteristic, expected)
	})
}

func characteristicHasStructure(t *testing.T, concept model.Characteristic, expected model.Characteristic) {
	expected = removeIdsCharacteristic(expected)
	concept = removeIdsCharacteristic(concept)
	if !reflect.DeepEqual(concept, expected) {
		t.Fatal(concept, expected)
	}
}

func characteristicWithIds(t *testing.T, char model.Characteristic) {
	if char.Id == "" {
		t.Fatal(char)
	}
	for i, characteristic := range char.SubCharacteristics {
		t.Run("characteristics subCharacteristics "+strconv.Itoa(i), func(t *testing.T) {
			subcharacteristicWithId(t, characteristic.Id)
		})
	}
}

func subcharacteristicWithId(t *testing.T, characteristicId string) {
	if characteristicId == "" {
		t.Fatal(characteristicId)
	}
}

func removeIdsCharacteristic(characteristic model.Characteristic) model.Characteristic {
	characteristic.Id = ""
	for i, ch := range characteristic.SubCharacteristics {
		characteristic.SubCharacteristics[i] = removeIdsCharacteristic(ch)
	}
	return characteristic
}
