/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ReturnBalance struct {
	Key
	Attributes ReturnBalanceAttributes `json:"attributes"`
}
type ReturnBalanceResponse struct {
	Data     ReturnBalance `json:"data"`
	Included Included      `json:"included"`
}

type ReturnBalanceListResponse struct {
	Data     []ReturnBalance `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
}

// MustCheckToken - returns CheckToken from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustReturnBalance(key Key) *ReturnBalance {
	var returnBalance ReturnBalance
	if c.tryFindEntry(key, &returnBalance) {
		return &returnBalance
	}
	return nil
}
