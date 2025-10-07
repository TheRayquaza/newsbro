import { ENV } from "../env";
import ArticlesApi from './articles/src/api/ArticlesApi';
import FeedbackApi from './articles/src/api/FeedbackApi';
import UserApi from './users/src/api/UserApi';
import AuthApi  from './users/src/api/AuthApi';
import RepoAccountSrcApiDtoLoginRequest from './users/src/model/RepoAccountSrcApiDtoLoginRequest';
import ApiClient  from "./users/src/ApiClient";
// import { RepoAccountSrcApiDtoRegisterRequest } from './auth/src/model/RepoAccountSrcApiDtoRegisterRequest';


class ApiService {
  constructor() {
    // ---- Initialize API clients ----
    this.articlesApi = new ArticlesApi(new ApiClient({ basePath: ENV.ARTICLE_BASE_URL + '/api/v1/' }));
    this.feedbackApi = new FeedbackApi(new ApiClient({ basePath: ENV.ARTICLE_BASE_URL + '/api/v1/' }));
    this.authApi = new AuthApi(new ApiClient({ basePath: ENV.ACCOUNT_BASE_URL + '/api/v1/' }));
    this.userApi = new UserApi(new ApiClient({ basePath: ENV.ACCOUNT_BASE_URL + '/api/v1/' }));
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

  // async oauthLogin() {
  //   return new Promise((resolve, reject) => {
  //     this.authApi.authOauthLoginGet((error, data) => {
  //       if (error) return reject(error);
  //       resolve(data);
  //     });
  //   });
  // }

  forgeOAuthUrl() {
    return `${ENV.API_BASE_URL}/api/v1/auth/oauth/login`;
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
    return new Promise((resolve, reject) => {
      this.authApi.authRegisterPost(userData, (error, data) => {
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
    return new Promise((resolve, reject) => {
      this.articlesApi.articlesPost(articleData, `Bearer ${token}`, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async updateArticle(id, articleData, token) {
    return new Promise((resolve, reject) => {
      this.articlesApi.articlesIdPut(id, articleData, `Bearer ${token}`, (error, data) => {
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
    return new Promise((resolve, reject) => {
      this.articlesApi.articlesIngestionPost(articleData, `Bearer ${token}`, (error, data) => {
        if (error) return reject(error);
        resolve(data);
      });
    });
  }

  async getSubcategories(token) {
    return new Promise((resolve, reject) => {
      this.articlesApi.articlesSubcategoriesGet(`Bearer ${token}`, (error, data) => {
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
    return new Promise((resolve, reject) => {
      this.feedbackApi.articlesIdFeedbackPost(articleId, feedbackData, `Bearer ${token}`, (error, data) => {
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
