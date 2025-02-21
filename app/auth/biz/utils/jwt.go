package utils

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	arj  *ARJWT
	once sync.Once
)

type ARJWT struct {
	// 密钥，用以加密 JWT
	Key []byte

	// 定义 access token 过期时间（单位：分钟）即当颁发 access token 后，多少分钟后 access token 过期
	AccessExpireTime int64

	// 定义 refresh token 过期时间（单位：分钟）即当颁发 refresh token 后，多少分钟后 refresh token 过期
	RefreshExpireTime int64

	// token 的签发者
	Issuer string
}

type JWTCustomClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewARJWT() *ARJWT {
	once.Do(func() {
		// arj = &ARJWT{
		// 	Key:               []byte(conf.GetConf().JWT.Secret),
		// 	AccessExpireTime:  conf.GetConf().JWT.AccessExpireTime,
		// 	RefreshExpireTime: conf.GetConf().JWT.RefreshExpireTime,
		// 	Issuer:            conf.GetConf().JWT.Issuer,
		// }
		arj = &ARJWT{
			Key:               []byte("gomall"),
			AccessExpireTime:  10,
			RefreshExpireTime: 30,
			Issuer:            "gomall",
		}
	})
	return arj
}

func (arj *ARJWT) GenerateToken(userId int64, username string) (accessToken, refreshToken string, err error) {
	// 生成 access token
	mc := JWTCustomClaims{
		UserID:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(arj.AccessExpireTime) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    arj.Issuer,
		},
	}

	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, mc).SignedString(arj.Key)
	if err != nil {
		log.Printf("generate access token failed: %v \n", err)
		return "", "", err
	}

	// 生成 refresh token
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(arj.RefreshExpireTime) * time.Minute)),
		Issuer:    arj.Issuer,
	}).SignedString(arj.Key)
	if err != nil {
		log.Printf("generate refresh token failed: %v \n", err)
		return "", "", err
	}
	return
}

func (arj *ARJWT) ParseAccessToken(tokenString string) (*JWTCustomClaims, error) {
	claims := JWTCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims,
		func(token *jwt.Token) (interface{}, error) {
			return arj.Key, nil
		},
	)

	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)
		if ok {
			switch validationErr.Errors {
			case jwt.ValidationErrorMalformed:
				return nil, errors.New("请求令牌格式有误")
			case jwt.ValidationErrorExpired:
				return nil, errors.New("令牌已过期")
			}
		}
		return nil, errors.New("请求令牌无效")
	}

	if _, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, errors.New("请求令牌无效")
}

func (arj *ARJWT) RefreshToken(accessToken, refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	// 先判断 refresh token 是否有效
	if _, err = jwt.Parse(refreshToken,
		func(token *jwt.Token) (interface{}, error) {
			return arj.Key, nil
		},
	); err != nil {
		return
	}

	// 从旧的 access token 中解析出 JWTCustomClaims 数据出来
	claims := JWTCustomClaims{}
	_, err = jwt.ParseWithClaims(accessToken, &claims,
		func(token *jwt.Token) (interface{}, error) {
			return arj.Key, nil
		},
	)
	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)
		// 当 access token 是过期错误，并且 refresh token 没有过期时就创建一个新的 access token 和 refresh token
		if ok && validationErr.Errors == jwt.ValidationErrorExpired {
			// 重新生成新的 access token 和 refresh token
			return arj.GenerateToken(claims.UserID, claims.Username)
		}
	}

	return accessToken, refreshToken, errors.New("access token still valid")
}
