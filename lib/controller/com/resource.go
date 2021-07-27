package com

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
)

func getResourceFromService(token string, endpoint string, id string, result interface{}) (err error, code int) {
	req, err := http.NewRequest("GET", endpoint+"/"+url.PathEscape(id), nil)
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

func validateResource(token string, endpoints []string, resource interface{}) (err error, code int) {
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
			log.Println("WARNING: validation error", endpoint, resource, resp.StatusCode, err)
			debug.PrintStack()
			return err, resp.StatusCode
		}
	}
	return nil, http.StatusOK
}
