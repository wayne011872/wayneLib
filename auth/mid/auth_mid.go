package mid

import (
	"github.com/gin-gonic/gin"
	"github.com/wayne011872/api-toolkit/auth"
	"github.com/wayne011872/api-toolkit/errors"
)

func NewGinInterAuthMid(address string) auth.GinAuthMidInter {
	// authSDK := authClient.New(address)
	return &interAuthMiddle{
		auth: auth.NewGinBearAuthMid(true),
		// authSDK: authSDK,
	}
}

func (lm *interAuthMiddle) GetName() string {
	return "auth"
}

type interAuthMiddle struct {
	errors.CommonApiErrorHandler
	auth auth.GinAuthMidInter
	// authSDK authClient.AuthClient
}

func (am *interAuthMiddle) AddAuthPath(path string, method string, isAuth bool, group []auth.ApiPerm) {
	am.auth.AddAuthPath(path, method, isAuth, group)
}
func (am *interAuthMiddle) IsAuth(path string, method string) bool {
	return am.auth.IsAuth(path, method)
}
func (am *interAuthMiddle) HasPerm(path, method string, perm []string) bool {
	return am.auth.HasPerm(path, method, perm)
}

func (am *interAuthMiddle) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: implement
		c.Next()
	}
}
