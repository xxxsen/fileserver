package s3base

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/xxxsen/common/errs"
)

//FROM: https://github.com/dovefi/go_aws4/blob/20d83dfd088ad1bd4a4723cac30e45b8bfcce224/signature.go#L182

const (
	iSO8601FormatDateTime = "20060102T150405Z" // 时间原点
	iSO8601FormatDate     = "20060102"
)

const (
	v4CredPrefix         = "Credential="
	v4SignedHeaderPrefix = "SignedHeaders="
	v4SignaturePrefix    = "Signature="
	v4SignAlgorithm      = "AWS4-HMAC-SHA256"
	s3Authorization      = "Authorization"
)

type S3V4Verify struct {
	r  *http.Request
	ak string
	sk string

	parsed *V4Signature
}

func NewS3V4Verify(r *http.Request, ak string, sk string, parsed *V4Signature) *S3V4Verify {
	return &S3V4Verify{
		r:      r,
		ak:     ak,
		sk:     sk,
		parsed: parsed,
	}
}

type V4Signature struct {
	AKey          string
	Date          string
	Region        string
	RequestType   string
	Service       string
	Algorithm     string
	SignedHeaders []string
	Signature     string
	Contentmd5    string
	Contentsha256 string
}

func ParseV4Signature(r *http.Request) (*V4Signature, bool, error) {
	auz := r.Header.Get(s3Authorization)
	if len(auz) == 0 {
		return nil, false, nil
	}
	items := strings.Split(auz, ",")
	if len(items) != 3 {
		return nil, false, errs.New(errs.ErrParam, "invalid authorization part, auz:%s", auz)
	}
	v4Data := &V4Signature{}
	credPart := strings.TrimSpace(items[0])
	if err := v4Data.parseCredPart(credPart); err != nil {
		return nil, false, errs.Wrap(errs.ErrParam, "decode cred part fail", err)
	}
	signedHeaderPart := strings.TrimSpace(items[1])
	if err := v4Data.parseSignedHeaderPart(signedHeaderPart); err != nil {
		return nil, false, errs.Wrap(errs.ErrParam, "decode signature header part fail", err)
	}
	signaturePart := strings.TrimSpace(items[2])
	if err := v4Data.parseSignaturePart(signaturePart); err != nil {
		return nil, false, errs.Wrap(errs.ErrParam, "decode signature part fail", err)
	}
	if err := v4Data.parseExtraPart(r); err != nil {
		return nil, false, errs.Wrap(errs.ErrParam, "decode extra part fail", err)
	}
	return v4Data, true, nil
}

func (d *V4Signature) parseExtraPart(r *http.Request) error {
	d.Contentmd5 = r.Header.Get("Content-Md5")
	d.Contentsha256 = r.Header.Get("X-Amz-Content-Sha256")
	d.Date = r.Header.Get("X-Amz-Date")
	return nil
}

func (d *V4Signature) parseCredPart(part string) error {
	items := strings.Split(part, " ")
	if len(items) != 2 {
		return errs.New(errs.ErrParam, "invalid cred part, part:%s", part)
	}
	algo := strings.TrimSpace(items[0])
	d.Algorithm = algo
	cred := strings.TrimSpace(items[1])
	if !strings.HasPrefix(cred, v4CredPrefix) {
		return errs.New(errs.ErrParam, "invalid cred prefix, cred:%s", cred)
	}
	cred = cred[len(v4CredPrefix):]
	parts := strings.Split(cred, "/")
	if len(parts) != 5 {
		return errs.New(errs.ErrParam, "invalid cred part, need 5, part:%s", cred)
	}
	d.AKey = parts[0]
	d.Date = parts[1]
	d.Region = parts[2]
	d.Service = parts[3]
	d.RequestType = parts[4]
	return nil
}

