package notificationhubs

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

// NamespacesClient is the azure NotificationHub client
type NamespacesClient struct {
	ManagementClient
}

// NewNamespacesClient creates an instance of the NamespacesClient client.
func NewNamespacesClient(subscriptionID string) NamespacesClient {
	return NewNamespacesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewNamespacesClientWithBaseURI creates an instance of the NamespacesClient
// client.
func NewNamespacesClientWithBaseURI(baseURI string, subscriptionID string) NamespacesClient {
	return NamespacesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CheckAvailability checks the availability of the given service namespace
// across all Windows Azure subscriptions. This is useful because the domain
// name is created based on the service namespace name.
//
// parameters is the namespace name.
func (client NamespacesClient) CheckAvailability(parameters CheckAvailabilityParameters) (result CheckAvailabilityResource, err error) {
	req, err := client.CheckAvailabilityPreparer(parameters)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "CheckAvailability", nil, "Failure preparing request")
	}

	resp, err := client.CheckAvailabilitySender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "CheckAvailability", resp, "Failure sending request")
	}

	result, err = client.CheckAvailabilityResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "CheckAvailability", resp, "Failure responding to request")
	}

	return
}

// CheckAvailabilityPreparer prepares the CheckAvailability request.
func (client NamespacesClient) CheckAvailabilityPreparer(parameters CheckAvailabilityParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.NotificationHubs/checkNamespaceAvailability", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// CheckAvailabilitySender sends the CheckAvailability request. The method will close the
// http.Response Body if it receives an error.
func (client NamespacesClient) CheckAvailabilitySender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// CheckAvailabilityResponder handles the response to the CheckAvailability request. The method always
// closes the http.Response Body.
func (client NamespacesClient) CheckAvailabilityResponder(resp *http.Response) (result CheckAvailabilityResource, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// CreateOrUpdate creates/Updates a service namespace. Once created, this
// namespace's resource manifest is immutable. This operation is idempotent.
//
// resourceGroupName is the name of the resource group. namespaceName is the
// namespace name. parameters is parameters supplied to create a Namespace
// Resource.
func (client NamespacesClient) CreateOrUpdate(resourceGroupName string, namespaceName string, parameters NamespaceCreateOrUpdateParameters) (result NamespaceResource, err error) {
	req, err := client.CreateOrUpdatePreparer(resourceGroupName, namespaceName, parameters)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "CreateOrUpdate", nil, "Failure preparing request")
	}

	resp, err := client.CreateOrUpdateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "CreateOrUpdate", resp, "Failure sending request")
	}

	result, err = client.CreateOrUpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "CreateOrUpdate", resp, "Failure responding to request")
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client NamespacesClient) CreateOrUpdatePreparer(resourceGroupName string, namespaceName string, parameters NamespaceCreateOrUpdateParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"namespaceName":     autorest.Encode("path", namespaceName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NotificationHubs/namespaces/{namespaceName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client NamespacesClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client NamespacesClient) CreateOrUpdateResponder(resp *http.Response) (result NamespaceResource, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// CreateOrUpdateAuthorizationRule creates an authorization rule for a
// namespace
//
// resourceGroupName is the name of the resource group. namespaceName is the
// namespace name. authorizationRuleName is aauthorization Rule Name.
// parameters is the shared access authorization rule.
func (client NamespacesClient) CreateOrUpdateAuthorizationRule(resourceGroupName string, namespaceName string, authorizationRuleName string, parameters SharedAccessAuthorizationRuleCreateOrUpdateParameters) (result SharedAccessAuthorizationRuleResource, err error) {
	req, err := client.CreateOrUpdateAuthorizationRulePreparer(resourceGroupName, namespaceName, authorizationRuleName, parameters)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "CreateOrUpdateAuthorizationRule", nil, "Failure preparing request")
	}

	resp, err := client.CreateOrUpdateAuthorizationRuleSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "CreateOrUpdateAuthorizationRule", resp, "Failure sending request")
	}

	result, err = client.CreateOrUpdateAuthorizationRuleResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "CreateOrUpdateAuthorizationRule", resp, "Failure responding to request")
	}

	return
}

