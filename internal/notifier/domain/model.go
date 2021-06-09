package domain

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type MailMessage struct {
	Sender      User   `json:"sender"`
	To          []User `json:"to"`
	Subject     string `json:"subject"`
	HtmlContent string `json:"htmlContent"`
}
