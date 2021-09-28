package controllers

import "net/http"

func getSession(r *http.Request) (string, error) {
	c, err := r.Cookie("session")
	if err != nil {
		return "", err
	}
	return c.Value, nil
}
