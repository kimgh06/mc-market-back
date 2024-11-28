package user

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/godruoyi/go-snowflake"
	"io"
	"maple/internal/api"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
	"time"
)

type CreateUserBody struct {
	Nickname string `json:"nickname"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateUserResponse struct {
	SurgeResponse interface{}  `json:"surge_response"`
	User          *schema.User `json:"user"`
}

type SurgeUserResponse struct {
	ID uint64 `json:"id"`

	Email    *string `json:"email"`
	Username *string `json:"username"`

	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	LastSignIn *time.Time `json:"last_sign_in"`
}

func createUser(ctx *gin.Context) {
	a := api.Get(ctx)

	body := CreateUserBody{}
	err := ctx.ShouldBind(&body)
	if err != nil {
		// failed to bind body, abort
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON(err.Error()))
		return
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.FailedValidate.MakeJSON(err.Error()))
		return
	}

	createUserResponse, err := createSurgeUser(a, body.Username, body.Password)
	if err != nil {
		println(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedAPI.WithJSON(createUserResponse, err.Error()))
		return
	}

	marshalled, err := json.Marshal(createUserResponse)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedAPI.WithJSON(createUserResponse, err.Error()))
		return
	}
	var unmarshalled SurgeUserResponse
	err = json.Unmarshal(marshalled, &unmarshalled)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedAPI.WithJSON(createUserResponse, err.Error()))
		return
	}

	snowflakeId := snowflake.ParseID(unmarshalled.ID)

	mapleUser, err := a.Queries.CreateUser(ctx, schema.CreateUserParams{
		ID:        int64(unmarshalled.ID),
		Nickname:  sql.NullString{String: body.Nickname, Valid: body.Nickname != ""},
		CreatedAt: snowflakeId.GenerateTime(),
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, CreateUserResponse{
		SurgeResponse: unmarshalled,
		User:          mapleUser,
	})
}

func createSurgeUser(api *api.MapleAPI, username string, password string) (interface{}, error) {
	bodyJson, err := json.Marshal(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: username,
		Password: password,
	})
	fmt.Printf("[user/create#createSurgeUser] Reached P1")

	bodyBuffer := bytes.NewBuffer(bodyJson)
	fmt.Printf("[user/create#createSurgeUser] Reached P2")

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/sign_up/credentials", api.Config.Surge.URL), bodyBuffer)
	if err != nil {
		return nil, err
	}
	fmt.Printf("[user/create#createSurgeUser] Reached P3")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Surge-Service-Key", api.Config.Surge.ServiceKey)

	result, err := api.SurgeHTTP.Do(req)
	if err != nil {
		fmt.Printf("[user/create#createSurgeUser] Error: %+v\n", result)
		fmt.Printf("[user/create#createSurgeUser] ------- DUMP START ------")
		fmt.Printf("[user/create#createSurgeUser] bodyJson %+v", bodyJson)
		fmt.Printf("[user/create#createSurgeUser] bodyBuffer %+v", bodyBuffer)
		fmt.Printf("[user/create#createSurgeUser] req %+v", req)
		fmt.Printf("[user/create#createSurgeUser] result %+v", result)
		fmt.Printf("[user/create#createSurgeUser] -------  DUMP END  ------")
		if result != nil {
			return unmarshalledBody(result)
		} else {
			return nil, err
		}
	}
	fmt.Printf("[user/create#createSurgeUser] Reached P4")

	return unmarshalledBody(result)
}

func unmarshalledBody(result *http.Response) (interface{}, error) {
	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	var bodyUnmarshalled interface{}

	err = json.Unmarshal(body, &bodyUnmarshalled)
	if err != nil {
		return nil, err
	}

	return bodyUnmarshalled, err
}
