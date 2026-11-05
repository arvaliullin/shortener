import { useRef, useState } from 'react'
import { ChatContainer, type ChatStatus, type TChatMessage, type TSubmitData } from '@gravity-ui/aikit'
import { chatTurns, isChatCanceled, sendChat } from '../../api/chat'
import './chat-page.css'

// Страница чата. Пока агент отвечает, список показывает ожидание.
export default function ChatPage() {
  const [messages, setMessages] = useState<TChatMessage[]>([])
  const [status, setStatus] = useState<ChatStatus>('ready')
  const generation = useRef(0)
  const request = useRef<AbortController | null>(null)

  // stopRequest прерывает текущий запрос к агенту.
  function stopRequest() {
    request.current?.abort()
    request.current = null
  }

  // requestReply запрашивает ответ по уже собранной истории.
  async function requestReply(history: TChatMessage[]) {
    const currentGeneration = generation.current
    stopRequest()
    const controller = new AbortController()
    request.current = controller
    setStatus('streaming_loading')

    try {
      const content = await sendChat(chatTurns(history), controller.signal)
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
      setStatus('ready')
    } catch (error) {
      if (currentGeneration !== generation.current || isChatCanceled(error)) {
        return
      }
      setStatus('error')
    } finally {
      if (request.current === controller) {
        request.current = null
      }
    }
  }

  // sendMessage добавляет реплику пользователя и ждет ответ.
  async function sendMessage(data: TSubmitData) {
    const history = [
      ...messages,
      {
        id: crypto.randomUUID(),
        role: 'user' as const,
        content: data.content,
        timestamp: new Date().toISOString(),
      },
    ]
    setMessages(history)
    await requestReply(history)
  }

  // retry повторяет последний запрос без новой реплики пользователя.
  async function retry() {
    await requestReply(messages)
  }

  // cancel останавливает ожидание и возвращает поле ввода.
  async function cancel() {
    stopRequest()
    setStatus('ready')
  }

  // createChat сбрасывает диалог и текущий запрос.
  function createChat() {
    generation.current += 1
    stopRequest()
    setMessages([])
    setStatus('ready')
  }

  return (
    <main className="chat-page">
      <ChatContainer
        messages={messages}
        status={status}
        onSendMessage={sendMessage}
        onCancel={cancel}
        onRetry={() => {
          void retry()
        }}
        onCreateChat={createChat}
        showHistory={false}
        showNewChat
        messageListConfig={{
          loaderMessage: 'Готовлю ответ. Это может занять минуту.',
        }}
        texts={{
          headerTitle: 'Чат',
          emptyStateTitle: 'Чат',
          emptyStateDescription: 'Спросите про ссылки или попросите создать короткую',
          promptPlaceholder: 'Сообщение',
          submitButtonCancelableText: 'Остановить',
          submitCancelTooltip: 'Остановить',
          errorText: 'Не удалось получить ответ. Попробуйте ещё раз.',
        }}
      />
    </main>
  )
}
