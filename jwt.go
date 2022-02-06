package main

import (
	"app/models"
	"errors"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 参考
// https://taisablog.com/archives/go-jwt
func createJwtMiddleware(secretKey, identityKey string, db *gorm.DB, tokenTimeout time.Duration) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		// Realm:       "test zone",
		Key:     []byte(secretKey),
		Timeout: tokenTimeout,
		// 期限が切れてから自動でリフレッシュできる期間。実質のトークンの期限はTimeout + MaxRefresh
		// デフォルト0
		// MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		// ログイン時にTokenを発行するLoginHandlerではAuthenticatorとPayloadFuncが呼ばれる
		// PayloadFuncはペイロードに含めるクレームを設定
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				// これでこんな感じの内容のJWTがクライアントに返される
				// header: {"alg":"HS256","typ":"JWT"}
				// payload: {"email":"test2@example.com","exp":1644053776,"orig_iat":16440501}
				return jwt.MapClaims{
					identityKey: v.Email,
				}
			}
			return jwt.MapClaims{}
		},
		// MiddlewareFuncで呼ばれる
		// クレームからログインIDを取得する
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.User{
				Email: claims[identityKey].(string),
			}
		},
		// ログイン時にTokenを発行するLoginHandlerではAuthenticatorとPayloadFuncが呼ばれる
		// Authenticatorはログイン認証の為の関数
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var ul models.UserLogin
			if err := c.ShouldBind(&ul); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			var u models.User
			tx := db.First(
				&u,
				models.User{Email: ul.Email},
			)
			if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
				return nil, jwt.ErrFailedAuthentication
			}

			if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordDigest), []byte(ul.Password)); err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return &u, nil
		},
		// MiddlewareFuncで呼ばれる
		// トークンのユーザ情報からの認証
		// dataはIdentityHandlerの戻り値
		Authorizator: func(data interface{}, c *gin.Context) bool {
			v, ok := data.(*models.User)
			if !ok {
				return false
			}

			var u models.User
			// identityKeyで探して、見つかればok
			tx := db.First(&u, v)
			return !errors.Is(tx.Error, gorm.ErrRecordNotFound)
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		// TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		// TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
}
