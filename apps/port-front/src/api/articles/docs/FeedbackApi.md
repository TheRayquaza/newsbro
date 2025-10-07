# ArticleNewsbroCc.FeedbackApi

All URIs are relative to *http://article.newsbro.cc*

Method | HTTP request | Description
------------- | ------------- | -------------
[**articlesIdFeedbackDelete**](FeedbackApi.md#articlesIdFeedbackDelete) | **DELETE** /articles/{id}/feedback | Delete feedback for an article
[**articlesIdFeedbackGet**](FeedbackApi.md#articlesIdFeedbackGet) | **GET** /articles/{id}/feedback | Get feedback for a specific article
[**articlesIdFeedbackPost**](FeedbackApi.md#articlesIdFeedbackPost) | **POST** /articles/{id}/feedback | Create feedback for an article
[**feedbackAllGet**](FeedbackApi.md#feedbackAllGet) | **GET** /feedback/all | Get all feedback
[**feedbackCsvGet**](FeedbackApi.md#feedbackCsvGet) | **GET** /feedback/csv | Export all feedback to CSV
[**feedbackIngestPost**](FeedbackApi.md#feedbackIngestPost) | **POST** /feedback/ingest | Trigger feedback ingestion
[**feedbackMyGet**](FeedbackApi.md#feedbackMyGet) | **GET** /feedback/my | Get user&#39;s feedback history
[**feedbackStatsGet**](FeedbackApi.md#feedbackStatsGet) | **GET** /feedback/stats | Get feedback statistics



## articlesIdFeedbackDelete

> articlesIdFeedbackDelete(id, authorization)

Delete feedback for an article

Delete user feedback for a specific article

### Example

```javascript
import ArticleNewsbroCc from 'article_newsbro_cc';

let apiInstance = new ArticleNewsbroCc.FeedbackApi();
let id = 56; // Number | Article ID
let authorization = "'Bearer <Add access token here>'"; // String | Insert your access token
apiInstance.articlesIdFeedbackDelete(id, authorization, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **Number**| Article ID | 
 **authorization** | **String**| Insert your access token | [default to &#39;Bearer &lt;Add access token here&gt;&#39;]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## articlesIdFeedbackGet

> RepoArticleSrcApiDtoFeedbackStatsResponse articlesIdFeedbackGet(id, authorization)

Get feedback for a specific article

Get feedback statistics and user&#39;s feedback for a specific article

### Example

```javascript
import ArticleNewsbroCc from 'article_newsbro_cc';

let apiInstance = new ArticleNewsbroCc.FeedbackApi();
let id = 56; // Number | Article ID
let authorization = "'Bearer <Add access token here>'"; // String | Insert your access token
apiInstance.articlesIdFeedbackGet(id, authorization, (error, data, response) => {
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
 **id** | **Number**| Article ID | 
 **authorization** | **String**| Insert your access token | [default to &#39;Bearer &lt;Add access token here&gt;&#39;]

### Return type

[**RepoArticleSrcApiDtoFeedbackStatsResponse**](RepoArticleSrcApiDtoFeedbackStatsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## articlesIdFeedbackPost

> {String: Object} articlesIdFeedbackPost(id, authorization, feedback)

Create feedback for an article

Create or update user feedback for a specific article

### Example

```javascript
import ArticleNewsbroCc from 'article_newsbro_cc';

let apiInstance = new ArticleNewsbroCc.FeedbackApi();
let id = 56; // Number | Article ID
let authorization = "'Bearer <Add access token here>'"; // String | Insert your access token
let feedback = new ArticleNewsbroCc.RepoArticleSrcApiDtoFeedbackRequest(); // RepoArticleSrcApiDtoFeedbackRequest | Feedback data
apiInstance.articlesIdFeedbackPost(id, authorization, feedback, (error, data, response) => {
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
 **id** | **Number**| Article ID | 
 **authorization** | **String**| Insert your access token | [default to &#39;Bearer &lt;Add access token here&gt;&#39;]
 **feedback** | [**RepoArticleSrcApiDtoFeedbackRequest**](RepoArticleSrcApiDtoFeedbackRequest.md)| Feedback data | 

### Return type

**{String: Object}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## feedbackAllGet

> [RepoArticleSrcApiDtoFeedbackResponse] feedbackAllGet(authorization, opts)

Get all feedback

Get all feedback with pagination (Admin only)

### Example

```javascript
import ArticleNewsbroCc from 'article_newsbro_cc';

let apiInstance = new ArticleNewsbroCc.FeedbackApi();
let authorization = "'Bearer <Add access token here>'"; // String | Insert your access token
let opts = {
  'page': 1, // Number | Page number
  'limit': 10 // Number | Items per page
};
apiInstance.feedbackAllGet(authorization, opts, (error, data, response) => {
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
 **page** | **Number**| Page number | [optional] [default to 1]
 **limit** | **Number**| Items per page | [optional] [default to 10]

### Return type

[**[RepoArticleSrcApiDtoFeedbackResponse]**](RepoArticleSrcApiDtoFeedbackResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## feedbackCsvGet

> File feedbackCsvGet(authorization, opts)

Export all feedback to CSV

Export all article feedback data to CSV format (Admin only)

### Example

```javascript
import ArticleNewsbroCc from 'article_newsbro_cc';

let apiInstance = new ArticleNewsbroCc.FeedbackApi();
let authorization = "'Bearer <Add access token here>'"; // String | Insert your access token
let opts = {
  'startDate': "startDate_example", // String | Start date (YYYY-MM-DD)
  'endDate': "endDate_example" // String | End date (YYYY-MM-DD)
};
apiInstance.feedbackCsvGet(authorization, opts, (error, data, response) => {
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
 **startDate** | **String**| Start date (YYYY-MM-DD) | [optional] 
 **endDate** | **String**| End date (YYYY-MM-DD) | [optional] 

### Return type

**File**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: text/csv


## feedbackIngestPost

> {String: Object} feedbackIngestPost(authorization, ingestion)

Trigger feedback ingestion

Trigger ingestion of feedback data for articles within a specified date range (Admin only)

### Example

```javascript
import ArticleNewsbroCc from 'article_newsbro_cc';

let apiInstance = new ArticleNewsbroCc.FeedbackApi();
let authorization = "'Bearer <Add access token here>'"; // String | Insert your access token
let ingestion = new ArticleNewsbroCc.RepoArticleSrcApiDtoFeedbackTriggerIngestionRequest(); // RepoArticleSrcApiDtoFeedbackTriggerIngestionRequest | Ingestion date range
apiInstance.feedbackIngestPost(authorization, ingestion, (error, data, response) => {
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
 **ingestion** | [**RepoArticleSrcApiDtoFeedbackTriggerIngestionRequest**](RepoArticleSrcApiDtoFeedbackTriggerIngestionRequest.md)| Ingestion date range | 

### Return type

**{String: Object}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: */*


## feedbackMyGet

> [RepoArticleSrcApiDtoFeedbackResponse] feedbackMyGet(authorization, opts)

Get user&#39;s feedback history

Get all feedback given by the authenticated user

### Example

```javascript
import ArticleNewsbroCc from 'article_newsbro_cc';

let apiInstance = new ArticleNewsbroCc.FeedbackApi();
let authorization = "'Bearer <Add access token here>'"; // String | Insert your access token
let opts = {
  'page': 1, // Number | Page number
  'limit': 10 // Number | Items per page
};
apiInstance.feedbackMyGet(authorization, opts, (error, data, response) => {
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
 **page** | **Number**| Page number | [optional] [default to 1]
 **limit** | **Number**| Items per page | [optional] [default to 10]

### Return type

[**[RepoArticleSrcApiDtoFeedbackResponse]**](RepoArticleSrcApiDtoFeedbackResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## feedbackStatsGet

> [RepoArticleSrcApiDtoFeedbackStatsResponse] feedbackStatsGet(authorization, opts)

Get feedback statistics

Get aggregated feedback statistics for all articles (Admin only)

### Example

```javascript
import ArticleNewsbroCc from 'article_newsbro_cc';

let apiInstance = new ArticleNewsbroCc.FeedbackApi();
let authorization = "'Bearer <Add access token here>'"; // String | Insert your access token
let opts = {
  'page': 1, // Number | Page number
  'limit': 10 // Number | Items per page
};
apiInstance.feedbackStatsGet(authorization, opts, (error, data, response) => {
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
 **page** | **Number**| Page number | [optional] [default to 1]
 **limit** | **Number**| Items per page | [optional] [default to 10]

### Return type

[**[RepoArticleSrcApiDtoFeedbackStatsResponse]**](RepoArticleSrcApiDtoFeedbackStatsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

