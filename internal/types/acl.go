package types

// UserACL represents an organization ACL entry.
type UserACL struct {
	OrgID       int64    `json:"orgId"`
	OrgName     string   `json:"orgName"`
	Currency    string   `json:"currency"`
	PaymentModel string  `json:"paymentModel"`
	RoleNames   []string `json:"roleNames"`
	TimeZone    string   `json:"timeZone"`
}

// MeDetail represents the caller's user details.
type MeDetail struct {
	UserID      int64  `json:"userId"`
	ParentOrgID int64  `json:"parentOrgId"`
}
