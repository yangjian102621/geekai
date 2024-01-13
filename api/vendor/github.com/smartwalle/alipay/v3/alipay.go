package alipay

import (
	"context"
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/smartwalle/ngx"
	"github.com/smartwalle/nsign"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/smartwalle/ncrypto"
)

var (
	ErrBadResponse          = errors.New("alipay: bad response")
	ErrSignNotFound         = errors.New("alipay: sign content not found")
	ErrAliPublicKeyNotFound = errors.New("alipay: alipay public key not found")
)

const (
	kAliPayPublicKeySN = "alipay-public-key"
	kAppAuthToken      = "app_auth_token"
	kEmptyBizContent   = "{}"
)

type Client struct {
	mu               sync.Mutex
	isProduction     bool
	appId            string
	host             string
	notifyVerifyHost string
	Client           *http.Client
	location         *time.Location
	onReceivedData   func(method string, data []byte)

	// 内容加密
	needEncrypt    bool
	encryptIV      []byte
	encryptType    string
	encryptKey     []byte
	encryptPadding ncrypto.Padding

	appCertSN     string
	aliRootCertSN string
	aliCertSN     string

	// 签名和验签
	signer    Signer
	verifiers map[string]Verifier
}

type Signer interface {
	SignValues(values url.Values, opts ...nsign.SignOption) ([]byte, error)

	SignBytes(data []byte, opts ...nsign.SignOption) ([]byte, error)
}

type Verifier interface {
	VerifyValues(values url.Values, signature []byte, opts ...nsign.SignOption) error

	VerifyBytes(data []byte, signature []byte, opts ...nsign.SignOption) error
}

type OptionFunc func(c *Client)

func WithTimeLocation(location *time.Location) OptionFunc {
	return func(c *Client) {
		c.location = location
	}
}

func WithHTTPClient(client *http.Client) OptionFunc {
	return func(c *Client) {
		if client != nil {
			c.Client = client
		}
	}
}

func WithSandboxGateway(gateway string) OptionFunc {
	return func(c *Client) {
		if gateway == "" {
			gateway = kNewSandboxGateway
		}
		if !c.isProduction {
			c.host = gateway
		}
	}
}

func WithProductionGateway(gateway string) OptionFunc {
	return func(c *Client) {
		if gateway == "" {
			gateway = kProductionGateway
		}
		if c.isProduction {
			c.host = gateway
		}
	}
}

func WithNewSandboxGateway() OptionFunc {
	return WithSandboxGateway(kNewSandboxGateway)
}

func WithPastSandboxGateway() OptionFunc {
	return WithSandboxGateway(kPastSandboxGateway)
}

// New 初始化支付宝客户端
//
// appId - 支付宝应用 id
//
// privateKey - 应用私钥，开发者自己生成
//
// isProduction - 是否为生产环境，传 false 的时候为沙箱环境，用于开发测试，正式上线的时候需要改为 true
func New(appId, privateKey string, isProduction bool, opts ...OptionFunc) (nClient *Client, err error) {
	priKey, err := ncrypto.DecodePrivateKey([]byte(privateKey)).PKCS1().RSAPrivateKey()
	if err != nil {
		priKey, err = ncrypto.DecodePrivateKey([]byte(privateKey)).PKCS8().RSAPrivateKey()
		if err != nil {
			return nil, err
		}
	}
	nClient = &Client{}
	nClient.isProduction = isProduction
	nClient.appId = appId

	if nClient.isProduction {
		nClient.host = kProductionGateway
		nClient.notifyVerifyHost = kProductionMAPIGateway
	} else {
		nClient.host = kNewSandboxGateway
		nClient.notifyVerifyHost = kNewSandboxGateway
	}
	nClient.Client = http.DefaultClient
	nClient.location = time.Local

	nClient.signer = nsign.New(nsign.WithMethod(nsign.NewRSAMethod(crypto.SHA256, priKey, nil)))
	nClient.verifiers = make(map[string]Verifier)

	for _, opt := range opts {
		if opt != nil {
			opt(nClient)
		}
	}

	return nClient, nil
}

func (this *Client) IsProduction() bool {
	return this.isProduction
}

