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
	"github.com/keithnyc/probicon/probe"
	"github.com/spf13/cobra"
)

var probeConf probe.ProbeConfig

// probeCmd represents the probe command
var probeCmd = &cobra.Command{
	Use:   "probe",
	Short: "Probe an address for health status",
	Long: `When given an address, such as a URL, probe conducts
a health check on the given endpoint for signs of life. 

If using simple output mode (-s), the response time (in milliseconds) will be returned 
if the address endpoint is responding, else a 0 will be returned. This is
useful for combining probe with other utilities.`,

	Run: func(cmd *cobra.Command, args []string) {
		probe.ExecProbe(&probeConf)
	},
}

func init() {
	rootCmd.AddCommand(probeCmd)
	probeConf = probe.ProbeConfig{}

	probeCmd.Flags().StringVarP(&probeConf.Address, "address", "a", "", "Address to probe, e.g. https://example.com")
	probeCmd.Flags().BoolVarP(&probeConf.SimpleOutput, "simpleoutput", "s", false, "Use simple output where only the response time is reported - useful for piping results to other commands. A value of 0 indicates the host is down.")
	probeCmd.Flags().IntVarP(&probeConf.Timeout, "timeout", "t", 10, "Timeout value in seconds")
	probeCmd.Flags().IntVarP(&probeConf.Repeat, "repeat", "r", 0, "Repeat check x times (default is 0). Use -1 to repeat until stopped")
	probeCmd.Flags().IntVarP(&probeConf.Delay, "delay", "d", 0, "Delay, in seconds, before running probe. When used with the -r flag, this delay is repeated between each execution")
	probeCmd.Flags().StringVarP(&probeConf.ExpectsCode, "code", "c", "", "Specify expected HTTP status code (e.g. 200). By default, HTTP status is ignored")
	probeCmd.Flags().StringVarP(&probeConf.ExpectsString, "value", "v", "", "Specify expected string to be found in the response content")

	probeCmd.MarkFlagRequired("address")
}
