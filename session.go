package funcy

// SetCookie in HTTP response.
func (cl *Client) SetCookie(name, value string) {
	cl.Session.Values[name] = value
}

// ClearCookie in HTTP response.
func (cl *Client) ClearCookie(name string) {
	delete(cl.Session.Values, name)
}

// GetCookie from HTTP request.
func (cl *Client) GetCookie(name string) string {
	if val, ok := cl.Session.Values[name].(string); ok {
		return val
	}
	return ""
}

// Save session.
func (cl *Client) Save() {
	cl.Session.Save(cl.R, cl.W)
}
