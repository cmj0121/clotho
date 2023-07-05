// Get the GitHub user information.
package github

type User struct {
	// The general information of the user.
	Login      string `json:"login"`
	Id         int    `json:"id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	Company    string `json:"company"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`

	// The raw information of the user.
	raw map[string]interface{}
}

// The raw information of the user from GitHub.
func (user User) Raw() map[string]interface{} {
	return user.raw
}
