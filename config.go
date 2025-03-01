package funcy

import "context"

// GetConfig returns a key's value from the config table.
func (cl *Client) GetConfig(key string) string {
	var value string
	err := cl.Conn.QueryRow(context.Background(), `select value from config where key = $1`, key).Scan(&value)
	if err != nil {
		return ""
	}

	return value
}
