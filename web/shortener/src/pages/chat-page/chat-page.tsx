import { useRef, useState } from 'react'
import { ChatContainer, type TChatMessage, type TSubmitData } from '@gravity-ui/aikit'
import { chatTurns, sendChat } from '../../api/chat'
import './chat-page.css'

// Страница чата. Сообщение и история уходят агенту.
export default function ChatPage() {
  const [messages, setMessages] = useState<TChatMessage[]>([])
  const generation = useRef(0)

  // Отправляет историю и добавляет ответ агента.
  async function sendMessage(data: TSubmitData) {
    const currentGeneration = generation.current
    const timestamp = new Date().toISOString()
    const userMessage: TChatMessage = {
      id: crypto.randomUUID(),
      role: 'user',
      content: data.content,
      timestamp,
    }
    const history = [...messages, userMessage]
    setMessages(history)

    let content = 'Не удалось получить ответ'
    try {
      content = await sendChat(chatTurns(history))
    } catch {
      content = 'Не удалось получить ответ'
    }
    if (currentGeneration !== generation.current) {
      return
    }

    setMessages((current) => [
      ...current,
      {
        id: crypto.randomUUID(),
        role: 'assistant',
        content,
        timestamp: new Date().toISOString(),
      },
    ])
  }

  return (
    <main className="chat-page">
      <ChatContainer
        messages={messages}
        onSendMessage={sendMessage}
        onCreateChat={() => {
          generation.current += 1
          setMessages([])
        }}
        showHistory={false}
        showNewChat
        texts={{
          headerTitle: 'Чат',
          emptyStateTitle: 'Чат',
          emptyStateDescription: 'Спросите про ссылки или попросите создать короткую',
          promptPlaceholder: 'Сообщение',
        }}
      />
    </main>
  )
}
