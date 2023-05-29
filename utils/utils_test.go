package utils_test

import (
	"bilibili-api-go/utils"
	"net/url"
	"testing"
)

func TestSignParams(t *testing.T) {
	// mid=475210&token=&platform=web&web_location=1550101&w_rid=222c5b33ef4f9e941071797e6f82a01d&wts=1685374758
	var ts int64 = 1685374758
	params := url.Values{}
	params.Add("mid", "475210")
	params.Add("token", "")
	params.Add("platform", "web")
	params.Add("web_location", "1550101")
	expected := "222c5b33ef4f9e941071797e6f82a01d"
	utils.SignParams("9a16e2304d794393b733badf79f09804", "69c7f4b06e3f44449a54ad06ca05676c", ts, &params)
	if params.Get("w_rid") != expected {
		t.Errorf("Expected %s, got %s", expected, params.Get("w_rid"))
	}
}
