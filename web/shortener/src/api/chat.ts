import axios from 'axios'
import { api } from './client'
import type { TChatMessage } from '@gravity-ui/aikit'

export type ChatTurn = {
  role: 'user' | 'assistant'
  content: string
}

// sendChat отправляет историю и возвращает текст ответа агента.
export async function sendChat(messages: ChatTurn[], signal?: AbortSignal): Promise<string> {
  const response = await api.post<{ content: string }>('/chat', { messages }, { timeout: 180000, signal })
  if (!response.data.content) {
    throw new Error('empty chat response')
  }
  return response.data.content
}

// isChatCanceled отличает отмену запроса от сбоя агента.
export function isChatCanceled(error: unknown) {
  return axios.isCancel(error)
}

// chatTurns оставляет только текстовые реплики пользователя и ассистента.
export function chatTurns(messages: TChatMessage[]): ChatTurn[] {
  return messages.flatMap((message) => {
    if ((message.role !== 'user' && message.role !== 'assistant') || typeof message.content !== 'string') {
      return []
    }
    return [{ role: message.role, content: message.content }]
  })
}
