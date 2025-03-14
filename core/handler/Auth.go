package handler

import (
	"agenda-kaki-go/core/lib"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/shareed2k/goth_fiber"
	"golang.org/x/crypto/bcrypt"
)

type Authentication struct {
	C           *fiber.Ctx
	Res         *lib.SendResponse
	sessionName string
}

type SessionsOptions struct {
	CookiesKey string
	MaxAge     int
	HttpOnly   bool
	Secure     bool
}

func Auth(c *fiber.Ctx) *Authentication {
	return &Authentication{
		C:           c,
		Res:         &lib.SendResponse{Ctx: c},
		sessionName: "user_session",
	}
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
	return err == nil
}

func NewAuth(session *sessions.CookieStore) {
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleUrl := os.Getenv("GOOGLE_URL_OAUTH")
	gothic.Store = session
	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, googleUrl),
	)
}

func NewCookieStore(opts SessionsOptions) *sessions.CookieStore {
	// Create a new cookie store
	store := sessions.NewCookieStore([]byte(opts.CookiesKey))
	store.MaxAge(opts.MaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = opts.HttpOnly
	store.Options.Secure = opts.Secure
	return store
}

func SessionOpts() SessionsOptions {
	return SessionsOptions{
		CookiesKey: os.Getenv("COOKIES_KEY"),
		MaxAge: func() int {
			intVal, err := strconv.Atoi(os.Getenv("MAX_AGE"))
			if err != nil {
				log.Fatalf("sessionsOptions MaxAge: %v", err)
			}
			return intVal
		}(),
		HttpOnly: func() bool {
			b, err := strconv.ParseBool(os.Getenv("HTTP_ONLY_COOKIE"))
			if err != nil {
				log.Fatalf("sessionsOptions httpOnly: %v", err)
			}
			return b
		}(),
		Secure: func() bool {
			b, err := strconv.ParseBool(os.Getenv("SECURE_COOKIE"))
			if err != nil {
				log.Fatalf("sessionsOptions secure: %v", err)
			}
			return b
		}(),
	}
}

func (a *Authentication) StoreUserSession(us goth.User) error {
	// Store the user session
	err := goth_fiber.StoreInSession("_gothic_session", us.AccessToken, a.C)
	if err != nil {
		log.Println("Error storing user session", err)
	}
	return nil
}

func (a *Authentication) WhoAreYou() error {
	// Check if the user is authenticated
	_, err := goth_fiber.GetFromSession("_gothic_session", a.C)
	if err != nil {
		a.Res.Http401(fmt.Errorf("user not authenticated"))
	}
	return a.C.Next()
}
