package user

import (
	"database/sql"
	"fmt"
	"maple/internal/api"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeleteUserResponse struct {
	Success bool `json:"success"`
}

func deleteUser(ctx *gin.Context) {
	a := api.Get(ctx)

	// Get user ID from URL parameter
	userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidParameter.MakeJSON("invalid user id"))
		return
	}

	// First, get the user to ensure they exist
	user, err := a.Queries.GetUserById(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.AbortWithStatusJSON(http.StatusNotFound, perrors.NotFound.MakeJSON("user not found"))
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	if user.ID != userID {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidParameter.MakeJSON("user id does not match"))
		return
	}

	// Delete user from local database
	err = a.Queries.DeleteUser(ctx, schema.DeleteUserParams{ID: userID})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, DeleteUserResponse{Success: true})
}

func deleteSurgeUser(api *api.MapleAPI, userID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/users/%s", api.Config.Surge.URL, userID), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Surge-Service-Key", api.Config.Surge.ServiceKey)

	result, err := api.SurgeHTTP.Do(req)
	if err != nil {
		return err
	}
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete surge user: status code %d", result.StatusCode)
	}

	return nil
}
