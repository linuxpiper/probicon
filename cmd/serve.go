/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/keithnyc/probicon/serve"

	"github.com/spf13/cobra"
)

var serveConf serve.ServerConfig

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run as HTTP Server",
	Long: `
Run as an HTTP Server and launch probes when invoked. 

When running in serve mode, probicon will listen for requests to /probe and, when invoked
with a GET, will execute a probe and return with the status of the specified endpoint address.

The GET request to /probe should include an address at a bare minimum as part of the query string. Valid
query string parameters are:
	address      - Endpoint to probe. Required.
	timeout      - Timeout in seconds for the address response. Defaults to 10 seconds.
	expectscode  - Expected HTTP status code (e.g. 200). Ignored if not specified.
	expectsvalue - Expected case-sensitive string value to be found within the response body. Ignored if not specified.

The /probe handler returns a simple JSON document in the following format:
{
	"IsDown":<true / false>,
	"ErrorMessage":<message string>,
	"ResponseTime":<response time in milliseconds>,
	"ResponseCode":<http status code>
}

Return Example:
{"IsDown":false,"ErrorMessage":"","ResponseTime":411,"ResponseCode":"200 OK"}

This example shows that the endpoint was up, had no error messages, responded in 411
milliseconds and returned a 200 OK HTTP status code.

Handler Examples:

Run a simple probe to see if the endpoint responds:
http://127.0.0.1:8080/probe?address=https://google.com

Run a simple probe to see if the endpoing response within 20 seconds:
http://127.0.0.1:8080/probe?address=https://google.com&timeout=20

Run a probe and ensure HTTP 200 is returned within 20 seconds
http://127.0.0.1:8080/probe?address=https://google.com&timeout=20&expectscode=200

Run a probe and ensure Username is found within the webpage body content:
http://127.0.0.1:8080/probe?address=https://google.com&timeout=20&expectscode=200&expectsvalue=Username

`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")
		serve.StartServer(&serveConf)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveConf = serve.ServerConfig{}
	serveCmd.Flags().IntVarP(&serveConf.ListenPort, "listenport", "p", 8080, "Port to listen on")
}
