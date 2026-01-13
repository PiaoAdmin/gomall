package jwt

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/elastic/pkcs8"
	"github.com/golang-jwt/jwt/v4"
)

// MapClaims type that uses the map[string]interface{} for JSON decoding
// This is the default claims type if you don't supply one
// MapClaims 使用 map[string]interface{} 进行 JSON 解码
// 如果你不提供 claims 类型，这是默认的 claims 类型
type MapClaims map[string]interface{}

// HertzJWTMiddleware provides a Json-Web-Token authentication implementation. On failure, a 401 HTTP response
// is returned. On success, the wrapped middleware is called, and the userID is made available as
// c.Get("userID").(string).
// Users can get a token by posting a json request to LoginHandler. The token then needs to be passed in
// the Authentication header. Example: Authorization:Bearer XXX_TOKEN_XXX
// HertzJWTMiddleware 提供 Json-Web-Token 认证实现。失败时返回 401 HTTP 响应。
// 成功时，调用包装的中间件，userID 可通过 c.Get("userID").(string) 获取。
// 用户可以通过向 LoginHandler 发送 json 请求来获取 token。然后需要在 Authentication header 中传递 token。
// 例如：Authorization:Bearer XXX_TOKEN_XXX
type HertzJWTMiddleware struct {
	// Realm name to display to the user. Required.
	// Realm 名称，展示给用户。必须。
	Realm string

	// signing algorithm - possible values are HS256, HS384, HS512, RS256, RS384 or RS512
	// Optional, default is HS256.
	// 签名算法 - 可能的值有 HS256, HS384, HS512, RS256, RS384 或 RS512
	// 可选，默认为 HS256。
	SigningAlgorithm string

	// Secret key used for signing. Required.
	// 用于签名的密钥。必须。
	Key []byte

	// Callback to retrieve key used for signing. Setting KeyFunc will bypass
	// all other key settings
	// 回调函数，用于获取签名的密钥。设置 KeyFunc 将绕过所有其他密钥设置
	KeyFunc func(token *jwt.Token) (interface{}, error)

	// Duration that a jwt token is valid. Optional, defaults to one hour.
	// jwt token 的有效期。可选，默认为一小时。
	Timeout time.Duration
	// Callback function that will override the default Timeout duration
	// Optional, defaults to return one hour
	// 回调函数，用于覆盖默认的 Timeout 时长
	// 可选，默认返回一小时
	TimeoutFunc func(claims jwt.MapClaims) time.Duration

	// This field allows clients to refresh their token until MaxRefresh has passed.
	// Note that clients can refresh their token in the last moment of MaxRefresh.
	// This means that the maximum validity timespan for a token is TokenTime + MaxRefresh.
	// Optional, defaults to 0 meaning not refreshable.
	// 此字段允许客户端在 MaxRefresh 过去之前刷新其 token。
	// 注意，客户端可以在 MaxRefresh 的最后时刻刷新 token。
	// 这意味着 token 的最大有效期是 TokenTime + MaxRefresh。
	// 可选，默认为 0，表示不可刷新。
	MaxRefresh time.Duration

	// Callback function that should perform the authentication of the user based on login info.
	// Must return user data as user identifier, it will be stored in Claim Array. Required.
	// Check error (e) to determine the appropriate error message.
	// 回调函数，应根据登录信息执行用户认证。
	// 必须返回用户数据作为用户标识符，它将存储在 Claim Array 中。必须。
	// 检查错误 (e) 以确定适当的错误消息。
	Authenticator func(ctx context.Context, c *app.RequestContext) (interface{}, error)

	// Callback function that should perform the authorization of the authenticated user. Called
	// only after an authentication success. Must return true on success, false on failure.
	// Optional, default to success.
	// 回调函数，应执行已认证用户的授权。仅在认证成功后调用。
	// 成功返回 true，失败返回 false。
	// 可选，默认为成功。
	Authorizator func(data interface{}, ctx context.Context, c *app.RequestContext) bool

	// Callback function that will be called during login.
	// Using this function it is possible to add additional payload data to the webtoken.
	// The data is then made available during requests via c.Get("JWT_PAYLOAD").
	// Note that the payload is not encrypted.
	// The attributes mentioned on jwt.io can't be used as keys for the map.
	// Optional, by default no additional data will be set.
	// 登录期间调用的回调函数。
	// 使用此函数可以向 webtoken 添加额外的 payload 数据。
	// 然后可以通过 c.Get("JWT_PAYLOAD") 在请求期间获取数据。
	// 注意 payload 未加密。
	// jwt.io 上提到的属性不能用作 map 的键。
	// 可选，默认不设置额外数据。
	PayloadFunc func(data interface{}) MapClaims

	// User can define own Unauthorized func.
	// 用户可以定义自己的 Unauthorized 函数。
	Unauthorized func(ctx context.Context, c *app.RequestContext, code int, message string)

	// User can define own LoginResponse func.
	// 用户可以定义自己的 LoginResponse 函数。
	LoginResponse func(ctx context.Context, c *app.RequestContext, code int, message string, time time.Time)

	// User can define own LogoutResponse func.
	// 用户可以定义自己的 LogoutResponse 函数。
	LogoutResponse func(ctx context.Context, c *app.RequestContext, code int)

	// User can define own RefreshResponse func.
	// 用户可以定义自己的 RefreshResponse 函数。
	RefreshResponse func(ctx context.Context, c *app.RequestContext, code int, message string, time time.Time)

	// Set the identity handler function
	// 设置身份处理函数
	IdentityHandler func(ctx context.Context, c *app.RequestContext) interface{}

	// Set the identity key
	// 设置身份键
	IdentityKey string

	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	// - "param:<name>"
	// - "form:<name>"
	// TokenLookup 是 "<source>:<name>" 形式的字符串，用于从请求中提取 token。
	// 可选。默认值为 "header:Authorization"。
	// 可能的值：
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	// - "param:<name>"
	// - "form:<name>"
	TokenLookup string

	// TokenHeadName is a string in the header. Default value is "Bearer"
	// TokenHeadName 是 header 中的字符串。默认值为 "Bearer"
	TokenHeadName string

	// WithoutDefaultTokenHeadName allow set empty TokenHeadName
	// WithoutDefaultTokenHeadName 允许设置空的 TokenHeadName
	WithoutDefaultTokenHeadName bool

	// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
	// TimeFunc 提供当前时间。你可以覆盖它以使用其他时间值。这对测试或如果你的服务器使用与 token 不同的时区很有用。
	TimeFunc func() time.Time

	// HTTP Status messages for when something in the JWT middleware fails.
	// Check error (e) to determine the appropriate error message.
	// JWT 中间件失败时的 HTTP 状态消息。
	// 检查错误 (e) 以确定适当的错误消息。
	HTTPStatusMessageFunc func(e error, ctx context.Context, c *app.RequestContext) string

	// Private key file for asymmetric algorithms
	// 非对称算法的私钥文件
	PrivKeyFile string

	// Private Key bytes for asymmetric algorithms
	//
	// Note: PrivKeyFile takes precedence over PrivKeyBytes if both are set
	// 非对称算法的私钥字节
	// 注意：如果同时设置，PrivKeyFile 优先于 PrivKeyBytes
	PrivKeyBytes []byte

	// Public key file for asymmetric algorithms
	// 非对称算法的公钥文件
	PubKeyFile string

	// Private key passphrase
	// 私钥密码短语
	PrivateKeyPassphrase string

	// Public key bytes for asymmetric algorithms.
	//
	// Note: PubKeyFile takes precedence over PubKeyBytes if both are set
	// 非对称算法的公钥字节。
	// 注意：如果同时设置，PubKeyFile 优先于 PubKeyBytes
	PubKeyBytes []byte

	// Private key
	// 私钥
	privKey *rsa.PrivateKey

	// Public key
	// 公钥
	pubKey *rsa.PublicKey

	// Optionally return the token as a cookie
	// 可选地将 token 作为 cookie 返回
	SendCookie bool

	// Duration that a cookie is valid. Optional, by default equals to Timeout value.
	// cookie 的有效期。可选，默认等于 Timeout 值。
	CookieMaxAge time.Duration

	// Allow insecure cookies for development over http
	// 允许在 http 开发中使用不安全的 cookie
	SecureCookie bool

	// Allow cookies to be accessed client side for development
	// 允许客户端访问 cookie 进行开发
	CookieHTTPOnly bool

	// Allow cookie domain change for development
	// 允许更改 cookie 域进行开发
	CookieDomain string

	// SendAuthorization allow return authorization header for every request
	// SendAuthorization 允许每个请求返回 authorization header
	SendAuthorization bool

	// Disable abort() of context.
	// 禁用 context 的 abort()。
	DisabledAbort bool

	// CookieName allow cookie name change for development
	// CookieName 允许更改 cookie 名称进行开发
	CookieName string

	// CookieSameSite allow use protocol.CookieSameSite cookie param
	// CookieSameSite 允许使用 protocol.CookieSameSite cookie 参数
	CookieSameSite protocol.CookieSameSite

	// ParseOptions allow to modify jwt's parser methods
	// ParseOptions 允许修改 jwt 的解析方法
	ParseOptions []jwt.ParserOption
}

