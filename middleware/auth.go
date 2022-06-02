package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"qxy-dy/serializer"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	ErrorServerBusy = "server is busy"
	ErrorReLogin    = "relogin"
)

type JWTClaims struct {
	jwt.StandardClaims
	Id       uint64 `json:"id"`
	Password string `json:"password"`
	Username string `json:"username"`
}

var (
	Secret     = "liuliumei"    //salt
	ExpireTime = 3600 * 24 * 30 //token expire time
)

// AuthRequired 需要登录
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenstr, ok := c.GetPostForm("token")
		if !ok {
			tokenstr = c.Query("token")
		}
		userInfo, err := verifyAction(tokenstr)
		if err == nil {
			// 这里设置了user的值，之后可以通过c.Get("user")拿到用户信息
			c.Set("MyId", strconv.FormatUint(userInfo.Id, 10))
			fmt.Printf("userInfo:%#v", userInfo)
			c.Next()
			return
		} else {
			c.JSON(200, serializer.Response{
				StatusCode: 1,
				StatusMsg:  "用户token校验失败,请登陆",
			})
			c.Abort()
		}
		// if user, _ := c.Get("user"); user != nil {
		// 	if _, ok := user.(*model.User); ok {
		// 		c.Next()
		// 		return
		// 	}
		// }
		// 在上面写校验用户token的逻辑，之后成功调用Next()

		// c.JSON(200, serializer.CheckLogin())
		// c.Abort()
	}
}

//generate jwt token
func GenToken(claims *JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(Secret))
	if err != nil {
		return "", errors.New(ErrorServerBusy)
	}
	return signedToken, nil
}

//验证jwt token
func verifyAction(strToken string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(strToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil {
		return nil, errors.New(ErrorServerBusy)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New(ErrorReLogin)
	}
	if err := token.Claims.Valid(); err != nil {
		return nil, errors.New(ErrorReLogin)
	}

	fmt.Println("verify")
	return claims, nil
}

func refresh(c *gin.Context) {
	strToken := c.Param("token")
	claims, err := verifyAction(strToken)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	claims.ExpiresAt = time.Now().Unix() + (claims.ExpiresAt - claims.IssuedAt)
	signedToken, err := GenToken(claims)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	c.String(http.StatusOK, signedToken, ", ", claims.ExpiresAt)
}
func GetIdByToken(tokenString string) (uint64, error) {
	var claims JWTClaims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil {
		return 0, errors.New("获取id失败")
	}
	return claims.Id, nil
}
