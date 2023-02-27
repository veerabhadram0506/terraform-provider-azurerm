// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appconfiguration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2022-05-01/configurationstores"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/sdk/1.0/appconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KeyResource struct{}

var _ sdk.ResourceWithCustomizeDiff = KeyResource{}

var _ sdk.ResourceWithStateMigration = KeyResource{}

const (
	KeyTypeVault        = "vault"
	KeyTypeKV           = "kv"
	VaultKeyContentType = "application/vnd.microsoft.appconfig.keyvaultref+json;charset=utf-8"
)

type KeyResourceModel struct {
	ConfigurationStoreId string                 `tfschema:"configuration_store_id"`
	Key                  string                 `tfschema:"key"`
	ContentType          string                 `tfschema:"content_type"`
	Etag                 string                 `tfschema:"etag"`
	Label                string                 `tfschema:"label"`
	Value                string                 `tfschema:"value"`
	Locked               bool                   `tfschema:"locked"`
	Tags                 map[string]interface{} `tfschema:"tags"`
	Type                 string                 `tfschema:"type"`
	VaultKeyReference    string                 `tfschema:"vault_key_reference"`
}

type VaultKeyReference struct {
	URI string `json:"uri"`
}

func (k KeyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"configuration_store_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: configurationstores.ValidateConfigurationStoreID,
		},
		"key": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"content_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},
		"etag": {
			Type:     pluginsdk.TypeString,
			Computed: true,
			Optional: true,
		},
		"label": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"value": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},
		"locked": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
		"type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "kv",
			ValidateFunc: validation.StringInSlice([]string{KeyTypeVault, KeyTypeKV}, false),
		},
		"vault_key_reference": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},
		"tags": tags.Schema(),
	}
}

func (k KeyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (k KeyResource) ModelObject() interface{} {
	return &KeyResourceModel{}
}

func (k KeyResource) ResourceType() string {
	return "azurerm_app_configuration_key"
}

func (k KeyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model KeyResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClient(ctx, model.ConfigurationStoreId)
			if err != nil {
				return err
			}
			if client == nil {
				return fmt.Errorf("app configuration %q was not found", model.ConfigurationStoreId)
			}

			appCfgKeyResourceID := parse.AppConfigurationKeyId{
				ConfigurationStoreId: model.ConfigurationStoreId,
				Key:                  model.Key,
				Label:                model.Label,
			}

			// from https://learn.microsoft.com/en-us/azure/azure-app-configuration/concept-enable-rbac#azure-built-in-roles-for-azure-app-configuration
			// allow up to 15 min for role permission to be done propagated
			metadata.Logger.Infof("[DEBUG] Waiting for App Configuration Key %q read permission to be done propagated", model.Key)
			stateConf := &pluginsdk.StateChangeConf{
				Pending:      []string{"Forbidden"},
				Target:       []string{"Error", "Exists"},
				Refresh:      appConfigurationGetKeyRefreshFunc(ctx, client, model.Key, model.Label),
				PollInterval: 20 * time.Second,
				Timeout:      15 * time.Minute,
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for App Configuration Key %q read permission to be propagated: %+v", model.Key, err)
			}

			kv, err := client.GetKeyValue(ctx, model.Key, model.Label, "", "", "", []string{})
			if err != nil {
				if v, ok := err.(autorest.DetailedError); ok {
					if !utils.ResponseWasNotFound(autorest.Response{Response: v.Response}) {
						return fmt.Errorf("checking for presence of existing %s: %+v", appCfgKeyResourceID, err)
					}
				} else {
					return fmt.Errorf("while checking for key's %q existence: %+v", model.Key, err)
				}
			} else if kv.Response.StatusCode == 200 {
				return tf.ImportAsExistsError(k.ResourceType(), appCfgKeyResourceID.ID())
			}

			entity := appconfiguration.KeyValue{
				Key:   utils.String(model.Key),
				Label: utils.String(model.Label),
				Tags:  tags.Expand(model.Tags),
			}

			switch model.Type {
			case KeyTypeKV:
				entity.ContentType = utils.String(model.ContentType)
				entity.Value = utils.String(model.Value)
			case KeyTypeVault:
				entity.ContentType = utils.String(VaultKeyContentType)
				ref, err := json.Marshal(VaultKeyReference{URI: model.VaultKeyReference})
				if err != nil {
					return fmt.Errorf("while encoding vault key reference: %+v", err)
				}
				entity.Value = utils.String(string(ref))
			}

			if _, err = client.PutKeyValue(ctx, model.Key, model.Label, &entity, "", ""); err != nil {
				return err
			}

			if model.Locked {
				_, err = client.PutLock(ctx, model.Key, model.Label, "", "")
				if err != nil {
					return fmt.Errorf("while locking key/label pair %q/%q: %+v", model.Key, model.Label, err)
				}
			}

			if appCfgKeyResourceID.Label == "" {
				// We set an empty label as %00 in the resource ID
				// Otherwise it breaks the ID parsing logic
				appCfgKeyResourceID.Label = "%00"
			}
			metadata.SetID(appCfgKeyResourceID)
			return nil
		},
		Timeout: 45 * time.Minute,
	}
}