var (
	// ErrMissingSecretKey indicates Secret key is required
	// ErrMissingSecretKey 表示需要 Secret key
	ErrMissingSecretKey = errors.New("secret key is required")

	// ErrForbidden when HTTP status 403 is given
	// ErrForbidden 当 HTTP 状态为 403 时
	ErrForbidden = errors.New("you don't have permission to access this resource")

	// ErrMissingAuthenticatorFunc indicates Authenticator is required
	// ErrMissingAuthenticatorFunc 表示需要 Authenticator
	ErrMissingAuthenticatorFunc = errors.New("HertzJWTMiddleware.Authenticator func is undefined")

	// ErrMissingLoginValues indicates a user tried to authenticate without username or password
	// ErrMissingLoginValues 表示用户尝试在没有用户名或密码的情况下进行身份验证
	ErrMissingLoginValues = errors.New("missing Username or Password")

	// ErrFailedAuthentication indicates authentication failed, could be faulty username or password
	// ErrFailedAuthentication 表示身份验证失败，可能是用户名或密码错误
	ErrFailedAuthentication = errors.New("incorrect Username or Password")

	// ErrFailedTokenCreation indicates JWT Token failed to create, reason unknown
	// ErrFailedTokenCreation 表示 JWT Token 创建失败，原因未知
	ErrFailedTokenCreation = errors.New("failed to create JWT Token")

	// ErrExpiredToken indicates JWT token has expired. Can't refresh.
	// ErrExpiredToken 表示 JWT token 已过期。无法刷新。
	ErrExpiredToken = errors.New("token is expired") // in practice, this is generated from the jwt library not by us

	// ErrEmptyAuthHeader can be thrown if authing with a HTTP header, the Auth header needs to be set
	// ErrEmptyAuthHeader 如果使用 HTTP header 进行身份验证，则需要设置 Auth header
	ErrEmptyAuthHeader = errors.New("auth header is empty")

	// ErrMissingExpField missing exp field in token
	// ErrMissingExpField token 中缺少 exp 字段
	ErrMissingExpField = errors.New("missing exp field")

	// ErrWrongFormatOfExp field must be float64 format
	// ErrWrongFormatOfExp 字段必须是 float64 格式
	ErrWrongFormatOfExp = errors.New("exp must be float64 format")

	// ErrInvalidAuthHeader indicates auth header is invalid, could for example have the wrong Realm name
	// ErrInvalidAuthHeader 表示 auth header 无效，例如可能有错误的 Realm 名称
	ErrInvalidAuthHeader = errors.New("auth header is invalid")

	// ErrEmptyQueryToken can be thrown if authing with URL Query, the query token variable is empty
	// ErrEmptyQueryToken 如果使用 URL Query 进行身份验证，则 query token 变量为空
	ErrEmptyQueryToken = errors.New("query token is empty")

	// ErrEmptyCookieToken can be thrown if authing with a cookie, the token cookie is empty
	// ErrEmptyCookieToken 如果使用 cookie 进行身份验证，则 token cookie 为空
	ErrEmptyCookieToken = errors.New("cookie token is empty")

	// ErrEmptyParamToken can be thrown if authing with parameter in path, the parameter in path is empty
	// ErrEmptyParamToken 如果使用路径中的参数进行身份验证，则路径中的参数为空
	ErrEmptyParamToken = errors.New("parameter token is empty")

	// ErrEmptyFormToken can be thrown if authing with post form, the form token is empty
	// ErrEmptyFormToken 如果使用 post form 进行身份验证，则 form token 为空
	ErrEmptyFormToken = errors.New("form token is empty")

	// ErrInvalidSigningAlgorithm indicates signing algorithm is invalid, needs to be HS256, HS384, HS512, RS256, RS384 or RS512
	// ErrInvalidSigningAlgorithm 表示签名算法无效，需要是 HS256, HS384, HS512, RS256, RS384 或 RS512
	ErrInvalidSigningAlgorithm = errors.New("invalid signing algorithm")

	// ErrNoPrivKeyFile indicates that the given private key is unreadable
	// ErrNoPrivKeyFile 表示给定的私钥不可读
	ErrNoPrivKeyFile = errors.New("private key file unreadable")

	// ErrNoPubKeyFile indicates that the given public key is unreadable
	// ErrNoPubKeyFile 表示给定的公钥不可读
	ErrNoPubKeyFile = errors.New("public key file unreadable")

	// ErrInvalidPrivKey indicates that the given private key is invalid
	// ErrInvalidPrivKey 表示给定的私钥无效
	ErrInvalidPrivKey = errors.New("private key invalid")

	// ErrInvalidPubKey indicates the the given public key is invalid
	// ErrInvalidPubKey 表示给定的公钥无效
	ErrInvalidPubKey = errors.New("public key invalid")

	// IdentityKey default identity key
	// IdentityKey 默认身份键
	IdentityKey = "identity"
)

