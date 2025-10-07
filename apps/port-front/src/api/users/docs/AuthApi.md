# AccountNewsbroCc.AuthApi

All URIs are relative to *http://account.newsbro.cc*

Method | HTTP request | Description
------------- | ------------- | -------------
[**authLoginPost**](AuthApi.md#authLoginPost) | **POST** /auth/login | Login a user
[**authOauthCallbackGet**](AuthApi.md#authOauthCallbackGet) | **GET** /auth/oauth/callback | OAuth Callback
[**authOauthLoginGet**](AuthApi.md#authOauthLoginGet) | **GET** /auth/oauth/login | OAuth Login
[**authRefreshPost**](AuthApi.md#authRefreshPost) | **POST** /auth/refresh | Refresh JWT token
[**authRegisterPost**](AuthApi.md#authRegisterPost) | **POST** /auth/register | Register a new user



## authLoginPost

> RepoAccountSrcApiDtoLoginResponse authLoginPost(loginRequest)

Login a user

Login a user with email and password

### Example

```javascript
import AccountNewsbroCc from 'account_newsbro_cc';

let apiInstance = new AccountNewsbroCc.AuthApi();
let loginRequest = new AccountNewsbroCc.RepoAccountSrcApiDtoLoginRequest(); // RepoAccountSrcApiDtoLoginRequest | Login Request
apiInstance.authLoginPost(loginRequest, (error, data, response) => {
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
 **loginRequest** | [**RepoAccountSrcApiDtoLoginRequest**](RepoAccountSrcApiDtoLoginRequest.md)| Login Request | 

### Return type

[**RepoAccountSrcApiDtoLoginResponse**](RepoAccountSrcApiDtoLoginResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## authOauthCallbackGet

> authOauthCallbackGet(code, state)

OAuth Callback

Handle OAuth provider callback

### Example

```javascript
import AccountNewsbroCc from 'account_newsbro_cc';

let apiInstance = new AccountNewsbroCc.AuthApi();
let code = "code_example"; // String | Authorization Code
let state = "state_example"; // String | State Parameter
apiInstance.authOauthCallbackGet(code, state, (error, data, response) => {
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
 **code** | **String**| Authorization Code | 
 **state** | **String**| State Parameter | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## authOauthLoginGet

> authOauthLoginGet()

OAuth Login

Redirect to OAuth provider for login

### Example

```javascript
import AccountNewsbroCc from 'account_newsbro_cc';

let apiInstance = new AccountNewsbroCc.AuthApi();
apiInstance.authOauthLoginGet((error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
});
```

### Parameters

This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## authRefreshPost

> RepoAccountSrcApiDtoLoginResponse authRefreshPost(authorization)

Refresh JWT token

Refresh JWT token using refresh token

### Example

```javascript
import AccountNewsbroCc from 'account_newsbro_cc';

let apiInstance = new AccountNewsbroCc.AuthApi();
let authorization = "'Bearer <Add access token here>'"; // String | Insert your access token
apiInstance.authRefreshPost(authorization, (error, data, response) => {
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

[**RepoAccountSrcApiDtoLoginResponse**](RepoAccountSrcApiDtoLoginResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## authRegisterPost

> authRegisterPost(registerRequest)

Register a new user

Register a new user with email and password

### Example

```javascript
import AccountNewsbroCc from 'account_newsbro_cc';

let apiInstance = new AccountNewsbroCc.AuthApi();
let registerRequest = new AccountNewsbroCc.RepoAccountSrcApiDtoRegisterRequest(); // RepoAccountSrcApiDtoRegisterRequest | Register Request
apiInstance.authRegisterPost(registerRequest, (error, data, response) => {
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
 **registerRequest** | [**RepoAccountSrcApiDtoRegisterRequest**](RepoAccountSrcApiDtoRegisterRequest.md)| Register Request | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