// SetEncryptKey 接口内容加密密钥 https://opendocs.alipay.com/common/02mse3
func (this *Client) SetEncryptKey(key string) error {
	if key == "" {
		this.needEncrypt = false
		return nil
	}

	var data, err = base64.StdEncoding.DecodeString(key)
	if err != nil {
		return err
	}
	this.needEncrypt = true
	this.encryptIV = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	this.encryptType = "AES"
	this.encryptKey = data
	this.encryptPadding = ncrypto.PKCS7Padding{}
	return nil
}

func (this *Client) loadVerifier(sn string, pub *rsa.PublicKey) Verifier {
	this.aliCertSN = sn
	var verifier = nsign.New(nsign.WithMethod(nsign.NewRSAMethod(crypto.SHA256, nil, pub)))
	this.verifiers[this.aliCertSN] = verifier
	return verifier
}

// LoadAliPayPublicKey 加载支付宝公钥
func (this *Client) LoadAliPayPublicKey(aliPublicKey string) error {
	var pub *rsa.PublicKey
	var err error
	if len(aliPublicKey) < 0 {
		return ErrAliPublicKeyNotFound
	}

	pub, err = ncrypto.DecodePublicKey([]byte(aliPublicKey)).PKIX().RSAPublicKey()
	if err != nil {
		return err
	}

	this.mu.Lock()
	this.loadVerifier(kAliPayPublicKeySN, pub)
	this.mu.Unlock()
	return nil
}

// LoadAppPublicCert 加载应用公钥证书
//
// Deprecated: use LoadAppCertPublicKey instead.
func (this *Client) LoadAppPublicCert(s string) error {
	return this.LoadAppCertPublicKey(s)
}

// LoadAppPublicCertFromFile 加载应用公钥证书
//
// Deprecated: use LoadAppCertPublicKeyFromFile instead.
func (this *Client) LoadAppPublicCertFromFile(filename string) error {
	return this.LoadAppCertPublicKeyFromFile(filename)
}

func (this *Client) loadAppCertPublicKey(b []byte) error {
	cert, err := ncrypto.DecodeCertificate(b)
	if err != nil {
		return err
	}
	this.appCertSN = getCertSN(cert)
	return nil
}

// LoadAppCertPublicKey 加载应用公钥证书
func (this *Client) LoadAppCertPublicKey(s string) error {
	return this.loadAppCertPublicKey([]byte(s))
}

// LoadAppCertPublicKeyFromFile 从文件加载应用公钥证书
func (this *Client) LoadAppCertPublicKeyFromFile(filename string) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return this.loadAppCertPublicKey(b)
}

// LoadAliPayPublicCert 加载支付宝公钥证书
//
// Deprecated: use LoadAlipayCertPublicKey instead.
func (this *Client) LoadAliPayPublicCert(s string) error {
	return this.LoadAlipayCertPublicKey(s)
}

// LoadAliPayPublicCertFromFile 加载支付宝公钥证书
//
// Deprecated: use LoadAlipayCertPublicKeyFromFile instead.
func (this *Client) LoadAliPayPublicCertFromFile(filename string) error {
	return this.LoadAlipayCertPublicKeyFromFile(filename)
}

// loadAlipayCertPublicKey 加载支付宝公钥证书
func (this *Client) loadAlipayCertPublicKey(b []byte) error {
	cert, err := ncrypto.DecodeCertificate(b)
	if err != nil {
		return err
	}
	pub, ok := cert.PublicKey.(*rsa.PublicKey)
	if ok == false {
		return nil
	}

	this.mu.Lock()
	this.loadVerifier(getCertSN(cert), pub)
	this.mu.Unlock()
	return nil
}

// LoadAlipayCertPublicKey 支付宝公钥证书
func (this *Client) LoadAlipayCertPublicKey(s string) error {
	return this.loadAlipayCertPublicKey([]byte(s))
}

// LoadAlipayCertPublicKeyFromFile 从文件支付宝公钥证书
func (this *Client) LoadAlipayCertPublicKeyFromFile(filename string) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return this.loadAlipayCertPublicKey(b)
}

