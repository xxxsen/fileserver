package s3base

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/xxxsen/common/errs"
)

// Whitelist resource list that will be used in query string for signature-V2 calculation.
//
// This list should be kept alphabetically sorted, do not hastily edit.
var resourceList = []string{
	"acl",
	"cors",
	"delete",
	"encryption",
	"legal-hold",
	"lifecycle",
	"location",
	"logging",
	"notification",
	"partNumber",
	"policy",
	"requestPayment",
	"response-cache-control",
	"response-content-disposition",
	"response-content-encoding",
	"response-content-language",
	"response-content-type",
	"response-expires",
	"retention",
	"select",
	"select-type",
	"tagging",
	"torrent",
	"uploadId",
	"uploads",
	"versionId",
	"versioning",
	"versions",
	"website",
}

const (
	authorization   = "Authorization"
	amzAccessKeyID  = "AWSAccessKeyId"
	signV2Algorithm = "AWS"
	date            = "Date"
	contentMD5      = "Content-Md5"
	contentType     = "Content-Type"
	signV4Algorithm = "AWS4-HMAC-SHA256"
)

func getReqAccessKeyV2(r *http.Request) (string, error) {
	if accessKey := r.Form.Get(amzAccessKeyID); accessKey != "" {
		return accessKey, nil
	}

	// below is V2 Signed Auth header format, splitting on `space` (after the `AWS` string).
	// Authorization = "AWS" + " " + AWSAccessKeyId + ":" + Signature
	authFields := strings.Split(r.Header.Get(authorization), " ")
	if len(authFields) != 2 {
		return "", errs.New(errs.ErrParam, "missing field")
	}

	// Then will be splitting on ":", this will seprate `AWSAccessKeyId` and `Signature` string.
	keySignFields := strings.Split(strings.TrimSpace(authFields[1]), ":")
	if len(keySignFields) != 2 {
		return "", errs.New(errs.ErrParam, "missing field")
	}

	return keySignFields[0], nil
}

func validateV2AuthHeader(r *http.Request) (string, error) {
	v2Auth := r.Header.Get(authorization)
	if v2Auth == "" {
		return "", errs.New(errs.ErrParam, "nil auth")
	}

	// Verify if the header algorithm is supported or not.
	if !strings.HasPrefix(v2Auth, signV2Algorithm) {
		return "", errs.New(errs.ErrParam, "signature version not supported")
	}

	ak, apiErr := getReqAccessKeyV2(r)
	if apiErr != nil {
		return "", apiErr
	}

	return ak, nil
}

func IsRequestSignatureV2(r *http.Request) bool {
	return (!strings.HasPrefix(r.Header.Get(authorization), signV4Algorithm) &&
		strings.HasPrefix(r.Header.Get(authorization), signV2Algorithm))
}

func S3AuthV2(r *http.Request, ak string, sk string) error {
	v2Auth := r.Header.Get(authorization)
	accessKey, apiError := validateV2AuthHeader(r)
	if apiError != nil {
		return errs.Wrap(errs.ErrParam, "api error", apiError)
	}
	if accessKey != ak {
		return errs.New(errs.ErrParam, "ak not match")
	}

	// r.RequestURI will have raw encoded URI as sent by the client.
	tokens := strings.SplitN(r.RequestURI, "?", 2)
	encodedResource := tokens[0]
	encodedQuery := ""
	if len(tokens) == 2 {
		encodedQuery = tokens[1]
	}

	unescapedQueries, err := unescapeQueries(encodedQuery)
	if err != nil {
		return errs.Wrap(errs.ErrParam, "invalid query param", err)
	}
	prefix := fmt.Sprintf("%s %s:", signV2Algorithm, accessKey)
	if !strings.HasPrefix(v2Auth, prefix) {
		return errs.New(errs.ErrParam, "signature not match")
	}
	v2Auth = v2Auth[len(prefix):]
	expectedAuth := signatureV2(ak, sk, r.Method, encodedResource, strings.Join(unescapedQueries, "&"), r.Header)
	if !compareSignatureV2(v2Auth, expectedAuth) {
		return errs.New(errs.ErrParam, "signature not match")
	}
	return nil
}

