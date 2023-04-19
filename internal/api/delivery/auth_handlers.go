package delivery

import (
	"context"
	"main/internal/constants"
	proto "main/internal/microservices/auth/proto"
	"main/internal/models"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
	"github.com/stroiman/go-automapper"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authHandler struct {
	logger *zap.SugaredLogger

	authMicroservice proto.AuthClient
}

// nolint:golint
func NewAuthHandler(logger *zap.SugaredLogger, auth proto.AuthClient) *authHandler {
	return &authHandler{authMicroservice: auth, logger: logger}
}

func (a *authHandler) Register(router *echo.Echo) {
	router.POST(constants.SignupURL, a.SignUp())
	router.POST(constants.LoginURL, a.LogIn())
	router.DELETE(constants.LogoutURL, a.LogOut())
	router.GET(constants.ConfirmEmailURL, a.ConfirmEmail())
}

// nolint:cyclop
func (a *authHandler) ParseError(ctx echo.Context, requestID string, err error) error {
	if getErr, ok := status.FromError(err); ok {
		// nolint:exhaustive
		switch getErr.Code() {
		case codes.Internal:
			a.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError),
			)
			resp, errMarshal := easyjson.Marshal(&models.Response{
				Status:  http.StatusInternalServerError,
				Message: getErr.Message(),
			})
			if errMarshal != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusInternalServerError, resp)
		case codes.NotFound:
			a.logger.Info(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusNotFound),
			)
			resp, errMarshal := easyjson.Marshal(&models.Response{
				Status:  http.StatusNotFound,
				Message: getErr.Message(),
			})
			if errMarshal != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusNotFound, resp)
		case codes.InvalidArgument:
			a.logger.Info(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusBadRequest),
			)
			resp, errMarshal := easyjson.Marshal(&models.Response{
				Status:  http.StatusBadRequest,
				Message: getErr.Message(),
			})
			if errMarshal != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusBadRequest, resp)
		case codes.Unavailable:
			a.logger.Info(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError),
			)
			resp, errMarshal := easyjson.Marshal(&models.Response{
				Status:  http.StatusInternalServerError,
				Message: getErr.Message(),
			})
			if errMarshal != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusInternalServerError, resp)
		case codes.Unauthenticated:
			a.logger.Info(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusUnauthorized),
			)
			resp, err := easyjson.Marshal(&models.Response{
				Status:  http.StatusUnauthorized,
				Message: getErr.Message(),
			})
			if err != nil {
				return ctx.NoContent(http.StatusUnauthorized)
			}
			return ctx.JSONBlob(http.StatusUnauthorized, resp)
		}

	}
	return nil
}

func (a *authHandler) ConfirmEmail() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestID, ok := ctx.Get("REQUEST_ID").(string)
		if !ok {
			a.logger.Error(
				zap.String("ERROR", constants.NoRequestID),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			resp, err := easyjson.Marshal(&models.Response{
				Status:  http.StatusInternalServerError,
				Message: constants.NoRequestID,
			})
			if err != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusInternalServerError, resp)
		}

		hash := strings.ReplaceAll(ctx.Request().Header.Get("Req"), "/confirm?hash=", "")

		a.logger.Info(hash)

		data := &proto.Hash{Hash: hash}
		_, err := a.authMicroservice.ConfirmEmail(context.Background(), data)
		if err != nil {
			return a.ParseError(ctx, requestID, err)
		}

		a.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)

		resp, err := easyjson.Marshal(&models.Response{
			Status:  http.StatusMovedPermanently,
			Message: constants.EmailConfirmed,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (a *authHandler) LogIn() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userData := models.LogInUserDTO{
			Email:    "",
			Password: "",
		}

		requestID, ok := ctx.Get("REQUEST_ID").(string)
		if !ok {
			a.logger.Error(
				zap.String("ERROR", constants.NoRequestID),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			resp, err := easyjson.Marshal(&models.Response{
				Status:  http.StatusInternalServerError,
				Message: constants.NoRequestID,
			})
			if err != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusInternalServerError, resp)
		}

		if err := ctx.Bind(&userData); err != nil {
			a.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError),
			)
			resp, errMarshal := easyjson.Marshal(&models.Response{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			if errMarshal != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusInternalServerError, resp)
		}

		data := &proto.LogInData{
			Email:    "",
			Password: "",
		}
		automapper.MapLoose(userData, data)

		session, err := a.authMicroservice.LogIn(context.Background(), data)
		if err != nil {
			return a.ParseError(ctx, requestID, err)
		}

		cookie := http.Cookie{
			Name:       "Session_cookie",
			Value:      session.Cookie,
			Path:       "",
			Domain:     "",
			Expires:    time.Now().Add(30 * 24 * time.Hour),
			RawExpires: "",
			MaxAge:     0,
			Secure:     false,
			HttpOnly:   true,
			SameSite:   0,
			Raw:        "",
			Unparsed:   nil,
		}

		ctx.SetCookie(&cookie)

		a.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)
		resp, err := easyjson.Marshal(&models.Response{
			Status:  http.StatusOK,
			Message: constants.UserCanBeLoggedIn,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (a *authHandler) SignUp() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userData := models.CreateUserDTO{
			Name:     "",
			Email:    "",
			Password: "",
		}

		requestID, ok := ctx.Get("REQUEST_ID").(string)
		if !ok {
			a.logger.Error(
				zap.String("ERROR", constants.NoRequestID),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			resp, err := easyjson.Marshal(&models.Response{
				Status:  http.StatusInternalServerError,
				Message: constants.NoRequestID,
			})
			if err != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusInternalServerError, resp)
		}

		if err := ctx.Bind(&userData); err != nil {
			a.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError),
			)
			resp, errMarshal := easyjson.Marshal(&models.Response{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			if errMarshal != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusInternalServerError, resp)
		}

		data := &proto.SignUpData{
			Name:     "",
			Email:    "",
			Password: "",
			Date:     "",
		}
		automapper.MapLoose(userData, data)
		link, err := a.authMicroservice.SignUp(context.Background(), data)
		if err != nil {
			return a.ParseError(ctx, requestID, err)
		}

		a.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusCreated),
		)

		resp, err := easyjson.Marshal(&models.Response{
			Status:  http.StatusOK,
			Message: link.Link,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (a *authHandler) LogOut() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestID, ok := ctx.Get("REQUEST_ID").(string)
		if !ok {
			a.logger.Error(
				zap.String("ERROR", constants.NoRequestID),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			resp, err := easyjson.Marshal(&models.Response{
				Status:  http.StatusInternalServerError,
				Message: constants.NoRequestID,
			})
			if err != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusInternalServerError, resp)
		}

		cookie, err := ctx.Cookie("Session_cookie")
		if err != nil {
			a.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError),
			)

			resp, errMarshal := easyjson.Marshal(&models.Response{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			if errMarshal != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusInternalServerError, resp)
		}

		data := &proto.Cookie{Cookie: cookie.Value}
		_, err = a.authMicroservice.LogOut(context.Background(), data)
		if err != nil {
			return a.ParseError(ctx, requestID, err)
		}

		cookie.Expires = time.Now().AddDate(0, 0, -1)
		ctx.SetCookie(cookie)

		a.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)
		resp, err := easyjson.Marshal(&models.Response{
			Status:  http.StatusOK,
			Message: constants.UserIsLoggedOut,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}
