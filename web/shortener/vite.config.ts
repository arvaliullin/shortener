import fs from 'node:fs'
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig(({ command }) => ({
  plugins: [react()],
  server: {
    https: command === 'serve' ? {
      cert: fs.readFileSync(new URL('../../certs/localhost.pem', import.meta.url)),
      key: fs.readFileSync(new URL('../../certs/localhost-key.pem', import.meta.url)),
    } : undefined,
    proxy: {
      // /api/chat идет в агента и должен стоять раньше общего /api.
      '/api/chat': {
        target: 'http://localhost:8090',
        timeout: 180000,
        proxyTimeout: 180000,
      },
      '/api': 'http://localhost:8080',
      '^/[A-Za-z0-9]{8}$': 'http://localhost:8080',
    },
  },
}))
