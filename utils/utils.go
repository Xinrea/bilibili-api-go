package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
)

var SIGN_TABLE = []int{
	46, 47, 18, 2, 53, 8, 23, 32, 15, 50, 10, 31, 58, 3, 45, 35, 27, 43, 5, 49, 33, 9, 42,
	19, 29, 28, 14, 39, 12, 38, 41, 13, 37, 48, 7, 16, 24, 55, 40, 61, 26, 17, 0, 1, 60,
	51, 30, 4, 22, 25, 54, 21, 56, 59, 6, 63, 57, 62, 11, 36, 20, 34, 44, 52,
}

// SignParams using img and sub which comes from nav api to sign params
func SignParams(img string, sub string, ts int64, params *url.Values) {
	originStr := img + sub
	encodedBytes := []byte{}
	for _, i := range SIGN_TABLE {
		if i < len(originStr) {
			encodedBytes = append(encodedBytes, originStr[i])
		}
	}
	if len(encodedBytes) > 32 {
		encodedBytes = encodedBytes[:32]
	}
	params.Add("wts", fmt.Sprint(ts))
	paramKeys := []string{}
	for key := range *params {
		paramKeys = append(paramKeys, key)
	}
	sort.Strings(paramKeys)
	paramBuilder := []byte{}
	// Concat with ordered key
	for i, key := range paramKeys {
		paramBuilder = append(paramBuilder, []byte(key)...)
		paramBuilder = append(paramBuilder, '=')
		paramBuilder = append(paramBuilder, []byte(url.QueryEscape(params.Get(key)))...)
		if i != len(paramKeys)-1 {
			paramBuilder = append(paramBuilder, '&')
		}
	}
	paramBuilder = append(paramBuilder, encodedBytes...)
	md5bytes := md5.Sum(paramBuilder)
	signed := hex.EncodeToString(md5bytes[:])
	params.Add("w_rid", signed)
}
