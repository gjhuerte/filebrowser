package share

// CreateBody is the request body for creating a share link.
type CreateBody struct {
	// Password protects the link when set.
	Password string `json:"password"`
	// Expires is how long the link lasts, in Unit. Empty means forever.
	Expires string `json:"expires"`
	// Unit is "seconds", "minutes", "hours" (the default) or "days".
	Unit string `json:"unit"`
}

// Link is the information needed to build a shareable link.
type Link struct {
	Hash         string `json:"hash" storm:"id,index"`
	Path         string `json:"path" storm:"index"`
	UserID       uint   `json:"userID"`
	Expire       int64  `json:"expire"`
	PasswordHash string `json:"password_hash,omitempty"`
	// Token is a random value that will only be set when PasswordHash is set. It is
	// URL-Safe and is used to download links in password-protected shares via a
	// query arg.
	Token string `json:"token,omitempty"`
}
