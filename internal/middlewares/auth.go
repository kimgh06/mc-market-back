package middlewares

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"maple/internal/api"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
	"strconv"
	"strings"
)

type SurgeTokenClaims struct {
	jwt.RegisteredClaims

	Email    *string `json:"email"`
	Username *string `json:"username"`
}

func (c *SurgeTokenClaims) ParseID() *uint64 {
	id, err := strconv.ParseUint(c.Subject, 10, 64)
	if err != nil {
		return nil
	}

	return &id
}

func getRequestTokenString(c *gin.Context) (string, error) {
	authorization := c.GetHeader("Authorization")

	var found bool
	if authorization, found = strings.CutPrefix(authorization, "bearer "); !found {
		if authorization, found = strings.CutPrefix(authorization, "Bearer "); !found {
			return "", perrors.UnauthorizedMissingBearer
		}
	}

	return authorization, nil
}

// RequireAuthentication resolves user from request and aborts if the user is unauthorized
func RequireAuthentication(a *api.MapleAPI) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := getRequestTokenString(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, perrors.UnauthorizedMissingBearer.MakeJSON())
			return
		}

		parsed, err := jwt.Parse(tokenString, a.JWKS.GetKeyFuncForParse())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, perrors.UnauthorizedFailedParse.MakeJSON(err.Error()))
			return
		}

		// Parse subject to UUID
		subject, err := parsed.Claims.GetSubject()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, perrors.UnauthorizedFailedSubject.MakeJSON(err.Error()))
			return
		}

		subjectID, err := strconv.ParseUint(subject, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, perrors.UnauthorizedFailedSubject.MakeJSON(err.Error()))
			return
		}

		user, err := a.Queries.GetUserById(c, int64(subjectID))
		if err != nil {
			if !errors.Is(sql.ErrNoRows, err) {
				c.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
				return
			}
		}

		c.Set(api.ContextTokenKey.String(), parsed)
		c.Set(api.ContextClaimKey.String(), decodeClaim(parsed))
		c.Set(api.ContextUserKey.String(), user)

		c.Next()
	}
}

// UseAuthentication resolves user but doesn't abort even if the user is unauthorized
func UseAuthentication(a *api.MapleAPI) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := getRequestTokenString(c)
		if err != nil {
			return
		}

		parsed, err := jwt.Parse(tokenString, a.JWKS.GetKeyFuncForParse())
		if err != nil {
			return
		}

		// Parse subject to UUID
		subject, err := parsed.Claims.GetSubject()
		if err != nil {
			return
		}

		subjectID, err := strconv.ParseUint(subject, 10, 64)
		if err != nil {
			return
		}

		user, err := a.Queries.GetUserById(c, int64(subjectID))
		if err != nil {
			if !errors.Is(sql.ErrNoRows, err) {
				return
			}
		}

		c.Set(api.ContextTokenKey.String(), parsed)
		c.Set(api.ContextClaimKey.String(), decodeClaim(parsed))
		c.Set(api.ContextUserKey.String(), user)

		c.Next()
	}
}

// GetToken reads the JWT token from the context.
func GetToken(ctx *gin.Context) *jwt.Token {
	obj, exists := ctx.Get(api.ContextTokenKey.String())
	if obj == nil || !exists {
		return nil
	}

	return obj.(*jwt.Token)
}

func GetClaims(ctx *gin.Context) *SurgeTokenClaims {
	obj, exists := ctx.Get(api.ContextClaimKey.String())
	if obj == nil || !exists {
		return nil
	}

	return obj.(*SurgeTokenClaims)
}

func GetUser(ctx *gin.Context) *schema.User {
	obj, exists := ctx.Get(api.ContextUserKey.String())
	if obj == nil || !exists {
		return nil
	}
	return obj.(*schema.User)
}

func GetUserID(ctx *gin.Context) uint64 {
	return uint64(GetUser(ctx).ID)
}

func decodeClaim(token *jwt.Token) *SurgeTokenClaims {
	marshalled, err := json.Marshal(token.Claims)
	if err != nil {
		return nil
	}

	var decoded SurgeTokenClaims

	if json.Unmarshal(marshalled, &decoded) != nil {
		return nil
	}

	return &decoded
}
