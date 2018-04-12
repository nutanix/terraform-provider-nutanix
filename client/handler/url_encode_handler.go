package handler

import (
	"net/url"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/private/protocol/query/queryutil"
)

// URLEncodeMarshalHander encodes the body to url encode
func URLEncodeMarshalHander(v interface{}, action, version string) (string, error) {
	body := url.Values{
		"Action":  {action},
		"Version": {version},
	}
	if err := queryutil.Parse(body, v, true); err != nil {
		return "", awserr.New("SerializationError", "failed encoding query request", err)
	}

	return body.Encode(), nil
}
