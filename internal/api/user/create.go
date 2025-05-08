package user

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"maple/internal/api"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/godruoyi/go-snowflake"
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
	ID *string `json:"id"`

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

	if body.Nickname != "" {
		if len(body.Nickname) > 32 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.FailedValidate.MakeJSON("nickname must be 32 characters or less"))
			return
		}

		// Check for duplicate nickname
		existingUser, err := a.Queries.GetUserByNickname(ctx, body.Nickname)
		if err != nil && err != sql.ErrNoRows {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
			return
		}
		fmt.Println(existingUser)
		if existingUser != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.FailedValidate.MakeJSON("nickname already exists"))
			return
		}
	}


	createUserResponse, err := createSurgeUser(a, body.Username, body.Password)
	if err != nil || createUserResponse.ID == nil || *createUserResponse.ID == "0" {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedAPI.MakeJSON(err.Error()))
		return
	}

	fmt.Printf("body %+v\n", body)
	fmt.Printf("createUserResponse: %+v\n", createUserResponse)

	uintedId, err := strconv.ParseUint(*createUserResponse.ID, 10, 64)
	snowflakeId := snowflake.ParseID(uintedId)

	mapleUser, err := a.Queries.CreateUser(ctx, schema.CreateUserParams{
		ID:        int64(uintedId),
		Nickname:  body.Nickname,
		CreatedAt: snowflakeId.GenerateTime(),
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, CreateUserResponse{
		SurgeResponse: createUserResponse,
		User:          mapleUser,
	})
}

func createSurgeUser(api *api.MapleAPI, username string, password string) (*SurgeUserResponse, error) {
	bodyJson, err := json.Marshal(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: username,
		Password: password,
	})
	fmt.Printf("[user/create#createSurgeUser] Reached P1\n")

	bodyBuffer := bytes.NewBuffer(bodyJson)
	fmt.Printf("[user/create#createSurgeUser] Reached P2\n")

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/sign_up/credentials", api.Config.Surge.URL), bodyBuffer)
	if err != nil {
		return nil, err
	}
	fmt.Printf("[user/create#createSurgeUser] Reached P3\n")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Surge-Service-Key", api.Config.Surge.ServiceKey)

	result, err := api.SurgeHTTP.Do(req)
	if err != nil {
		fmt.Printf("[user/create#createSurgeUser] Error: %+v\n", result)
		fmt.Printf("[user/create#createSurgeUser] ------- DUMP START ------\n")
		fmt.Printf("[user/create#createSurgeUser] bodyJson %+v\n", bodyJson)
		fmt.Printf("[user/create#createSurgeUser] bodyBuffer %+v\n", bodyBuffer)
		fmt.Printf("[user/create#createSurgeUser] req %+v\n", req)
		fmt.Printf("[user/create#createSurgeUser] result %+v\n", result)
		fmt.Printf("[user/create#createSurgeUser] -------  DUMP END  ------\n")
		if result != nil {
			return unmarshalledBody(result)
		} else {
			return nil, err
		}
	}
	fmt.Printf("[user/create#createSurgeUser] Reached P4\n")

	return unmarshalledBody(result)
}

func unmarshalledBody(result *http.Response) (*SurgeUserResponse, error) {
	defer result.Body.Close()

	unmarshalled := new(SurgeUserResponse)

	if err := json.NewDecoder(result.Body).Decode(unmarshalled); err != nil {
		return nil, err
	}

	return unmarshalled, nil
}
