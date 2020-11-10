/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type GetBalance struct {
	Key
	Attributes GetBalanceAttributes `json:"attributes"`
}
type GetBalanceResponse struct {
	Data     GetBalance `json:"data"`
	Included Included   `json:"included"`
}

type GetBalanceListResponse struct {
	Data     []GetBalance `json:"data"`
	Included Included     `json:"included"`
	Links    *Links       `json:"links"`
}

// MustCheckToken - returns CheckToken from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustGetBalance(key Key) *GetBalance {
	var createBalance GetBalance
	if c.tryFindEntry(key, &createBalance) {
		return &createBalance
	}
	return nil
}
