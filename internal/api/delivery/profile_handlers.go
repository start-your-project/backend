package delivery

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"main/internal/constants"
	"main/internal/csrf"
	profile "main/internal/microservices/profile/proto"
	"main/internal/models"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"code.sajari.com/docconv"
	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type profileHandler struct {
	logger *zap.SugaredLogger

	profileMicroservice profile.ProfileClient
}

// nolint:golint
func NewProfileHandler(logger *zap.SugaredLogger, profile profile.ProfileClient) *profileHandler {
	return &profileHandler{profileMicroservice: profile, logger: logger}
}

func (p *profileHandler) Register(router *echo.Echo) {
	router.GET(constants.ProfileURL, p.GetUserProfile())
	router.PUT(constants.EditURL, p.EditProfile())
	router.PUT(constants.AvatarURL, p.EditAvatar())
	router.GET(constants.CsrfURL, p.GetCsrf())
	router.POST(constants.AddLikeURL, p.AddLike())
	router.DELETE(constants.RemoveLikeURL, p.RemoveLike())
	router.GET(constants.LikesURL, p.GetFavorites())
	router.POST(constants.Resume, p.Resume())
	router.POST(constants.FinishURL, p.Finish())
	router.DELETE(constants.CancelURL, p.Cancel())
	router.POST(constants.FinishedURL, p.GetFinished())
	router.POST(constants.LetterURL, p.Letter())
}

