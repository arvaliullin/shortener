import { api } from './client'
import type { ShortenResponse, UrlItem } from '../types/shorten'

export async function shortenUrl(url: string): Promise<string> {
  const response = await api.post<ShortenResponse>('/shorten', { url })
  if (!response.data.result) {
    throw new Error('shorten failed')
  }
  return response.data.result
}

export async function fetchUrls(): Promise<UrlItem[]> {
  const response = await api.get<UrlItem[]>('/urls')
  return response.data
}

export async function updateUrl(id: string, url: string): Promise<UrlItem> {
  const response = await api.put<UrlItem>(`/urls/${encodeURIComponent(id)}`, { url })
  return response.data
}

export async function deleteUrl(id: string): Promise<void> {
  await api.delete(`/urls/${encodeURIComponent(id)}`)
}
