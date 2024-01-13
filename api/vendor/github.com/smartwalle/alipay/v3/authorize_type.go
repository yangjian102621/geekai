package alipay

const (
	kProductionPublicAppAuthorize = "https://openauth.alipay.com/oauth2/publicAppAuthorize.htm"
	kSandboxPublicAppAuthorize    = "https://openauth.alipaydev.com/oauth2/publicAppAuthorize.htm"
)

const (
	kProductionAppToAppAuth = "https://openauth.alipay.com/oauth2/appToAppAuth.htm"
	kSandboxAppToAppAuth    = "https://openauth.alipaydev.com/oauth2/appToAppAuth.htm"
)

// SystemOauthToken 换取授权访问令牌接口请求参数 https://docs.open.alipay.com/api_9/alipay.system.oauth.token
type SystemOauthToken struct {
	AuxParam
	AppAuthToken string `json:"-"` // 可选
	GrantType    string `json:"-"` // 值为 authorization_code 时，代表用code换取；值为refresh_token时，代表用refresh_token换取
	Code         string `json:"-"`
	RefreshToken string `json:"-"`
}

func (this SystemOauthToken) APIName() string {
	return "alipay.system.oauth.token"
}

func (this SystemOauthToken) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	m["grant_type"] = this.GrantType
	if this.Code != "" {
		m["code"] = this.Code
	}
	if this.RefreshToken != "" {
		m["refresh_token"] = this.RefreshToken
	}
	return m
}

// SystemOauthTokenRsp 换取授权访问令牌接口请求参数
type SystemOauthTokenRsp struct {
	Error
	UserId       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	ReExpiresIn  int64  `json:"re_expires_in"`
	AuthStart    string `json:"auth_start"`
	OpenId       string `json:"open_id"`
}

// UserInfoShare 支付宝会员授权信息查询接口请求参数 https://docs.open.alipay.com/api_2/alipay.user.info.share
type UserInfoShare struct {
	AuxParam
	AppAuthToken string `json:"-"` // 可选
	AuthToken    string `json:"-"` // 是
}

func (this UserInfoShare) APIName() string {
	return "alipay.user.info.share"
}

func (this UserInfoShare) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	m["auth_token"] = this.AuthToken
	return m
}

// UserInfoShareRsp 支付宝会员授权信息查询接口响应参数
type UserInfoShareRsp struct {
	Error
	AuthNo             string `json:"auth_no"`
	UserId             string `json:"user_id"`
	Avatar             string `json:"avatar"`
	Province           string `json:"province"`
	City               string `json:"city"`
	NickName           string `json:"nick_name"`
	IsStudentCertified string `json:"is_student_certified"`
	UserType           string `json:"user_type"`
	UserStatus         string `json:"user_status"`
	IsCertified        string `json:"is_certified"`
	Gender             string `json:"gender"`
}

// OpenAuthTokenApp 换取应用授权令牌请求参数 https://docs.open.alipay.com/api_9/alipay.open.auth.token.app
type OpenAuthTokenApp struct {
	AuxParam
	GrantType    string `json:"grant_type"` // 值为 authorization_code 时，代表用code换取；值为refresh_token时，代表用refresh_token换取
	Code         string `json:"code"`
	RefreshToken string `json:"refresh_token"`
}

func (this OpenAuthTokenApp) APIName() string {
	return "alipay.open.auth.token.app"
}

func (this OpenAuthTokenApp) Params() map[string]string {
	var m = make(map[string]string)
	m["grant_type"] = this.GrantType
	if this.Code != "" {
		m["code"] = this.Code
	}
	if this.RefreshToken != "" {
		m["refresh_token"] = this.RefreshToken
	}
	return m
}

// OpenAuthTokenAppRsp 换取应用授权令牌响应参数 新版返回值 参见 https://opendocs.alipay.com/open/20160728150111277227/intro
type OpenAuthTokenAppRsp struct {
	Error
	Tokens []*OpenAuthToken `json:"tokens"`
}

type OpenAuthToken struct {
	AppAuthToken    string `json:"app_auth_token"`    // 授权令牌信息
	AppRefreshToken string `json:"app_refresh_token"` // 令牌信息
	AuthAppId       string `json:"auth_app_id"`       // 授权方应用id
	ExpiresIn       int64  `json:"expires_in"`        // 令牌有效期
	ReExpiresIn     int64  `json:"re_expires_in"`     // 有效期
	UserId          string `json:"user_id"`           // 支付宝用户标识
}

// OpenAuthTokenAppQuery 查询某个应用授权AppAuthToken的授权信息 https://opendocs.alipay.com/isv/04hgcp?pathHash=7ea21afe
type OpenAuthTokenAppQuery struct {
	AuxParam
	AppAuthToken string `json:"app_auth_token"` // 必选 应用授权令牌
}

func (this OpenAuthTokenAppQuery) APIName() string {
	return "alipay.open.auth.token.app.query"
}

func (this OpenAuthTokenAppQuery) Params() map[string]string {
	return nil
}

type OpenAuthTokenAppQueryRsp struct {
	Error
	UserId      string   `json:"user_id"`
	AuthAppId   string   `json:"auth_app_id"`
	ExpiresIn   int64    `json:"expires_in"`
	AuthMethods []string `json:"auth_methods"`
	AuthStart   string   `json:"auth_start"`
	AuthEnd     string   `json:"auth_end"`
	Status      string   `json:"status"`
	IsByAppAuth bool     `json:"is_by_app_auth"`
}

// AccountAuth 支付宝登录时, 帮客户端做参数签名, 返回授权请求信息字串接口请求参数 https://docs.open.alipay.com/218/105327/
type AccountAuth struct {
	AuxParam
	Pid      string `json:"pid"`
	TargetId string `json:"target_id"`
	AuthType string `json:"auth_type"`
}

func (this AccountAuth) APIName() string {
	return "alipay.open.auth.sdk.code.get"
}

func (this AccountAuth) Params() map[string]string {
	var m = make(map[string]string)
	m["apiname"] = "com.alipay.account.auth"
	m["app_name"] = "mc"
	m["biz_type"] = "openservice"
	m["pid"] = this.Pid
	m["product_id"] = "APP_FAST_LOGIN"
	m["scope"] = "kuaijie"
	m["target_id"] = this.TargetId
	m["auth_type"] = this.AuthType
	return m
}

// OpenAuthAppAuthInviteCreate ISV向商户发起应用授权邀约 https://opendocs.alipay.com/isv/06evao?pathHash=f46ecafa
type OpenAuthAppAuthInviteCreate struct {
	AuxParam
	AppAuthToken string `json:"-"`            // 可选
	AuthAppId    string `json:"auth_app_id"`  // 必选 指定授权的商户appid
	RedirectURL  string `json:"redirect_url"` // 可选 授权回调地址，用于返回应用授权码
	State        string `json:"state"`        // 可选 自定义参数，授权后回调时透传回服务商。对应的值必须为 base64 编码。
}

func (this OpenAuthAppAuthInviteCreate) APIName() string {
	return "alipay.open.auth.appauth.invite.create"
}

func (this OpenAuthAppAuthInviteCreate) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = this.AppAuthToken
	return m
}
