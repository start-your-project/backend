package constants

import (
	"errors"
	"main/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

const (
	DefaultImage           = "default_avatar.webp"
	UserObjectsBucketName  = "avatars"
	SessionRequired        = "Session required"
	UserIsUnauthorized     = "User is unauthorized"
	NoRequestID            = "No RequestID in context"
	UserIsLoggedOut        = "User is logged out"
	UserCanBeLoggedIn      = "User can be logged in"
	EmailConfirmed         = "Email confirmed"
	FileTypeIsNotSupported = "File type is not supported"
	ProfileIsEdited        = "Profile is edited"
	LikeIsAdded            = "Like is added"
	LikeIsRemoved          = "Like is removed"
	Finished               = "Finished"
	Canceled               = "Canceled"
)

var (
	ErrLetter                = errors.New("at least one letter is required")
	ErrNum                   = errors.New("at least one digit is required")
	ErrCount                 = errors.New("at least eight characters long is required")
	ErrBan                   = errors.New("password uses unavailable symbols")
	ErrWrongData             = errors.New("wrong data")
	ErrEmailAlreadyConfirmed = errors.New("email is already confirmed")
	ErrEmailIsNotUnique      = errors.New("email is not unique")
	ErrNoEmailLink           = errors.New("no email link")
	ErrEmailIsNotConfirmed   = errors.New("email is not confirmed")
	ErrTryAgain              = errors.New("error try again")
)

const (
	SignupURL       = "/api/v1/signup"
	LoginURL        = "/api/v1/login"
	LogoutURL       = "/api/v1/logout"
	ConfirmEmailURL = "/confirm"
	ProfileURL      = "/api/v1/profile"
	EditURL         = "/api/v1/edit"
	AvatarURL       = "/api/v1/avatar"
	CsrfURL         = "/api/v1/csrf"
	GetTechnologies = "/api/v1/technologies"
	AddLikeURL      = "/api/v1/like"
	RemoveLikeURL   = "/api/v1/dislike"
	LikesURL        = "/api/v1/likes"
	TopPosition     = "/api/v1/top"
	Recommendation  = "/api/v1/recommend"
	Professions     = "/api/v1/professions"
	GetProfessions  = "/api/v1/list"
	Resume          = "/api/v1/resume"
	FinishURL       = "/api/v1/finish"
	CancelURL       = "/api/v1/cancel"
	FinishedURL     = "/api/v1/finished"
	LetterURL       = "/api/v1/letter"
)

var (
	ImageTypes = map[string]interface{}{
		"image/jpeg": nil,
		"image/png":  nil,
	}
)

func RespError(ctx echo.Context, logger *zap.SugaredLogger, requestID, errorMsg string, status int) error {
	logger.Error(
		zap.String("ID", requestID),
		zap.String("ERROR", errorMsg),
		zap.Int("ANSWER STATUS", status),
	)
	resp, err := easyjson.Marshal(&models.Response{
		Status:  status,
		Message: errorMsg,
	})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSONBlob(status, resp)
}

func DefaultUserChecks(ctx echo.Context, logger *zap.SugaredLogger) (int64, string, error) {
	requestID, okey := ctx.Get("REQUEST_ID").(string)
	if !okey {
		err := RespError(ctx, logger, requestID, NoRequestID, http.StatusInternalServerError)
		if err != nil {
			return 0, "", err
		}
		return 0, "", errors.New("")
	}

	userID, ok := ctx.Get("USER_ID").(int64)
	if !ok {
		err := RespError(ctx, logger, requestID, SessionRequired, http.StatusBadRequest)
		if err != nil {
			return 0, "", err
		}
		return 0, "", errors.New("")
	}

	if userID == -1 {
		err := RespError(ctx, logger, requestID, UserIsUnauthorized, http.StatusUnauthorized)
		if err != nil {
			return 0, "", err
		}
		return userID, "", errors.New("")
	}
	return userID, requestID, nil
}
