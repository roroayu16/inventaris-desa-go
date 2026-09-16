package main

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Flash struct {
	Type        string
	Message     string
	Duration    int
	RedirectURL string
}

func SetFlash(
	w http.ResponseWriter,
	flashType string,
	message string,
) {

	SetFlashWithDuration(
		w,
		flashType,
		message,
		5000,
	)
}

func SetFlashWithDuration(
	w http.ResponseWriter,
	flashType string,
	message string,
	duration int,
) {

	value := url.QueryEscape(
		flashType + "|" + message + "|" + strconv.Itoa(duration),
	)

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

func SetFlashWithRedirect(
	w http.ResponseWriter,
	flashType string,
	message string,
	duration int,
	redirectURL string,
) {

	value := url.QueryEscape(
		flashType + "|" +
			message + "|" +
			strconv.Itoa(duration) + "|" +
			redirectURL,
	)

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

	parts := strings.SplitN(value, "|", 4)

	if len(parts) < 2 {
		return nil
	}

	duration := 5000

	if len(parts) == 3 {
		parsedDuration, err := strconv.Atoi(parts[2])

		if err == nil && parsedDuration > 0 {
			duration = parsedDuration
		}
	}

	redirectURL := ""

	if len(parts) == 4 {
		redirectURL = parts[3]
	}

	return &Flash{
		Type:        parts[0],
		Message:     parts[1],
		Duration:    duration,
		RedirectURL: redirectURL,
	}
}
