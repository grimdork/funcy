package funcy

import (
	"context"
	"net/http"
	"os"

	ll "github.com/grimdork/loglines"
	"golang.org/x/crypto/bcrypt"
)

const (
	userSessions = "sessions"
	adminSessons = "admin_sessions"

	getPasswordSQL      = `select password from users where name = $1;`
	getAdminPasswordSQL = `select password from users where name = $1 and admin = true;`
	insertSessionSQL    = `insert into sessions (user_id, token, expires_at)
select id, $2, now() + interval '1 day'
from users where name = $1;
`
	insertAdminSessionSQL = `insert into admin_sessions (user_id, token, expires_at)
select id, $2, now() + interval '1 day'
from users where name = $1;
`
	validateSessionSQL = `select 1 from sessions s
join users u on s.user_id = u.id
where u.name = $1 and s.token = $2 and s.expires_at > now();
`
	validateAdminSessionSQL = `select 1 from admin_sessions s
join users u on s.user_id = u.id
where u.name = $1 and u.admin = true and s.token = $2 and s.expires_at > now();
`
	invalidateUserSessionsSQL  = `delete from sessions where user_id = (select id from users where name = $1);`
	invalidateAdminSessionsSQL = `delete from admin_sessions where user_id = (select id from users where name = $1);`
)

// Authenticate checks the user's password and updates the session if valid, and the user is an admin.
func (cl *Client) Authenticate(username, password string, asAdmin bool) bool {
	if username == "" || password == "" {
		return false
	}

	var hashedPassword string
	sql := getPasswordSQL
	if asAdmin {
		sql = getAdminPasswordSQL
	}
	err := cl.Conn.QueryRow(context.Background(), sql, username).Scan(&hashedPassword)
	if err != nil {
		ll.Msg("Failed to get password: %s", err.Error())
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		ll.Msg("Password mismatch: %s", err.Error())
		return false
	}

	cl.SetCookie("username", username)
	token := GenerateToken(username)
	cl.SetCookie("token", token)
	sql = insertSessionSQL
	if asAdmin {
		sql = insertAdminSessionSQL
	}
	_, err = cl.Conn.Exec(context.Background(), sql, username, token)
	if err != nil {
		ll.Msg("Failed to create session: %s", err.Error())
		return false
	}

	return true
}

// IsAuthenticated checks if the user is logged in.
func (cl *Client) IsAuthenticated(asAdmin bool) bool {
	username := cl.GetCookie("username")
	token := cl.GetCookie("token")
	if username == "" || token == "" {
		return false
	}

	sql := validateSessionSQL
	if asAdmin {
		sql = validateAdminSessionSQL
	}
	var valid int
	err := cl.Conn.QueryRow(context.Background(), sql, username, token).Scan(&valid)
	if err != nil || valid != 1 {
		return false
	}

	return true
}

// IsAdmin checks if the user is an admin.
func (cl *Client) IsAdmin() bool {
	return cl.CheckAdmin(cl.GetCookie("username"))
}

// CheckAdmin checks if a user is an admin.
func (cl *Client) CheckAdmin(username string) bool {
	if username == "" {
		return false
	}

	var admin bool
	err := cl.Conn.QueryRow(context.Background(), `select admin from users where name = $1`, username).Scan(&admin)
	if err != nil {
		return false
	}

	return admin
}

// InvalidateSessions removes all sessions for a user.
func (cl *Client) InvalidateSessions(username string) {
	_, err := cl.Conn.Exec(context.Background(), invalidateUserSessionsSQL, username)
	if err != nil {
		ll.Msg("Failed to invalidate user sessions: %s", err.Error())
	}
}

// InvalidateAdminSessions removes all admin sessions for a user.
func (cl *Client) InvalidateAdminSessions(username string) {
	_, err := cl.Conn.Exec(context.Background(), invalidateAdminSessionsSQL, username)
	if err != nil {
		ll.Msg("Failed to invalidate admin sessions: %s", err.Error())
	}
}

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

// ClearSession removes the session cookies of the current user.
func (cl *Client) ClearSession() {
	cl.ClearCookie("username")
	cl.ClearCookie("token")
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
