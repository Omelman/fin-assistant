/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type CompleteRecovery struct {
	Key
	Attributes CompleteRecoveryAttributes `json:"attributes"`
}
type CompleteRecoveryResponse struct {
	Data     CompleteRecovery `json:"data"`
	Included Included         `json:"included"`
}

type CompleteRecoveryListResponse struct {
	Data     []CompleteRecovery `json:"data"`
	Included Included           `json:"included"`
	Links    *Links             `json:"links"`
}

// MustCompleteRecovery - returns CompleteRecovery from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCompleteRecovery(key Key) *CompleteRecovery {
	var completeRecovery CompleteRecovery
	if c.tryFindEntry(key, &completeRecovery) {
		return &completeRecovery
	}
	return nil
}