// New for check error with HertzJWTMiddleware
// New 用于检查 HertzJWTMiddleware 的错误
func New(m *HertzJWTMiddleware) (*HertzJWTMiddleware, error) {
	if err := m.MiddlewareInit(); err != nil {
		return nil, err
	}

	return m, nil
}

func (mw *HertzJWTMiddleware) readKeys() error {
	err := mw.privateKey()
	if err != nil {
		return err
	}
	err = mw.publicKey()
	if err != nil {
		return err
	}
	return nil
}

func (mw *HertzJWTMiddleware) privateKey() error {
	var keyData []byte
	if mw.PrivKeyFile == "" {
		keyData = mw.PrivKeyBytes
	} else {
		filecontent, err := os.ReadFile(mw.PrivKeyFile)
		if err != nil {
			return ErrNoPrivKeyFile
		}
		keyData = filecontent
	}

	if mw.PrivateKeyPassphrase != "" {
		key, err := pkcs8.ParsePKCS8PrivateKey(keyData, []byte(mw.PrivateKeyPassphrase))
		if err != nil {
			return ErrInvalidPrivKey
		}

		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return ErrInvalidPrivKey
		}

		mw.privKey = rsaKey
		return nil
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return ErrInvalidPrivKey
	}
	mw.privKey = key
	return nil
}

