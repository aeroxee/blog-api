package handlers

import (
	"net/http"
	"strconv"
)

func getQueryInteger(r *http.Request, key string, defaultValue int) int {
	result := r.URL.Query().Get(key)
	if result == "" {
		return defaultValue
	}

	resultInt, err := strconv.Atoi(result)
	if err != nil {
		return defaultValue
	}
	return resultInt
}

func getQueryString(r *http.Request, key, defaultValue string) string {
	result := r.URL.Query().Get(key)
	if result == "" {
		return defaultValue
	}
	return result
}

func getQueryBool(r *http.Request, key string, defaultValue bool) bool {
	result := r.URL.Query().Get(key)
	if result == "" {
		return defaultValue
	}

	resultBool, err := strconv.ParseBool(result)
	if err != nil {
		return defaultValue
	}

	return resultBool
}