func (d *V4Signature) parseSignedHeaderPart(part string) error {
	part = strings.TrimSpace(part)
	if !strings.HasPrefix(part, v4SignedHeaderPrefix) {
		return errs.New(errs.ErrParam, "invalid sign header, should containsm signed headers prefix, part:%s", part)
	}
	part = part[len(v4SignedHeaderPrefix):]
	headers := strings.Split(part, ";")
	for _, h := range headers {
		d.SignedHeaders = append(d.SignedHeaders, strings.TrimSpace(h))
	}
	return nil
}

func (d *V4Signature) parseSignaturePart(part string) error {
	if !strings.HasPrefix(part, v4SignaturePrefix) {
		return errs.New(errs.ErrParam, "invalid signature prefix, signature:%s", part)
	}
	part = part[len(v4SignaturePrefix):]
	d.Signature = strings.TrimSpace(part)
	return nil
}

// 规范uri
func (s *S3V4Verify) canonicalURI(r *http.Request) string {
	path := r.URL.RequestURI()
	if r.URL.RawQuery != "" {
		path = path[:len(path)-len(r.URL.RawQuery)-1]
	}
	slash := strings.HasSuffix(path, "/")
	path = filepath.Clean(path)
	if path != "/" && slash {
		path += "/"
	}
	return path
}

// 构建规范的请求字符串
func (s *S3V4Verify) canonicalRequest(r *http.Request) (string, error) {
	var err error
	httpMethod := strings.ToUpper(r.Method)
	canonURI := s.canonicalURI(r)
	signHeader := s.canonicalSignHeaders(r)

	canonQuery, err := s.canonicalQueryString(r)
	if err != nil {
		return "", err
	}

	canonHeader, err := s.canonicalHeaders(r)
	if err != nil {
		return "", err
	}
	payload, err := s.hashPayload(r)
	if err != nil {
		return "", errs.Wrap(errs.ErrIO, "hash payload fail", err)
	}
	return httpMethod + "\n" +
		canonURI + "\n" +
		canonQuery + "\n" +
		canonHeader + "\n" +
		signHeader + "\n" +
		payload, nil
}

// 将秘钥加入到sign中
func (s *S3V4Verify) signKSecret(t time.Time) []byte {
	kDate := s.gHmac([]byte("AWS4"+s.sk), []byte(t.Format(iSO8601FormatDate)))
	kRegion := s.gHmac(kDate, []byte(s.parsed.Region))
	kService := s.gHmac(kRegion, []byte(s.parsed.Service))
	kSigning := s.gHmac(kService, []byte(s.parsed.RequestType))
	return kSigning
}