func (mw *HertzJWTMiddleware) publicKey() error {
	var keyData []byte
	if mw.PubKeyFile == "" {
		keyData = mw.PubKeyBytes
	} else {
		filecontent, err := os.ReadFile(mw.PubKeyFile)
		if err != nil {
			return ErrNoPubKeyFile
		}
		keyData = filecontent
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return ErrInvalidPubKey
	}
	mw.pubKey = key
	return nil
}

func (mw *HertzJWTMiddleware) usingPublicKeyAlgo() bool {
	switch mw.SigningAlgorithm {
	case "RS256", "RS512", "RS384":
		return true
	}
	return false
}

// MiddlewareInit initialize jwt configs.
// MiddlewareInit 初始化 jwt 配置。
func (mw *HertzJWTMiddleware) MiddlewareInit() error {
	if mw.TokenLookup == "" {
		mw.TokenLookup = "header:Authorization"
	}

	if mw.SigningAlgorithm == "" {
		mw.SigningAlgorithm = "HS256"
	}

	if mw.Timeout == 0 {
		mw.Timeout = time.Hour
	}

	if mw.TimeoutFunc == nil {
		mw.TimeoutFunc = func(claims jwt.MapClaims) time.Duration {
			return mw.Timeout
		}
	}

	if mw.TimeFunc == nil {
		mw.TimeFunc = time.Now
	}

	mw.TokenHeadName = strings.TrimSpace(mw.TokenHeadName)
	if len(mw.TokenHeadName) == 0 && !mw.WithoutDefaultTokenHeadName {
		mw.TokenHeadName = "Bearer"
	}

	if mw.Authorizator == nil {
		mw.Authorizator = func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
			return true
		}
	}

	if mw.Unauthorized == nil {
		mw.Unauthorized = func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(code, map[string]interface{}{
				"code":    code,
				"message": message,
			})
		}
	}

	if mw.LoginResponse == nil {
		mw.LoginResponse = func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code":   http.StatusOK,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		}
	}

	if mw.LogoutResponse == nil {
		mw.LogoutResponse = func(ctx context.Context, c *app.RequestContext, code int) {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code": http.StatusOK,
			})
		}
	}

	if mw.RefreshResponse == nil {
		mw.RefreshResponse = func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code":   http.StatusOK,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		}
	}

	if mw.IdentityKey == "" {
		mw.IdentityKey = IdentityKey
	}

	if mw.IdentityHandler == nil {
		mw.IdentityHandler = func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := ExtractClaims(ctx, c)
			return claims[mw.IdentityKey]
		}
	}

	if mw.HTTPStatusMessageFunc == nil {
		mw.HTTPStatusMessageFunc = func(e error, ctx context.Context, c *app.RequestContext) string {
			return e.Error()
		}
	}

	if mw.Realm == "" {
		mw.Realm = "hertz jwt"
	}

	if mw.CookieMaxAge == 0 {
		mw.CookieMaxAge = mw.Timeout
	}

	if mw.CookieName == "" {
		mw.CookieName = "jwt"
	}

	// bypass other key settings if KeyFunc is set
	// 如果设置了 KeyFunc，则绕过其他密钥设置
	if mw.KeyFunc != nil {
		return nil
	}

	if mw.usingPublicKeyAlgo() {
		return mw.readKeys()
	}

	if mw.Key == nil {
		return ErrMissingSecretKey
	}
	return nil
}

