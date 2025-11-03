# RepoFeedServiceApi.FeedApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**feedGet**](FeedApi.md#feedGet) | **GET** /feed | Get user feed
[**feedModelsGet**](FeedApi.md#feedModelsGet) | **GET** /feed/models | Get user feed models



## feedGet

> RepoFeedSrcApiDtoFeedResponse feedGet(authorization, opts)

Get user feed

Get the feed of the authenticated user

### Example

```javascript
import RepoFeedServiceApi from 'repo_feed_service_api';

let apiInstance = new RepoFeedServiceApi.FeedApi();
let authorization = "'Bearer <Add access token here>'"; // String | Insert your access token
let opts = {
  'model': "model_example" // String | Model name (optional, uses default if not specified)
};
apiInstance.feedGet(authorization, opts, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **authorization** | **String**| Insert your access token | [default to &#39;Bearer &lt;Add access token here&gt;&#39;]
 **model** | **String**| Model name (optional, uses default if not specified) | [optional] 

### Return type

[**RepoFeedSrcApiDtoFeedResponse**](RepoFeedSrcApiDtoFeedResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## feedModelsGet

> [String] feedModelsGet(authorization)

Get user feed models

Get the models available for user feed

### Example

```javascript
import RepoFeedServiceApi from 'repo_feed_service_api';

let apiInstance = new RepoFeedServiceApi.FeedApi();
let authorization = "'Bearer <Add access token here>'"; // String | Insert your access token
apiInstance.feedModelsGet(authorization, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **authorization** | **String**| Insert your access token | [default to &#39;Bearer &lt;Add access token here&gt;&#39;]

### Return type

**[String]**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

