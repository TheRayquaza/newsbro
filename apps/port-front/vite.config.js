import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import fs from 'fs';
import path from 'path';

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  optimizeDeps: {
    include: ['superagent'],
  },
  /*
  server: {
    https: {
      key: fs.readFileSync(path.resolve(__dirname, 'localhost-key.pem')),
      cert: fs.readFileSync(path.resolve(__dirname, 'localhost-cert.pem')),
    },
    proxy: {
      '/api/v1/users': {
        target: 'https://account.newsbro.cc',
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
        target: 'https://account.newsbro.cc',
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
      '/api/v1/article': {
        target: 'https://article.newsbro.cc',
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
        target: 'https://article.newsbro.cc',
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
    */
  //},
});
