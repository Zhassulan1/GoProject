package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/Zhassulan1/Go_Project/pkg/clinic-api/model"
)

type contextKey string

// userContextKey is used as a key for getting and setting user information in the request
// context.
const userContextKey = contextKey("user")

// contextSetUser returns a new copy of the request with the provided User struct added to the
// context.
func (app *application) contextSetUser(r *http.Request, user *model.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// contextGetUser retrieves the User struct from the request context. The only time that
// this helper should be used is when we logically expect there to be a User struct value
// in the context, and if it doesn't exist it will firmly be an 'unexpected' error, upon we panic.
func (app *application) contextGetUser(r *http.Request) (*model.User, error) {
	var err error
	user, ok := r.Context().Value(userContextKey).(*model.User)

	log.Print("\n\nIs there user\n\n")
	authToken := r.Header.Get("Authorization")

	if len(authToken) < 8 {
		log.Print("Error getting user \nAuth: ", authToken)
		return model.AnonymousUser, errors.New("invalid or missing authentication token")
	}

	if ok {
		return user, nil
	}

	user, err = app.models.Users.GetUserByToken(authToken[7:])
	if err != nil {
		log.Print(err, "\n\n authToken[7:] = ", authToken[7:], "\n\n")
		return model.AnonymousUser, errors.New("invalid or missing authentication token")
	}

	log.Print("Auth: ", authToken)

	return user, nil
}
