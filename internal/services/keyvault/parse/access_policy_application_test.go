// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = AccessPolicyApplicationId{}

func TestAccessPolicyApplicationIDFormatter(t *testing.T) {
	actual := NewAccessPolicyApplicationID("12345678-1234-9876-4563-123456789012", "resGroup1", "vault1", "object1", "application1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/vaults/vault1/objectId/object1/applicationId/application1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestAccessPolicyApplicationID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *AccessPolicyApplicationId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Error: true,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Error: true,
		},

		{
			// missing ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},

		{
			// missing VaultName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/",
			Error: true,
		},

		{
			// missing value for VaultName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/vaults/",
			Error: true,
		},

		{
			// missing ObjectIdName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/vaults/vault1/",
			Error: true,
		},

		{
			// missing value for ObjectIdName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/vaults/vault1/objectId/",
			Error: true,
		},

		{
			// missing ApplicationIdName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/vaults/vault1/objectId/object1/",
			Error: true,
		},

		{
			// missing value for ApplicationIdName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/vaults/vault1/objectId/object1/applicationId/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/vaults/vault1/objectId/object1/applicationId/application1",
			Expected: &AccessPolicyApplicationId{
				SubscriptionId:    "12345678-1234-9876-4563-123456789012",
				ResourceGroup:     "resGroup1",
				VaultName:         "vault1",
				ObjectIdName:      "object1",
				ApplicationIdName: "application1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.KEYVAULT/VAULTS/VAULT1/OBJECTID/OBJECT1/APPLICATIONID/APPLICATION1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := AccessPolicyApplicationID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
		if actual.VaultName != v.Expected.VaultName {
			t.Fatalf("Expected %q but got %q for VaultName", v.Expected.VaultName, actual.VaultName)
		}
		if actual.ObjectIdName != v.Expected.ObjectIdName {
			t.Fatalf("Expected %q but got %q for ObjectIdName", v.Expected.ObjectIdName, actual.ObjectIdName)
		}
		if actual.ApplicationIdName != v.Expected.ApplicationIdName {
			t.Fatalf("Expected %q but got %q for ApplicationIdName", v.Expected.ApplicationIdName, actual.ApplicationIdName)
		}
	}
}
