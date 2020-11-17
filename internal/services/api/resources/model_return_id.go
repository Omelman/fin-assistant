/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ReturnId struct {
	Key
	Attributes ReturnIdAttributes `json:"attributes"`
}
type ReturnIdResponse struct {
	Data     ReturnId `json:"data"`
	Included Included `json:"included"`
}

type ReturnIdListResponse struct {
	Data     []ReturnId `json:"data"`
	Included Included   `json:"included"`
	Links    *Links     `json:"links"`
}

// MustCheckToken - returns CheckToken from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustReturnBalance(key Key) *ReturnId {
	var returnId ReturnId
	if c.tryFindEntry(key, &returnId) {
		return &returnId
	}
	return nil
}
