The following changes are made in go_sdk repo.

1. Added insecure flag in APIClient struct (api_client.go):

```
type APIClient struct {
    Insecure bool
}
```

2.A new argument is supplied to "prepareRequest" function which represents the insecure flag of APIClient. And correspondingly, request is modified to skip verification while transport.Changes are made at two places in api_client.go

```
//Add the following in imported modules.
import(
    ...
    "crypto/tls"
    ...
)

//In CallAPI function where prepareRequest is called
request := prepareRequest(postBody, headerParams, queryParams, formParams, fileName, fileBytes,c.Insecure)

//In function definition
func prepareRequest(postBody interface{},headerParams map[string]string,queryParams url.Values,formParams map[string]string,fileName string,fileBytes []byte,insecure bool) *resty.Request {
    var jsontype interface{}
    request := resty.R()
    if insecure==true {
        restycli := resty.New()
        config := &tls.Config{InsecureSkipVerify: true}
        restycli = restycli.SetTLSClientConfig(config)
        request = restycli.R()
    }
    .....
    .....
}

```

3. In NewConfiguration function, APIClient is initialised with Insecure flag as false.(configuration.go)

```
func NewConfiguration() *Configuration {
    return &Configuration{
        ...
        ...
        APIClient:     APIClient{Insecure: false},
    }
}

```