// CreateOrUpdateAuthorizationRulePreparer prepares the CreateOrUpdateAuthorizationRule request.
func (client NamespacesClient) CreateOrUpdateAuthorizationRulePreparer(resourceGroupName string, namespaceName string, authorizationRuleName string, parameters SharedAccessAuthorizationRuleCreateOrUpdateParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"authorizationRuleName": autorest.Encode("path", authorizationRuleName),
		"namespaceName":         autorest.Encode("path", namespaceName),
		"resourceGroupName":     autorest.Encode("path", resourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NotificationHubs/namespaces/{namespaceName}/AuthorizationRules/{authorizationRuleName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// CreateOrUpdateAuthorizationRuleSender sends the CreateOrUpdateAuthorizationRule request. The method will close the
// http.Response Body if it receives an error.
func (client NamespacesClient) CreateOrUpdateAuthorizationRuleSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// CreateOrUpdateAuthorizationRuleResponder handles the response to the CreateOrUpdateAuthorizationRule request. The method always
// closes the http.Response Body.
func (client NamespacesClient) CreateOrUpdateAuthorizationRuleResponder(resp *http.Response) (result SharedAccessAuthorizationRuleResource, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete deletes an existing namespace. This operation also removes all
// associated notificationHubs under the namespace. This method may poll for
// completion. Polling can be canceled by passing the cancel channel
// argument. The channel will be used to cancel polling and any outstanding
// HTTP requests.
//
// resourceGroupName is the name of the resource group. namespaceName is the
// namespace name.
func (client NamespacesClient) Delete(resourceGroupName string, namespaceName string, cancel <-chan struct{}) (result autorest.Response, err error) {
	req, err := client.DeletePreparer(resourceGroupName, namespaceName, cancel)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "Delete", nil, "Failure preparing request")
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "Delete", resp, "Failure sending request")
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "Delete", resp, "Failure responding to request")
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client NamespacesClient) DeletePreparer(resourceGroupName string, namespaceName string, cancel <-chan struct{}) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"namespaceName":     autorest.Encode("path", namespaceName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NotificationHubs/namespaces/{namespaceName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{Cancel: cancel})
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client NamespacesClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoPollForAsynchronous(client.PollingDelay))
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client NamespacesClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}

// DeleteAuthorizationRule deletes a namespace authorization rule
//
// resourceGroupName is the name of the resource group. namespaceName is the
// namespace name. authorizationRuleName is authorization Rule Name.
func (client NamespacesClient) DeleteAuthorizationRule(resourceGroupName string, namespaceName string, authorizationRuleName string) (result autorest.Response, err error) {
	req, err := client.DeleteAuthorizationRulePreparer(resourceGroupName, namespaceName, authorizationRuleName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "DeleteAuthorizationRule", nil, "Failure preparing request")
	}

	resp, err := client.DeleteAuthorizationRuleSender(req)
	if err != nil {
		result.Response = resp
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "DeleteAuthorizationRule", resp, "Failure sending request")
	}

	result, err = client.DeleteAuthorizationRuleResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "DeleteAuthorizationRule", resp, "Failure responding to request")
	}

	return
}

// DeleteAuthorizationRulePreparer prepares the DeleteAuthorizationRule request.
func (client NamespacesClient) DeleteAuthorizationRulePreparer(resourceGroupName string, namespaceName string, authorizationRuleName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"authorizationRuleName": autorest.Encode("path", authorizationRuleName),
		"namespaceName":         autorest.Encode("path", namespaceName),
		"resourceGroupName":     autorest.Encode("path", resourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NotificationHubs/namespaces/{namespaceName}/AuthorizationRules/{authorizationRuleName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// DeleteAuthorizationRuleSender sends the DeleteAuthorizationRule request. The method will close the
// http.Response Body if it receives an error.
func (client NamespacesClient) DeleteAuthorizationRuleSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// DeleteAuthorizationRuleResponder handles the response to the DeleteAuthorizationRule request. The method always
// closes the http.Response Body.
func (client NamespacesClient) DeleteAuthorizationRuleResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get returns the description for the specified namespace.
//
// resourceGroupName is the name of the resource group. namespaceName is the
// namespace name.
func (client NamespacesClient) Get(resourceGroupName string, namespaceName string) (result NamespaceResource, err error) {
	req, err := client.GetPreparer(resourceGroupName, namespaceName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "Get", nil, "Failure preparing request")
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "Get", resp, "Failure sending request")
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client NamespacesClient) GetPreparer(resourceGroupName string, namespaceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"namespaceName":     autorest.Encode("path", namespaceName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NotificationHubs/namespaces/{namespaceName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client NamespacesClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client NamespacesClient) GetResponder(resp *http.Response) (result NamespaceResource, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetAuthorizationRule gets an authorization rule for a namespace by name.
//
// resourceGroupName is the name of the resource group. namespaceName is the
// namespace name authorizationRuleName is authorization rule name.
func (client NamespacesClient) GetAuthorizationRule(resourceGroupName string, namespaceName string, authorizationRuleName string) (result SharedAccessAuthorizationRuleResource, err error) {
	req, err := client.GetAuthorizationRulePreparer(resourceGroupName, namespaceName, authorizationRuleName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "GetAuthorizationRule", nil, "Failure preparing request")
	}

	resp, err := client.GetAuthorizationRuleSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "GetAuthorizationRule", resp, "Failure sending request")
	}

	result, err = client.GetAuthorizationRuleResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "GetAuthorizationRule", resp, "Failure responding to request")
	}

	return
}

