package appservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WindowsFunctionAppDataSource struct{}

type WindowsFunctionAppDataSourceModel struct {
	Name               string `tfschema:"name"`
	ResourceGroup      string `tfschema:"resource_group_name"`
	Location           string `tfschema:"location"`
	ServicePlanId      string `tfschema:"service_plan_id"`
	StorageAccountName string `tfschema:"storage_account_name"`

	StorageAccountKey string `tfschema:"storage_account_access_key"`
	StorageUsesMSI    bool   `tfschema:"storage_uses_managed_identity"`

	AppSettings               map[string]string                      `tfschema:"app_settings"`
	AuthSettings              []helpers.AuthSettings                 `tfschema:"auth_settings"`
	Backup                    []helpers.Backup                       `tfschema:"backup"`
	BuiltinLogging            bool                                   `tfschema:"builtin_logging_enabled"`
	ClientCertEnabled         bool                                   `tfschema:"client_certificate_enabled"`
	ClientCertMode            string                                 `tfschema:"client_certificate_mode"`
	ConnectionStrings         []helpers.ConnectionString             `tfschema:"connection_string"`
	DailyMemoryTimeQuota      int                                    `tfschema:"daily_memory_time_quota"`
	Enabled                   bool                                   `tfschema:"enabled"`
	FunctionExtensionsVersion string                                 `tfschema:"functions_extension_version"`
	ForceDisableContentShare  bool                                   `tfschema:"content_share_force_disabled"`
	HttpsOnly                 bool                                   `tfschema:"https_only"`
	Identity                  []helpers.Identity                     `tfschema:"identity"`
	SiteConfig                []helpers.SiteConfigWindowsFunctionApp `tfschema:"site_config"`
	Tags                      map[string]string                      `tfschema:"tags"`

	CustomDomainVerificationId    string   `tfschema:"custom_domain_verification_id"`
	DefaultHostname               string   `tfschema:"default_hostname"`
	Kind                          string   `tfschema:"kind"`
	OutboundIPAddresses           string   `tfschema:"outbound_ip_addresses"`
	OutboundIPAddressList         []string `tfschema:"outbound_ip_address_list"`
	PossibleOutboundIPAddresses   string   `tfschema:"possible_outbound_ip_addresses"`
	PossibleOutboundIPAddressList []string `tfschema:"possible_outbound_ip_address_list"`

	SiteCredentials []helpers.SiteCredential `tfschema:"site_credential"`
}

var _ sdk.DataSource = WindowsFunctionAppDataSource{}

func (d WindowsFunctionAppDataSource) ModelObject() interface{} {
	return &WindowsFunctionAppDataSourceModel{}
}

func (d WindowsFunctionAppDataSource) ResourceType() string {
	return "azurerm_windows_function_app"
}

func (d WindowsFunctionAppDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.WebAppName,
		},

		"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),
	}
}

func (d WindowsFunctionAppDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": location.SchemaComputed(),

		"service_plan_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"storage_account_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"storage_account_access_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"storage_uses_managed_identity": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"app_settings": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"auth_settings": helpers.AuthSettingsSchemaComputed(),

		"backup": helpers.BackupSchemaComputed(),

		"builtin_logging_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"client_certificate_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"client_certificate_mode": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"connection_string": helpers.ConnectionStringSchemaComputed(),

		"daily_memory_time_quota": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"content_share_force_disabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"functions_extension_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"https_only": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"custom_domain_verification_id": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"default_hostname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kind": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"outbound_ip_addresses": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"outbound_ip_address_list": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"possible_outbound_ip_addresses": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"possible_outbound_ip_address_list": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"site_credential": helpers.SiteCredentialSchema(),

		"site_config": helpers.SiteConfigSchemaWindowsFunctionAppComputed(),

		"identity": helpers.IdentitySchemaComputed(),

		"tags": tags.SchemaDataSource(),
	}
}

