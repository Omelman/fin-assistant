/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type CreateTransactionAttributes struct {
	Date        string `json:"date,omitempty"`
	Description string `json:"description,omitempty"`
	Amount      int    `json:"amount,string,omitempty"`
	Category    string `json:"category,omitempty"`
	Include     bool   `json:"include,string,bool,omitempty"`
	BalaceId    int    `json:"balance_id,string,omitempty"`
}
