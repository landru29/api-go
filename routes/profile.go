package routes

// ProfileFacebook define a user profile
type ProfileFacebook struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Hd        string `json:"hd"`
	Locale    string `json:"locale"`
	Name      string `json:"name"`
	Picture   struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}

// ProfileGoogle define a user profile
type ProfileGoogle struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
	Hd            string `json:"hd"`
	Locale        string `json:"locale"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

// Profile define a user profile
type Profile struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
	Hd        string
	Locale    string
	Name      string
}
