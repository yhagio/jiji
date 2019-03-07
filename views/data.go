package views

import (
	"jiji/models"
	"net/http"
	"time"
)

type Data struct {
	Alert *Alert
	User  *models.User
	Yield interface{}
}

// RedirectAlert accepts all the normal params for an
// http.Redirect and performs a redirect, but only after
// persisting the provided alert in a cookie so that it can
// be displayed when the new page is loaded.
func RedirectAlert(w http.ResponseWriter, r *http.Request, urlStr string, statusCode int, alert Alert) {
	persistAlert(w, alert)
	http.Redirect(w, r, urlStr, statusCode)
}

func persistAlert(w http.ResponseWriter, alert Alert) {
	// We don't want alerts showing up days later. If the
	// user doesnt load the redirect in 5 minutes we will just expire it.
	expiresAt := time.Now().Add(5 * time.Minute)

	alertLevel := http.Cookie{
		Name:     "alert_level",
		Value:    alert.Level,
		Expires:  expiresAt,
		HttpOnly: true,
	}

	message := http.Cookie{
		Name:     "alert_message",
		Value:    alert.Message,
		Expires:  expiresAt,
		HttpOnly: true,
	}

	http.SetCookie(w, &alertLevel)
	http.SetCookie(w, &message)
}

func clearAlert(w http.ResponseWriter) {
	alertLevel := http.Cookie{
		Name:     "alert_level",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}

	message := http.Cookie{
		Name:     "alert_message",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}

	http.SetCookie(w, &alertLevel)
	http.SetCookie(w, &message)
}

func getAlert(r *http.Request) *Alert {
	// If either cookie is missing we will assume the alert
	// is invalid and return nil
	alertLevel, err := r.Cookie("alert_level")
	if err != nil {
		return nil
	}
	message, err := r.Cookie("alert_message")
	if err != nil {
		return nil
	}
	alert := Alert{
		Level:   alertLevel.Value,
		Message: message.Value,
	}
	return &alert
}
