/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Balance struct {
	Key
	Attributes BalanceAttributes `json:"attributes"`
}
type BalanceResponse struct {
	Data     Balance  `json:"data"`
	Included Included `json:"included"`
}

type BalanceListResponse struct {
	Data     []Balance `json:"data"`
	Included Included  `json:"included"`
	Links    *Links    `json:"links"`
}

// MustCheckToken - returns CheckToken from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustBalance(key Key) *Balance {
	var createBalance Balance
	if c.tryFindEntry(key, &createBalance) {
		return &createBalance
	}
	return nil
}
