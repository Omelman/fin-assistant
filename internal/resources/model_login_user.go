/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type LoginUser struct {
	Key
	Attributes LoginUserAttributes `json:"attributes"`
}
type LoginUserResponse struct {
	Data     LoginUser `json:"implementation"`
	Included Included  `json:"included"`
}

type LoginUserListResponse struct {
	Data     []LoginUser `json:"implementation"`
	Included Included    `json:"included"`
	Links    *Links      `json:"links"`
}

// MustLoginUser - returns LoginUser from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustLoginUser(key Key) *LoginUser {
	var loginUser LoginUser
	if c.tryFindEntry(key, &loginUser) {
		return &loginUser
	}
	return nil
}