func (k KeyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			resourceID, err := parse.KeyId(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			// We set an empty label as %00 in the ID to make the ID validator happy
			// but in reality the label is just an empty string
			if resourceID.Label == "%00" {
				resourceID.Label = ""
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClient(ctx, resourceID.ConfigurationStoreId)
			if err != nil {
				return err
			}
			if client == nil {
				// if the parent AppConfiguration is gone, all the data will be too
				return metadata.MarkAsGone(resourceID)
			}

			kv, err := client.GetKeyValue(ctx, resourceID.Key, resourceID.Label, "", "", "", []string{})
			if err != nil {
				if v, ok := err.(autorest.DetailedError); ok {
					if utils.ResponseWasNotFound(autorest.Response{Response: v.Response}) {
						return metadata.MarkAsGone(resourceID)
					}
				} else {
					return fmt.Errorf("while checking for key's %q existence: %+v", resourceID.Key, err)
				}
				return fmt.Errorf("while checking for key's %q existence: %+v", resourceID.Label, err)
			}

			model := KeyResourceModel{
				ConfigurationStoreId: resourceID.ConfigurationStoreId,
				Key:                  utils.NormalizeNilableString(kv.Key),
				ContentType:          utils.NormalizeNilableString(kv.ContentType),
				Etag:                 utils.NormalizeNilableString(kv.Etag),
				Label:                utils.NormalizeNilableString(kv.Label),
				Tags:                 tags.Flatten(kv.Tags),
			}

			if utils.NormalizeNilableString(kv.ContentType) != VaultKeyContentType {
				model.Type = KeyTypeKV
				model.Value = utils.NormalizeNilableString(kv.Value)
			} else {
				var ref VaultKeyReference
				refBytes := []byte(utils.NormalizeNilableString(kv.Value))
				err := json.Unmarshal(refBytes, &ref)
				if err != nil {
					return fmt.Errorf("while unmarshalling vault reference: %+v", err)
				}

				model.Type = KeyTypeVault
				model.VaultKeyReference = ref.URI
				model.ContentType = VaultKeyContentType
				model.Value = utils.NormalizeNilableString(kv.Value)
			}

			if kv.Locked != nil {
				model.Locked = *kv.Locked
			}
			return metadata.Encode(&model)
		},
		Timeout: 5 * time.Minute,
	}
}

func (k KeyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			resourceID, err := parse.KeyId(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClient(ctx, resourceID.ConfigurationStoreId)
			if err != nil {
				return err
			}
			if client == nil {
				return fmt.Errorf("app configuration %q was not found", resourceID.ConfigurationStoreId)
			}

			var model KeyResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			if metadata.ResourceData.HasChange("value") || metadata.ResourceData.HasChange("content_type") || metadata.ResourceData.HasChange("tags") || metadata.ResourceData.HasChange("type") || metadata.ResourceData.HasChange("vault_key_reference") {
				entity := appconfiguration.KeyValue{
					Key:   utils.String(model.Key),
					Label: utils.String(model.Label),
					Tags:  tags.Expand(model.Tags),
				}

				switch model.Type {
				case KeyTypeKV:
					entity.ContentType = utils.String(model.ContentType)
					entity.Value = utils.String(model.Value)
				case KeyTypeVault:
					entity.ContentType = utils.String(VaultKeyContentType)
					ref, err := json.Marshal(VaultKeyReference{URI: model.VaultKeyReference})
					if err != nil {
						return fmt.Errorf("while encoding vault key reference: %+v", err)
					}
					entity.Value = utils.String(string(ref))
				}
				if _, err = client.PutKeyValue(ctx, model.Key, model.Label, &entity, "", ""); err != nil {
					return fmt.Errorf("while updating key/label pair %s/%s: %+v", model.Key, model.Label, err)
				}
			}

			if metadata.ResourceData.HasChange("locked") {
				if model.Locked {
					if _, err = client.PutLock(ctx, model.Key, model.Label, "", ""); err != nil {
						return fmt.Errorf("while locking key/label pair %s/%s: %+v", model.Key, model.Label, err)
					}
				} else {
					if _, err = client.DeleteLock(ctx, model.Key, model.Label, "", ""); err != nil {
						return fmt.Errorf("while unlocking key/label pair %s/%s: %+v", model.Key, model.Label, err)
					}
				}
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (k KeyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			resourceID, err := parse.KeyId(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClient(ctx, resourceID.ConfigurationStoreId)
			if err != nil {
				return err

			}
			if client == nil {
				return fmt.Errorf("app configuration %q was not found", resourceID.ConfigurationStoreId)
			}

			decodedKey, err := url.QueryUnescape(resourceID.Key)
			if err != nil {
				return fmt.Errorf("while decoding key of resource ID: %+v", err)
			}

			decodedLabel, err := url.QueryUnescape(resourceID.Label)
			if err != nil {
				return fmt.Errorf("while decoding label of resource ID: %+v", err)
			}

			if _, err = client.DeleteLock(ctx, decodedKey, decodedLabel, "", ""); err != nil {
				return fmt.Errorf("while unlocking key/label pair %s/%s: %+v", decodedKey, resourceID.Label, err)
			}

			_, err = client.DeleteKeyValue(ctx, decodedKey, resourceID.Label, "")
			if err != nil {
				return fmt.Errorf("while removing key %q from App Configuration Store %q: %+v", decodedKey, resourceID.ConfigurationStoreId, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (k KeyResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff
			keyType := rd.Get("type").(string)
			if keyType == KeyTypeVault {
				contentType := rd.Get("content_type").(string)
				if rd.HasChange("content_type") && contentType != VaultKeyContentType {
					return fmt.Errorf("vault reference key %q cannot have content type other than %q (found %q)", rd.Get("key").(string), VaultKeyContentType, contentType)
				}

				value := rd.Get("value").(string)
				var v VaultKeyReference
				if rd.HasChange("value") {
					if err := json.Unmarshal([]byte(value), &v); err != nil {
						return fmt.Errorf("while validating attribute 'value' (%q): %+v", value, err)
					}
					if v.URI == "" {
						return fmt.Errorf("invalid data in 'value' contents: URI cannot be empty")
					}
				}
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (k KeyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.AppConfigurationKeyID
}

func (k KeyResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.KeyResourceV0ToV1{},
		},
	}
}
