package middleware

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/realjf/goframe/utils"
	"github.com/realjf/goframe/utils/conv"
)

const (
	SUCCESS                        = 200
	INVALID_PARAMS                 = 99
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 98
	ERROR_AUTH_CHECK_TOKEN_FAIL    = 97

	ISSUER = "realjf"
)

var (
	jwtSecret []byte
)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var code int
		var data interface{}

		code = SUCCESS
		token := r.Header.Get("token")
		if token == "" {
			code = INVALID_PARAMS
		} else {
			_, err := ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					code = ERROR_AUTH_CHECK_TOKEN_TIMEOUT
				default:
					code = ERROR_AUTH_CHECK_TOKEN_FAIL
				}
			}
		}

		if code != SUCCESS {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(conv.Bytes(gin.H{
				"code":    code,
				"message": "",
				"data":    data,
			}))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		utils.EncodeMD5(username),
		utils.EncodeMD5(password),
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    ISSUER,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
