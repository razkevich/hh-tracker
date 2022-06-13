import {SelfLink} from '../common/http-data-structures'

export interface LogResponseBody {
  data: LogResponseData[]
  meta: LogMeta
  links: SelfLink
  relationships?: LogRel
}

export interface LogResponseData {
  id?: string
  links: SelfLink
  store_id: string
  user: LogResponseUserData
  time: string
  event_type: string
  delta: object
  type: string
}

export interface LogMeta {
  timestamps: {
    created_at: string
    updated_at: string
  }
}

export interface LogResponseUserData {
  id: string
  type: string
  role: string
  name: string
}

export interface LogRel {
  resource_path: {
    url: string
  }
}

export interface LogReadParams {
  'store_id'?: string
  'page-limit'?: string
  'page-offset'?: string
  'X-Moltin-Settings-page_length'?: string
  'ep-internal-search-json'?: string
}
