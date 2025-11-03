const mode = import.meta.env.MODE;

const ENVIRONMENTS = {
  development: {
    ACCOUNT_BASE_URL: "http://localhost:5173", // thanks to vite proxy
    ARTICLE_BASE_URL: "http://localhost:5173", // thanks to vite proxy
    FEED_BASE_URL: "http://localhost:5173", // thanks to vite proxy
  },
  production: {
    ACCOUNT_BASE_URL: "https://account.newsbro.cc",
    ARTICLE_BASE_URL: "https://article.newsbro.cc",
    FEED_BASE_URL: "https://feed.newsbro.cc",
  },
};

export const ENV = ENVIRONMENTS[mode] || ENVIRONMENTS.production;
