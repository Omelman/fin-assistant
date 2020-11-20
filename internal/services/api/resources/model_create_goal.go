/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type CreateGoal struct {
	Key
	Attributes CreateGoalAttributes `json:"attributes"`
}
type CreateGoalResponse struct {
	Data     CreateGoal `json:"data"`
	Included Included   `json:"included"`
}

type CreateGoalListResponse struct {
	Data     []CreateGoal `json:"data"`
	Included Included     `json:"included"`
	Links    *Links       `json:"links"`
}

// MustCheckToken - returns CheckToken from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCreateGoal(key Key) *CreateGoal {
	var createTransaction CreateGoal
	if c.tryFindEntry(key, &createTransaction) {
		return &createTransaction
	}
	return nil
}
