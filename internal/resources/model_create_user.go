/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type CreateUser struct {
	Key
	Attributes CreateUserAttributes `json:"attributes"`
}
type CreateUserResponse struct {
	Data     CreateUser `json:"data"`
	Included Included   `json:"included"`
}

type CreateUserListResponse struct {
	Data     []CreateUser `json:"data"`
	Included Included     `json:"included"`
	Links    *Links       `json:"links"`
}

// MustCreateUser - returns CreateUser from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCreateUser(key Key) *CreateUser {
	var createUser CreateUser
	if c.tryFindEntry(key, &createUser) {
		return &createUser
	}
	return nil
}
