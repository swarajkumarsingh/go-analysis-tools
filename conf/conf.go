package conf

import "fmt"

var Stage string = Dev
var Dev string = "dev"
var Prod string = "prod"
var PORT string = fmt.Sprintf("localhost:%s", "8080")

var ClientENV = ""
// DDAgentHost is Hostname for Datadog agent
var DDAgentHost string = "localhost"
const DDServiceName = "go-analysis-tool"