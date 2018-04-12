package handler

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestUnmarshalXML(t *testing.T) {
	buf := bytes.NewReader([]byte("<OperationNameResponse><Str>myname</Str><FooNum>123</FooNum><FalseBool>false</FalseBool><TrueBool>true</TrueBool><Float>1.2</Float><Double>1.3</Double><Long>200</Long><Char>a</Char><RequestId>request-id</RequestId></OperationNameResponse>"))
	res := &http.Response{StatusCode: 200, Body: ioutil.NopCloser(buf), Header: http.Header{}}
	var v struct{}
	if err := UnmarshalXML(&v, res); err != nil {
		t.Fatalf("err: %s", err)
	}
}
