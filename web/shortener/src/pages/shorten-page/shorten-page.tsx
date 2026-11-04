import { useState, type FormEvent } from 'react'
import { Alert, Button, Card, ClipboardButton, Flex, Link, Text, TextInput } from '@gravity-ui/uikit'
import { shortenUrl } from '../../api/urls'
import { COPY_FAILED_MESSAGE, INVALID_URL_MESSAGE, SHORTEN_FAILED_MESSAGE } from '../../consts'
import { httpUrl } from '../../utils/http-url'
import './shorten-page.css'

export default function ShortenPage() {
  const [url, setUrl] = useState('')
  const [result, setResult] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError('')

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

  return (
    <main className="shorten-page">
      <Card className="shorten-page__panel" type="container" view="raised" size="l" spacing={{ p: 6 }}>
        <Flex direction="column" gap={6}>
          <Flex direction="column" gap={1}>
            <Text variant="header-1">Сокращатель ссылок</Text>
            <Text variant="body-1" color="secondary">
              Вставьте длинный адрес и получите короткую ссылку
            </Text>
          </Flex>

          <form noValidate onSubmit={submit}>
            <Flex direction="column" gap={4}>
              <Flex direction="column" gap={2}>
                <Text as="label" variant="body-2" htmlFor="shorten-url">
                  Длинный адрес
                </Text>
                <TextInput
                  id="shorten-url"
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
              </Flex>
              <Button view="action" size="l" type="submit" loading={loading} width="max">
                Сократить
              </Button>
            </Flex>
          </form>

          {result ? (
            <Flex direction="column" gap={3}>
              <Text variant="subheader-2">Короткая ссылка</Text>
              <Flex alignItems="center" gap={2}>
                <Flex grow={1} minWidth={0}>
                  <Link className="shorten-page__link" href={result} target="_blank" rel="noreferrer">
                    {result}
                  </Link>
                </Flex>
                <Flex shrink={0}>
                  <ClipboardButton
                    text={result}
                    view="outlined"
                    size="l"
                    onCopy={(_text, copied) => {
                      setError(copied ? '' : COPY_FAILED_MESSAGE)
                    }}
                  >
                    Скопировать
                  </ClipboardButton>
                </Flex>
              </Flex>
            </Flex>
          ) : null}

          {error && result ? (
            <Alert theme="danger" title="Ошибка" message={error} onClose={() => setError('')} />
          ) : null}
        </Flex>
      </Card>
    </main>
  )
}
