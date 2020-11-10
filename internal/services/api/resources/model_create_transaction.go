/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type CreateTransaction struct {
	Key
	Attributes CreateTransactionAttributes `json:"attributes"`
}
type CreateTransactionResponse struct {
	Data     CreateTransaction `json:"data"`
	Included Included          `json:"included"`
}

type CreateTransactionListResponse struct {
	Data     []Balance `json:"data"`
	Included Included  `json:"included"`
	Links    *Links    `json:"links"`
}

// MustCheckToken - returns CheckToken from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCreateTransaction(key Key) *CreateTransaction {
	var createTransaction CreateTransaction
	if c.tryFindEntry(key, &createTransaction) {
		return &createTransaction
	}
	return nil
}
