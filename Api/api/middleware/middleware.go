package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"api/api/token"
	"api/config"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type JwtRoleAuth struct {
	enforcer   *casbin.Enforcer
	jwtHandler token.JWTHandler
}

func NewAuth(enforce *casbin.Enforcer) gin.HandlerFunc {
	auth := JwtRoleAuth{
		enforcer: enforce,
	}

	return func(ctx *gin.Context) {
		allow, err := auth.CheckPermission(ctx)
		if err != nil {
			valid, _ := err.(jwt.ValidationError)
			if valid.Errors == jwt.ValidationErrorExpired {
				ctx.AbortWithStatusJSON(http.StatusForbidden, "Invalid token !!!")
			} else {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Access token expired or invalid")
			}
			return
		}

		if !allow {
			ctx.AbortWithStatusJSON(http.StatusForbidden, "Permission denied")
			return
		}

		ctx.Next()
	}
}

func (a *JwtRoleAuth) GetRole(r *gin.Context) (string, error) {
	var (
		claims jwt.MapClaims
		err    error
	)

	jwtToken := r.Request.Header.Get("Authorization")
	if jwtToken == "" || !strings.HasPrefix(jwtToken, "Bearer ") {
		return "unauthorized", nil
	}
	jwtToken = strings.TrimPrefix(jwtToken, "Bearer ")

	a.jwtHandler.Token = jwtToken
	a.jwtHandler.SigningKey = config.Load().TokenKey
	claims, err = a.jwtHandler.ExtractClaims()
	if err != nil {
		log.Println("Error while extracting claims: ", err)
		return "unauthorized", err
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "unauthorized", fmt.Errorf("role claim not found or not a string")
	}

	return role, nil
}

func (a *JwtRoleAuth) CheckPermission(r *gin.Context) (bool, error) {
	role, err := a.GetRole(r)
	if err != nil {
		log.Println("Error while getting role from token: ", err)
		return false, err
	}

	method := r.Request.Method
	path := r.FullPath()

	log.Printf("Checking permissions: role=%s, path=%s, method=%s\n", role, path, method)

	allowed, err := a.enforcer.Enforce(role, path, method)
	if err != nil {
		log.Println("Error while checking permission: ", err)
		return false, err
	}

	return allowed, nil
}
