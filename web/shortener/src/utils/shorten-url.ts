import type { ShortenResponse } from '../types/shorten'

export async function shortenUrl(url: string): Promise<string> {
  const response = await fetch('/api/shorten', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ url }),
  })

  if (!response.ok) {
    throw new Error('shorten failed')
  }

  const body = (await response.json()) as ShortenResponse
  if (!body.result) {
    throw new Error('shorten failed')
  }

  return body.result
}
