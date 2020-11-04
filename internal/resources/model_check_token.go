/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type CheckToken struct {
	Key
	Attributes CheckTokenAttributes `json:"attributes"`
}
type CheckTokenResponse struct {
	Data     CheckToken `json:"implementation"`
	Included Included   `json:"included"`
}

type CheckTokenListResponse struct {
	Data     []CheckToken `json:"implementation"`
	Included Included     `json:"included"`
	Links    *Links       `json:"links"`
}

// MustCheckToken - returns CheckToken from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCheckToken(key Key) *CheckToken {
	var checkToken CheckToken
	if c.tryFindEntry(key, &checkToken) {
		return &checkToken
	}
	return nil
}
