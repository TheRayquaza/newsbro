import { ENV } from "../env";
import ArticlesApi from './articles/src/api/ArticlesApi';
import FeedbackApi from './articles/src/api/FeedbackApi';
import UserApi from './users/src/api/UserApi';
import AuthApi from './users/src/api/AuthApi';

import RepoAccountSrcApiDtoLoginRequest from './users/src/model/RepoAccountSrcApiDtoLoginRequest';
import RepoAccountSrcApiDtoRegisterRequest from './users/src/model/RepoAccountSrcApiDtoRegisterRequest';
import RepoArticleSrcApiDtoArticleCreateRequest from './articles/src/model/RepoArticleSrcApiDtoArticleCreateRequest';
import RepoArticleSrcApiDtoArticleUpdateRequest from './articles/src/model/RepoArticleSrcApiDtoArticleUpdateRequest';
import RepoArticleSrcApiDtoArticleTriggerIngestionRequest from './articles/src/model/RepoArticleSrcApiDtoArticleTriggerIngestionRequest';
import RepoArticleSrcApiDtoFeedbackRequest from './articles/src/model/RepoArticleSrcApiDtoFeedbackRequest';

import AccountApiClient from "./users/src/ApiClient";
import ArticlesApiClient from "./articles/src/ApiClient";


class ApiService {
  constructor() {
    // ---- Initialize API clients ----
    this.articlesApi = new ArticlesApi(new ArticlesApiClient(ENV.ARTICLE_BASE_URL));
    this.feedbackApi = new FeedbackApi(new ArticlesApiClient(ENV.ARTICLE_BASE_URL));
    this.authApi = new AuthApi(new AccountApiClient(ENV.ACCOUNT_BASE_URL));
    this.userApi = new UserApi(new AccountApiClient(ENV.ACCOUNT_BASE_URL));
  }

