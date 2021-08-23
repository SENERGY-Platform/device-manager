package com

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/device-manager/lib/auth"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
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

func validateResource(token auth.Token, endpoints []string, resource interface{}) (err error, code int) {
	for _, endpoint := range endpoints {
		b := new(bytes.Buffer)
		err = json.NewEncoder(b).Encode(resource)
		if err != nil {
			debug.PrintStack()
			return err, http.StatusInternalServerError
		}
		req, err := http.NewRequest("PUT", endpoint, b)
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
			log.Println("WARNING: validation error", endpoint, resource, resp.StatusCode, err)
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

func (this *Com) ResourcesEffectedByUserDelete(token auth.Token, resource string) (deleteResourceIds []string, deleteUserFromResourceIds []string, err error) {
	rights := "a"
	limit := 1000
	lastCount := limit
	lastElement := PermSearchElement{}
	for lastCount == limit {
		query := url.Values{}
		query.Add("limit", strconv.Itoa(limit))
		query.Add("sort", "name.asc")
		query.Add("rights", rights)
		if lastElement.Id == "" {
			query.Add("offset", "0")
		} else {
			name, err := json.Marshal(lastElement.Name)
			if err != nil {
				return deleteResourceIds, deleteUserFromResourceIds, err
			}
			query.Add("after.sort_field_value", string(name))
			query.Add("after.id", lastElement.Id)
		}
		temp := []PermSearchElement{}
		err = this.queryResourceInPermsearch(token, resource, query, &temp)
		if err != nil {
			return deleteResourceIds, deleteUserFromResourceIds, err
		}
		lastCount = len(temp)
		lastElement = temp[lastCount-1]
		for _, element := range temp {
			if len(element.PermissionHolders.AdminUsers) > 1 {
				deleteUserFromResourceIds = append(deleteUserFromResourceIds, element.Id)
			} else {
				deleteResourceIds = append(deleteResourceIds, element.Id)
			}
		}
	}
	return deleteResourceIds, deleteUserFromResourceIds, err
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
