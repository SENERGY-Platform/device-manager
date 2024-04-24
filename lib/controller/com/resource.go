package com

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/device-manager/lib/auth"
	"github.com/SENERGY-Platform/device-manager/lib/config"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func getResourceFromService(token auth.Token, endpoint string, id string, result interface{}) (err error, code int) {
	req, err := http.NewRequest("GET", endpoint+"/"+url.PathEscape(id), nil)
	if err != nil {
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}
	req.Header.Set("Authorization", token.Token)
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
		log.Println("WARNING: unable to get resource", endpoint, id, resp.StatusCode, err)
		//debug.PrintStack()
		return err, resp.StatusCode
	}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

func validateResources(token auth.Token, config config.Config, endpoints []string, resource interface{}) (err error, code int) {
	if config.DisableValidation {
		return nil, http.StatusOK
	}
	for _, endpoint := range endpoints {
		err, code = validateResource(token, config, "PUT", endpoint, resource)
		if err != nil {
			return err, code
		}
	}
	return nil, http.StatusOK
}

func validateResource(token auth.Token, config config.Config, method string, endpoint string, resource interface{}) (err error, code int) {
	if config.DisableValidation {
		return nil, http.StatusOK
	}
	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(resource)
	if err != nil {
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}
	req, err := http.NewRequest(method, endpoint, b)
	if err != nil {
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}
	req.Header.Set("Authorization", token.Token)
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		err = errors.New(buf.String())
		log.Println("WARNING: validation error", endpoint, resource, resp.StatusCode, err)
		debug.PrintStack()
		return err, resp.StatusCode
	}
	return nil, http.StatusOK
}

func validateResourceDelete(token auth.Token, config config.Config, endpoints []string, id string) (err error, code int) {
	if config.DisableValidation {
		return nil, http.StatusOK
	}
	for _, endpoint := range endpoints {
		req, err := http.NewRequest("DELETE", endpoint+"/"+url.PathEscape(id)+"?dry-run=true", nil)
		if err != nil {
			debug.PrintStack()
			return err, http.StatusInternalServerError
		}
		req.Header.Set("Authorization", token.Token)
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
			log.Println("WARNING: validation error", endpoint, resp.StatusCode, err)
			debug.PrintStack()
			return err, resp.StatusCode
		}
	}
	return nil, http.StatusOK
}

type PermSearchElement struct {
	Id                string            `json:"id"`
	Name              string            `json:"name"`
	Shared            bool              `json:"shared"`
	Creator           string            `json:"creator"`
	PermissionHolders PermissionHolders `json:"permission_holders"`
}

type PermissionHolders struct {
	AdminUsers   []string `json:"admin_users"`
	ReadUsers    []string `json:"read_users"`
	WriteUsers   []string `json:"write_users"`
	ExecuteUsers []string `json:"execute_users"`
}

var ResourcesEffectedByUserDelete_BATCH_SIZE = 1000

func (this *Com) ResourcesEffectedByUserDelete(token auth.Token, resource string) (deleteResourceIds []string, deleteUserFromResourceIds []string, err error) {
	err = this.iterateResource(token, resource, ResourcesEffectedByUserDelete_BATCH_SIZE, "a", func(element PermSearchElement) {
		if len(element.PermissionHolders.AdminUsers) > 1 {
			deleteUserFromResourceIds = append(deleteUserFromResourceIds, element.Id)
		} else {
			deleteResourceIds = append(deleteResourceIds, element.Id)
		}
	})
	if err != nil {
		return
	}
	userid := token.GetUserId()
	err = this.iterateResource(token, resource, ResourcesEffectedByUserDelete_BATCH_SIZE, "r", func(element PermSearchElement) {
		if !contains(element.PermissionHolders.AdminUsers, userid) {
			deleteUserFromResourceIds = append(deleteUserFromResourceIds, element.Id)
		}
	})
	if err != nil {
		return
	}
	err = this.iterateResource(token, resource, ResourcesEffectedByUserDelete_BATCH_SIZE, "w", func(element PermSearchElement) {
		if !contains(element.PermissionHolders.AdminUsers, userid) &&
			!contains(element.PermissionHolders.ReadUsers, userid) {
			deleteUserFromResourceIds = append(deleteUserFromResourceIds, element.Id)
		}
	})
	if err != nil {
		return
	}
	err = this.iterateResource(token, resource, ResourcesEffectedByUserDelete_BATCH_SIZE, "x", func(element PermSearchElement) {
		if !contains(element.PermissionHolders.AdminUsers, userid) &&
			!contains(element.PermissionHolders.ReadUsers, userid) &&
			!contains(element.PermissionHolders.WriteUsers, userid) {
			deleteUserFromResourceIds = append(deleteUserFromResourceIds, element.Id)
		}
	})
	if err != nil {
		return
	}
	return deleteResourceIds, deleteUserFromResourceIds, err
}

func (this *Com) iterateResource(token auth.Token, resource string, batchsize int, rights string, handler func(element PermSearchElement)) (err error) {
	lastCount := batchsize
	lastElement := PermSearchElement{}
	for lastCount == batchsize {
		query := url.Values{}
		query.Add("limit", strconv.Itoa(batchsize))
		query.Add("sort", "name.asc")
		query.Add("rights", rights)
		if lastElement.Id == "" {
			query.Add("offset", "0")
		} else {
			name, err := json.Marshal(lastElement.Name)
			if err != nil {
				return err
			}
			query.Add("after.sort_field_value", string(name))
			query.Add("after.id", lastElement.Id)
		}
		temp := []PermSearchElement{}
		err = this.queryResourceInPermsearch(token, resource, query, &temp)
		if err != nil {
			return err
		}
		lastCount = len(temp)
		if lastCount > 0 {
			lastElement = temp[lastCount-1]
		}
		for _, element := range temp {
			handler(element)
		}
	}
	return err
}

func (this *Com) queryResourceInPermsearch(token auth.Token, resource string, query url.Values, result interface{}) (err error) {
	req, err := http.NewRequest("GET", this.config.PermissionsUrl+"/v3/resources/"+resource+"?"+query.Encode(), nil)
	if err != nil {
		debug.PrintStack()
		return err
	}
	req.Header.Set("Authorization", token.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		debug.PrintStack()
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		err = errors.New(buf.String())
		log.Println("ERROR: queryResourceInPermsearch()", resource, resp.StatusCode, err)
		debug.PrintStack()
		return err
	}
	err = json.NewDecoder(resp.Body).Decode(result)
	return
}

func contains(list []string, value string) bool {
	for _, element := range list {
		if element == value {
			return true
		}
	}
	return false
}

const Seperator = "$"

func removeIdModifier(id string) string {
	return strings.SplitN(id, Seperator, 2)[0]
}

func removeIdModifiers(ids []string) (result []string) {
	for _, id := range ids {
		result = append(result, removeIdModifier(id))
	}
	return result
}

func PreventIdModifier(id string) error {
	if strings.Contains(id, Seperator) {
		return errors.New("no edit on ids with " + Seperator + " part allowed")
	}
	return nil
}

func RemoveDuplicates[T comparable](slice []T) []T {
	keys := make(map[T]bool)
	result := []T{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			result = append(result, entry)
		}
	}
	return result
}
