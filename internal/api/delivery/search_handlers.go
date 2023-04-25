package delivery

import (
	"bufio"
	"context"
	"crypto/tls"
	"main/internal/constants"
	search "main/internal/microservices/search/proto"
	"main/internal/models"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

type searchHandler struct {
	logger *zap.SugaredLogger

	searchMicroservice search.SearchClient
}

// nolint:golint
func NewSearchHandler(logger *zap.SugaredLogger, search search.SearchClient) *searchHandler {
	return &searchHandler{searchMicroservice: search, logger: logger}
}

func (a *searchHandler) Register(router *echo.Echo) {
	router.GET(constants.GetTechnologies, a.GetTechnologies())
	router.GET(constants.TopPosition, a.GetTop())
	router.GET(constants.Recommendation, a.Recommendation())
	router.POST(constants.Professions, a.Professions())
	router.POST(constants.GetProfessions, a.GetProfessions())
	router.GET(constants.GetTechSearch, a.TechSearch())
}

// nolint:cyclop
func (a *searchHandler) GetTechnologies() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		_, ok := ctx.Get("REQUEST_ID").(string)
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

		searchText := ctx.QueryParam("search_text")

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		res, err := http.Get(os.Getenv("HOST_SEARCH") + searchText)
		if err != nil {
			a.logger.Error(
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			resp, err := easyjson.Marshal(&models.Response{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			if err != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusInternalServerError, resp)
		}
		defer res.Body.Close()

		scanner := bufio.NewScanner(res.Body)
		scanner.Scan()
		singleRequestForm := scanner.Text()

		if res.Status != "200 OK" {
			status := res.Status[:3]
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
				return ctx.JSONBlob(statusInt, []byte(singleRequestForm))
			}
		}

		profession := &models.Profession{
			Profession: "",
			InBase:     "",
		}
		err = easyjson.Unmarshal([]byte(singleRequestForm), profession)
		if err != nil {
			a.logger.Error(
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

		if profession.InBase == "1" {
			data := &search.SearchText{Text: profession.Profession}
			technologies, err := a.searchMicroservice.GetTechnologies(context.Background(), data)
			if err != nil {
				a.logger.Error(
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

			if technologies.Technology != nil {
				positionResult := models.PositionData{
					JobName:          profession.Profession,
					TechnologyNumber: len(technologies.Technology),
					Additional:       make([]models.Technology, 0),
				}
				for _, technology := range technologies.Technology {
					positionResult.Additional = append(positionResult.Additional, models.Technology{
						TechnologyName:  technology.Name,
						Distance:        technology.Distance,
						Professionalism: technology.Professionalism,
						HardSkill:       technology.HardSkill,
					})
				}
				a.logger.Info(positionResult)

				resp, errMarshal := easyjson.Marshal(&models.ResponseTechnologies{
					Status:       http.StatusOK,
					PositionData: positionResult,
				})
				if errMarshal != nil {
					return ctx.NoContent(http.StatusInternalServerError)
				}
				return ctx.JSONBlob(http.StatusOK, resp)
			}

			res, err = http.Get(os.Getenv("HOST_TECHNOLOGIES") + profession.Profession)
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
			defer res.Body.Close()

			scanner = bufio.NewScanner(res.Body)
			scanner.Scan()
			positionRes := scanner.Text()

			if res.Status != "200 OK" {
				status := res.Status[:3]
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
					return ctx.JSONBlob(statusInt, []byte(positionRes))
				}
			}

			technologies, err = a.searchMicroservice.GetTechnologies(context.Background(), data)
			if err != nil {
				a.logger.Error(
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

			if technologies.Technology != nil {
				positionResult := models.PositionData{
					JobName:          profession.Profession,
					TechnologyNumber: len(technologies.Technology),
					Additional:       make([]models.Technology, 0),
				}
				for _, technology := range technologies.Technology {
					positionResult.Additional = append(positionResult.Additional, models.Technology{
						TechnologyName:  technology.Name,
						Distance:        technology.Distance,
						Professionalism: technology.Professionalism,
						HardSkill:       technology.HardSkill,
					})
				}

				a.logger.Info(positionResult)
				resp, errMarshal := easyjson.Marshal(&models.ResponseTechnologies{
					Status:       http.StatusOK,
					PositionData: positionResult,
				})
				if errMarshal != nil {
					return ctx.NoContent(http.StatusInternalServerError)
				}
				return ctx.JSONBlob(http.StatusOK, resp)
			}

			resp, errMarshal := easyjson.Marshal(&models.Response{
				Status:  http.StatusInternalServerError,
				Message: constants.NoTechnologies,
			})
			if errMarshal != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusInternalServerError, resp)
		} else {
			res, err = http.Get(os.Getenv("HOST_TECHNOLOGIES") + profession.Profession)
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
			defer res.Body.Close()

			scanner = bufio.NewScanner(res.Body)
			scanner.Scan()
			positionRes := scanner.Text()

			if res.Status != "200 OK" {
				status := res.Status[:3]
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
					return ctx.JSONBlob(statusInt, []byte(positionRes))
				}
			}

			positionResult := models.PositionData{
				JobName:          profession.Profession,
				TechnologyNumber: 0,
				Additional:       make([]models.Technology, 0),
			}

			err = positionResult.UnmarshalJSON([]byte(positionRes))
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

			resp, errMarshal := easyjson.Marshal(&models.ResponseTechnologies{
				Status:       http.StatusOK,
				PositionData: positionResult,
			})
			if errMarshal != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusOK, resp)
		}
	}
}

func (a *searchHandler) GetTop() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		_, ok := ctx.Get("REQUEST_ID").(string)
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

		positions, err := a.searchMicroservice.GetTop(context.Background(), &search.Empty{})
		if err != nil {
			a.logger.Error(
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

		positionResult := models.ResponseTop{
			Status: http.StatusOK,
			Top:    make([]models.Profession, 0),
		}

		for _, position := range positions.Position {
			positionResult.Top = append(positionResult.Top, models.Profession{
				Profession: position.Name,
			})
		}

		positionResult.Status = http.StatusOK
		resp, err := easyjson.Marshal(&positionResult)
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (a *searchHandler) Recommendation() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		_, ok := ctx.Get("REQUEST_ID").(string)
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

		searchText := ctx.QueryParam("search_text")

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		res, err := http.Get(os.Getenv("HOST_RECOMMEND") + searchText)
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		defer res.Body.Close()

		scanner := bufio.NewScanner(res.Body)
		scanner.Scan()
		singleRequestForm := scanner.Text()

		professions := &models.Professions{Profession: make([]string, 0)}
		err = easyjson.Unmarshal([]byte(singleRequestForm), professions)
		if err != nil {
			a.logger.Error(
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

		resp, errMarshal := easyjson.Marshal(&models.ResponseProfessions{
			Status:      http.StatusOK,
			Professions: professions.Profession,
		})
		if errMarshal != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (a *searchHandler) Professions() echo.HandlerFunc {
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

		techs := models.SearchTechs{SearchText: ""}
		if err := ctx.Bind(&techs); err != nil {
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

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		res, err := http.Get(os.Getenv("HOST_PROFESSIONS") + url.QueryEscape(techs.SearchText))
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
		defer res.Body.Close()

		scanner := bufio.NewScanner(res.Body)
		scanner.Scan()
		respProfessions := scanner.Text()

		a.logger.Info(respProfessions)

		if res.Status != "200 OK" {
			status := res.Status[:3]
			statusInt, errAt := strconv.Atoi(status)
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
				return ctx.JSONBlob(statusInt, []byte(respProfessions))
			}
		}

		professions := &models.RespProfessions{
			Techs:      "",
			JobNumber:  0,
			Additional: make([]models.RespProfession, 0),
		}

		err = easyjson.Unmarshal([]byte(respProfessions), professions)
		if err != nil {
			a.logger.Error(
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

		resp, errMarshal := easyjson.Marshal(&models.ResponseProfessionsWithTechnology{
			Status:      http.StatusOK,
			Professions: professions,
		})
		if errMarshal != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)

	}
}

func (a *searchHandler) GetProfessions() echo.HandlerFunc {
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

		technology := models.Technology{TechnologyName: ""}
		if err := ctx.Bind(&technology); err != nil {
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

		data := search.GetTechnology{Name: technology.TechnologyName}

		positions, err := a.searchMicroservice.GetPositions(context.Background(), &data)
		if err != nil {
			a.logger.Error(
				zap.String("ID", requestID),
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

		positionResult := models.ResponseTop{
			Status:      http.StatusOK,
			Top:         make([]models.Profession, 0),
			TipsToLearn: positions.TipsToLearn,
		}

		for _, position := range positions.Position {
			positionResult.Top = append(positionResult.Top, models.Profession{
				Profession: position.Name,
			})
		}

		positionResult.Status = http.StatusOK
		resp, err := easyjson.Marshal(&positionResult)
		if err != nil {
			a.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			res, errMarshal := easyjson.Marshal(&models.Response{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			if errMarshal != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusInternalServerError, res)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (a *searchHandler) TechSearch() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		_, ok := ctx.Get("REQUEST_ID").(string)
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

		searchText := ctx.QueryParam("search_text")

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		res, err := http.Get(os.Getenv("HOST_TECH") + searchText)
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		defer res.Body.Close()

		scanner := bufio.NewScanner(res.Body)
		scanner.Scan()
		response := scanner.Text()

		if res.Status != "200 OK" {
			status := res.Status[:3]
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
				return ctx.JSONBlob(statusInt, []byte(response))
			}
		}

		techs := &models.Techs{Techs: make([]string, 0)}
		err = easyjson.Unmarshal([]byte(response), techs)
		if err != nil {
			a.logger.Error(
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

		resp, errMarshal := easyjson.Marshal(&models.ResponseTechs{
			Status: http.StatusOK,
			Techs:  techs.Techs,
		})
		if errMarshal != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}
