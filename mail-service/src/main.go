package main

import "mail-service/src/api"

func main() {
	app := api.Config{
		WebPort: "80",
		Mailer:  api.CreateMailer(),
	}
	app.Serve()
}
