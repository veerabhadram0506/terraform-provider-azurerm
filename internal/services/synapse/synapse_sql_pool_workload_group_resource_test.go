// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package synapse_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SynapseWorkloadGroupResource struct{}

func TestAccSynapseWorkloadGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool_workload_group", "test")
	r := SynapseWorkloadGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseWorkloadGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool_workload_group", "test")
	r := SynapseWorkloadGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseWorkloadGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool_workload_group", "test")
	r := SynapseWorkloadGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSynapseWorkloadGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool_workload_group", "test")
	r := SynapseWorkloadGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SynapseWorkloadGroupResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SqlPoolWorkloadGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Synapse.SQLPoolWorkloadGroupClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.SqlPoolName, id.WorkloadGroupName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r SynapseWorkloadGroupResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_sql_pool_workload_group" "test" {
  name                             = "acctestWG%s"
  sql_pool_id                      = azurerm_synapse_sql_pool.test.id
  max_resource_percent             = 100
  min_resource_percent             = 0
  min_resource_percent_per_request = 3
}
`, template, data.RandomString)
}

func (r SynapseWorkloadGroupResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_sql_pool_workload_group" "import" {
  name                             = azurerm_synapse_sql_pool_workload_group.test.name
  sql_pool_id                      = azurerm_synapse_sql_pool_workload_group.test.sql_pool_id
  max_resource_percent             = 100
  min_resource_percent             = 0
  min_resource_percent_per_request = 3
}
`, config)
}

func (r SynapseWorkloadGroupResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_sql_pool_workload_group" "test" {
  name                               = "acctestWG%s"
  sql_pool_id                        = azurerm_synapse_sql_pool.test.id
  importance                         = "normal"
  max_resource_percent               = 100
  min_resource_percent               = 0
  max_resource_percent_per_request   = 3
  min_resource_percent_per_request   = 3
  query_execution_timeout_in_seconds = 0
}
`, template, data.RandomString)
}

func (r SynapseWorkloadGroupResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-synapse-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%d"
  storage_account_id = azurerm_storage_account.test.id
}

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_synapse_sql_pool" "test" {
  name                 = "acctestSP%s"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  sku_name             = "DW100c"
  create_mode          = "Default"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger, data.RandomString)
}
