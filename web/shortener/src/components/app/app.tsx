import { useState } from 'react'
import { ThemeProvider, type Theme } from '@gravity-ui/uikit'
import ShortenPage from '../../pages/shorten-page/shorten-page'
import './app.css'

export default function App() {
  const [theme, setTheme] = useState<Theme>('system')

  return (
    <ThemeProvider theme={theme}>
      <div className="app">
        <ShortenPage onThemeChange={setTheme} />
      </div>
    </ThemeProvider>
  )
}
