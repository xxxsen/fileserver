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
	wrap                  = "\n"
	unSignPayload         = "UNSIGNED-PAYLOAD" // 如果不对BODY进行哈希

)

type parsedData struct {
	ak            string
	date          string
	region        string
	requestType   string
	service       string
	algorithm     string
	signedHeaders []string
	signature     string
	contentmd5    string
	contentsha256 string
}

type s3VerifyV4 struct {
	r  *http.Request
	ak string
	sk string

	parsed *parsedData
}

func (s *s3VerifyV4) parseCred(part string) error {
	items := strings.Split(strings.TrimSpace(part), " ")
	if len(items) != 2 {
		return errs.New(errs.ErrParam, "invalid cred part, need 2, part:%s", part)
	}
	s.parsed.algorithm = strings.TrimSpace(items[0])

	cred := strings.TrimSpace(items[1])
	credPrefix := "Credential="
	if !strings.HasPrefix(cred, credPrefix) {
		return errs.New(errs.ErrParam, "invalid cred prefix, cred:%s", cred)
	}
	cred = cred[len(credPrefix):]

	parts := strings.Split(cred, "/")
	if len(parts) != 5 {
		return errs.New(errs.ErrParam, "invalid cred part, need 5, part:%s", cred)
	}
	s.parsed.ak = parts[0]
	s.parsed.date = parts[1]
	s.parsed.region = parts[2]
	s.parsed.service = parts[3]
	s.parsed.requestType = parts[4]
	return nil
}

func (s *s3VerifyV4) parseSignHeaders(part string) error {
	signPrefix := "SignedHeaders="
	part = strings.TrimSpace(part)
	if !strings.HasPrefix(part, signPrefix) {
		return errs.New(errs.ErrParam, "invalid sign header, should containsm signed headers prefix, part:%s", part)
	}
	part = part[len(signPrefix):]
	headers := strings.Split(part, ";")
	for _, h := range headers {
		s.parsed.signedHeaders = append(s.parsed.signedHeaders, strings.TrimSpace(h))
	}
	return nil
}

func (s *s3VerifyV4) parseSignature(part string) error {
	part = strings.TrimSpace(part)
	signaturePrefix := "Signature="
	if !strings.HasPrefix(part, signaturePrefix) {
		return errs.New(errs.ErrParam, "invalid signature prefix, signature:%s", part)
	}
	part = part[len(signaturePrefix):]
	s.parsed.signature = strings.TrimSpace(part)
	return nil
}

func (s *s3VerifyV4) parseV4(r *http.Request) error {
	auz := r.Header.Get(authorization)
	if len(auz) == 0 {
		return errs.New(errs.ErrParam, "not found authorization")
	}
	items := strings.Split(auz, ",")
	if len(items) != 3 {
		return errs.New(errs.ErrParam, "invalid authorization part, need:3, get:%d, auz:%s", len(items), auz)
	}
	s.parsed = &parsedData{}
	if err := s.parseCred(items[0]); err != nil {
		return errs.Wrap(errs.ErrParam, "decode cred fail", err)
	}
	if err := s.parseSignHeaders(items[1]); err != nil {
		return errs.Wrap(errs.ErrParam, "decode sign headers fail", err)
	}
	if err := s.parseSignature(items[2]); err != nil {
		return errs.Wrap(errs.ErrParam, "decode signature fail", err)
	}
	s.parsed.contentmd5 = r.Header.Get("Content-Md5")
	s.parsed.contentsha256 = r.Header.Get("X-Amz-Content-Sha256")
	s.parsed.date = r.Header.Get("X-Amz-Date")
	return nil
}