// GetAuthorizationRulePreparer prepares the GetAuthorizationRule request.
func (client NamespacesClient) GetAuthorizationRulePreparer(resourceGroupName string, namespaceName string, authorizationRuleName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"authorizationRuleName": autorest.Encode("path", authorizationRuleName),
		"namespaceName":         autorest.Encode("path", namespaceName),
		"resourceGroupName":     autorest.Encode("path", resourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NotificationHubs/namespaces/{namespaceName}/AuthorizationRules/{authorizationRuleName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetAuthorizationRuleSender sends the GetAuthorizationRule request. The method will close the
// http.Response Body if it receives an error.
func (client NamespacesClient) GetAuthorizationRuleSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetAuthorizationRuleResponder handles the response to the GetAuthorizationRule request. The method always
// closes the http.Response Body.
func (client NamespacesClient) GetAuthorizationRuleResponder(resp *http.Response) (result SharedAccessAuthorizationRuleResource, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetLongRunningOperationStatus the Get Operation Status operation returns
// the status of the specified operation. After calling an asynchronous
// operation, you can call Get Operation Status to determine whether the
// operation has succeeded, failed, or is still in progress.
//
// operationStatusLink is location value returned by the Begin operation.
func (client NamespacesClient) GetLongRunningOperationStatus(operationStatusLink string) (result autorest.Response, err error) {
	req, err := client.GetLongRunningOperationStatusPreparer(operationStatusLink)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "GetLongRunningOperationStatus", nil, "Failure preparing request")
	}

	resp, err := client.GetLongRunningOperationStatusSender(req)
	if err != nil {
		result.Response = resp
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "GetLongRunningOperationStatus", resp, "Failure sending request")
	}

	result, err = client.GetLongRunningOperationStatusResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "GetLongRunningOperationStatus", resp, "Failure responding to request")
	}

	return
}

// GetLongRunningOperationStatusPreparer prepares the GetLongRunningOperationStatus request.
func (client NamespacesClient) GetLongRunningOperationStatusPreparer(operationStatusLink string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"operationStatusLink": operationStatusLink,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/{operationStatusLink}", pathParameters))
	return preparer.Prepare(&http.Request{})
}

// GetLongRunningOperationStatusSender sends the GetLongRunningOperationStatus request. The method will close the
// http.Response Body if it receives an error.
func (client NamespacesClient) GetLongRunningOperationStatusSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetLongRunningOperationStatusResponder handles the response to the GetLongRunningOperationStatus request. The method always
// closes the http.Response Body.
func (client NamespacesClient) GetLongRunningOperationStatusResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusNotFound, http.StatusAccepted, http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}

// List lists the available namespaces within a resourceGroup.
//
// resourceGroupName is the name of the resource group. If resourceGroupName
// value is null the method lists all the namespaces within subscription
func (client NamespacesClient) List(resourceGroupName string) (result NamespaceListResult, err error) {
	req, err := client.ListPreparer(resourceGroupName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "List", nil, "Failure preparing request")
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "List", resp, "Failure sending request")
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client NamespacesClient) ListPreparer(resourceGroupName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NotificationHubs/namespaces", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client NamespacesClient) ListSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client NamespacesClient) ListResponder(resp *http.Response) (result NamespaceListResult, err error) {
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
func (client NamespacesClient) ListNextResults(lastResults NamespaceListResult) (result NamespaceListResult, err error) {
	req, err := lastResults.NamespaceListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "List", nil, "Failure preparing next results request request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "List", resp, "Failure sending next results request request")
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "List", resp, "Failure responding to next results request request")
	}

	return
}

// ListAll lists all the available namespaces within the subscription
// irrespective of the resourceGroups.
func (client NamespacesClient) ListAll() (result NamespaceListResult, err error) {
	req, err := client.ListAllPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "ListAll", nil, "Failure preparing request")
	}

	resp, err := client.ListAllSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "ListAll", resp, "Failure sending request")
	}

	result, err = client.ListAllResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "ListAll", resp, "Failure responding to request")
	}

	return
}

