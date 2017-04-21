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

package ecsclient

//go:generate mockgen -destination=mocks/ec2/ec2_mocks.go github.com/aws/aws-sdk-go/service/ec2/ec2iface EC2API
//go:generate mockgen -destination=mocks/ecs/ecs_mocks.go github.com/aws/aws-sdk-go/service/ecs/ecsiface ECSAPI
//go:generate mockgen -destination=mocks/client_mocks.go github.com/luxbet/ecs-task-kite/lib/ecsclient AugmentedTask,AugmentedContainer
