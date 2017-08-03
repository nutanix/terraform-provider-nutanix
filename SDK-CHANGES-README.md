The following changes should be made in your SDK repo.

1. Added insecure flag in APIClient struct (api_client.go):

```
type APIClient struct {
    Insecure bool
}
```

2.A new argument is supplied to "prepareRequest" function which represents the insecure flag of APIClient. And correspondingly, request is modified to skip verification while transport.Changes are made at two places in api_client.go

```
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

3.In VmResources, struct field "MemorySizeMib" is modified as follows (vm_resources.go)_

```
MemorySizeMib int64 `json:"memory_size_mb,omitempty" bson:"memory_size_mb,omitempty"`
```

4.Changed the "Ip" of "IpAddress" struct to "Address" as Get API returns json with ip in the address field (ip_address.go)

```
// Address string.
Address string `json:"ip,omitempty" bson:"ip,omitempty"`
```


