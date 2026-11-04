import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { configure } from '@gravity-ui/uikit'
import App from './components/app/app'
import '@gravity-ui/uikit/styles/fonts.css'
import '@gravity-ui/uikit/styles/styles.css'
import './styles/theme.css'
import './index.css'

configure({ lang: 'ru' })

const root = document.getElementById('root')
if (!root) {
  throw new Error('root element is missing')
}

createRoot(root).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
