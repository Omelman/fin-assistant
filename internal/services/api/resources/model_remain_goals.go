/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type RemainGoals struct {
	Key
	Attributes RemainGoalsAttributes `json:"attributes"`
}
type RemainGoalsResponse struct {
	Data     RemainGoals `json:"data"`
	Included Included    `json:"included"`
}

type RemainGoalsListResponse struct {
	Data     []RemainGoals `json:"data"`
	Included Included      `json:"included"`
	Links    *Links        `json:"links"`
}

// MustCheckToken - returns CheckToken from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustRemainGoals(key Key) *RemainGoals {
	var createBalance RemainGoals
	if c.tryFindEntry(key, &createBalance) {
		return &createBalance
	}
	return nil
}
