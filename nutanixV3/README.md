# Go API client for nutanix

## Overview
This GO package is automatically generated, Details of this project:

- API version: 3.0.0
- Package version: 1.0.0
- Build date: 2017-06-19T20:38:34.869+05:30
- Build package: nutanix

## Requirements.
Go : go version >= 1.7.4

## Download
You can download this Go API client [here](sdk/go_sdk.tgz).

## Installation & Usage-Go
Download and extract gosdk packages in src/nutanix folder. To install the package use command :

go get -v ./...


~~~ go
import (
        "nutanix"
)
~~~

## Getting Started
Please follow the [Installation & Usage](#go-api-client-for-nutanix-installation--usage-go) and then run the following:

> Sample Code:

~~~ go
import (
        "nutanix"
)
func main() {

    nutanix.Configuration.Username := "YOUR_USERNAME"
    nutanix.Configuration.Password := "YOUR_PASSWORD"
	nutanix.Configuration.BasePath := "URL"

    // create an instance of the API class
    api_instance := nutanix.AppblueprintApi
    getEntitiesRequest := nutanix.AppBlueprintListMetadata() // AppBlueprintListMetadata |

    // List the App Blueprints
    appblueprintslistpost_response, , api_response, err := api_instance.AppBlueprintsListPost(getEntitiesRequest)
    if err != nil {
        fmt.Println("Error : ", err)
    } else {
        fmt.Println("Api response : ", *api_response)
        fmt.Println("AppBlueprintsListPost response : ", *appblueprintslistpost_response)
    }
}
~~~
