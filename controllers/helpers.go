package controllers

import (
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
)

func parseURLParams(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return parseValues(r.Form, dst)
}

func parseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return parseValues(r.PostForm, dst)
}

func parseValues(values url.Values, dst interface{}) error {
	dec := schema.NewDecoder()
	// Call the IgnoreUnkownKeys function to tell schema's decoder
	// to ignore the CSRF token key
	dec.IgnoreUnknownKeys(true)
	err := dec.Decode(dst, values)
	if err != nil {
		return err
	}
	return nil
}