func (d WindowsFunctionAppDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var functionApp WindowsFunctionAppModel
			if err := metadata.Decode(&functionApp); err != nil {
				return err
			}

			id := parse.NewFunctionAppID(subscriptionId, functionApp.ResourceGroup, functionApp.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("Windows %s not found", id)
				}
				return fmt.Errorf("checking for presence of existing Windows %s: %+v", id, err)
			}

			if existing.SiteProperties == nil {
				return fmt.Errorf("reading properties of Windows %s", id)
			}
			props := *existing.SiteProperties

			functionApp.Name = id.SiteName
			functionApp.ResourceGroup = id.ResourceGroup
			functionApp.ServicePlanId = utils.NormalizeNilableString(props.ServerFarmID)
			functionApp.Location = location.NormalizeNilable(existing.Location)
			functionApp.Enabled = utils.NormaliseNilableBool(existing.Enabled)
			functionApp.ClientCertMode = string(existing.ClientCertMode)
			functionApp.DailyMemoryTimeQuota = int(utils.NormaliseNilableInt32(props.DailyMemoryTimeQuota))
			functionApp.Tags = tags.ToTypedObject(existing.Tags)
			functionApp.Kind = utils.NormalizeNilableString(existing.Kind)

			appSettingsResp, err := client.ListApplicationSettings(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading App Settings for Windows %s: %+v", id, err)
			}

			connectionStrings, err := client.ListConnectionStrings(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Connection String information for Windows %s: %+v", id, err)
			}

			siteCredentialsFuture, err := client.ListPublishingCredentials(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("listing Site Publishing Credential information for Windows %s: %+v", id, err)
			}

			if err := siteCredentialsFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for Site Publishing Credential information for Windows %s: %+v", id, err)
			}
			siteCredentials, err := siteCredentialsFuture.Result(*client)
			if err != nil {
				return fmt.Errorf("reading Site Publishing Credential information for Windows %s: %+v", id, err)
			}

			auth, err := client.GetAuthSettings(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading Auth Settings for Windows %s: %+v", id, err)
			}

			backup, err := client.GetBackupConfiguration(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				if !utils.ResponseWasNotFound(backup.Response) {
					return fmt.Errorf("reading Backup Settings for Windows %s: %+v", id, err)
				}
			}

			logs, err := client.GetDiagnosticLogsConfiguration(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("reading logs configuration for Windows %s: %+v", id, err)
			}

			if identity := helpers.FlattenIdentity(existing.Identity); identity != nil {
				functionApp.Identity = identity
			}

			configResp, err := client.GetConfiguration(ctx, id.ResourceGroup, id.SiteName)
			if err != nil {
				return fmt.Errorf("making Read request on AzureRM Function App Configuration %q: %+v", id.SiteName, err)
			}

			siteConfig, err := helpers.FlattenSiteConfigWindowsFunctionApp(configResp.SiteConfig)
			if err != nil {
				return fmt.Errorf("reading Site Config for Windows %s: %+v", id, err)
			}

			functionApp.SiteConfig = []helpers.SiteConfigWindowsFunctionApp{*siteConfig}

			functionApp.unpackWindowsFunctionAppSettings(appSettingsResp)

			functionApp.ConnectionStrings = helpers.FlattenConnectionStrings(connectionStrings)

			functionApp.SiteCredentials = helpers.FlattenSiteCredentials(siteCredentials)

			functionApp.AuthSettings = helpers.FlattenAuthSettings(auth)

			functionApp.Backup = helpers.FlattenBackupConfig(backup)

			functionApp.SiteConfig[0].AppServiceLogs = helpers.FlattenFunctionAppAppServiceLogs(logs)

			functionApp.HttpsOnly = utils.NormaliseNilableBool(existing.HTTPSOnly)

			functionApp.ClientCertEnabled = utils.NormaliseNilableBool(existing.ClientCertEnabled)

			metadata.SetID(id)

			return metadata.Encode(&functionApp)
		},
	}
}
