// Copyright 2014-2015 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.-

package taskhelpers

import (
	"fmt"

	"github.com/luxbet/ecs-task-kite/lib/ecsclient"
)

// ContainerPorts returns all of the ports that a given container within the
// tasks is listening on.
func ContainerPorts(tasks []ecsclient.AugmentedTask, containerName string, protocol string) []uint16 {
	// dedupe map to return the minimal array
	seenPorts := make(map[uint16]bool)
	output := make([]uint16, 0, len(tasks)/2)
	for _, task := range tasks {
		container := task.Container(containerName)
		if container == nil {
			continue
		}
		if !container.Running() {
			continue
		}
		ports := container.ContainerPorts(protocol)
		for _, port := range ports {
			if _, ok := seenPorts[port]; !ok {
				output = append(output, port)
				seenPorts[port] = true
			}
		}
	}
	return output
}

// FilterIPPort returns the "ip:port" pair for the given containerName within
// all tasks where the given container is known to be running.
func FilterIPPort(tasks []ecsclient.AugmentedTask, containerName string, containerPort uint16, publicIP bool) []string {
	output := make([]string, 0, len(tasks)/2)
	for _, task := range tasks {
		container := task.Container(containerName)
		if container == nil {
			continue
		}
		if !container.Running() {
			continue
		}
		hostPort := container.ResolvePort(containerPort)
		if hostPort == 0 {
			continue
		}
		var taskIP string
		if publicIP {
			taskIP = task.PublicIP()
		} else {
			taskIP = task.PrivateIP()
		}
		if taskIP == "" {
			continue
		}
		output = append(output, fmt.Sprintf("%s:%d", taskIP, hostPort))
	}
	return output
}
