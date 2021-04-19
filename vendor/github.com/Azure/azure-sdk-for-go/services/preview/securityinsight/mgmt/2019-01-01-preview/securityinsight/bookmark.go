package securityinsight

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// BookmarkClient is the API spec for Microsoft.SecurityInsights (Azure Security Insights) resource provider
type BookmarkClient struct {
	BaseClient
}

// NewBookmarkClient creates an instance of the BookmarkClient client.
func NewBookmarkClient(subscriptionID string) BookmarkClient {
	return NewBookmarkClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewBookmarkClientWithBaseURI creates an instance of the BookmarkClient client using a custom endpoint.  Use this
// when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewBookmarkClientWithBaseURI(baseURI string, subscriptionID string) BookmarkClient {
	return BookmarkClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Expand expand an bookmark
// Parameters:
// resourceGroupName - the name of the resource group within the user's subscription. The name is case
// insensitive.
// operationalInsightsResourceProvider - the namespace of workspaces resource provider-
// Microsoft.OperationalInsights.
// workspaceName - the name of the workspace.
// bookmarkID - bookmark ID
// parameters - the parameters required to execute an expand operation on the given bookmark.
func (client BookmarkClient) Expand(ctx context.Context, resourceGroupName string, operationalInsightsResourceProvider string, workspaceName string, bookmarkID string, parameters BookmarkExpandParameters) (result BookmarkExpandResponse, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/BookmarkClient.Expand")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.Pattern, Rule: `^[0-9A-Fa-f]{8}-([0-9A-Fa-f]{4}-){3}[0-9A-Fa-f]{12}$`, Chain: nil}}},
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "resourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._\(\)]+$`, Chain: nil}}},
		{TargetValue: workspaceName,
			Constraints: []validation.Constraint{{Target: "workspaceName", Name: validation.MaxLength, Rule: 90, Chain: nil},
				{Target: "workspaceName", Name: validation.MinLength, Rule: 1, Chain: nil}}}}); err != nil {
		return result, validation.NewError("securityinsight.BookmarkClient", "Expand", err.Error())
	}

	req, err := client.ExpandPreparer(ctx, resourceGroupName, operationalInsightsResourceProvider, workspaceName, bookmarkID, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securityinsight.BookmarkClient", "Expand", nil, "Failure preparing request")
		return
	}

	resp, err := client.ExpandSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "securityinsight.BookmarkClient", "Expand", resp, "Failure sending request")
		return
	}

	result, err = client.ExpandResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securityinsight.BookmarkClient", "Expand", resp, "Failure responding to request")
		return
	}

	return
}

// ExpandPreparer prepares the Expand request.
func (client BookmarkClient) ExpandPreparer(ctx context.Context, resourceGroupName string, operationalInsightsResourceProvider string, workspaceName string, bookmarkID string, parameters BookmarkExpandParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"bookmarkId":                          autorest.Encode("path", bookmarkID),
		"operationalInsightsResourceProvider": autorest.Encode("path", operationalInsightsResourceProvider),
		"resourceGroupName":                   autorest.Encode("path", resourceGroupName),
		"subscriptionId":                      autorest.Encode("path", client.SubscriptionID),
		"workspaceName":                       autorest.Encode("path", workspaceName),
	}

	const APIVersion = "2019-01-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{operationalInsightsResourceProvider}/workspaces/{workspaceName}/providers/Microsoft.SecurityInsights/bookmarks/{bookmarkId}/expand", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ExpandSender sends the Expand request. The method will close the
// http.Response Body if it receives an error.
func (client BookmarkClient) ExpandSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ExpandResponder handles the response to the Expand request. The method always
// closes the http.Response Body.
func (client BookmarkClient) ExpandResponder(resp *http.Response) (result BookmarkExpandResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
