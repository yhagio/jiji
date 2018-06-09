package views

import (
	"jiji/models"
)

type Data struct {
	Alert *Alert
	User  *models.User
	Yield interface{}
}
