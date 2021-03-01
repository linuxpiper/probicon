package probe

import (
	"fmt"
	"github.com/pterm/pterm"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ProbeConfig struct {
	Address       string
	SimpleOutput  bool
	Repeat        int
	Delay         int
	Timeout       int
	ExpectsCode   string
	ExpectsString string
}

type ProbeResponse struct {
	IsDown       bool   `json:"IsDown"`
	ErrorMessage string `json:"ErrorMessage"`
	ResponseTime int64  `json:"ResponseTime"`
	ResponseCode string `json:"ResponseCode"`
}

func StatusCheck(address string, timeoutSeconds int, delaySeconds int, expectsCode string, expectsValue string) ProbeResponse {
	var probeResponse = ProbeResponse{}
	address = sanitizeHttp(address)

	if delaySeconds > 0 {
		time.Sleep(time.Duration(delaySeconds) * time.Second)
	}

	start := time.Now()
	timeout := time.Duration(timeoutSeconds * int(time.Second))
	client := http.Client{Timeout: timeout}

	response, err := client.Get(address)

	elapsed := time.Since(start)

	if err != nil {
		probeResponse.IsDown = true
		probeResponse.ResponseCode = ""
		switch err := err.(type) {
		case *url.Error:
			probeResponse.ErrorMessage = "URL Error: " + err.Error()
		case net.Error:
			probeResponse.ErrorMessage = "Network Error: " + err.Error()
		default:
			probeResponse.ErrorMessage = "Error: " + err.Error()
		}
	} else {
		probeResponse.IsDown = false
		probeResponse.ResponseCode = response.Status
		mili := elapsed.Nanoseconds() / int64(time.Millisecond)
		probeResponse.ResponseTime = mili

		if len(expectsCode) > 1 {
			expectedCode, _ := strconv.Atoi(expectsCode)

			if expectedCode != response.StatusCode {
				probeResponse.IsDown = true
				probeResponse.ErrorMessage = "Incorrect HTTP status code returned: " + response.Status + " was returned but expected " + expectsCode
			}
		}

		if len(expectsValue) > 1 {
			contents, _ := ioutil.ReadAll(response.Body)
			defer response.Body.Close()
			if !strings.Contains(string(contents), expectsValue) {
				probeResponse.IsDown = true
				probeResponse.ErrorMessage = "Expected value " + expectsValue + " not found in response body"
			}
		}
	}

	return probeResponse
}

func printSimpleOutput(response ProbeResponse) {
	if response.IsDown {
		fmt.Println("down,'" + response.ErrorMessage + "'," + response.ResponseCode)
		return
	}

	fmt.Println("up," + strconv.FormatInt(response.ResponseTime, 10) + "," + response.ResponseCode)
}

func printOutput(command *ProbeConfig, response ProbeResponse) {
	var responseTime = strconv.FormatInt(response.ResponseTime, 10)
	if response.IsDown {
		pterm.Error.Println("Host is DOWN\n" + response.ErrorMessage)
		//pterm.Error.Println(response.ErrorMessage)
	} else {
		pterm.Success.Println("Respone Time: " + responseTime + "ms\tResponse Code: " + response.ResponseCode)
		//pterm.Println("Response Code: " + response.ResponseCode)
	}
}

func sanitizeHttp(address string) string {
	if strings.HasPrefix(address, "http") {
		return address
	}
	address = "http://" + address

	return address
}

func ExecProbe(command *ProbeConfig) {
	var repeater = 0

	if command.Repeat < 0 { //infinite repeat
		repeater = -2
	}

	if !command.SimpleOutput {
		pterm.DefaultSection.Println("Started Probe: " + command.Address)
	}

	if command.Repeat != 0 && command.Delay == 0 {
		command.Delay = 2
		if !command.SimpleOutput {
			pterm.Warning.Println("Defaulting to 2 second delay between repeats (use --delay to specify a delay time)")
		}
	}

	if command.Repeat == 0 {
		command.Repeat = 1
	}

	for repeater < (command.Repeat) {
		if !command.SimpleOutput {
			spinnerSuccess, _ := pterm.DefaultSpinner.Start("Running Probe")
			spinnerSuccess.RemoveWhenDone = true

			if command.Delay > 0 {
				spinnerSuccess.Text = "Running Probe with " + strconv.Itoa(command.Delay) + " second delay "
				if command.Repeat < 0 {
					spinnerSuccess.Text += "\t[Continuous Mode ON - Ctrl+C to Quit]"
				}
			}

			var probeResponse = StatusCheck(command.Address, command.Timeout, command.Delay, command.ExpectsCode, command.ExpectsString)
			spinnerSuccess.Success("-----")
			printOutput(command, probeResponse)
		} else {
			var probeResponse = StatusCheck(command.Address, command.Timeout, command.Delay, command.ExpectsCode, command.ExpectsString)
			printSimpleOutput(probeResponse)
		}
		if command.Repeat >= 0 {
			repeater++
		}
	}

	if !command.SimpleOutput {
		pterm.DefaultSection.Println("Probe Complete")
	}

}
