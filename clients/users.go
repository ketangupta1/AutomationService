package clients

// UserSession represents the response after a successful authentication.
type UserSession struct {
	UserProfile
	UserSessionTokens
}

// UserSessionTokens represents response after renew access token.
type UserSessionTokens struct {
	AccessToken  string `json:"jwtToken"`
	RefreshToken string `json:"refreshToken"`
	FeedToken    string `json:"feedToken"`
}

// UserProfile represents a user's personal and financial profile.
type UserProfile struct {
	ClientCode    string   `json:"clientcode"`
	UserName      string   `json:"name"`
	Email         string   `json:"email"`
	Phone         string   `json:"mobileno"`
	Broker        string   `json:"broker"`
	Products      []string `json:"products"`
	LastLoginTime string   `json:"lastlogintime"`
	Exchanges     []string `json:"exchanges"`
}
