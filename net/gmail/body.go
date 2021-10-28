package gmail

import "github.com/matcornic/hermes/v2"

// Create beautiful html email

type Template struct {
	Signature string
}

func MakeBody() (string, error) {
	// Configure hermes by setting a theme and your product info
	h := hermes.Hermes{
		// Optional Theme
		// Theme: new(Default)
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: "Hermes",
			Link: "https://example-hermes.com/",
			// Optional product logo
			Logo: "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
		},
	}

	email := hermes.Email{
		Body: hermes.Body{
			Name: "Jon Snow",
			Intros: []string{
				"Welcome to Hermes! We're very excited to have you on board.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started with Hermes, please click here:",
					Button: hermes.Button{
						Color: "#22BC66", // Optional action button color
						Text:  "Confirm your account",
						Link:  "https://hermes-example.com/confirm?token=d9729feb74992cc3482b350163a1a010",
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}

	// Generate an HTML email with the provided contents (for modern clients)
	return h.GenerateHTML(email)
	/*
		// Generate the plaintext version of the e-mail (for clients that do not support xHTML)
		emailText, err := h.GeneratePlainText(email)
		if err != nil {
			panic(err) // Tip: Handle error with something else than a panic ;)
		}*/
}

// https://github.com/toikarin/go-email2html/

func MimeConvHtml() {

}