// 规范uri
func (s *s3VerifyV4) canonicalURI(r *http.Request) string {
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
func (s *s3VerifyV4) canonicalRequest(r *http.Request) (string, error) {
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
func (s *s3VerifyV4) signKSecret(t time.Time) []byte {
	kDate := s.gHmac([]byte("AWS4"+s.sk), []byte(t.Format(iSO8601FormatDate)))
	kRegion := s.gHmac(kDate, []byte(s.parsed.region))
	kService := s.gHmac(kRegion, []byte(s.parsed.service))
	kSigning := s.gHmac(kService, []byte(s.parsed.requestType))
	return kSigning
}

func (s *s3VerifyV4) gSha256(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (s *s3VerifyV4) gHmac(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

// 规范化header
func (s *s3VerifyV4) canonicalHeaders(r *http.Request) (string, error) {
	// 将header的值转为字符串，字符串内连续的空格替换为一个空格
	canonValue := func(values []string) (string, error) {
		var vals string
		// 2个以上的空格替换为一个空格，这个很容易忘记处理
		re, err := regexp.Compile("\\s{2,}")
		if err != nil {
			return "", err
		}

		for i, v := range values {
			vals += strings.TrimSpace(re.ReplaceAllString(v, " "))
			if i < (len(values) - 1) {
				vals += "; "
			}
		}
		return vals, nil
	}

	// 获取keys并排序
	keys := s.parsed.signedHeaders
	sort.Strings(keys)
	var canonHeader string
	for _, k := range keys {
		v := r.Header.Values(k)
		if strings.EqualFold(k, "host") {
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
func (s *s3VerifyV4) canonicalSignHeaders(r *http.Request) string {
	header := s.parsed.signedHeaders
	keys := make([]string, 0, len(header))
	for _, k := range header {
		keys = append(keys, strings.ToLower(k))
	}
	sort.Strings(keys)
	return strings.Join(keys, ";")
}

// 将url中的请求参数按照参数名
// 升序排序, 如果一个参数有多个值，按照值排序
func (s *s3VerifyV4) canonicalQueryString(r *http.Request) (string, error) {
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
func (s *s3VerifyV4) finalSign(t time.Time, r *http.Request) (string, error) {
	signKey := s.signKSecret(t)
	strToSign, err := s.stringToSign(t, r)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", s.gHmac(signKey, []byte(strToSign))), nil
}

// 创建待签名字符串
func (s *s3VerifyV4) stringToSign(t time.Time, r *http.Request) (string, error) {
	credentialScope := fmt.Sprintf("%s/%s/%s/%s",
		t.Format(iSO8601FormatDate),
		s.parsed.region, s.parsed.service,
		s.parsed.requestType)

	hr, err := s.hashRequest(r)
	if err != nil {
		return "", err
	}

	return s.parsed.algorithm + "\n" +
		t.Format(iSO8601FormatDateTime) + "\n" +
		credentialScope + "\n" +
		hr, nil
}

func (s *s3VerifyV4) hashRequest(r *http.Request) (string, error) {
	canonReq, err := s.canonicalRequest(r)
	if err != nil {
		return "", err
	}
	return s.gSha256([]byte(canonReq)), nil
}

// 规范request body 进行sha256哈希
func (s *s3VerifyV4) hashPayload(r *http.Request) (string, error) {
	if strings.EqualFold(s.parsed.contentsha256, unSignPayload) {
		// 如果request 的body 为空，那就用空字符串代替
		return s.gSha256([]byte("")), nil
	}
	return s.parsed.contentsha256, nil
}

func (s *s3VerifyV4) authorization(t time.Time, r *http.Request) (string, error) {
	var err error
	credentialScope := fmt.Sprintf("%s/%s/%s/%s",
		t.Format(iSO8601FormatDate),
		s.parsed.region, s.parsed.service,
		s.parsed.requestType)

	signHeaders := s.canonicalSignHeaders(r)
	sign, err := s.finalSign(t, r)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		s.parsed.algorithm,
		s.ak,
		credentialScope,
		signHeaders,
		sign), nil
}

func (s *s3VerifyV4) createSignature(r *http.Request) (string, error) {
	//t, err = time.Parse(http.TimeFormat, date)
	t, err := time.Parse(iSO8601FormatDateTime, s.parsed.date)
	if err != nil {
		return "", err
	}

	auth, err := s.authorization(t, r)
	if err != nil {
		return "", err
	}
	return auth, nil
}

func (s *s3VerifyV4) verify() (bool, error) {
	if s.parsed.algorithm != signV4Algorithm {
		return false, errs.New(errs.ErrParam, "algo not match, need:%s, get:%s", signV4Algorithm, s.parsed.algorithm)
	}
	sign, err := s.createSignature(s.r)
	if err != nil {
		return false, err
	}
	if sign != s.r.Header.Get(authorization) {
		return false, nil
	}
	return true, nil
}

func IsRequestSignatureV4(r *http.Request) bool {
	auz := r.Header.Get(authorization)
	if len(auz) == 0 {
		return false
	}
	if !strings.HasPrefix(auz, signV4Algorithm) {
		return false
	}
	return true
}

func S3AuthV4(r *http.Request, ak string, sk string) (bool, error) {
	s := &s3VerifyV4{
		r:  r,
		ak: ak,
		sk: sk,
	}
	if err := s.parseV4(r); err != nil {
		return false, errs.Wrap(errs.ErrParam, "parse signature fail", err)
	}
	return s.verify()
}
