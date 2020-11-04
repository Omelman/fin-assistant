/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type AskRecovery struct {
	Key
	Attributes AskRecoveryAttributes `json:"attributes"`
}
type AskRecoveryResponse struct {
	Data     AskRecovery `json:"data"`
	Included Included    `json:"included"`
}

type AskRecoveryListResponse struct {
	Data     []AskRecovery `json:"data"`
	Included Included      `json:"included"`
	Links    *Links        `json:"links"`
}

// MustAskRecovery - returns AskRecovery from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustAskRecovery(key Key) *AskRecovery {
	var askRecovery AskRecovery
	if c.tryFindEntry(key, &askRecovery) {
		return &askRecovery
	}
	return nil
}
