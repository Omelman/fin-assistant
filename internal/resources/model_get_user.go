/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type GetUser struct {
	Key
	Attributes *GetUserAttributes `json:"attributes,omitempty"`
}
type GetUserResponse struct {
	Data     GetUser  `json:"data"`
	Included Included `json:"included"`
}

type GetUserListResponse struct {
	Data     []GetUser `json:"data"`
	Included Included  `json:"included"`
	Links    *Links    `json:"links"`
}

// MustGetUser - returns GetUser from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustGetUser(key Key) *GetUser {
	var getUser GetUser
	if c.tryFindEntry(key, &getUser) {
		return &getUser
	}
	return nil
}
