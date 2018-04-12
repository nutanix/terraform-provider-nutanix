package handler

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestUnmarshallErrorHandler(t *testing.T) {
	data := ` <Response>
						<Errors>
						<Error>
						<Code where='Code'>MissingParameter</Code>
						<Message where='Message'>Mensaje</Message>
						</Error>
						</Errors>
						<RequestID where='RequesID'></RequestID>
						</Response>`
	test := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte(data))),
	}

	if err := UnmarshalErrorHandler(test); err == nil {

		t.Fatalf("err: %s", err)
	}
}
