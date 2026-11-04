import { createAsyncThunk, createSlice } from '@reduxjs/toolkit'
import { deleteUrl, fetchUrls, updateUrl } from '../../api/urls'
import {
  LINK_DELETE_FAILED_MESSAGE,
  LINK_SAVE_FAILED_MESSAGE,
  LINKS_LOAD_FAILED_MESSAGE,
} from '../../consts'
import type { UrlItem } from '../../types/shorten'

type LinksState = {
  items: UrlItem[]
  loading: boolean
  error: string
}

const initialState: LinksState = {
  items: [],
  loading: false,
  error: '',
}

export const loadLinks = createAsyncThunk('links/load', () => fetchUrls())

export const saveLink = createAsyncThunk(
  'links/save',
  ({ id, url }: { id: string; url: string }) => updateUrl(id, url),
)

export const removeLink = createAsyncThunk('links/remove', async (id: string) => {
  await deleteUrl(id)
  return id
})

const linksSlice = createSlice({
  name: 'links',
  initialState,
  reducers: {
    clearLinksError(state) {
      state.error = ''
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(loadLinks.pending, (state) => {
        state.loading = true
        state.error = ''
      })
      .addCase(loadLinks.fulfilled, (state, action) => {
        state.loading = false
        state.items = action.payload
      })
      .addCase(loadLinks.rejected, (state) => {
        state.loading = false
        state.error = LINKS_LOAD_FAILED_MESSAGE
      })
      .addCase(saveLink.fulfilled, (state, action) => {
        state.error = ''
        state.items = state.items.map((item) =>
          item.id === action.payload.id ? action.payload : item,
        )
      })
      .addCase(saveLink.rejected, (state) => {
        state.error = LINK_SAVE_FAILED_MESSAGE
      })
      .addCase(removeLink.fulfilled, (state, action) => {
        state.error = ''
        state.items = state.items.filter((item) => item.id !== action.payload)
      })
      .addCase(removeLink.rejected, (state) => {
        state.error = LINK_DELETE_FAILED_MESSAGE
      })
  },
})

export const { clearLinksError } = linksSlice.actions
export default linksSlice.reducer
