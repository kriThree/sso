package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	storage "sso/internal/storage"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}
type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appId int) (models.App, error)
}

var (
	ErrInvalidatedCredentials = errors.New("invalid credentials")
)

// New creates a new instance of the Auth struct.
func New(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *Auth {

	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

// Login logs in a user.
func (a Auth) Login(ctx context.Context, email string, password string, appId int64) (token string, err error) {

	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("attempting to login user")

	user, err := a.userProvider.User(context.TODO(), email)

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("user not found")
			return "", fmt.Errorf("%s: %w", op, ErrInvalidatedCredentials)
		}

		log.Error("failed getting user", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Error("invalid credentials", err)

		return "", fmt.Errorf("%s: %w", op, ErrInvalidatedCredentials)
	}

	app, err := a.appProvider.App(ctx, int(appId))

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully logged in user")

	token, err = jwt.NewToken(user, app, a.tokenTTL)

	if err != nil {
		log.Error("failed generating token", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// RegisterNewUser registers a new user.
func (a Auth) RegisterNewUser(ctx context.Context, email string, password string) (userId int64, err error) {

	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Error("failed generating password hash", err)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)

	if err != nil {
		log.Error("failed saving user", err)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("registering new user")

	return id, nil
}

// IsAdmin checks if a user is an admin.
func (a Auth) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op = "Auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("userId", userId),
	)

	log.Info("checking if user is admin")

	IsAdmin, err := a.userProvider.IsAdmin(ctx, userId)

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checking if user is admin", slog.Bool("isAdmin", IsAdmin))

	return IsAdmin, nil

}
