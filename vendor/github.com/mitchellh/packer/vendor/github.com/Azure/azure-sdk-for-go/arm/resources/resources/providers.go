package resources

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator 0.17.0.0
// Changes may cause incorrect behavior and will be lost if the code is
// regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"net/http"
)

// ProvidersClient is the client for the Providers methods of the Resources
// service.
type ProvidersClient struct {
	ManagementClient
}

// NewProvidersClient creates an instance of the ProvidersClient client.
func NewProvidersClient(subscriptionID string) ProvidersClient {
	return NewProvidersClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewProvidersClientWithBaseURI creates an instance of the ProvidersClient
// client.
func NewProvidersClientWithBaseURI(baseURI string, subscriptionID string) ProvidersClient {
	return ProvidersClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Get gets a resource provider.
//
// resourceProviderNamespace is namespace of the resource provider.
func (client ProvidersClient) Get(resourceProviderNamespace string) (result Provider, err error) {
	req, err := client.GetPreparer(resourceProviderNamespace)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "resources.ProvidersClient", "Get", nil, "Failure preparing request")
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "resources.ProvidersClient", "Get", resp, "Failure sending request")
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.ProvidersClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client ProvidersClient) GetPreparer(resourceProviderNamespace string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceProviderNamespace": autorest.Encode("path", resourceProviderNamespace),
		"subscriptionId":            autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/{resourceProviderNamespace}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client ProvidersClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client ProvidersClient) GetResponder(resp *http.Response) (result Provider, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List gets a list of resource providers.
//
// top is query parameters. If null is passed returns all deployments.
func (client ProvidersClient) List(top *int32) (result ProviderListResult, err error) {
	req, err := client.ListPreparer(top)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "resources.ProvidersClient", "List", nil, "Failure preparing request")
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "resources.ProvidersClient", "List", resp, "Failure sending request")
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.ProvidersClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client ProvidersClient) ListPreparer(top *int32) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}
	if top != nil {
		queryParameters["$top"] = autorest.Encode("query", *top)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client ProvidersClient) ListSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client ProvidersClient) ListResponder(resp *http.Response) (result ProviderListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListNextResults retrieves the next set of results, if any.
func (client ProvidersClient) ListNextResults(lastResults ProviderListResult) (result ProviderListResult, err error) {
	req, err := lastResults.ProviderListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "resources.ProvidersClient", "List", nil, "Failure preparing next results request request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "resources.ProvidersClient", "List", resp, "Failure sending next results request request")
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.ProvidersClient", "List", resp, "Failure responding to next results request request")
	}

	return
}

// Register registers provider to be used with a subscription.
//
// resourceProviderNamespace is namespace of the resource provider.
func (client ProvidersClient) Register(resourceProviderNamespace string) (result Provider, err error) {
	req, err := client.RegisterPreparer(resourceProviderNamespace)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "resources.ProvidersClient", "Register", nil, "Failure preparing request")
	}

	resp, err := client.RegisterSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "resources.ProvidersClient", "Register", resp, "Failure sending request")
	}

	result, err = client.RegisterResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.ProvidersClient", "Register", resp, "Failure responding to request")
	}

	return
}

// RegisterPreparer prepares the Register request.
func (client ProvidersClient) RegisterPreparer(resourceProviderNamespace string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceProviderNamespace": autorest.Encode("path", resourceProviderNamespace),
		"subscriptionId":            autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/{resourceProviderNamespace}/register", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// RegisterSender sends the Register request. The method will close the
// http.Response Body if it receives an error.
func (client ProvidersClient) RegisterSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// RegisterResponder handles the response to the Register request. The method always
// closes the http.Response Body.
func (client ProvidersClient) RegisterResponder(resp *http.Response) (result Provider, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Unregister unregisters provider from a subscription.
//
// resourceProviderNamespace is namespace of the resource provider.
func (client ProvidersClient) Unregister(resourceProviderNamespace string) (result Provider, err error) {
	req, err := client.UnregisterPreparer(resourceProviderNamespace)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "resources.ProvidersClient", "Unregister", nil, "Failure preparing request")
	}

	resp, err := client.UnregisterSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "resources.ProvidersClient", "Unregister", resp, "Failure sending request")
	}

	result, err = client.UnregisterResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.ProvidersClient", "Unregister", resp, "Failure responding to request")
	}

	return
}

// UnregisterPreparer prepares the Unregister request.
func (client ProvidersClient) UnregisterPreparer(resourceProviderNamespace string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceProviderNamespace": autorest.Encode("path", resourceProviderNamespace),
		"subscriptionId":            autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/{resourceProviderNamespace}/unregister", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// UnregisterSender sends the Unregister request. The method will close the
// http.Response Body if it receives an error.
func (client ProvidersClient) UnregisterSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// UnregisterResponder handles the response to the Unregister request. The method always
// closes the http.Response Body.
func (client ProvidersClient) UnregisterResponder(resp *http.Response) (result Provider, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
