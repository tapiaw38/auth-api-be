package user

import (
	"context"
	"errors"
	"sync"
	"time"

	user_repo "github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories/user"
	"github.com/tapiaw38/auth-api-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/auth-api-be/internal/platform/errors"
	"github.com/tapiaw38/auth-api-be/internal/platform/errors/mappings"
)

const mobileRedirectURI = "https://app.practiq.com.ar/auth/mobile-callback"

type (
	GoogleMobileLoginUsecase interface {
		Complete(context.Context, string, string) apperrors.ApplicationError
		Poll(string) (*GoogleMobileSessionResult, bool)
	}

	googleMobileLoginUsecase struct {
		contextFactory appcontext.Factory
		sessions       map[string]storedGoogleMobileSession
		mutex          sync.Mutex
	}

	GoogleMobileSessionStatus string

	GoogleMobileSessionResult struct {
		Status       GoogleMobileSessionStatus `json:"status"`
		Token        string                    `json:"token,omitempty"`
		RefreshToken string                    `json:"refresh_token,omitempty"`
		Data         UserOutputData            `json:"data,omitempty"`
		Message      string                    `json:"message,omitempty"`
	}

	storedGoogleMobileSession struct {
		result    GoogleMobileSessionResult
		createdAt time.Time
	}
)

const (
	MobileSessionPending   GoogleMobileSessionStatus = "pending"
	MobileSessionDone      GoogleMobileSessionStatus = "done"
	MobileSessionError     GoogleMobileSessionStatus = "error"
	googleMobileSessionTTL                           = 5 * time.Minute
)

func NewGoogleMobileLoginUsecase(contextFactory appcontext.Factory) GoogleMobileLoginUsecase {
	return &googleMobileLoginUsecase{
		contextFactory: contextFactory,
		sessions:       make(map[string]storedGoogleMobileSession),
	}
}

func (u *googleMobileLoginUsecase) Complete(ctx context.Context, code string, state string) apperrors.ApplicationError {
	app := u.contextFactory()
	token, err := app.Integrations.SSO.ExchangeCode(ctx, code, mobileRedirectURI)
	if err != nil {
		return u.storeError(state, apperrors.NewApplicationError(mappings.UserLoginGoogleExchangeError, err))
	}

	userInfo, err := app.Integrations.SSO.GetUserInfo(ctx, token)
	if err != nil {
		return u.storeError(state, apperrors.NewApplicationError(mappings.UserLoginGoogleUserInfoError, err))
	}

	userID, appErr := resolveGoogleUser(ctx, app, userInfo)
	if appErr != nil {
		return u.storeError(state, appErr)
	}

	if userID == nil {
		return u.storeError(state, apperrors.NewApplicationError(mappings.UserLoginUserNotFoundError, errors.New("user not found after google login")))
	}

	resolvedUser, appErr := app.Repositories.User.Get(ctx, user_repo.GetFilterOptions{ID: *userID})
	if appErr != nil {
		return u.storeError(state, appErr)
	}

	if resolvedUser == nil {
		return u.storeError(state, apperrors.NewApplicationError(mappings.UserLoginUserNotFoundError, errors.New("user not found")))
	}

	issued, appErr := issueSession(ctx, app, resolvedUser)
	if appErr != nil {
		return u.storeError(state, appErr)
	}

	u.store(state, GoogleMobileSessionResult{
		Status:       MobileSessionDone,
		Token:        issued.accessToken,
		RefreshToken: issued.refreshToken,
		Data:         toUserOutputData(resolvedUser),
	})
	return nil
}

func (u *googleMobileLoginUsecase) Poll(state string) (*GoogleMobileSessionResult, bool) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	session, found := u.sessions[state]
	if !found || time.Since(session.createdAt) > googleMobileSessionTTL {
		delete(u.sessions, state)
		return nil, false
	}

	delete(u.sessions, state)
	return &session.result, true
}

func (u *googleMobileLoginUsecase) storeError(state string, appErr apperrors.ApplicationError) apperrors.ApplicationError {
	u.store(state, GoogleMobileSessionResult{
		Status:  MobileSessionError,
		Message: "No se pudo completar el inicio de sesión con Google",
	})
	return appErr
}

func (u *googleMobileLoginUsecase) store(state string, result GoogleMobileSessionResult) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	for key, session := range u.sessions {
		if time.Since(session.createdAt) > googleMobileSessionTTL {
			delete(u.sessions, key)
		}
	}
	u.sessions[state] = storedGoogleMobileSession{result: result, createdAt: time.Now()}
}