// LoadAliPayRootCert 加载支付宝根证书
func (this *Client) LoadAliPayRootCert(s string) error {
	var certStrList = strings.Split(s, kCertificateEnd)
	var certSNList = make([]string, 0, len(certStrList))
	for _, certStr := range certStrList {
		certStr = strings.Replace(certStr, kCertificateBegin, "", 1)
		var cert, _ = ncrypto.DecodeCertificate([]byte(certStr))
		if cert != nil && (cert.SignatureAlgorithm == x509.SHA256WithRSA || cert.SignatureAlgorithm == x509.SHA1WithRSA) {
			certSNList = append(certSNList, getCertSN(cert))
		}
	}
	this.aliRootCertSN = strings.Join(certSNList, "_")
	return nil
}

// LoadAliPayRootCertFromFile 加载支付宝根证书
func (this *Client) LoadAliPayRootCertFromFile(filename string) error {
	b, err := os.ReadFile(filename)

	if err != nil {
		return err
	}

	return this.LoadAliPayRootCert(string(b))
}

func (this *Client) URLValues(param Param) (value url.Values, err error) {
	var values = url.Values{}
	values.Add(kFieldAppId, this.appId)
	values.Add(kFieldMethod, param.APIName())
	values.Add(kFieldFormat, kFormat)
	values.Add(kFieldCharset, kCharset)
	values.Add(kFieldSignType, kSignTypeRSA2)
	values.Add(kFieldTimestamp, time.Now().In(this.location).Format(kTimeFormat))
	values.Add(kFieldVersion, kVersion)
	if this.appCertSN != "" {
		values.Add(kFieldAppCertSN, this.appCertSN)
	}
	if this.aliRootCertSN != "" {
		values.Add(kFieldAliPayRootCertSN, this.aliRootCertSN)
	}

	jsonBytes, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	var content = string(jsonBytes)
	if content != kEmptyBizContent {
		if this.needEncrypt && param.NeedEncrypt() {
			jsonBytes, err = ncrypto.AESCBCEncrypt(jsonBytes, this.encryptKey, this.encryptIV, this.encryptPadding)
			if err != nil {
				return nil, err
			}
			content = base64.StdEncoding.EncodeToString(jsonBytes)
			values.Add(kFieldEncryptType, this.encryptType)
		}
		values.Add(kFieldBizContent, content)
	}

	var params = param.Params()
	for k, v := range params {
		if k == kAppAuthToken && v == "" {
			continue
		}
		values.Add(k, v)
	}

	signature, err := this.sign(values)
	if err != nil {
		return nil, err
	}

	values.Add(kFieldSign, signature)
	return values, nil
}

func (this *Client) doRequest(method string, param Param, result interface{}) (err error) {
	var req = ngx.NewRequest(method, this.host, ngx.WithClient(this.Client))
	req.SetContentType(kContentType)
	if param != nil {
		var values url.Values
		values, err = this.URLValues(param)
		if err != nil {
			return err
		}
		req.SetParams(values)

		var files = param.FileParams()
		for _, file := range files {
			req.AddFile(file.Name, file.Filename, file.Filepath)
		}
	}

	rsp, err := req.Do(context.Background())
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	bodyBytes, err := io.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	var apiName = param.APIName()
	var bizFieldName = strings.Replace(apiName, ".", "_", -1) + kResponseSuffix

	return this.decode(bodyBytes, bizFieldName, param.NeedVerify(), result)
}

func (this *Client) decode(data []byte, bizFieldName string, needVerifySign bool, result interface{}) (err error) {
	var raw = make(map[string]json.RawMessage)
	if err = json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var signBytes = raw[kFieldSign]
	var certBytes = raw[kFieldAlyPayCertSN]
	var bizBytes = raw[bizFieldName]
	var errBytes = raw[kErrorResponse]

	if len(certBytes) > 1 {
		certBytes = certBytes[1 : len(certBytes)-1]
	}
	if len(signBytes) > 1 {
		signBytes = signBytes[1 : len(signBytes)-1]
	}

	if len(bizBytes) == 0 {
		if len(errBytes) > 0 {
			var rErr *Error
			if err = json.Unmarshal(errBytes, &rErr); err != nil {
				return err
			}
			return rErr
		}
		return ErrBadResponse
	}

	// 数据解密
	var plaintext []byte
	if plaintext, err = this.decrypt(bizBytes); err != nil {
		return err
	}

	// 验证签名
	if needVerifySign {
		if this.onReceivedData != nil {
			this.onReceivedData(bizFieldName, plaintext)
		}

		if len(signBytes) == 0 {
			// 没有签名数据，返回的内容一般为错误信息
			var rErr *Error
			if err = json.Unmarshal(plaintext, &rErr); err != nil {
				return err
			}
			return rErr
		}

		// 验证签名
		if err = this.verify(string(certBytes), bizBytes, signBytes); err != nil {
			return err
		}
	}

	if err = json.Unmarshal(plaintext, result); err != nil {
		return err
	}
	return nil
}