// Escape encodedQuery string into unescaped list of query params, returns error
// if any while unescaping the values.
func unescapeQueries(encodedQuery string) (unescapedQueries []string, err error) {
	for _, query := range strings.Split(encodedQuery, "&") {
		var unescapedQuery string
		unescapedQuery, err = url.QueryUnescape(query)
		if err != nil {
			return nil, err
		}
		unescapedQueries = append(unescapedQueries, unescapedQuery)
	}
	return unescapedQueries, nil
}

// compareSignatureV2 returns true if and only if both signatures
// are equal. The signatures are expected to be base64 encoded strings
// according to the AWS S3 signature V2 spec.
func compareSignatureV2(sig1, sig2 string) bool {
	// Decode signature string to binary byte-sequence representation is required
	// as Base64 encoding of a value is not unique:
	// For example "aGVsbG8=" and "aGVsbG8=\r" will result in the same byte slice.
	signature1, err := base64.StdEncoding.DecodeString(sig1)
	if err != nil {
		return false
	}
	signature2, err := base64.StdEncoding.DecodeString(sig2)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(signature1, signature2) == 1
}

func calculateSignatureV2(stringToSign string, secret string) string {
	hm := hmac.New(sha1.New, []byte(secret))
	hm.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(hm.Sum(nil))
}

// Return the signature v2 of a given request.
func signatureV2(ak, sk string, method string, encodedResource string, encodedQuery string, headers http.Header) string {
	stringToSign := getStringToSignV2(method, encodedResource, encodedQuery, headers, "")
	signature := calculateSignatureV2(stringToSign, sk)
	return signature
}

// Return canonical headers.
func canonicalizedAmzHeadersV2(headers http.Header) string {
	var keys []string
	keyval := make(map[string]string, len(headers))
	for key := range headers {
		lkey := strings.ToLower(key)
		if !strings.HasPrefix(lkey, "x-amz-") {
			continue
		}
		keys = append(keys, lkey)
		keyval[lkey] = strings.Join(headers[key], ",")
	}
	sort.Strings(keys)
	var canonicalHeaders []string
	for _, key := range keys {
		canonicalHeaders = append(canonicalHeaders, key+":"+keyval[key])
	}
	return strings.Join(canonicalHeaders, "\n")
}

// Return string to sign under two different conditions.
// - if expires string is set then string to sign includes date instead of the Date header.
// - if expires string is empty then string to sign includes date header instead.
func getStringToSignV2(method string, encodedResource, encodedQuery string, headers http.Header, expires string) string {
	canonicalHeaders := canonicalizedAmzHeadersV2(headers)
	if len(canonicalHeaders) > 0 {
		canonicalHeaders += "\n"
	}

	date := expires // Date is set to expires date for presign operations.
	if date == "" {
		// If expires date is empty then request header Date is used.
		date = headers.Get(date)
	}

	// From the Amazon docs:
	//
	// StringToSign = HTTP-Verb + "\n" +
	// 	 Content-Md5 + "\n" +
	//	 Content-Type + "\n" +
	//	 Date/Expires + "\n" +
	//	 CanonicalizedProtocolHeaders +
	//	 CanonicalizedResource;
	stringToSign := strings.Join([]string{
		method,
		headers.Get(contentMD5),
		headers.Get(contentType),
		date,
		canonicalHeaders,
	}, "\n")

	return stringToSign + canonicalizedResourceV2(encodedResource, encodedQuery)
}

// Return canonical resource string.
func canonicalizedResourceV2(encodedResource, encodedQuery string) string {
	queries := strings.Split(encodedQuery, "&")
	keyval := make(map[string]string)
	for _, query := range queries {
		key := query
		val := ""
		index := strings.Index(query, "=")
		if index != -1 {
			key = query[:index]
			val = query[index+1:]
		}
		keyval[key] = val
	}

	var canonicalQueries []string
	for _, key := range resourceList {
		val, ok := keyval[key]
		if !ok {
			continue
		}
		if val == "" {
			canonicalQueries = append(canonicalQueries, key)
			continue
		}
		canonicalQueries = append(canonicalQueries, key+"="+val)
	}

	// The queries will be already sorted as resourceList is sorted, if canonicalQueries
	// is empty strings.Join returns empty.
	canonicalQuery := strings.Join(canonicalQueries, "&")
	if canonicalQuery != "" {
		return encodedResource + "?" + canonicalQuery
	}
	return encodedResource
}
