const mode = import.meta.env.MODE;

const ENVIRONMENTS = {
  development: {
    ACCOUNT_BASE_URL: "http://localhost:5173", // thanks to vite proxy
    ARTICLE_BASE_URL: "http://localhost:8081", // thanks to vite proxy
  },
  production: {
    ACCOUNT_BASE_URL: "https://account.newsbro.cc",
    ARTICLE_BASE_URL: "https://article.newsbro.cc",
  },
};

export const ENV = ENVIRONMENTS[mode] || ENVIRONMENTS.production;
