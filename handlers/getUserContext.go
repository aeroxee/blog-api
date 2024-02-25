package handlers

import (
	"net/http"

	"github.com/Aeroxee/blog-api/auth"
	"github.com/Aeroxee/blog-api/models"
)

func getUserContext(r *http.Request) (models.User, error) {
	context := r.Context().Value(&auth.UserAuth{}).(auth.Claims)
	user, err := models.NewUserModel(models.GetDB()).GetUserByID(context.Credential.UserID)
	return user, err
}
