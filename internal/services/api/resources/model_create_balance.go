/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type CreateBalance struct {
	Key
	Attributes CreateBalanceAttributes `json:"attributes"`
}
type CreateBalanceResponse struct {
	Data     CreateBalance `json:"data"`
	Included Included      `json:"included"`
}

type CreateBalanceListResponse struct {
	Data     []CreateBalance `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
}

// MustCheckToken - returns CheckToken from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCreateBalance(key Key) *CreateBalance {
	var createBalance CreateBalance
	if c.tryFindEntry(key, &createBalance) {
		return &createBalance
	}
	return nil
}
