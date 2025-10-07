# ArticleNewsbroCc.ArticlesApi

All URIs are relative to *http://article.newsbro.cc*

Method | HTTP request | Description
------------- | ------------- | -------------
[**articlesCategoriesGet**](ArticlesApi.md#articlesCategoriesGet) | **GET** /articles/categories | Get categories
[**articlesGet**](ArticlesApi.md#articlesGet) | **GET** /articles | Get all articles
[**articlesIdDelete**](ArticlesApi.md#articlesIdDelete) | **DELETE** /articles/{id} | Delete article
[**articlesIdGet**](ArticlesApi.md#articlesIdGet) | **GET** /articles/{id} | Get article by ID
[**articlesIdPut**](ArticlesApi.md#articlesIdPut) | **PUT** /articles/{id} | Update article
[**articlesIngestionPost**](ArticlesApi.md#articlesIngestionPost) | **POST** /articles/ingestion | Trigger article ingestion
[**articlesPost**](ArticlesApi.md#articlesPost) | **POST** /articles | Create a new article
[**articlesSubcategoriesGet**](ArticlesApi.md#articlesSubcategoriesGet) | **GET** /articles/subcategories | Get subcategories



## articlesCategoriesGet

> [String] articlesCategoriesGet(authorization)

Get categories

Get all available article categories

### Example

```javascript
import ArticleNewsbroCc from 'article_newsbro_cc';

let apiInstance = new ArticleNewsbroCc.ArticlesApi();
let authorization = "'Bearer <Add access token here>'"; // String | Insert your access token
apiInstance.articlesCategoriesGet(authorization, (error, data, response) => {
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


## articlesGet

> [RepoArticleSrcApiDtoArticleResponse] articlesGet(authorization, opts)

Get all articles

Get a list of articles with optional filtering and pagination

### Example

```javascript
import ArticleNewsbroCc from 'article_newsbro_cc';

let apiInstance = new ArticleNewsbroCc.ArticlesApi();
let authorization = "'Bearer <Add access token here>'"; // String | Insert your access token
let opts = {
  'category': "category_example", // String | Filter by category
  'subcategory': "subcategory_example", // String | Filter by subcategory
  'search': "search_example", // String | Search in title and abstract
  'limit': 10, // Number | Limit number of results
  'offset': 0 // Number | Offset for pagination
};
apiInstance.articlesGet(authorization, opts, (error, data, response) => {
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
 **category** | **String**| Filter by category | [optional] 
 **subcategory** | **String**| Filter by subcategory | [optional] 
 **search** | **String**| Search in title and abstract | [optional] 
 **limit** | **Number**| Limit number of results | [optional] [default to 10]
 **offset** | **Number**| Offset for pagination | [optional] [default to 0]

### Return type

[**[RepoArticleSrcApiDtoArticleResponse]**](RepoArticleSrcApiDtoArticleResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## articlesIdDelete

> articlesIdDelete(id, authorization)

Delete article

Delete an article by ID

### Example

```javascript
import ArticleNewsbroCc from 'article_newsbro_cc';

let apiInstance = new ArticleNewsbroCc.ArticlesApi();
let id = 56; // Number | Article ID
let authorization = "'Bearer <Add access token here>'"; // String | Insert your access token
apiInstance.articlesIdDelete(id, authorization, (error, data, response) => {
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
- **Accept**: */*


## articlesIdGet

> RepoArticleSrcApiDtoArticleResponse articlesIdGet(id, authorization)

Get article by ID

Get a specific article by its ID

### Example

```javascript
import ArticleNewsbroCc from 'article_newsbro_cc';

let apiInstance = new ArticleNewsbroCc.ArticlesApi();
let id = 56; // Number | Article ID
let authorization = "'Bearer <Add access token here>'"; // String | Insert your access token
apiInstance.articlesIdGet(id, authorization, (error, data, response) => {
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

[**RepoArticleSrcApiDtoArticleResponse**](RepoArticleSrcApiDtoArticleResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## articlesIdPut

> RepoArticleSrcApiDtoArticleResponse articlesIdPut(id, authorization, article)

Update article

Update an existing article by ID

### Example

```javascript
import ArticleNewsbroCc from 'article_newsbro_cc';

let apiInstance = new ArticleNewsbroCc.ArticlesApi();
let id = 56; // Number | Article ID
let authorization = "'Bearer <Add access token here>'"; // String | Insert your access token
let article = new ArticleNewsbroCc.RepoArticleSrcApiDtoArticleUpdateRequest(); // RepoArticleSrcApiDtoArticleUpdateRequest | Article update data
apiInstance.articlesIdPut(id, authorization, article, (error, data, response) => {
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
 **article** | [**RepoArticleSrcApiDtoArticleUpdateRequest**](RepoArticleSrcApiDtoArticleUpdateRequest.md)| Article update data | 

### Return type

[**RepoArticleSrcApiDtoArticleResponse**](RepoArticleSrcApiDtoArticleResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## articlesIngestionPost

> {String: Object} articlesIngestionPost(authorization, article)

Trigger article ingestion

Trigger ingestion of articles published within a specified date range

### Example

```javascript
import ArticleNewsbroCc from 'article_newsbro_cc';

let apiInstance = new ArticleNewsbroCc.ArticlesApi();
let authorization = "'Bearer <Add access token here>'"; // String | Insert your access token
let article = new ArticleNewsbroCc.RepoArticleSrcApiDtoArticleTriggerIngestionRequest(); // RepoArticleSrcApiDtoArticleTriggerIngestionRequest | Article update data
apiInstance.articlesIngestionPost(authorization, article, (error, data, response) => {
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
 **article** | [**RepoArticleSrcApiDtoArticleTriggerIngestionRequest**](RepoArticleSrcApiDtoArticleTriggerIngestionRequest.md)| Article update data | 

### Return type

**{String: Object}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## articlesPost

> RepoArticleSrcApiDtoArticleResponse articlesPost(authorization, article)

Create a new article

Create a new article with the provided data

### Example

```javascript
import ArticleNewsbroCc from 'article_newsbro_cc';

let apiInstance = new ArticleNewsbroCc.ArticlesApi();
let authorization = "'Bearer <Add access token here>'"; // String | Insert your access token
let article = new ArticleNewsbroCc.RepoArticleSrcApiDtoArticleCreateRequest(); // RepoArticleSrcApiDtoArticleCreateRequest | Article data
apiInstance.articlesPost(authorization, article, (error, data, response) => {
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
 **article** | [**RepoArticleSrcApiDtoArticleCreateRequest**](RepoArticleSrcApiDtoArticleCreateRequest.md)| Article data | 

### Return type

[**RepoArticleSrcApiDtoArticleResponse**](RepoArticleSrcApiDtoArticleResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## articlesSubcategoriesGet

> [String] articlesSubcategoriesGet(authorization, opts)

Get subcategories

Get subcategories for a specific category or all subcategories

### Example

```javascript
import ArticleNewsbroCc from 'article_newsbro_cc';

let apiInstance = new ArticleNewsbroCc.ArticlesApi();
let authorization = "'Bearer <Add access token here>'"; // String | Insert your access token
let opts = {
  'category': "category_example" // String | Filter subcategories by category
};
apiInstance.articlesSubcategoriesGet(authorization, opts, (error, data, response) => {
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
 **category** | **String**| Filter subcategories by category | [optional] 

### Return type

**[String]**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

