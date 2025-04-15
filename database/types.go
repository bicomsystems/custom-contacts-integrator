package database

type Contact struct {
	Id        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Company   string  `json:"company"`
	Type      string  `json:"type"`
	FbUserID  string  `json:"fb_user_id"`
	Phones    []Phone `json:"phones"`
	Emails    []Email `json:"emails"`
}

type Phone struct {
	Number string `json:"number"`
	Label  string `json:"label"`
}

type Email struct {
	Email string `json:"email"`
	Label string `json:"label"`
}

type Config struct {
	Phones []Phone `json:"phones"`
	Emails []Email `json:"emails"`
}