func (s *S3V4Verify) gSha256(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (s *S3V4Verify) gHmac(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

// 规范化header
func (s *S3V4Verify) canonicalHeaders(r *http.Request) (string, error) {
	// 2个以上的空格替换为一个空格，这个很容易忘记处理
	re, err := regexp.Compile(`\s{2,}`)
	if err != nil {
		return "", err
	}
	// 将header的值转为字符串，字符串内连续的空格替换为一个空格
	canonValue := func(values []string) (string, error) {
		var vals string
		for i, v := range values {
			vals += strings.TrimSpace(re.ReplaceAllString(v, " "))
			if i < (len(values) - 1) {
				vals += "; "
			}
		}
		return vals, nil
	}
	keys := s.parsed.SignedHeaders
	var canonHeader string
	for _, k := range keys {
		v := r.Header.Values(k)
		if len(v) == 0 && strings.EqualFold(k, "host") {
			v = []string{r.Host}
		}
		canVal, err := canonValue(v)
		if err != nil {
			return "", err
		}
		canonHeader += strings.ToLower(k) + ":" + canVal + "\n"
	}
	return canonHeader, nil
}

// 这里将加入到签名的header名称拼凑起来
// 目的是为了让服务端知道是基于哪几个header进行签名
func (s *S3V4Verify) canonicalSignHeaders(r *http.Request) string {
	return strings.Join(s.parsed.SignedHeaders, ";")
}

// 将url中的请求参数按照参数名
// 升序排序, 如果一个参数有多个值，按照值排序
func (s *S3V4Verify) canonicalQueryString(r *http.Request) (string, error) {
	err := r.ParseForm()
	if err != nil {
		return "", err
	}
	queryString := r.Form
	kvs := make([]string, 0, len(queryString))
	for k, vs := range queryString {
		for _, v := range vs {
			if v == "" {
				kvs = append(kvs, url.QueryEscape(k))
			} else {
				kvs = append(kvs, url.QueryEscape(k)+"="+url.QueryEscape(v))
			}
		}
	}
	sort.Strings(kvs)
	return strings.Join(kvs, "&"), nil
}

// 得出最终的签名结果
func (s *S3V4Verify) finalSign(t time.Time, r *http.Request) (string, error) {
	signKey := s.signKSecret(t)
	strToSign, err := s.stringToSign(t, r)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", s.gHmac(signKey, []byte(strToSign))), nil
}

// 创建待签名字符串
func (s *S3V4Verify) stringToSign(t time.Time, r *http.Request) (string, error) {
	credentialScope := fmt.Sprintf("%s/%s/%s/%s",
		t.Format(iSO8601FormatDate),
		s.parsed.Region, s.parsed.Service,
		s.parsed.RequestType)

	hr, err := s.hashRequest(r)
	if err != nil {
		return "", err
	}

	return s.parsed.Algorithm + "\n" +
		t.Format(iSO8601FormatDateTime) + "\n" +
		credentialScope + "\n" +
		hr, nil
}

func (s *S3V4Verify) hashRequest(r *http.Request) (string, error) {
	canonReq, err := s.canonicalRequest(r)
	if err != nil {
		return "", err
	}
	return s.gSha256([]byte(canonReq)), nil
}

// 规范request body 进行sha256哈希
func (s *S3V4Verify) hashPayload(r *http.Request) (string, error) {
	return s.parsed.Contentsha256, nil
}

func (s *S3V4Verify) authorization(t time.Time, r *http.Request) (string, error) {
	var err error
	credentialScope := fmt.Sprintf("%s/%s/%s/%s",
		t.Format(iSO8601FormatDate),
		s.parsed.Region, s.parsed.Service,
		s.parsed.RequestType)

	signHeaders := s.canonicalSignHeaders(r)
	sign, err := s.finalSign(t, r)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		s.parsed.Algorithm,
		s.ak,
		credentialScope,
		signHeaders,
		sign), nil
}

func (s *S3V4Verify) createSignature(r *http.Request) (string, error) {
	t, err := time.Parse(iSO8601FormatDateTime, s.parsed.Date)
	if err != nil {
		return "", err
	}

	auth, err := s.authorization(t, r)
	if err != nil {
		return "", err
	}
	return auth, nil
}

func (s *S3V4Verify) Verify() (bool, error) {
	if s.parsed.Algorithm != v4SignAlgorithm {
		return false, errs.New(errs.ErrParam, "algo not match, need:%s, get:%s", v4SignAlgorithm, s.parsed.Algorithm)
	}
	sign, err := s.createSignature(s.r)
	if err != nil {
		return false, err
	}
	if sign != s.r.Header.Get(s3Authorization) {
		return false, nil
	}
	return true, nil
}

func IsRequestSignatureV4(r *http.Request) bool {
	auz := r.Header.Get(s3Authorization)
	if len(auz) == 0 {
		return false
	}
	if !strings.HasPrefix(auz, v4SignAlgorithm) {
		return false
	}
	return true
}

func S3AuthV4(r *http.Request, ak string, sk string, parsed *V4Signature) (bool, error) {
	s := &S3V4Verify{
		r:      r,
		ak:     ak,
		sk:     sk,
		parsed: parsed,
	}
	return s.Verify()
}