// nolint:cyclop
func (p *profileHandler) ParseError(ctx echo.Context, requestID string, err error) error {
	if getErr, ok := status.FromError(err); ok {
		// nolint:exhaustive
		switch getErr.Code() {
		case codes.Unavailable:
			p.logger.Info(
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
		case codes.InvalidArgument:
			p.logger.Info(
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
		case codes.PermissionDenied:
			p.logger.Info(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusForbidden),
			)
			resp, errMarshal := easyjson.Marshal(&models.Response{
				Status:  http.StatusForbidden,
				Message: getErr.Message(),
			})
			if errMarshal != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusForbidden, resp)
		case codes.AlreadyExists:
			p.logger.Info(
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
		default:
			p.logger.Error(
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
		}
	}
	return nil
}

func (p *profileHandler) GetUserProfile() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID, requestID, err := constants.DefaultUserChecks(ctx, p.logger)
		if err != nil {
			return err
		}
		data := &profile.UserID{ID: userID}
		userData, err := p.profileMicroservice.GetUserProfile(context.Background(), data)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		p.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
			userData,
		)

		profileData := models.ProfileUserDTO{
			Name:   userData.Name,
			Email:  userData.Email,
			Avatar: userData.Avatar,
		}

		sanitizer := bluemonday.UGCPolicy()
		profileData.Name = sanitizer.Sanitize(profileData.Name)
		resp, err := easyjson.Marshal(&models.ResponseUserProfile{
			Status:   http.StatusOK,
			UserData: &profileData,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

// nolint:cyclop
func (p *profileHandler) EditAvatar() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID, requestID, err := constants.DefaultUserChecks(ctx, p.logger)
		if err != nil {
			return err
		}

		file, err := ctx.FormFile("file")
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}

		src, err := file.Open()
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}

		buffer := make([]byte, file.Size)
		_, err = src.Read(buffer)
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}
		err = src.Close()
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}

		file, err = ctx.FormFile("file")
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}
		src, err = file.Open()
		defer func(src multipart.File) {
			err = src.Close()
			if err != nil {
				return
			}
		}(src)
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}

		fileType := http.DetectContentType(buffer)

		// Validate File Type
		if _, ex := constants.ImageTypes[fileType]; !ex {
			return constants.RespError(ctx, p.logger, requestID, constants.FileTypeIsNotSupported, http.StatusBadRequest)
		}

		uploadData := &profile.UploadInputFile{
			ID:          userID,
			File:        buffer,
			Size:        file.Size,
			ContentType: fileType,
		}

		fileName, err := p.profileMicroservice.UploadAvatar(context.Background(), uploadData)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		editData := &profile.EditAvatarData{
			ID:     userID,
			Avatar: fileName.Name,
		}

		_, err = p.profileMicroservice.EditAvatar(context.Background(), editData)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		p.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)
		resp, err := easyjson.Marshal(&models.Response{
			Status:  http.StatusOK,
			Message: constants.ProfileIsEdited,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (p *profileHandler) EditProfile() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID, requestID, err := constants.DefaultUserChecks(ctx, p.logger)
		if err != nil {
			return err
		}

		userData := models.EditProfileDTO{Name: "", Password: ""}

		if err = ctx.Bind(&userData); err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusBadRequest)
		}

		data := &profile.EditProfileData{
			ID:       userID,
			Name:     userData.Name,
			Password: userData.Password,
		}

		_, err = p.profileMicroservice.EditProfile(context.Background(), data)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		p.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)

		resp, err := easyjson.Marshal(&models.Response{
			Status:  http.StatusOK,
			Message: constants.ProfileIsEdited,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (p *profileHandler) GetCsrf() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestID, ok := ctx.Get("REQUEST_ID").(string)
		if !ok {
			return constants.RespError(ctx, p.logger, requestID, constants.NoRequestID, http.StatusInternalServerError)
		}

		cookie, err := ctx.Cookie("Session_cookie")
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}

		token, err := csrf.Tokens.Create(cookie.Value, time.Now().Add(time.Hour).Unix())

		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}

		p.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)
		resp, err := easyjson.Marshal(&models.Response{
			Status:  http.StatusOK,
			Message: token,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (p *profileHandler) getLikeData(ctx echo.Context) (*profile.LikeData, string, error) {
	userID, requestID, err := constants.DefaultUserChecks(ctx, p.logger)
	if err != nil {
		return nil, "", err
	}

	positionName := models.LikeDTO{Name: ""}

	if err = ctx.Bind(&positionName); err != nil {
		return nil, "", constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
	}

	return &profile.LikeData{
		UserID:       userID,
		PositionName: positionName.Name,
	}, requestID, nil
}

// nolint:dupl
func (p *profileHandler) AddLike() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		data, requestID, err := p.getLikeData(ctx)
		if err != nil {
			return err
		}

		_, err = p.profileMicroservice.AddLike(context.Background(), data)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		p.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)

		resp, err := easyjson.Marshal(&models.Response{
			Status:  http.StatusOK,
			Message: constants.LikeIsAdded,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

// nolint:dupl
func (p *profileHandler) RemoveLike() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		data, requestID, err := p.getLikeData(ctx)
		if err != nil {
			return err
		}

		_, err = p.profileMicroservice.RemoveLike(context.Background(), data)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		p.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)
		resp, err := easyjson.Marshal(&models.Response{
			Status:  http.StatusOK,
			Message: constants.LikeIsRemoved,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (p *profileHandler) GetFavorites() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID, requestID, err := constants.DefaultUserChecks(ctx, p.logger)
		if err != nil {
			return err
		}

		data := &profile.UserID{ID: userID}
		favorites, err := p.profileMicroservice.GetFavorites(context.Background(), data)

		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		p.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)

		favoritesResponse := make([]models.Favorite, 0)
		for _, favorite := range favorites.Favorite {
			favoritesResponse = append(favoritesResponse, models.Favorite{
				ID:            favorite.PositionID,
				Name:          favorite.Name,
				CountAll:      favorite.CountAll,
				CountFinished: favorite.CountFinished,
			})
		}

		resp, err := easyjson.Marshal(&models.ResponseFavorites{
			Status:        http.StatusOK,
			FavoritesData: favoritesResponse,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

// nolint:cyclop
func (p *profileHandler) Resume() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestID, ok := ctx.Get("REQUEST_ID").(string)
		if !ok {
			p.logger.Error(
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

		req := models.ResumeRequest{
			CvText: "",
			NTech:  10,
			Role:   "",
		}

		if err := ctx.Bind(&req); err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusBadRequest)
		}

		file, err := ctx.FormFile("file")
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}

		src, err := file.Open()
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}

		pdf, _, err := docconv.ConvertPDF(src)
		if err != nil {
			resp, errMarshal := easyjson.Marshal(&models.Response{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			if errMarshal != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusInternalServerError, resp)
		}

		req.CvText = pdf

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		client := &http.Client{}
		request, err := http.NewRequest("GET", os.Getenv("HOST_SEARCH"), nil)
		if err != nil {
			resp, errMarshal := easyjson.Marshal(&models.Response{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			if errMarshal != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusInternalServerError, resp)
		}
		q := request.URL.Query()
		q.Add("query", req.Role)
		request.URL.RawQuery = q.Encode()
		request.Header.Add("User-Agent", "Localhost 1.0") // добавляем заголовок User-Agent
		resRole, err := client.Do(request)
		if err != nil {
			resp, errMarshal := easyjson.Marshal(&models.Response{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			if errMarshal != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusInternalServerError, resp)
		}

		defer resRole.Body.Close()

		scanner := bufio.NewScanner(resRole.Body)
		scanner.Scan()
		singleRequestForm := scanner.Text()

		if resRole.Status != "200 OK" {
			stat := resRole.Status[:3]
			statusInt, errAt := strconv.Atoi(stat)
			if errAt != nil {
				resp, errMarshal := easyjson.Marshal(&models.Response{
					Status:  http.StatusInternalServerError,
					Message: errAt.Error(),
				})
				if errMarshal != nil {
					return ctx.NoContent(http.StatusInternalServerError)
				}
				return ctx.JSONBlob(http.StatusInternalServerError, resp)
			} else {
				return ctx.JSONBlob(statusInt, []byte(singleRequestForm))
			}
		}

		profession := &models.Profession{
			Profession: "",
			InBase:     "",
		}
		err = easyjson.Unmarshal([]byte(singleRequestForm), profession)
		if err != nil {
			p.logger.Error(
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			resp, errMarshal := easyjson.Marshal(&models.Response{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			if errMarshal != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusInternalServerError, resp)
		}

		if profession.InBase == "0" {
			p.logger.Error(
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			resp, errMarshal := easyjson.Marshal(&models.Response{
				Status:  http.StatusNotAcceptable,
				Message: constants.NotInBase,
			})
			if errMarshal != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusNotAcceptable, resp)
		}

		req.Role = profession.Profession

		json, err := req.MarshalJSON()
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}

		//nolint:bodyclose
		response, err := http.Post(os.Getenv("HOST_RESUME"), "application/json", bytes.NewBuffer(json))
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				p.logger.Error(err)
			}
		}(response.Body)
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}

		if response.Status != "200 OK" {
			stat := response.Status[:3]
			statusInt, errAtoi := strconv.Atoi(stat)
			if errAtoi != nil {
				resp, errMarshal := easyjson.Marshal(&models.Response{
					Status:  http.StatusInternalServerError,
					Message: errAtoi.Error(),
				})
				if errMarshal != nil {
					return ctx.NoContent(http.StatusInternalServerError)
				}
				return ctx.JSONBlob(http.StatusInternalServerError, resp)
			} else {
				return ctx.JSONBlob(statusInt, body)
			}
		}

		recommends := &models.ResponseResume{
			Status:    http.StatusOK,
			Recommend: make([]models.Recommend, 0),
		}
		err = easyjson.Unmarshal(body, recommends)
		if err != nil {
			p.logger.Error(
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			resp, errMarshal := easyjson.Marshal(&models.Response{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			if errMarshal != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusInternalServerError, resp)
		}

		resp, errMarshal := easyjson.Marshal(recommends)
		if errMarshal != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

// nolint:dupl
func (p *profileHandler) Finish() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		data, requestID, err := p.getLikeData(ctx)
		if err != nil {
			return err
		}

		_, err = p.profileMicroservice.Finish(context.Background(), data)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		p.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)

		resp, err := easyjson.Marshal(&models.Response{
			Status:  http.StatusOK,
			Message: constants.Finished,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

// nolint:dupl
func (p *profileHandler) Cancel() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		data, requestID, err := p.getLikeData(ctx)
		if err != nil {
			return err
		}

		_, err = p.profileMicroservice.Cancel(context.Background(), data)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		p.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)
		resp, err := easyjson.Marshal(&models.Response{
			Status:  http.StatusOK,
			Message: constants.Canceled,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (p *profileHandler) GetFinished() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		data, requestID, err := p.getLikeData(ctx)
		if err != nil {
			return err
		}

		finished, err := p.profileMicroservice.GetFinished(context.Background(), data)

		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		p.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)

		resp, err := easyjson.Marshal(&models.ResponseFinished{
			Status:   http.StatusOK,
			Finished: finished.Names,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

// nolint:cyclop
func (p *profileHandler) Letter() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestID, ok := ctx.Get("REQUEST_ID").(string)
		if !ok {
			p.logger.Error(
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

		link := models.LinkDTO{
			Link: "",
		}

		if err := ctx.Bind(&link); err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}

		file, err := ctx.FormFile("file")
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}

		src, err := file.Open()
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}

		pdf, _, err := docconv.ConvertPDF(src)
		if err != nil {
			resp, errMarshal := easyjson.Marshal(&models.Response{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			if errMarshal != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusInternalServerError, resp)
		}

		req := models.LetterRequest{
			Resume:  pdf,
			Vacancy: link.Link,
		}

		json, err := req.MarshalJSON()
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		//nolint:bodyclose
		response, err := http.Post(os.Getenv("HOST_LETTER"), "application/json", bytes.NewBuffer(json))
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}

		if response.Status != "200 OK" {
			status := response.Status[:3]
			statusInt, errAtoi := strconv.Atoi(status)
			if errAtoi != nil {
				resp, errMarshal := easyjson.Marshal(&models.Response{
					Status:  http.StatusInternalServerError,
					Message: errAtoi.Error(),
				})
				if errMarshal != nil {
					return ctx.NoContent(http.StatusInternalServerError)
				}
				return ctx.JSONBlob(http.StatusInternalServerError, resp)
			} else {
				return ctx.JSONBlob(statusInt, body)
			}
		}

		letter := &models.ResponseLetter{
			Status:      http.StatusOK,
			CoverLetter: "",
		}
		err = easyjson.Unmarshal(body, letter)
		if err != nil {
			p.logger.Error(
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			resp, errMarshal := easyjson.Marshal(&models.Response{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			if errMarshal != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusInternalServerError, resp)
		}

		resp, errMarshal := easyjson.Marshal(letter)
		if errMarshal != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}