func (this *Client) decrypt(data []byte) ([]byte, error) {
	var plaintext = data
	if len(data) > 1 && data[0] == '"' {
		var ciphertext, err = base64decode(data[1 : len(data)-1])
		if err != nil {
			return nil, err
		}
		plaintext, err = ncrypto.AESCBCDecrypt(ciphertext, this.encryptKey, this.encryptIV, this.encryptPadding)
		if err != nil {
			return nil, err
		}
	}
	return plaintext, nil
}

func (this *Client) VerifySign(values url.Values) (err error) {
	var verifier Verifier
	if verifier, err = this.getVerifier(values.Get(kFieldAlyPayCertSN)); err != nil {
		return err
	}

	var signBytes []byte
	if signBytes, err = base64.StdEncoding.DecodeString(values.Get(kFieldSign)); err != nil {
		return err
	}

	return verifier.VerifyValues(values, signBytes, nsign.WithIgnore(kFieldSign, kFieldSignType, kFieldAlyPayCertSN))
}

func (this *Client) getVerifier(certSN string) (verifier Verifier, err error) {
	this.mu.Lock()
	defer this.mu.Unlock()

	if certSN == "" {
		certSN = this.aliCertSN
	}

	verifier = this.verifiers[certSN]

	if verifier == nil {
		if !this.isProduction {
			return nil, ErrAliPublicKeyNotFound
		}

		cert, err := this.downloadAliPayCert(certSN)
		if err != nil {
			return nil, err
		}

		pub, ok := cert.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, ErrAliPublicKeyNotFound
		}
		verifier = this.loadVerifier(getCertSN(cert), pub)
	}
	return verifier, nil
}

func (this *Client) CertDownload(param CertDownload) (result *CertDownloadRsp, err error) {
	err = this.doRequest(http.MethodPost, param, &result)
	return result, err
}

func (this *Client) downloadAliPayCert(certSN string) (cert *x509.Certificate, err error) {
	var param = CertDownload{}
	param.AliPayCertSN = certSN
	rsp, err := this.CertDownload(param)
	if err != nil {
		return nil, err
	}
	certBytes, err := base64.StdEncoding.DecodeString(rsp.AliPayCertContent)
	if err != nil {
		return nil, err
	}

	cert, err = ncrypto.DecodeCertificate(certBytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func (this *Client) sign(values url.Values) (signature string, err error) {
	sBytes, err := this.signer.SignValues(values)
	if err != nil {
		return "", err
	}
	signature = base64.StdEncoding.EncodeToString(sBytes)
	return signature, nil
}

func (this *Client) verify(certSN string, data, signature []byte) (err error) {
	var verifier Verifier
	if verifier, err = this.getVerifier(certSN); err != nil {
		return err
	}

	if signature, err = base64decode(signature); err != nil {
		return err
	}

	if err = verifier.VerifyBytes(data, signature); err != nil {
		return err
	}
	return nil
}

func (this *Client) Request(payload *Payload, result interface{}) (err error) {
	return this.doRequest(http.MethodPost, payload, result)
}

func (this *Client) OnReceivedData(fn func(method string, data []byte)) {
	this.onReceivedData = fn
}

func base64decode(data []byte) ([]byte, error) {
	var dBuf = make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(dBuf, data)
	return dBuf[:n], err
}

func getCertSN(cert *x509.Certificate) string {
	var value = md5.Sum([]byte(cert.Issuer.String() + cert.SerialNumber.String()))
	return hex.EncodeToString(value[:])
}