// MiddlewareFunc makes HertzJWTMiddleware implement the Middleware interface.
// MiddlewareFunc 使 HertzJWTMiddleware 实现 Middleware 接口。
func (mw *HertzJWTMiddleware) MiddlewareFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		mw.middlewareImpl(ctx, c)
	}
}

func (mw *HertzJWTMiddleware) middlewareImpl(ctx context.Context, c *app.RequestContext) {
	claims, err := mw.GetClaimsFromJWT(ctx, c)
	if err != nil {
		mw.unauthorized(ctx, c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err, ctx, c))
		return
	}

	switch v := claims["exp"].(type) {
	case nil:
		mw.unauthorized(ctx, c, http.StatusBadRequest, mw.HTTPStatusMessageFunc(ErrMissingExpField, ctx, c))
		return
	case float64:
		if int64(v) < mw.TimeFunc().Unix() {
			mw.unauthorized(ctx, c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(ErrExpiredToken, ctx, c))
			return
		}
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			mw.unauthorized(ctx, c, http.StatusBadRequest, mw.HTTPStatusMessageFunc(ErrWrongFormatOfExp, ctx, c))
			return
		}
		if n < mw.TimeFunc().Unix() {
			mw.unauthorized(ctx, c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(ErrExpiredToken, ctx, c))
			return
		}
	default:
		mw.Unauthorized(ctx, c, http.StatusBadRequest, mw.HTTPStatusMessageFunc(ErrWrongFormatOfExp, ctx, c))
	}

	c.Set("JWT_PAYLOAD", claims)
	identity := mw.IdentityHandler(ctx, c)

	if identity != nil {
		c.Set(mw.IdentityKey, identity)
	}

	if !mw.Authorizator(identity, ctx, c) {
		mw.unauthorized(ctx, c, http.StatusForbidden, mw.HTTPStatusMessageFunc(ErrForbidden, ctx, c))
		return
	}

	c.Next(ctx)
}

