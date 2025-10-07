# AccountNewsbroCc.UserApi

All URIs are relative to *http://account.newsbro.cc*

Method | HTTP request | Description
------------- | ------------- | -------------
[**usersGet**](UserApi.md#usersGet) | **GET** /users | Get all users
[**usersProfileGet**](UserApi.md#usersProfileGet) | **GET** /users/profile | Get user profile
[**usersProfilePut**](UserApi.md#usersProfilePut) | **PUT** /users/profile | Update user profile



## usersGet

> [RepoAccountSrcApiDtoUserResponse] usersGet(authorization, opts)

Get all users

Get a list of all users (admin only)

### Example

```javascript
import AccountNewsbroCc from 'account_newsbro_cc';

let apiInstance = new AccountNewsbroCc.UserApi();
let authorization = "'Bearer <Add access token here>'"; // String | Insert your access token
let opts = {
  'limit': 10, // Number | Limit
  'offset': 0 // Number | Offset
};
apiInstance.usersGet(authorization, opts, (error, data, response) => {
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
 **limit** | **Number**| Limit | [optional] [default to 10]
 **offset** | **Number**| Offset | [optional] [default to 0]

### Return type

[**[RepoAccountSrcApiDtoUserResponse]**](RepoAccountSrcApiDtoUserResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## usersProfileGet

> RepoAccountSrcApiDtoUserResponse usersProfileGet(authorization)

Get user profile

Get the profile of the authenticated user

### Example

```javascript
import AccountNewsbroCc from 'account_newsbro_cc';

let apiInstance = new AccountNewsbroCc.UserApi();
let authorization = "'Bearer <Add access token here>'"; // String | Insert your access token
apiInstance.usersProfileGet(authorization, (error, data, response) => {
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

[**RepoAccountSrcApiDtoUserResponse**](RepoAccountSrcApiDtoUserResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## usersProfilePut

> RepoAccountSrcApiDtoUserResponse usersProfilePut(authorization, updateRequest)

Update user profile

Update the profile of the authenticated user

### Example

```javascript
import AccountNewsbroCc from 'account_newsbro_cc';

let apiInstance = new AccountNewsbroCc.UserApi();
let authorization = "'Bearer <Add access token here>'"; // String | Insert your access token
let updateRequest = {key: null}; // {String: Object} | Update Request
apiInstance.usersProfilePut(authorization, updateRequest, (error, data, response) => {
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
 **updateRequest** | [**{String: Object}**](Object.md)| Update Request | 

### Return type

[**RepoAccountSrcApiDtoUserResponse**](RepoAccountSrcApiDtoUserResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

