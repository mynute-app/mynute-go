package handlers

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gorilla/sessions"
	gothfiber "github.com/luigiazoreng/goth_fiber"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"golang.org/x/crypto/bcrypt"
)

const (
	key    = "randomString"
	MaxAge = 86400 * 30
	IsProd = false
)

type Authentication struct {
	C           fiber.Ctx
	Res         *Res
	sessionName string
}

type SessionsOptions struct {
	CookiesKey string
	MaxAge     int
	HttpOnly   bool
	Secure     bool
}

func HashPassword(password string) (string, error) {
	// Generate a hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword, password string) bool {
	// Compare the password with the hashed one
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func NewAuth(session *sessions.CookieStore) {
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	gothic.Store = session
	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, "http://localhost:3000/auth/google/callback"),
	)
}

func NewCookieStore(opts SessionsOptions) *sessions.CookieStore {
	// Create a new cookie store
	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(opts.MaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = opts.HttpOnly
	store.Options.Secure = opts.Secure
	return store
}

func Auth(c fiber.Ctx) *Authentication {
	return &Authentication{
		C:           c,
		Res:         &Res{Ctx: c},
		sessionName: "user_session",
	}
}

func (a *Authentication) StoreUserSession() error {
	// Login the user
	session, err := gothfiber.SessionStore.Get(a.C)
	if err != nil {
		return a.Res.Http401(err)
	}
	log.Println("User id: ", session.ID())

	err = session.Save()
	if err != nil {
		return a.Res.Http401(err)
	}
	return nil
}

func (a *Authentication) WhoAreYou() error {
	// Check if the user is authenticated
	us, err := gothfiber.GetFromSession("_gothic_session", a.C)
	if err != nil {
		log.Println("User not authenticated", err)
		log.Println(us)
		return a.Res.Http401(err)
	}
	log.Println("User authenticated:", us)
	return a.C.Next()

}
