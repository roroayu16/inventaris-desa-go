package main

import (
	"net/http"
	"net/url"
	"strings"
)

type Flash struct {
	Type    string
	Message string
}

func SetFlash(
	w http.ResponseWriter,
	flashType string,
	message string,
) {
	value := url.QueryEscape(flashType + "|" + message)

	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "flash",
			Value:    value,
			Path:     "/",
			HttpOnly: true,
		},
	)
}

func GetFlash(r *http.Request, w http.ResponseWriter) *Flash {

	cookie, err := r.Cookie("flash")
	if err != nil {
		return nil
	}

	// Hapus cookie agar hanya tampil sekali
	http.SetCookie(w, &http.Cookie{
		Name:   "flash",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	value, _ := url.QueryUnescape(cookie.Value)

	parts := strings.SplitN(value, "|", 2)

	if len(parts) != 2 {
		return nil
	}

	return &Flash{
		Type:    parts[0],
		Message: parts[1],
	}
}
