// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package healthcare

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Service -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/services/service1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Workspace -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FhirService -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/fhirservices/service1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DicomService -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/dicomServices/service1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FhirService -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/fhirServices/service1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=MedTechService -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/iotConnectors/iotconnector1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=MedTechServiceFhirDestination -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/iotConnectors/iotconnector1/fhirDestinations/destination1 -rewrite=true
