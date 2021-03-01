package serve

import (
	"encoding/json"
	"fmt"
	"github.com/keithnyc/probicon/probe"
	"log"
	"net/http"
	"strconv"
)

type ServerConfig struct {
	ListenPort int
}

func StartServer(serverConf *ServerConfig) {
	http.HandleFunc("/probe", execProbe)
	log.Println("Awaiting Requests")
	http.ListenAndServe(fmt.Sprintf(":%d", serverConf.ListenPort), nil)

}

func execProbe(w http.ResponseWriter, r *http.Request) {
	var address, expectsCode, expectsValue string
	var timeoutSeconds int

	log.Println("Request Received. Processing...")
	addressBuffer, found := r.URL.Query()["address"]

	if !found {
		log.Println("Address parameter was not found")
		return
	}
	address = addressBuffer[0]

	timeoutBuffer, found := r.URL.Query()["timeout"]

	if !found {
		timeoutSeconds = 10
	} else {
		timeoutSeconds, err := strconv.Atoi(timeoutBuffer[0])
		if err != nil {
			timeoutSeconds = 10
		}
		if timeoutSeconds > 120{
			timeoutSeconds = 30
		}
		log.Printf("Using timeout of %d", timeoutSeconds)
	}
	expectsCodeBuffer, found := r.URL.Query()["expectscode"]

	if found {
		expectsCode = expectsCodeBuffer[0]
	}

	expectsValBuffer, found := r.URL.Query()["expectsvalue"]
	if found {
		expectsValue = expectsValBuffer[0]
	}

	var probeResponse = probe.StatusCheck(address, timeoutSeconds, 0, expectsCode, expectsValue)
	log.Println("Probe Response: ", probeResponse)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(probeResponse)
}
