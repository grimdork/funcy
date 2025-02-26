package funcy

import (
	"net/http"
	"os"
)

// SetCookie in HTTP response.
func (cl *Client) SetCookie(name, value string) {
	http.SetCookie(cl.W, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   os.Getenv("DOMAIN"),
		MaxAge:   86400,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearCookie in HTTP response.
func (cl *Client) ClearCookie(name string) {
	http.SetCookie(cl.W, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Domain:   os.Getenv("DOMAIN"),
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// GetCookie from HTTP request.
func (cl *Client) GetCookie(name string) string {
	cookie, err := cl.R.Cookie(name)
	if err == nil {
		return cookie.Value
	}

	return ""
}

// SetHeader in HTTP response.
func (cl *Client) SetHeader(name, value string) {
	cl.W.Header().Set(name, value)
}

// Redirect to URL.
func (cl *Client) Redirect(url string, code int) {
	http.Redirect(cl.W, cl.R, url, code)
}
