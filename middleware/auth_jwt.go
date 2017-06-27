package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/dgrijalva/jwt-go.v3"
)


type GinJWTMiddleware struct {
	Realm string
	SigningAlgorithm string
	Key []byte
	Timeout time.Duration
	MaxRefresh time.Duration
	Authenticator func(userID string, password string, c *gin.Context) (string, bool)
	PayloadFunc func(userID string) map[string]interface{}
	Unauthorized func(*gin.Context, int, string)
	IdentityHandler func(jwt.MapClaims) string
	GuardHandler func(claims jwt.MapClaims) string
	TimeFunc func() time.Time
}

type UserCredential struct {
	Userid string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Guard string `form:"guard" json:"guard" binding:""`
}

func (mw *GinJWTMiddleware) MiddlewareInit() error {
	if mw.SigningAlgorithm == "" {
		mw.SigningAlgorithm = "HS256"
	}

	if mw.Timeout == 0 {
		mw.Timeout = time.Hour
	}

	if mw.TimeFunc == nil {
		mw.TimeFunc = time.Now
	}

	if mw.Unauthorized == nil {
		mw.Unauthorized = func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		}
	}

	if mw.IdentityHandler == nil {
		mw.IdentityHandler = func(claims jwt.MapClaims) string {
			return claims["id"].(string)
		}
	}

	if mw.GuardHandler == nil {
		mw.GuardHandler = func(claims jwt.MapClaims) string {
			return claims["guard"].(string)
		}
	}

	if mw.Realm == "" {
		return errors.New("realm is required")
	}

	if mw.Key == nil {
		return errors.New("secret key is required")
	}

	return nil
}

func (mw *GinJWTMiddleware) LoginHandler(c *gin.Context) {

	// Initial middleware default setting.
	mw.MiddlewareInit()

	var loginVals UserCredential

	if c.BindJSON(&loginVals) != nil {
		mw.unauthorized(c, http.StatusBadRequest, "Missing Userid or Password")
		return
	}

	if mw.Authenticator == nil {
		mw.unauthorized(c, http.StatusInternalServerError, "Missing define authenticator func")
		return
	}

	userID, ok := mw.Authenticator(loginVals.Userid, loginVals.Password, c)


	if !ok {
		mw.unauthorized(c, http.StatusUnauthorized, "Incorrect Userid / Password")
		return
	}

	// Create the token
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(loginVals.Userid) {
			claims[key] = value
		}
	}

	if userID == "" {
		userID = loginVals.Userid
	}

	guard := loginVals.Guard

	if guard == "" {
		guard = "user"
	}


	expire := mw.TimeFunc().Add(mw.Timeout)
	claims["id"] = userID
	claims["guard"] = guard
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = mw.TimeFunc().Unix()

	tokenString, err := token.SignedString(mw.Key)

	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, "Create JWT Token faild")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":  tokenString,
		"expire": expire.Format(time.RFC3339),
	})
}

func (mw *GinJWTMiddleware) RefreshHandler(c *gin.Context) {
	token, _ := mw.parseToken(c)
	claims := token.Claims.(jwt.MapClaims)

	origIat := int64(claims["orig_iat"].(float64))

	if origIat < mw.TimeFunc().Add(-mw.MaxRefresh).Unix() {
		mw.unauthorized(c, http.StatusUnauthorized, "Token is expired.")
		return
	}

	// Create the token
	newToken := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	newClaims := newToken.Claims.(jwt.MapClaims)

	for key := range claims {
		newClaims[key] = claims[key]
	}

	expire := mw.TimeFunc().Add(mw.Timeout)
	newClaims["id"] = claims["id"]
	newClaims["guard"] = claims["guard"]
	newClaims["exp"] = expire.Unix()
	newClaims["orig_iat"] = origIat

	tokenString, err := newToken.SignedString(mw.Key)

	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, "Create JWT Token faild")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":  tokenString,
		"expire": expire.Format(time.RFC3339),
	})
}

func (mw *GinJWTMiddleware) TokenGenerator(userID string, guard string) string {
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(userID) {
			claims[key] = value
		}
	}

	claims["id"] = userID
	claims["guard"] = guard
	claims["exp"] = mw.TimeFunc().Add(mw.Timeout).Unix()
	claims["orig_iat"] = mw.TimeFunc().Unix()

	tokenString, _ := token.SignedString(mw.Key)

	return tokenString
}

func (mw *GinJWTMiddleware) MiddlewareFunc(guard string) gin.HandlerFunc {
	if err := mw.MiddlewareInit(); err != nil {
		return func(c *gin.Context) {
			mw.unauthorized(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	return func(c *gin.Context) {
		mw.middlewareImpl(c, guard)
		return
	}
}

func (mw *GinJWTMiddleware) middlewareImpl(c *gin.Context, guard string) {
	token, err := mw.parseToken(c)

	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, err.Error())
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	id := mw.IdentityHandler(claims)
	tokenGuard := mw.GuardHandler(claims)

	c.Set("JWT_PAYLOAD", claims)
	c.Set("userID", id)
	c.Set("guard", tokenGuard)

	if tokenGuard != guard {
		mw.unauthorized(c, http.StatusForbidden, "You don't have permission to access.")
		return
	}

	c.Next()
}

func (mw *GinJWTMiddleware) parseToken(c *gin.Context) (*jwt.Token, error) {
	var token string
	var err error

	parts := strings.Split("header:Authorization", ":")
	token, err = mw.jwtFromHeader(c, parts[1])

	if err != nil {
		return nil, err
	}

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(mw.SigningAlgorithm) != token.Method {
			return nil, errors.New("invalid signing algorithm")
		}

		return mw.Key, nil
	})
}

func (mw *GinJWTMiddleware) jwtFromHeader(c *gin.Context, key string) (string, error) {
	authHeader := c.Request.Header.Get(key)

	if authHeader == "" {
		return "", errors.New("auth header empty")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", errors.New("invalid auth header")
	}

	return parts[1], nil
}

func (mw *GinJWTMiddleware) unauthorized(c *gin.Context, code int, message string) {

	if mw.Realm == "" {
		mw.Realm = "gin jwt"
	}

	c.Header("WWW-Authenticate", "JWT realm="+mw.Realm)
	c.Abort()

	mw.Unauthorized(c, code, message)

	return
}