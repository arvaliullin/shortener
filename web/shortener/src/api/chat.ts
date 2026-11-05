import { api } from './client'
import type { TChatMessage } from '@gravity-ui/aikit'

export type ChatTurn = {
  role: 'user' | 'assistant'
  content: string
}

export async function sendChat(messages: ChatTurn[]): Promise<string> {
  const response = await api.post<{ content: string }>('/chat', { messages }, { timeout: 180000 })
  if (!response.data.content) {
    throw new Error('empty chat response')
  }
  return response.data.content
}

export function chatTurns(messages: TChatMessage[]): ChatTurn[] {
  return messages.flatMap((message) => {
    if ((message.role !== 'user' && message.role !== 'assistant') || typeof message.content !== 'string') {
      return []
    }
    return [{ role: message.role, content: message.content }]
  })
}
