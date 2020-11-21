/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Category struct {
	Key
	Attributes CategoryAttributes `json:"attributes"`
}
type CategoryResponse struct {
	Data     Category `json:"data"`
	Included Included `json:"included"`
}

type CategoryListResponse struct {
	Data     []Category `json:"data"`
	Included Included   `json:"included"`
	Links    *Links     `json:"links"`
}

// MustCheckToken - returns CheckToken from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCategory(key Key) *Category {
	var category Category
	if c.tryFindEntry(key, &category) {
		return &category
	}
	return nil
}