// GetClaimsFromJWT get claims from JWT token
// GetClaimsFromJWT 从 JWT token 获取 claims
func (mw *HertzJWTMiddleware) GetClaimsFromJWT(ctx context.Context, c *app.RequestContext) (MapClaims, error) {
	token, err := mw.ParseToken(ctx, c)
	if err != nil {
		return nil, err
	}

	if mw.SendAuthorization {
		if v, ok := c.Get("JWT_TOKEN"); ok {
			c.Header("Authorization", mw.TokenHeadName+" "+v.(string))
		}
	}

	claims := MapClaims{}
	for key, value := range token.Claims.(jwt.MapClaims) {
		claims[key] = value
	}

	return claims, nil
}

// LoginHandler can be used by clients to get a jwt token.
// Payload needs to be json in the form of {"username": "USERNAME", "password": "PASSWORD"}.
// Reply will be of the form {"token": "TOKEN"}.
// LoginHandler 可供客户端用于获取 jwt token。
// Payload 需要是 {"username": "USERNAME", "password": "PASSWORD"} 形式的 json。
// 回复将是 {"token": "TOKEN"} 形式。
func (mw *HertzJWTMiddleware) LoginHandler(ctx context.Context, c *app.RequestContext) {
	if mw.Authenticator == nil {
		mw.unauthorized(ctx, c,
			http.StatusInternalServerError,
			mw.HTTPStatusMessageFunc(ErrMissingAuthenticatorFunc, ctx, c))
		return
	}

	data, err := mw.Authenticator(ctx, c)
	if err != nil {
		mw.unauthorized(ctx, c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err, ctx, c))
		return
	}

	// Create the token
	// 创建 token
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(data) {
			claims[key] = value
		}
	}

	copyClaims := make(jwt.MapClaims, len(claims))
	for k, v := range claims {
		copyClaims[k] = v
	}

	expire := mw.TimeFunc().Add(mw.TimeoutFunc(copyClaims))
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = mw.TimeFunc().Unix()
	tokenString, err := mw.signedString(token)
	if err != nil {
		mw.unauthorized(ctx, c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(ErrFailedTokenCreation, ctx, c))
		return
	}

	// set cookie
	// 设置 cookie
	if mw.SendCookie {
		expireCookie := mw.TimeFunc().Add(mw.CookieMaxAge)
		maxage := int(expireCookie.Unix() - mw.TimeFunc().Unix())
		c.SetCookie(mw.CookieName, tokenString, maxage, "/", mw.CookieDomain, mw.CookieSameSite, mw.SecureCookie, mw.CookieHTTPOnly)
	}

	mw.LoginResponse(ctx, c, http.StatusOK, tokenString, expire)
}

// LogoutHandler can be used by clients to remove the jwt cookie (if set)
// LogoutHandler 可供客户端用于删除 jwt cookie（如果已设置）
func (mw *HertzJWTMiddleware) LogoutHandler(ctx context.Context, c *app.RequestContext) {
	// delete auth cookie
	// 删除 auth cookie
	if mw.SendCookie {
		c.SetCookie(mw.CookieName, "", -1, "/", mw.CookieDomain, mw.CookieSameSite, mw.SecureCookie, mw.CookieHTTPOnly)
	}

	mw.LogoutResponse(ctx, c, http.StatusOK)
}

func (mw *HertzJWTMiddleware) signedString(token *jwt.Token) (string, error) {
	var tokenString string
	var err error
	if mw.usingPublicKeyAlgo() {
		tokenString, err = token.SignedString(mw.privKey)
	} else {
		tokenString, err = token.SignedString(mw.Key)
	}
	return tokenString, err
}

// RefreshHandler can be used to refresh a token. The token still needs to be valid on refresh.
// Shall be put under an endpoint that is using the HertzJWTMiddleware.
// Reply will be of the form {"token": "TOKEN"}.
// RefreshHandler 可用于刷新 token。刷新时 token 仍需有效。
// 应放在使用 HertzJWTMiddleware 的端点下。
// 回复将是 {"token": "TOKEN"} 形式。
func (mw *HertzJWTMiddleware) RefreshHandler(ctx context.Context, c *app.RequestContext) {
	tokenString, expire, err := mw.RefreshToken(ctx, c)
	if err != nil {
		mw.unauthorized(ctx, c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err, ctx, c))
		return
	}

	mw.RefreshResponse(ctx, c, http.StatusOK, tokenString, expire)
}

