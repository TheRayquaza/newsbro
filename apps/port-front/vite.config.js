import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  optimizeDeps: {
    include: ['superagent'],
  },
  server: {
    proxy: {
      '/api/v1/users': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: true,
        cookieDomainRewrite: 'localhost',
        configure: (proxy) => {
          proxy.on('proxyRes', (proxyRes) => {
            const cookies = proxyRes.headers['set-cookie'];
            if (cookies) {
              proxyRes.headers['set-cookie'] = cookies.map((cookie) =>
                cookie
                  .replace(/Domain=[^;]+/i, 'Domain=localhost')
                  .replace(/; *Secure/gi, '')
              );
            }
          });
        },
      },
      '/api/v1/auth': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: true,
        cookieDomainRewrite: 'localhost',
        configure: (proxy) => {
          proxy.on('proxyRes', (proxyRes) => {
            const cookies = proxyRes.headers['set-cookie'];
            if (cookies) {
              proxyRes.headers['set-cookie'] = cookies.map((cookie) =>
                cookie
                  .replace(/Domain=[^;]+/i, 'Domain=localhost')
                  .replace(/; *Secure/gi, '')
              );
            }
          });
        },
      },
      '/api/v1/articles': {
        target: 'http://localhost:8081',
        changeOrigin: true,
        secure: true,
        cookieDomainRewrite: 'localhost',
        configure: (proxy) => {
          proxy.on('proxyRes', (proxyRes) => {
            const cookies = proxyRes.headers['set-cookie'];
            if (cookies) {
              proxyRes.headers['set-cookie'] = cookies.map((cookie) =>
                cookie
                  .replace(/Domain=[^;]+/i, 'Domain=localhost')
                  .replace(/; *Secure/gi, '')
              );
            }
          });
        },
      },
      '/api/v1/feedback': {
        target: 'http://localhost:8081',
        changeOrigin: true,
        secure: true,
        cookieDomainRewrite: 'localhost',
        configure: (proxy) => {
          proxy.on('proxyRes', (proxyRes) => {
            const cookies = proxyRes.headers['set-cookie'];
            if (cookies) {
              proxyRes.headers['set-cookie'] = cookies.map((cookie) =>
                cookie
                  .replace(/Domain=[^;]+/i, 'Domain=localhost')
                  .replace(/; *Secure/gi, '')
              );
            }
          });
        },
      },
      '/api/v1/rss': {
        target: 'http://localhost:8081',
        changeOrigin: true,
        secure: true,
        cookieDomainRewrite: 'localhost',
        configure: (proxy) => {
          proxy.on('proxyRes', (proxyRes) => {
            const cookies = proxyRes.headers['set-cookie'];
            if (cookies) {
              proxyRes.headers['set-cookie'] = cookies.map((cookie) =>
                cookie
                  .replace(/Domain=[^;]+/i, 'Domain=localhost')
                  .replace(/; *Secure/gi, '')
              );
            }
          });
        },
      },
    },
  }
});
