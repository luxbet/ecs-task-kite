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

package ecsclient_test

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	"github.com/luxbet/ecs-task-kite/lib/ecsclient"
	"github.com/luxbet/ecs-task-kite/lib/ecsclient/mocks/ec2"
	"github.com/luxbet/ecs-task-kite/lib/ecsclient/mocks/ecs"
)

const cluster = "testCluster"

func strptr(s string) *string {
	return &s
}

var pcluster = strptr(cluster)

func setup(t *testing.T) (*gomock.Controller, ecsclient.ECSSimpleClient, *mock_ecsiface.MockECSAPI, *mock_ec2iface.MockEC2API) {
	ctrl := gomock.NewController(t)
	mockecs := mock_ecsiface.NewMockECSAPI(ctrl)
	mockec2 := mock_ec2iface.NewMockEC2API(ctrl)
	ecsClient := ecsclient.New(cluster, "us-east-1", mockecs, mockec2)
	return ctrl, ecsClient, mockecs, mockec2
}

func TestListAllTasks(t *testing.T) {
	ctrl, ecsClient, mockecs, mockec2 := setup(t)
	defer ctrl.Finish()

	mockTaskArns := []*string{strptr("task1"), strptr("task2")}
	mockCIArns := []*string{strptr("ci1"), strptr("ci2")}
	mockEC2Ids := []*string{strptr("i-1"), strptr("i-2")}
	mockTasks := []*ecs.Task{
		&ecs.Task{
			TaskArn:              mockTaskArns[0],
			LastStatus:           strptr("RUNNING"),
			ContainerInstanceArn: mockCIArns[0],
		},
		&ecs.Task{
			TaskArn:              mockTaskArns[1],
			LastStatus:           strptr("RUNNING"),
			ContainerInstanceArn: mockCIArns[1],
		},
	}
	mockCIs := []*ecs.ContainerInstance{
		&ecs.ContainerInstance{
			ContainerInstanceArn: mockCIArns[0],
			Ec2InstanceId:        mockEC2Ids[0],
		},
		&ecs.ContainerInstance{
			ContainerInstanceArn: mockCIArns[1],
			Ec2InstanceId:        mockEC2Ids[1],
		},
	}
	mockEC2Instances := []*ec2.Instance{
		&ec2.Instance{
			InstanceId:      mockEC2Ids[0],
			PublicIpAddress: strptr("1.1.1.1"),
		},
		&ec2.Instance{
			InstanceId:      mockEC2Ids[1],
			PublicIpAddress: strptr("2.2.2.2"),
		},
	}
	gomock.InOrder(
		mockecs.EXPECT().ListTasksPages(&ecs.ListTasksInput{Cluster: pcluster}, gomock.Any()).Do(func(_, f interface{}) {
			f.(func(*ecs.ListTasksOutput, bool) bool)(&ecs.ListTasksOutput{TaskArns: mockTaskArns}, true)
		}).Return(nil),
		mockecs.EXPECT().DescribeTasks(&ecs.DescribeTasksInput{Cluster: pcluster, Tasks: mockTaskArns}).Return(
			&ecs.DescribeTasksOutput{
				Tasks: mockTasks,
			},
			nil,
		),
		mockecs.EXPECT().DescribeContainerInstances(describeContainerInstanceMatcher{&ecs.DescribeContainerInstancesInput{Cluster: pcluster, ContainerInstances: mockCIArns}}).Return(
			&ecs.DescribeContainerInstancesOutput{
				ContainerInstances: mockCIs,
			},
			nil,
		),
		mockec2.EXPECT().DescribeInstances(&ec2.DescribeInstancesInput{InstanceIds: mockEC2Ids}).Return(&ec2.DescribeInstancesOutput{
			Reservations: []*ec2.Reservation{
				&ec2.Reservation{Instances: mockEC2Instances},
			},
		},
			nil,
		),
	)
	tasks, err := ecsClient.Tasks(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	for i, task := range tasks {
		if !reflect.DeepEqual(task.ECSTask(), mockTasks[i]) {
			t.Fatal("Tasks did not match expected")
		}

		if !reflect.DeepEqual(task.EC2Instance(), mockEC2Instances[i]) {
			t.Fatal("Task's ec2 instance did not match expected")
		}
	}
}

type describeContainerInstanceMatcher struct {
	*ecs.DescribeContainerInstancesInput
}

// Checks for the same clusters and arns, ignoring order
func (lhs describeContainerInstanceMatcher) Matches(x interface{}) bool {
	rhs, ok := x.(*ecs.DescribeContainerInstancesInput)
	if !ok {
		return false
	}
	if *lhs.Cluster != *rhs.Cluster {
		return false
	}
	if len(lhs.ContainerInstances) != len(rhs.ContainerInstances) {
		return false
	}
	arns := make(map[string]bool)
	for _, arn := range lhs.ContainerInstances {
		arns[*arn] = true
	}
	for _, arn := range rhs.ContainerInstances {
		if _, ok := arns[*arn]; !ok {
			return false
		}
	}
	return true
}

func (describeContainerInstanceMatcher) String() string {
	return "Container Instance Describe Matcher"
}