  // -------------------- AUTH --------------------
  async login(email, password) {
    const loginRequest = new RepoAccountSrcApiDtoLoginRequest();
    loginRequest.email = email;
    loginRequest.password = password;

    return new Promise((resolve, reject) => {
      this.authApi.authLoginPost(loginRequest, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async oauthCallback(params) {
    return new Promise((resolve, reject) => {
      this.authApi.authOauthCallbackGet(params, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  forgeOAuthUrl() {
    return `${ENV.ACCOUNT_BASE_URL}/api/v1/auth/oauth/login`;
  }

  async refreshToken(token) {
    return new Promise((resolve, reject) => {
      this.authApi.authRefreshPost(`Bearer ${token}`, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async register(userData) {
    const registerRequest = new RepoAccountSrcApiDtoRegisterRequest();
    registerRequest.email = userData.email;
    registerRequest.password = userData.password;
    registerRequest.username = userData.username;
    registerRequest.first_name = userData.first_name;
    registerRequest.last_name = userData.last_name;

    return new Promise((resolve, reject) => {
      this.authApi.authRegisterPost(registerRequest, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async getUsers(token) {
    return new Promise((resolve, reject) => {
      this.userApi.usersGet(`Bearer ${token}`, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async getProfile(token) {
    return new Promise((resolve, reject) => {
      this.userApi.usersProfileGet(`Bearer ${token}`, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async updateProfile(profileData, token) {
    return new Promise((resolve, reject) => {
      console.log('Updating profile with data:', profileData);
      this.userApi.usersProfilePut(`Bearer ${token}`, profileData, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  // -------------------- ARTICLES --------------------
  async getCategories(token) {
    return new Promise((resolve, reject) => {
      this.articlesApi.articlesCategoriesGet(`Bearer ${token}`, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async getArticles(token, opts = {}) {
    return new Promise((resolve, reject) => {
      this.articlesApi.articlesGet(`Bearer ${token}`, opts, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async getArticleById(id, token) {
    return new Promise((resolve, reject) => {
      this.articlesApi.articlesIdGet(id, `Bearer ${token}`, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async createArticle(articleData, token) {
    const createRequest = new RepoArticleSrcApiDtoArticleCreateRequest();
    createRequest.title = articleData.title;
    createRequest.category = articleData.category;
    createRequest.subcategory = articleData.subcategory;
    createRequest.link = articleData.link;
    createRequest.published_at = articleData.published_at;
    createRequest.abstract = articleData.abstract;
    if (articleData.id !== undefined) {
      createRequest.id = articleData.id;
    }

    return new Promise((resolve, reject) => {
      this.articlesApi.articlesPost(createRequest, `Bearer ${token}`, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async updateArticle(id, articleData, token) {
    const updateRequest = new RepoArticleSrcApiDtoArticleUpdateRequest();
    updateRequest.id = id;
    if (articleData.title !== undefined) updateRequest.title = articleData.title;
    if (articleData.category !== undefined) updateRequest.category = articleData.category;
    if (articleData.subcategory !== undefined) updateRequest.subcategory = articleData.subcategory;
    if (articleData.link !== undefined) updateRequest.link = articleData.link;
    if (articleData.abstract !== undefined) updateRequest.abstract = articleData.abstract;

    return new Promise((resolve, reject) => {
      this.articlesApi.articlesIdPut(id, updateRequest, `Bearer ${token}`, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async deleteArticle(id, token) {
    return new Promise((resolve, reject) => {
      this.articlesApi.articlesIdDelete(id, `Bearer ${token}`, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async triggerIngestion(articleData, token) {
    const ingestionRequest = new RepoArticleSrcApiDtoArticleTriggerIngestionRequest();
    ingestionRequest.begin_date = articleData.begin_date;
    ingestionRequest.end_date = articleData.end_date;

    return new Promise((resolve, reject) => {
      this.articlesApi.articlesIngestionPost(ingestionRequest, `Bearer ${token}`, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async getSubcategories(token, category) {
    if (category === undefined) {
      return new Promise((resolve, reject) => {
        this.articlesApi.articlesSubcategoriesGet(`Bearer ${token}`, (error, data) => {
          if (error) return reject(error);
          resolve(data);
        });
      });
    }
    return new Promise((resolve, reject) => {
      this.articlesApi.articlesSubcategoriesGet(`Bearer ${token}`, category, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  // -------------------- FEEDBACK --------------------
  async getArticleFeedback(articleId, token) {
    return new Promise((resolve, reject) => {
      this.feedbackApi.articlesIdFeedbackGet(articleId, `Bearer ${token}`, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async createArticleFeedback(articleId, feedbackData, token) {
    const feedbackRequest = new RepoArticleSrcApiDtoFeedbackRequest();
    feedbackRequest.value = feedbackData.value;

    return new Promise((resolve, reject) => {
      this.feedbackApi.articlesIdFeedbackPost(articleId, `Bearer ${token}`, feedbackRequest, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async deleteArticleFeedback(articleId, token) {
    return new Promise((resolve, reject) => {
      this.feedbackApi.articlesIdFeedbackDelete(articleId, `Bearer ${token}`, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async getAllFeedback(token) {
    return new Promise((resolve, reject) => {
      this.feedbackApi.feedbackAllGet(`Bearer ${token}`, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async exportFeedbackCsv(token) {
    return new Promise((resolve, reject) => {
      this.feedbackApi.feedbackCsvGet(`Bearer ${token}`, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async triggerFeedbackIngestion(token) {
    return new Promise((resolve, reject) => {
      this.feedbackApi.feedbackIngestPost(`Bearer ${token}`, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async getMyFeedback(token) {
    return new Promise((resolve, reject) => {
      this.feedbackApi.feedbackMyGet(`Bearer ${token}`, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async getFeedbackStats(token) {
    return new Promise((resolve, reject) => {
      this.feedbackApi.feedbackStatsGet(`Bearer ${token}`, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }
}

const api = new ApiService();
export default api;