// RefreshToken refresh token and check if token is expired
// RefreshToken 刷新 token 并检查 token 是否过期
func (mw *HertzJWTMiddleware) RefreshToken(ctx context.Context, c *app.RequestContext) (string, time.Time, error) {
	claims, err := mw.CheckIfTokenExpire(ctx, c)
	if err != nil {
		return "", time.Now(), err
	}

	// Create the token
	// 创建 token
	newToken := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	newClaims := newToken.Claims.(jwt.MapClaims)
	copyClaims := make(jwt.MapClaims, len(claims))

	for k, v := range claims {
		newClaims[k] = claims[k]
		copyClaims[k] = v
	}

	expire := mw.TimeFunc().Add(mw.TimeoutFunc(copyClaims))
	newClaims["exp"] = expire.Unix()
	// Preserve the original orig_iat to maintain MaxRefresh window
	// 保留原始 orig_iat 以维护 MaxRefresh 窗口
	if origIat, exists := claims["orig_iat"]; exists {
		newClaims["orig_iat"] = origIat
	} else {
		// If orig_iat doesn't exist (backward compatibility), set it to current time
		// 如果 orig_iat 不存在（向后兼容），则将其设置为当前时间
		newClaims["orig_iat"] = mw.TimeFunc().Unix()
	}
	tokenString, err := mw.signedString(newToken)
	if err != nil {
		return "", time.Now(), err
	}

	// set cookie
	// 设置 cookie
	if mw.SendCookie {
		expireCookie := mw.TimeFunc().Add(mw.CookieMaxAge)
		maxage := int(expireCookie.Unix() - time.Now().Unix())
		c.SetCookie(mw.CookieName, tokenString, maxage, "/", mw.CookieDomain, mw.CookieSameSite, mw.SecureCookie, mw.CookieHTTPOnly)
	}

	return tokenString, expire, nil
}

// CheckIfTokenExpire check if token expire
// CheckIfTokenExpire 检查 token 是否过期
func (mw *HertzJWTMiddleware) CheckIfTokenExpire(ctx context.Context, c *app.RequestContext) (jwt.MapClaims, error) {
	token, err := mw.ParseToken(ctx, c)
	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)
		if !ok || validationErr.Errors != jwt.ValidationErrorExpired {
			return nil, err
		}
	}

	claims := token.Claims.(jwt.MapClaims)

	origIat := int64(claims["orig_iat"].(float64))

	if origIat < mw.TimeFunc().Add(-mw.MaxRefresh).Unix() {
		return nil, ErrExpiredToken
	}

	return claims, nil
}

// TokenGenerator method that clients can use to get a jwt token.
// TokenGenerator 客户端可用于获取 jwt token 的方法。
func (mw *HertzJWTMiddleware) TokenGenerator(data interface{}) (string, time.Time, error) {
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(data) {
			claims[key] = value
		}
	}

	copyClaims := make(jwt.MapClaims, len(claims))
	for k, v := range claims {
		copyClaims[k] = v
	}

	expire := mw.TimeFunc().UTC().Add(mw.TimeoutFunc(copyClaims))
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = mw.TimeFunc().Unix()
	tokenString, err := mw.signedString(token)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expire, nil
}