// ListAllPreparer prepares the ListAll request.
func (client NamespacesClient) ListAllPreparer() (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.NotificationHubs/namespaces", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListAllSender sends the ListAll request. The method will close the
// http.Response Body if it receives an error.
func (client NamespacesClient) ListAllSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListAllResponder handles the response to the ListAll request. The method always
// closes the http.Response Body.
func (client NamespacesClient) ListAllResponder(resp *http.Response) (result NamespaceListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListAllNextResults retrieves the next set of results, if any.
func (client NamespacesClient) ListAllNextResults(lastResults NamespaceListResult) (result NamespaceListResult, err error) {
	req, err := lastResults.NamespaceListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "ListAll", nil, "Failure preparing next results request request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListAllSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "ListAll", resp, "Failure sending next results request request")
	}

	result, err = client.ListAllResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "ListAll", resp, "Failure responding to next results request request")
	}

	return
}

// ListAuthorizationRules gets the authorization rules for a namespace.
//
// resourceGroupName is the name of the resource group. namespaceName is the
// namespace name
func (client NamespacesClient) ListAuthorizationRules(resourceGroupName string, namespaceName string) (result SharedAccessAuthorizationRuleListResult, err error) {
	req, err := client.ListAuthorizationRulesPreparer(resourceGroupName, namespaceName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "ListAuthorizationRules", nil, "Failure preparing request")
	}

	resp, err := client.ListAuthorizationRulesSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "ListAuthorizationRules", resp, "Failure sending request")
	}

	result, err = client.ListAuthorizationRulesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "ListAuthorizationRules", resp, "Failure responding to request")
	}

	return
}

// ListAuthorizationRulesPreparer prepares the ListAuthorizationRules request.
func (client NamespacesClient) ListAuthorizationRulesPreparer(resourceGroupName string, namespaceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"namespaceName":     autorest.Encode("path", namespaceName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NotificationHubs/namespaces/{namespaceName}/AuthorizationRules", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListAuthorizationRulesSender sends the ListAuthorizationRules request. The method will close the
// http.Response Body if it receives an error.
func (client NamespacesClient) ListAuthorizationRulesSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListAuthorizationRulesResponder handles the response to the ListAuthorizationRules request. The method always
// closes the http.Response Body.
func (client NamespacesClient) ListAuthorizationRulesResponder(resp *http.Response) (result SharedAccessAuthorizationRuleListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListAuthorizationRulesNextResults retrieves the next set of results, if any.
func (client NamespacesClient) ListAuthorizationRulesNextResults(lastResults SharedAccessAuthorizationRuleListResult) (result SharedAccessAuthorizationRuleListResult, err error) {
	req, err := lastResults.SharedAccessAuthorizationRuleListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "ListAuthorizationRules", nil, "Failure preparing next results request request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListAuthorizationRulesSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "ListAuthorizationRules", resp, "Failure sending next results request request")
	}

	result, err = client.ListAuthorizationRulesResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "ListAuthorizationRules", resp, "Failure responding to next results request request")
	}

	return
}

// ListKeys gets the Primary and Secondary ConnectionStrings to the namespace
//
// resourceGroupName is the name of the resource group. namespaceName is the
// namespace name. authorizationRuleName is the connection string of the
// namespace for the specified authorizationRule.
func (client NamespacesClient) ListKeys(resourceGroupName string, namespaceName string, authorizationRuleName string) (result ResourceListKeys, err error) {
	req, err := client.ListKeysPreparer(resourceGroupName, namespaceName, authorizationRuleName)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "ListKeys", nil, "Failure preparing request")
	}

	resp, err := client.ListKeysSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "ListKeys", resp, "Failure sending request")
	}

	result, err = client.ListKeysResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "notificationhubs.NamespacesClient", "ListKeys", resp, "Failure responding to request")
	}

	return
}

// ListKeysPreparer prepares the ListKeys request.
func (client NamespacesClient) ListKeysPreparer(resourceGroupName string, namespaceName string, authorizationRuleName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"authorizationRuleName": autorest.Encode("path", authorizationRuleName),
		"namespaceName":         autorest.Encode("path", namespaceName),
		"resourceGroupName":     autorest.Encode("path", resourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NotificationHubs/namespaces/{namespaceName}/AuthorizationRules/{authorizationRuleName}/listKeys", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListKeysSender sends the ListKeys request. The method will close the
// http.Response Body if it receives an error.
func (client NamespacesClient) ListKeysSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListKeysResponder handles the response to the ListKeys request. The method always
// closes the http.Response Body.
func (client NamespacesClient) ListKeysResponder(resp *http.Response) (result ResourceListKeys, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
