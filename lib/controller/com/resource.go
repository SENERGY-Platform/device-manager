package com

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/device-manager/lib/auth"
	"github.com/SENERGY-Platform/device-manager/lib/config"
	"github.com/SENERGY-Platform/permissions-v2/pkg/client"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"slices"
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

func getResourceFromServiceWithQueryParam(token auth.Token, endpoint string, id string, query url.Values, result interface{}) (err error, code int) {
	req, err := http.NewRequest("GET", endpoint+"/"+url.PathEscape(id)+"?"+query.Encode(), nil)
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
		Timeout: 10 * time.Second,
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

func containsOtherAdmin(m map[string]client.PermissionsMap, notThisKey string) bool {
	for k, v := range m {
		if k != notThisKey && v.Administrate {
			return true
		}
	}
	return false
}

var ResourcesEffectedByUserDelete_BATCH_SIZE int64 = 1000

func (this *Com) ResourcesEffectedByUserDelete(token auth.Token, resource string) (deleteResourceIds []string, deleteUserFromResource []client.Resource, err error) {
	userid := token.GetUserId()
	err = this.iterateResource(token, resource, ResourcesEffectedByUserDelete_BATCH_SIZE, client.Administrate, func(element client.Resource) {
		if containsOtherAdmin(element.UserPermissions, userid) {
			deleteUserFromResource = append(deleteUserFromResource, element)
		} else {
			deleteResourceIds = append(deleteResourceIds, element.Id)
		}
	})
	if err != nil {
		return
	}

	err = this.iterateResource(token, resource, ResourcesEffectedByUserDelete_BATCH_SIZE, client.Read, func(element client.Resource) {
		if !slices.ContainsFunc(deleteUserFromResource, func(resource client.Resource) bool {
			return resource.Id == element.Id
		}) {
			deleteUserFromResource = append(deleteUserFromResource, element)
		}
	})
	if err != nil {
		return
	}
	err = this.iterateResource(token, resource, ResourcesEffectedByUserDelete_BATCH_SIZE, client.Write, func(element client.Resource) {
		if !slices.ContainsFunc(deleteUserFromResource, func(resource client.Resource) bool {
			return resource.Id == element.Id
		}) {
			deleteUserFromResource = append(deleteUserFromResource, element)
		}
	})
	if err != nil {
		return
	}
	err = this.iterateResource(token, resource, ResourcesEffectedByUserDelete_BATCH_SIZE, client.Execute, func(element client.Resource) {
		if !slices.ContainsFunc(deleteUserFromResource, func(resource client.Resource) bool {
			return resource.Id == element.Id
		}) {
			deleteUserFromResource = append(deleteUserFromResource, element)
		}
	})
	if err != nil {
		return
	}
	return deleteResourceIds, deleteUserFromResource, err
}

func (this *Com) iterateResource(token auth.Token, resource string, batchsize int64, rights client.Permission, handler func(element client.Resource)) (err error) {
	lastCount := batchsize
	var offset int64 = 0
	for lastCount == batchsize {
		options := client.ListOptions{
			Limit:  batchsize,
			Offset: offset,
		}
		offset += batchsize
		ids, err, _ := this.perm.ListAccessibleResourceIds(token.Jwt(), resource, options, rights)
		if err != nil {
			return err
		}
		lastCount = int64(len(ids))
		for _, id := range ids {
			element, err, _ := this.perm.GetResource(client.InternalAdminToken, resource, id)
			if err != nil {
				return err
			}
			handler(element)
		}
	}
	return err
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
