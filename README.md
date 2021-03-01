# probicon
Simple and flexible probe / healthcheck utility

Probicon's primary function is to probe an endpoint address (currently only HTTP/S is supported) and determine whether that address is alive and well or not.

Probicon can run in 3 different modes:
1. Console Full
2. Console Simple
3. Server

*Full console* mode is appropriate when a human is running probicon and evaluating the output in a terminal. 

*Simple console* mode is useful when probicon is invoked by another program, or when the output needs to be easily parsed. The output is sent to stdout as comma
separated value(s) that can be easily piped or redirected to other programs (e.g. tail, sed, awk, etc.).

*Server* mode enables probicon's built in HTTP server which listens for GET requests and then executes the probe function based on query string parameters. 
The results are then returned in the response as a JSON structured value.

## Using probicon

Probicon can be executed from the command line with various arguments and flags appended.

To see what arguments are available, simply run:
`probicon -h`

To get specific information about an argument, run:
`probicon probe -h` or `probicon serve -h`

Only the `-a` or `--address` flag is required when using the `probe` argument.

To use console full console mode, invoke probicon from the command line using the `probe` argument and pass in an address. For example:

`probicon probe --address example.com`

Probicon will then execute a health check on the given address and report back the results. In this simple example, probicon is only checking that the address
that you gave it is responding. However, probicon is capable of check several other things as well to determine if the endpoint is healthy. For instance, running:

`probicon probe --address example.com --code 200 --value login`

will instruct probicon to not only check that the address is responding, but that it responds with an HTTP status code of 200 and that the word `login` is present
in the content of the response. This is useful in cases where a web server may be responding, but it's responding with a 500 error message and should be considered
unhealthy or down.

Probicon can be instructed to conduct its health check multiple times, or even forever, and you can include a delay in between those checks as well.

For example:

`probicon probe --address example.com --repeat 3 --delay 10`

will instruct probicon to repeat its check of example.com 3 times with a 10 second delay before each check. Passing a flag of `--repeat -1` will instruct 
probicon to continuously check the address until asked to stop (e.g. via ctrl+c).

## Serve Mode

When running probicon with the `serve` argument (i.e. `probicon serve`), an HTTP server will be started which will listen for requests on all local IP addresses
for requests to the /probe handler.

The /probe handler expects a GET request with, at a minimum, an `address` query string parameter.

The following are valid query string parameters:

| Param        	| Value Type 	| Description                                                  	|   	|   	|
|--------------	|------------	|--------------------------------------------------------------	|---	|---	|
| address      	| string     	| The address to run the health check against                  	|   	|   	|
| timeout      	| integer    	| The value, in seconds, to wait before timing the request out 	|   	|   	|
| expectsvalue 	| string     	| An expected value to be found in the response body           	|   	|   	|
| expectscode  	| integer    	| Expected HTTP status code, such as 200                       	|   	|   	|

Example handler usage:

`http://127.0.0.1:8080/probe?address=http://example.com&timeout=10&expectscode=200`

Example handler response:

`{"IsDown":false,"ErrorMessage":"","ResponseTime":304,"ResponseCode":"200 OK"}`
