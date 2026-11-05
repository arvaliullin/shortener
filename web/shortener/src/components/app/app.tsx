import { useEffect, useState } from 'react'
import { FaceRobot, Link as LinkIcon, ListUl, Moon, Plus, Sun } from '@gravity-ui/icons'
import { AsideHeader } from '@gravity-ui/navigation'
import { ActionTooltip, Button, Icon, ThemeProvider, useThemeValue, type Theme } from '@gravity-ui/uikit'
import { Provider } from 'react-redux'
import { BrowserRouter, Route, Routes, useLocation, useNavigate } from 'react-router-dom'
import { ensureAuth } from '../../auth'
import ChatPage from '../../pages/chat-page/chat-page'
import LinksPage from '../../pages/links-page/links-page'
import ShortenPage from '../../pages/shorten-page/shorten-page'
import { store } from '../../store'
import './app.css'

type AppLayoutProps = {
  onThemeChange: (theme: Theme) => void
}

function AppLayout({ onThemeChange }: AppLayoutProps) {
  const navigate = useNavigate()
  const location = useLocation()
  const themeValue = useThemeValue()

  function switchTheme() {
    onThemeChange(themeValue === 'dark' ? 'light' : 'dark')
  }

  return (
    <AsideHeader
      className="app"
      compact
      hideCollapseButton
      menuDensity="compact"
      onChangeCompact={() => undefined}
      logo={{
        text: 'Сокращатель',
        icon: LinkIcon,
        href: '/',
        onClick: (event) => {
          event.preventDefault()
          navigate('/')
        },
      }}
      menuItems={[
        {
          id: 'create',
          title: 'Создать',
          icon: Plus,
          current: location.pathname === '/',
          enableTooltip: true,
          onItemClick: () => navigate('/'),
        },
        {
          id: 'links',
          title: 'Ссылки',
          icon: ListUl,
          current: location.pathname === '/links',
          enableTooltip: true,
          onItemClick: () => navigate('/links'),
        },
        {
          id: 'chat',
          title: 'Чат',
          icon: FaceRobot,
          current: location.pathname === '/chat',
          enableTooltip: true,
          onItemClick: () => navigate('/chat'),
        },
      ]}
      renderContent={() => (
        <div className="app__content">
          <Routes>
            <Route path="/" element={<ShortenPage />} />
            <Route path="/links" element={<LinksPage />} />
            <Route path="/chat" element={<ChatPage />} />
          </Routes>
        </div>
      )}
      renderFooter={() => {
        const themeLabel = themeValue === 'dark' ? 'Светлая тема' : 'Темная тема'

        return (
          <div className="app__theme">
            <ActionTooltip title={themeLabel}>
              <Button view="flat" size="l" aria-label={themeLabel} onClick={switchTheme}>
                <Button.Icon>
                  <Icon data={themeValue === 'dark' ? Sun : Moon} />
                </Button.Icon>
              </Button>
            </ActionTooltip>
          </div>
        )
      }}
    />
  )
}

export default function App() {
  const [theme, setTheme] = useState<Theme>('system')
  const [authenticated, setAuthenticated] = useState(false)

  useEffect(() => {
    ensureAuth().then((ok) => setAuthenticated(ok)).catch(() => setAuthenticated(false))
  }, [])

  if (!authenticated) {
    return null
  }

  return (
    <Provider store={store}>
      <BrowserRouter>
        <ThemeProvider theme={theme}>
          <AppLayout onThemeChange={setTheme} />
        </ThemeProvider>
      </BrowserRouter>
    </Provider>
  )
}
