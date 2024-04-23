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

package config

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	LogLevel   string `json:"log_level"` //DEBUG | CALL | NONE
	ServerPort string `json:"server_port"`
	Debug      bool   `json:"debug"`

	EditForward string `json:"edit_forward"`

	PermissionsUrl      string `json:"permissions_url"`
	DeviceTopic         string `json:"device_topic"`
	DeviceTypeTopic     string `json:"device_type_topic"`
	DeviceGroupTopic    string `json:"device_group_topic"`
	ProtocolTopic       string `json:"protocol_topic"`
	HubTopic            string `json:"hub_topic"`
	ConceptTopic        string `json:"concept_topic"`
	CharacteristicTopic string `json:"characteristic_topic"`
	AspectTopic         string `json:"aspect_topic"`
	FunctionTopic       string `json:"function_topic"`
	DeviceClassTopic    string `json:"device_class_topic"`
	DeviceRepoUrl       string `json:"device_repo_url"`
	KafkaUrl            string `json:"kafka_url"`
	LocationTopic       string `json:"location_topic"`
	UserTopic           string `json:"user_topic"`
	GroupId             string `json:"group_id"`
	HttpClientTimeout   string `json:"http_client_timeout"`

	DisableValidation bool   `json:"disable_validation"`
	ConverterUrl      string `json:"converter_url"` //to validate concept conversions, may be empty to disable concept conversion validation
	HandleDoneWait    bool   `json:"handle_done_wait"`
}

// loads config from json in location and used environment variables (e.g KafkaUrl --> KAFKA_URL)
func Load(location string) (config Config, err error) {
	file, err := os.Open(location)
	if err != nil {
		return config, err
	}
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return config, err
	}
	handleEnvironmentVars(&config)
	setDefaultHttpClient(config)
	return config, nil
}

var camel = regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)")

func fieldNameToEnvName(s string) string {
	var a []string
	for _, sub := range camel.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}
		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}
	return strings.ToUpper(strings.Join(a, "_"))
}

// preparations for docker
func handleEnvironmentVars(config *Config) {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	configType := configValue.Type()
	for index := 0; index < configType.NumField(); index++ {
		fieldName := configType.Field(index).Name
		envName := fieldNameToEnvName(fieldName)
		envValue := os.Getenv(envName)
		if envValue != "" {
			fmt.Println("use environment variable: ", envName, " = ", envValue)
			if configValue.FieldByName(fieldName).Kind() == reflect.Int64 {
				i, _ := strconv.ParseInt(envValue, 10, 64)
				configValue.FieldByName(fieldName).SetInt(i)
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.String {
				configValue.FieldByName(fieldName).SetString(envValue)
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.Bool {
				b, _ := strconv.ParseBool(envValue)
				configValue.FieldByName(fieldName).SetBool(b)
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.Float64 {
				f, _ := strconv.ParseFloat(envValue, 64)
				configValue.FieldByName(fieldName).SetFloat(f)
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.Slice {
				val := []string{}
				for _, element := range strings.Split(envValue, ",") {
					val = append(val, strings.TrimSpace(element))
				}
				configValue.FieldByName(fieldName).Set(reflect.ValueOf(val))
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.Map {
				value := map[string]string{}
				for _, element := range strings.Split(envValue, ",") {
					keyVal := strings.Split(element, ":")
					key := strings.TrimSpace(keyVal[0])
					val := strings.TrimSpace(keyVal[1])
					value[key] = val
				}
				configValue.FieldByName(fieldName).Set(reflect.ValueOf(value))
			}
		}
	}
}

func setDefaultHttpClient(config Config) {
	var err error
	http.DefaultClient.Timeout, err = time.ParseDuration(config.HttpClientTimeout)
	if err != nil {
		log.Println("WARNING: invalid http timeout --> no timeouts\n", err)
	}
}
