import { ENV } from "../env";
import ArticlesApi from './articles/src/api/ArticlesApi';
import FeedbackApi from './articles/src/api/FeedbackApi';
import RSSApi from './articles/src/api/RSSApi';
import UserApi from './users/src/api/UserApi';
import AuthApi from './users/src/api/AuthApi';
import FeedApi from './feed/src/api/FeedApi';

import RepoAccountSrcApiDtoLoginRequest from './users/src/model/RepoAccountSrcApiDtoLoginRequest';
import RepoAccountSrcApiDtoRegisterRequest from './users/src/model/RepoAccountSrcApiDtoRegisterRequest';
import RepoArticleSrcApiDtoArticleCreateRequest from './articles/src/model/RepoArticleSrcApiDtoArticleCreateRequest';
import RepoArticleSrcApiDtoArticleUpdateRequest from './articles/src/model/RepoArticleSrcApiDtoArticleUpdateRequest';
import RepoArticleSrcApiDtoArticleTriggerIngestionRequest from './articles/src/model/RepoArticleSrcApiDtoArticleTriggerIngestionRequest';
import RepoArticleSrcApiDtoFeedbackRequest from './articles/src/model/RepoArticleSrcApiDtoFeedbackRequest';

import AccountApiClient from "./users/src/ApiClient";
import ArticlesApiClient from "./articles/src/ApiClient";
import FeedApiClient from "./feed/src/ApiClient";


class ApiService {
  constructor() {
    // ---- Initialize API clients ----
    this.articlesApi = new ArticlesApi(new ArticlesApiClient(ENV.ARTICLE_BASE_URL));
    this.feedbackApi = new FeedbackApi(new ArticlesApiClient(ENV.ARTICLE_BASE_URL));
    this.rssApi = new RSSApi(new ArticlesApiClient(ENV.ARTICLE_BASE_URL));
    this.authApi = new AuthApi(new AccountApiClient(ENV.ACCOUNT_BASE_URL));
    this.userApi = new UserApi(new AccountApiClient(ENV.ACCOUNT_BASE_URL));
    this.feedApi = new FeedApi(new FeedApiClient(ENV.FEED_BASE_URL));
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

  async refresh() {
    return new Promise((resolve, reject) => {
      this.authApi.authRefreshPost((error, data) => {
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

  async getUsers() {
    return new Promise((resolve, reject) => {
      this.userApi.usersGet((error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async getProfile() {
    return new Promise((resolve, reject) => {
      this.userApi.usersProfileGet((error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async updateProfile(profileData, ) {
    return new Promise((resolve, reject) => {
      this.userApi.usersProfilePut(profileData, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  // -------------------- ARTICLES --------------------
  async getCategories() {
    return new Promise((resolve, reject) => {
      this.articlesApi.articlesCategoriesGet((error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async getArticles(opts = {}) {
    return new Promise((resolve, reject) => {
      this.articlesApi.articlesGet(opts, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async getArticleHistory(opts = {}) {
    return new Promise((resolve, reject) => {
      this.articlesApi.articleHistoryGet(opts, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async getArticleById(id, ) {
    return new Promise((resolve, reject) => {
      this.articlesApi.articlesIdGet(id, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async createArticle(articleData, ) {
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
      this.articlesApi.articlesPost(createRequest, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async updateArticle(id, articleData, ) {
    const updateRequest = new RepoArticleSrcApiDtoArticleUpdateRequest();
    updateRequest.id = id;
    if (articleData.title !== undefined) updateRequest.title = articleData.title;
    if (articleData.category !== undefined) updateRequest.category = articleData.category;
    if (articleData.subcategory !== undefined) updateRequest.subcategory = articleData.subcategory;
    if (articleData.link !== undefined) updateRequest.link = articleData.link;
    if (articleData.abstract !== undefined) updateRequest.abstract = articleData.abstract;

    return new Promise((resolve, reject) => {
      this.articlesApi.articlesIdPut(id, updateRequest, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async deleteArticle(id) {
    return new Promise((resolve, reject) => {
      this.articlesApi.articlesIdDelete(id, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async triggerIngestion(articleData, ) {
    const ingestionRequest = new RepoArticleSrcApiDtoArticleTriggerIngestionRequest();
    ingestionRequest.begin_date = articleData.begin_date;
    ingestionRequest.end_date = articleData.end_date;

    return new Promise((resolve, reject) => {
      this.articlesApi.articlesIngestionPost(ingestionRequest, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async getSubcategories(category) {
    if (category === undefined) {
      return new Promise((resolve, reject) => {
        this.articlesApi.articlesSubcategoriesGet((error, data) => {
          if (error) return reject(error);
          resolve(data);
        });
      });
    }
    return new Promise((resolve, reject) => {
      var opts = { category: category };
      this.articlesApi.articlesSubcategoriesGet(opts, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  // -------------------- FEEDBACK --------------------
  async getArticleFeedback(articleId) {
    return new Promise((resolve, reject) => {
      this.feedbackApi.articlesIdFeedbackGet(articleId, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async createArticleFeedback(articleId, feedbackData) {
    const feedbackRequest = new RepoArticleSrcApiDtoFeedbackRequest();
    feedbackRequest.value = feedbackData.value;

    return new Promise((resolve, reject) => {
      this.feedbackApi.articlesIdFeedbackPost(articleId, feedbackRequest, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async deleteArticleFeedback(articleId) {
    return new Promise((resolve, reject) => {
      this.feedbackApi.articlesIdFeedbackDelete(articleId, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async getAllFeedback() {
    return new Promise((resolve, reject) => {
      this.feedbackApi.feedbackAllGet((error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async exportFeedbackCsv() {
    return new Promise((resolve, reject) => {
      this.feedbackApi.feedbackCsvGet((error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async triggerFeedbackIngestion() {
    return new Promise((resolve, reject) => {
      this.feedbackApi.feedbackIngestPost((error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async getMyFeedback() {
    return new Promise((resolve, reject) => {
      this.feedbackApi.feedbackMyGet((error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async getFeedbackStats() {
    return new Promise((resolve, reject) => {
      this.feedbackApi.feedbackStatsGet((error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  // -------------------- RSS --------------------
  async getRssTree(opts = {}) {
    return new Promise((resolve, reject) => {
      this.rssApi.rssTreeGet(opts, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  // -------------------- FEED --------------------
  async getFeedModels() {
    return new Promise((resolve, reject) => {
      this.feedApi.feedModelsGet((error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async getFeed(opts = {}) {
    return new Promise((resolve, reject) => {
      this.feedApi.feedGet(opts, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async removeArticleFromFeed(articleId) {
    return new Promise((resolve, reject) => {
      this.feedApi.feedIdDelete(articleId, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }
}

const api = new ApiService();
export default api;
