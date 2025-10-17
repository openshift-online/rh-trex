# \DefaultAPI

All URIs are relative to *http://localhost:8000*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiRhTrexV1DinosaursGet**](DefaultAPI.md#ApiRhTrexV1DinosaursGet) | **Get** /api/rh-trex/v1/dinosaurs | Returns a list of dinosaurs
[**ApiRhTrexV1DinosaursIdGet**](DefaultAPI.md#ApiRhTrexV1DinosaursIdGet) | **Get** /api/rh-trex/v1/dinosaurs/{id} | Get an dinosaur by id
[**ApiRhTrexV1DinosaursIdPatch**](DefaultAPI.md#ApiRhTrexV1DinosaursIdPatch) | **Patch** /api/rh-trex/v1/dinosaurs/{id} | Update an dinosaur
[**ApiRhTrexV1DinosaursPost**](DefaultAPI.md#ApiRhTrexV1DinosaursPost) | **Post** /api/rh-trex/v1/dinosaurs | Create a new dinosaur



## ApiRhTrexV1DinosaursGet

> DinosaurList ApiRhTrexV1DinosaursGet(ctx).Page(page).Size(size).Search(search).OrderBy(orderBy).Fields(fields).Execute()

Returns a list of dinosaurs

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	page := int32(56) // int32 | Page number of record list when record list exceeds specified page size (optional) (default to 1)
	size := int32(56) // int32 | Maximum number of records to return (optional) (default to 100)
	search := "search_example" // string | Specifies the search criteria. The syntax of this parameter is similar to the syntax of the _where_ clause of an SQL statement, using the names of the json attributes / column names of the account.  For example, in order to retrieve all the accounts with a username starting with `my`:  ```sql username like 'my%' ```  The search criteria can also be applied on related resource. For example, in order to retrieve all the subscriptions labeled by `foo=bar`,  ```sql subscription_labels.key = 'foo' and subscription_labels.value = 'bar' ```  If the parameter isn't provided, or if the value is empty, then all the accounts that the user has permission to see will be returned. (optional)
	orderBy := "orderBy_example" // string | Specifies the order by criteria. The syntax of this parameter is similar to the syntax of the _order by_ clause of an SQL statement, but using the names of the json attributes / column of the account. For example, in order to retrieve all accounts ordered by username:  ```sql username asc ```  Or in order to retrieve all accounts ordered by username _and_ first name:  ```sql username asc, firstName asc ```  If the parameter isn't provided, or if the value is empty, then no explicit ordering will be applied. (optional)
	fields := "fields_example" // string | Supplies a comma-separated list of fields to be returned. Fields of sub-structures and of arrays use <structure>.<field> notation. <stucture>.* means all field of a structure Example: For each Subscription to get id, href, plan(id and kind) and labels (all fields)  ``` ocm get subscriptions --parameter fields=id,href,plan.id,plan.kind,labels.* --parameter fetchLabels=true ``` (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.ApiRhTrexV1DinosaursGet(context.Background()).Page(page).Size(size).Search(search).OrderBy(orderBy).Fields(fields).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ApiRhTrexV1DinosaursGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRhTrexV1DinosaursGet`: DinosaurList
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ApiRhTrexV1DinosaursGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRhTrexV1DinosaursGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **int32** | Page number of record list when record list exceeds specified page size | [default to 1]
 **size** | **int32** | Maximum number of records to return | [default to 100]
 **search** | **string** | Specifies the search criteria. The syntax of this parameter is similar to the syntax of the _where_ clause of an SQL statement, using the names of the json attributes / column names of the account.  For example, in order to retrieve all the accounts with a username starting with &#x60;my&#x60;:  &#x60;&#x60;&#x60;sql username like &#39;my%&#39; &#x60;&#x60;&#x60;  The search criteria can also be applied on related resource. For example, in order to retrieve all the subscriptions labeled by &#x60;foo&#x3D;bar&#x60;,  &#x60;&#x60;&#x60;sql subscription_labels.key &#x3D; &#39;foo&#39; and subscription_labels.value &#x3D; &#39;bar&#39; &#x60;&#x60;&#x60;  If the parameter isn&#39;t provided, or if the value is empty, then all the accounts that the user has permission to see will be returned. | 
 **orderBy** | **string** | Specifies the order by criteria. The syntax of this parameter is similar to the syntax of the _order by_ clause of an SQL statement, but using the names of the json attributes / column of the account. For example, in order to retrieve all accounts ordered by username:  &#x60;&#x60;&#x60;sql username asc &#x60;&#x60;&#x60;  Or in order to retrieve all accounts ordered by username _and_ first name:  &#x60;&#x60;&#x60;sql username asc, firstName asc &#x60;&#x60;&#x60;  If the parameter isn&#39;t provided, or if the value is empty, then no explicit ordering will be applied. | 
 **fields** | **string** | Supplies a comma-separated list of fields to be returned. Fields of sub-structures and of arrays use &lt;structure&gt;.&lt;field&gt; notation. &lt;stucture&gt;.* means all field of a structure Example: For each Subscription to get id, href, plan(id and kind) and labels (all fields)  &#x60;&#x60;&#x60; ocm get subscriptions --parameter fields&#x3D;id,href,plan.id,plan.kind,labels.* --parameter fetchLabels&#x3D;true &#x60;&#x60;&#x60; | 

### Return type

[**DinosaurList**](DinosaurList.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRhTrexV1DinosaursIdGet

> Dinosaur ApiRhTrexV1DinosaursIdGet(ctx, id).Execute()

Get an dinosaur by id

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := "id_example" // string | The id of record

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.ApiRhTrexV1DinosaursIdGet(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ApiRhTrexV1DinosaursIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRhTrexV1DinosaursIdGet`: Dinosaur
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ApiRhTrexV1DinosaursIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The id of record | 

### Other Parameters

Other parameters are passed through a pointer to a apiApiRhTrexV1DinosaursIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**Dinosaur**](Dinosaur.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRhTrexV1DinosaursIdPatch

> Dinosaur ApiRhTrexV1DinosaursIdPatch(ctx, id).DinosaurPatchRequest(dinosaurPatchRequest).Execute()

Update an dinosaur

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := "id_example" // string | The id of record
	dinosaurPatchRequest := *openapiclient.NewDinosaurPatchRequest() // DinosaurPatchRequest | Updated dinosaur data

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.ApiRhTrexV1DinosaursIdPatch(context.Background(), id).DinosaurPatchRequest(dinosaurPatchRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ApiRhTrexV1DinosaursIdPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRhTrexV1DinosaursIdPatch`: Dinosaur
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ApiRhTrexV1DinosaursIdPatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The id of record | 

### Other Parameters

Other parameters are passed through a pointer to a apiApiRhTrexV1DinosaursIdPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **dinosaurPatchRequest** | [**DinosaurPatchRequest**](DinosaurPatchRequest.md) | Updated dinosaur data | 

### Return type

[**Dinosaur**](Dinosaur.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiRhTrexV1DinosaursPost

> Dinosaur ApiRhTrexV1DinosaursPost(ctx).Dinosaur(dinosaur).Execute()

Create a new dinosaur

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	dinosaur := *openapiclient.NewDinosaur("Species_example") // Dinosaur | Dinosaur data

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.ApiRhTrexV1DinosaursPost(context.Background()).Dinosaur(dinosaur).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.ApiRhTrexV1DinosaursPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiRhTrexV1DinosaursPost`: Dinosaur
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.ApiRhTrexV1DinosaursPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiRhTrexV1DinosaursPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **dinosaur** | [**Dinosaur**](Dinosaur.md) | Dinosaur data | 

### Return type

[**Dinosaur**](Dinosaur.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

