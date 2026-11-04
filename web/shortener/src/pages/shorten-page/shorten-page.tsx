import { useState, type FormEvent } from 'react'
import { Copy, CopyCheck, Moon, Sun } from '@gravity-ui/icons'
import {
  Alert,
  Button,
  Card,
  Flex,
  Icon,
  Link,
  Text,
  TextInput,
  useThemeValue,
  type Theme,
} from '@gravity-ui/uikit'
import { COPY_FAILED_MESSAGE, INVALID_URL_MESSAGE, SHORTEN_FAILED_MESSAGE } from '../../consts'
import { httpUrl } from '../../utils/http-url'
import { shortenUrl } from '../../utils/shorten-url'
import './shorten-page.css'

type ShortenPageProps = {
  onThemeChange: (theme: Theme) => void
}

export default function ShortenPage({ onThemeChange }: ShortenPageProps) {
  const themeValue = useThemeValue()
  const [url, setUrl] = useState('')
  const [result, setResult] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [copied, setCopied] = useState(false)

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError('')
    setCopied(false)

    const original = httpUrl(url)
    if (!original) {
      setResult('')
      setError(INVALID_URL_MESSAGE)
      return
    }

    setLoading(true)
    try {
      const shortened = await shortenUrl(original)
      setResult(shortened)
    } catch {
      setResult('')
      setError(SHORTEN_FAILED_MESSAGE)
    } finally {
      setLoading(false)
    }
  }

  async function copyResult() {
    try {
      await navigator.clipboard.writeText(result)
      setCopied(true)
      setError('')
    } catch {
      setCopied(false)
      setError(COPY_FAILED_MESSAGE)
    }
  }

  function switchTheme() {
    onThemeChange(themeValue === 'dark' ? 'light' : 'dark')
  }

  return (
    <main className="shorten-page">
      <Card className="shorten-page__panel" type="container" view="raised" size="l">
        <Flex direction="column" gap="spacing-5">
          <Flex alignItems="flex-start" justifyContent="space-between" gap="spacing-3" wrap>
            <Flex direction="column" gap="spacing-1">
              <Text variant="header-1">Сокращатель ссылок</Text>
              <Text variant="body-1" color="secondary">
                Вставьте длинный адрес и получите короткую ссылку
              </Text>
            </Flex>
            <Button
              view="flat"
              size="l"
              aria-label={themeValue === 'dark' ? 'Светлая тема' : 'Темная тема'}
              onClick={switchTheme}
            >
              <Button.Icon>
                <Icon data={themeValue === 'dark' ? Sun : Moon} />
              </Button.Icon>
            </Button>
          </Flex>

          <form noValidate onSubmit={submit}>
            <Flex direction="column" gap="spacing-3">
              <TextInput
                size="l"
                type="url"
                placeholder="https://example.com/very/long/path"
                value={url}
                onUpdate={setUrl}
                hasClear
                disabled={loading}
                validationState={error && !result ? 'invalid' : undefined}
                errorMessage={error && !result ? error : undefined}
              />
              <Button view="action" size="l" type="submit" loading={loading} width="max">
                Сократить
              </Button>
            </Flex>
          </form>

          {result ? (
            <Flex direction="column" gap="spacing-3">
              <Text variant="subheader-2">Короткая ссылка</Text>
              <Link className="shorten-page__link" href={result} target="_blank" rel="noreferrer">
                {result}
              </Link>
              <Button view="outlined" size="l" onClick={copyResult}>
                <Button.Icon>
                  <Icon data={copied ? CopyCheck : Copy} />
                </Button.Icon>
                {copied ? 'Скопировано' : 'Скопировать'}
              </Button>
            </Flex>
          ) : null}

          {error && result ? <Alert theme="danger" title="Ошибка" message={error} /> : null}
        </Flex>
      </Card>
    </main>
  )
}