func (mw *HertzJWTMiddleware) jwtFromHeader(_ context.Context, c *app.RequestContext, key string) (string, error) {
	authHeader := c.Request.Header.Get(key)

	if authHeader == "" {
		return "", ErrEmptyAuthHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !((len(parts) == 1 && mw.WithoutDefaultTokenHeadName && mw.TokenHeadName == "") ||
		(len(parts) == 2 && parts[0] == mw.TokenHeadName)) {
		return "", ErrInvalidAuthHeader
	}

	return parts[len(parts)-1], nil
}

func (mw *HertzJWTMiddleware) jwtFromQuery(_ context.Context, c *app.RequestContext, key string) (string, error) {
	token := c.Query(key)

	if token == "" {
		return "", ErrEmptyQueryToken
	}

	return token, nil
}

func (mw *HertzJWTMiddleware) jwtFromCookie(_ context.Context, c *app.RequestContext, key string) (string, error) {
	cookie := string(c.Cookie(key))

	if cookie == "" {
		return "", ErrEmptyCookieToken
	}

	return cookie, nil
}

func (mw *HertzJWTMiddleware) jwtFromParam(_ context.Context, c *app.RequestContext, key string) (string, error) {
	token := c.Param(key)

	if token == "" {
		return "", ErrEmptyParamToken
	}

	return token, nil
}

func (mw *HertzJWTMiddleware) jwtFromForm(_ context.Context, c *app.RequestContext, key string) (string, error) {
	token := c.PostForm(key)

	if token == "" {
		return "", ErrEmptyFormToken
	}

	return token, nil
}

// ParseToken parse jwt token from hertz context
// ParseToken 从 hertz context 解析 jwt token
func (mw *HertzJWTMiddleware) ParseToken(ctx context.Context, c *app.RequestContext) (*jwt.Token, error) {
	var token string
	var err error

	methods := strings.Split(mw.TokenLookup, ",")
	for _, method := range methods {
		if len(token) > 0 {
			break
		}
		parts := strings.Split(strings.TrimSpace(method), ":")
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		switch k {
		case "header":
			token, err = mw.jwtFromHeader(ctx, c, v)
		case "query":
			token, err = mw.jwtFromQuery(ctx, c, v)
		case "cookie":
			token, err = mw.jwtFromCookie(ctx, c, v)
		case "param":
			token, err = mw.jwtFromParam(ctx, c, v)
		case "form":
			token, err = mw.jwtFromForm(ctx, c, v)
		}
	}

	if err != nil {
		return nil, err
	}

	if mw.KeyFunc != nil {
		return jwt.Parse(token, mw.KeyFunc, mw.ParseOptions...)
	}

	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(mw.SigningAlgorithm) != t.Method {
			return nil, ErrInvalidSigningAlgorithm
		}
		if mw.usingPublicKeyAlgo() {
			return mw.pubKey, nil
		}

		// save token string if valid
		// 如果有效，保存 token 字符串
		c.Set("JWT_TOKEN", token)

		return mw.Key, nil
	}, mw.ParseOptions...)
}

// ParseTokenString parse jwt token string
// ParseTokenString 解析 jwt token 字符串
func (mw *HertzJWTMiddleware) ParseTokenString(token string) (*jwt.Token, error) {
	if mw.KeyFunc != nil {
		return jwt.Parse(token, mw.KeyFunc, mw.ParseOptions...)
	}

	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(mw.SigningAlgorithm) != t.Method {
			return nil, ErrInvalidSigningAlgorithm
		}
		if mw.usingPublicKeyAlgo() {
			return mw.pubKey, nil
		}

		return mw.Key, nil
	}, mw.ParseOptions...)
}

func (mw *HertzJWTMiddleware) unauthorized(ctx context.Context, c *app.RequestContext, code int, message string) {
	c.Header("WWW-Authenticate", "JWT realm="+mw.Realm)
	if !mw.DisabledAbort {
		c.Abort()
	}

	mw.Unauthorized(ctx, c, code, message)
}

// ExtractClaims help to extract the JWT claims
// ExtractClaims 帮助提取 JWT claims
func ExtractClaims(ctx context.Context, c *app.RequestContext) MapClaims {
	claims, exists := c.Get("JWT_PAYLOAD")
	if !exists {
		return make(MapClaims)
	}

	return claims.(MapClaims)
}

// ExtractClaimsFromToken help to extract the JWT claims from token
// ExtractClaimsFromToken 帮助从 token 中提取 JWT claims
func ExtractClaimsFromToken(token *jwt.Token) MapClaims {
	if token == nil {
		return make(MapClaims)
	}

	claims := MapClaims{}
	for key, value := range token.Claims.(jwt.MapClaims) {
		claims[key] = value
	}

	return claims
}

// GetToken help to get the JWT token string
// GetToken 帮助获取 JWT token 字符串
func GetToken(ctx context.Context, c *app.RequestContext) string {
	token, exists := c.Get("JWT_TOKEN")
	if !exists {
		return ""
	}

	return token.(string)
}
