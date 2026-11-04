import { useEffect, useState } from 'react'
import { ListUl } from '@gravity-ui/icons'
import {
  Alert,
  Button,
  Card,
  ClipboardButton,
  Dialog,
  Flex,
  Icon,
  Link,
  PlaceholderContainer,
  Skeleton,
  Text,
  TextInput,
} from '@gravity-ui/uikit'
import { useNavigate } from 'react-router-dom'
import { COPY_FAILED_MESSAGE, INVALID_URL_MESSAGE, LINK_SAVE_FAILED_MESSAGE } from '../../consts'
import { clearLinksError, loadLinks, removeLink, saveLink } from '../../store/slices/links-slice'
import { useAppDispatch, useAppSelector } from '../../store/hooks'
import type { UrlItem } from '../../types/shorten'
import { httpUrl } from '../../utils/http-url'
import './links-page.css'

export default function LinksPage() {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const items = useAppSelector((state) => state.links.items)
  const loading = useAppSelector((state) => state.links.loading)
  const error = useAppSelector((state) => state.links.error)
  const [editing, setEditing] = useState<UrlItem | null>(null)
  const [draft, setDraft] = useState('')
  const [draftError, setDraftError] = useState('')
  const [saving, setSaving] = useState(false)
  const [removing, setRemoving] = useState<UrlItem | null>(null)
  const [deleting, setDeleting] = useState(false)
  const [copyError, setCopyError] = useState('')

  useEffect(() => {
    dispatch(loadLinks())
  }, [dispatch])

  function openEdit(item: UrlItem) {
    setEditing(item)
    setDraft(item.url)
    setDraftError('')
  }

  async function confirmEdit() {
    if (!editing) {
      return
    }
    const original = httpUrl(draft)
    if (!original) {
      setDraftError(INVALID_URL_MESSAGE)
      return
    }

    setSaving(true)
    setDraftError('')
    try {
      await dispatch(saveLink({ id: editing.id, url: original })).unwrap()
      setEditing(null)
    } catch {
      setDraftError(LINK_SAVE_FAILED_MESSAGE)
    } finally {
      setSaving(false)
    }
  }

  async function confirmDelete() {
    if (!removing) {
      return
    }

    setDeleting(true)
    try {
      await dispatch(removeLink(removing.id)).unwrap()
      setRemoving(null)
    } catch {
      setRemoving(null)
    } finally {
      setDeleting(false)
    }
  }

  const pageError = copyError || error

  return (
    <Flex as="main" className="links-page" direction="column" gap={5} spacing={{ p: 6 }}>
      <Flex direction="column" gap={1}>
        <Text variant="header-1">Ссылки</Text>
        <Text variant="body-1" color="secondary">
          Сохраненные адреса. Короткий идентификатор при правке не меняется.
        </Text>
      </Flex>

      {pageError && !editing ? (
        <Alert
          theme="danger"
          title="Ошибка"
          message={pageError}
          onClose={() => {
            setCopyError('')
            dispatch(clearLinksError())
          }}
        />
      ) : null}

      {loading && items.length === 0 ? (
        <Card type="container" view="outlined">
          <Flex direction="column" gap={3} spacing={{ p: 5 }}>
            <Skeleton style={{ height: 'var(--g-spacing-5)', width: '72%' }} />
            <Skeleton style={{ height: 'var(--g-spacing-5)', width: '48%' }} />
            <Skeleton style={{ height: 'var(--g-spacing-5)', width: '64%' }} />
          </Flex>
        </Card>
      ) : null}

      {!loading && items.length === 0 ? (
        <Card type="container" view="outlined">
          <PlaceholderContainer
            direction="column"
            align="center"
            size="s"
            title="Пока нет сохраненных ссылок"
            description="Сократите адрес, и он появится в этом списке."
            image={<Icon data={ListUl} size={48} color="hint" />}
            actions={[
              {
                text: 'Создать',
                view: 'action',
                size: 'l',
                onClick: () => navigate('/'),
              },
            ]}
          />
        </Card>
      ) : null}

      {items.length > 0 ? (
        <Flex direction="column" gap={3}>
          {items.map((item) => (
            <Card key={item.id} type="container" view="outlined" spacing={{ p: 4 }}>
              <Flex direction="column" gap={4}>
                <Flex direction="column" gap={1}>
                  <Text variant="body-2" color="secondary">
                    Исходный адрес
                  </Text>
                  <Text className="links-page__url">{item.url}</Text>
                </Flex>
                <Flex direction="column" gap={1}>
                  <Text variant="body-2" color="secondary">
                    Короткая ссылка
                  </Text>
                  <Flex alignItems="center" gap={1} minWidth={0}>
                    <Link className="links-page__url" href={item.result} target="_blank" rel="noreferrer">
                      {item.result}
                    </Link>
                    <Flex shrink={0}>
                      <ClipboardButton
                        text={item.result}
                        view="flat"
                        size="s"
                        onCopy={(_text, copied) => {
                          setCopyError(copied ? '' : COPY_FAILED_MESSAGE)
                        }}
                      />
                    </Flex>
                  </Flex>
                </Flex>
                <Flex gap={2}>
                  <Button view="flat" onClick={() => openEdit(item)}>
                    Изменить
                  </Button>
                  <Button view="flat-danger" onClick={() => setRemoving(item)}>
                    Удалить
                  </Button>
                </Flex>
              </Flex>
            </Card>
          ))}
        </Flex>
      ) : null}

      <Dialog size="m" open={editing !== null} onClose={() => setEditing(null)}>
        <Dialog.Header caption="Изменить ссылку" />
        <Dialog.Body>
          <Flex direction="column" gap={2}>
            <Text as="label" variant="body-2" htmlFor="edit-url">
              Исходный адрес
            </Text>
            <TextInput
              id="edit-url"
              size="l"
              value={draft}
              onUpdate={setDraft}
              validationState={draftError ? 'invalid' : undefined}
              errorMessage={draftError || undefined}
            />
          </Flex>
        </Dialog.Body>
        <Dialog.Footer
          textButtonApply="Сохранить"
          textButtonCancel="Отмена"
          loading={saving}
          onClickButtonCancel={() => setEditing(null)}
          onClickButtonApply={confirmEdit}
        />
      </Dialog>

      <Dialog size="s" open={removing !== null} onClose={() => setRemoving(null)}>
        <Dialog.Header caption="Удалить ссылку" />
        <Dialog.Body>
          <Flex direction="column" gap={2}>
            <Text>Короткая ссылка перестанет открываться.</Text>
            {removing ? (
              <Text variant="body-2" color="secondary" className="links-page__url">
                {removing.result}
              </Text>
            ) : null}
          </Flex>
        </Dialog.Body>
        <Dialog.Footer
          textButtonApply="Удалить"
          textButtonCancel="Отмена"
          loading={deleting}
          onClickButtonCancel={() => setRemoving(null)}
          onClickButtonApply={confirmDelete}
          propsButtonApply={{ view: 'outlined-danger' }}
        />
      </Dialog>
    </Flex>
  )
}
