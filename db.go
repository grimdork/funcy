package funcy

import (
	"context"

	ll "github.com/grimdork/loglines"
	"golang.org/x/crypto/bcrypt"
)

const getPasswordSQL = `select password from users where name = $1 and admin = true;`
const insertSessionSQL = `insert into sessions (user_id, token, expires_at)
select id, $2, now() + interval '1 hour'
from users where name = $1;
`

// Authenticate checks the user's password and updates the session if valid, and the user is an admin.
func (cl *Client) Authenticate(username, password string) bool {
	if username == "" || password == "" {
		return false
	}

	var hashedPassword string
	err := cl.Conn.QueryRow(context.Background(), getPasswordSQL, username).Scan(&hashedPassword)
	if err != nil {
		cl.SetCookie("message", "Invalid credentials.")
		cl.Save()
		ll.Msg("failed to get password: %s", err.Error())
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		cl.SetCookie("message", "Invalid credentials.")
		cl.Save()
		ll.Msg("passwords do not match: %s", err.Error())
		return false
	}

	cl.SetCookie("username", username)
	token := GenerateToken(username)
	cl.SetCookie("token", token)
	cl.Save()
	_, err = cl.Conn.Exec(context.Background(), insertSessionSQL, username, token)
	if err != nil {
		cl.SetCookie("message", "Failed to create session: "+err.Error())
		cl.Save()
		ll.Msg("failed to create session: %s", err.Error())
		return false
	}

	return true
}

const validateSessionSQL = `select 1 from sessions s
join users u on s.user_id = u.id
where u.name = $1 and u.admin = true and s.token = $2 and s.expires_at > now();
`

// IsAuthenticated checks if the user is logged in.
func (cl *Client) IsAuthenticated() bool {
	username := cl.GetCookie("username")
	token := cl.GetCookie("token")
	if username == "" || token == "" {
		return false
	}

	var valid int
	err := cl.Conn.QueryRow(context.Background(), validateSessionSQL, username, token).Scan(&valid)
	if err != nil && valid != 1 {
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
