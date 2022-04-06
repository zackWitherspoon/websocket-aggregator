# polygon-websocket-aggregator
This project is my implementation of a 30 second aggregator to print a specific ticker's aggregate info.

## How to run this project
In order to run this project please follow the following steps: 

#### Check your go version 
````
//Check your go version
go version

// This project was built with go version 1.18
````

#### Clone the git repo
* the repo's url is https://github.com/zackWitherspoon/polygon-websocket-aggregator

#### To build the app
````
// Install dependencies
go mod download

//Add your API key under service/web_socket_service.go

// Build the application
go build main.go
````

#### To test the app
````
// Install the ginkgo CLI
go get -u github.com/onsi/ginkgo/ginkgo

// Install the gomega library for assertions
go get github.com/onsi/gomega/...

// run all ginkgo test in the directory
ginkgo -r -skipPackage=integration_test

// run all ginkgo test with coverage
ginkgo -r --cover
````

#### To run the app
````
// execute the following
go run main.go TICKER_NAME_HERE 
````