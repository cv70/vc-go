import { FormEvent, useMemo, useState } from 'react'
import './App.css'
import type { Founder } from './api/client/models/Founder'
import { ApiError } from './api/client/core/ApiError'
import { FounderService } from './api/client/services/FounderService'
import { handleApiErrorStatus } from './api/http'

function formatMoney(value: number | null | undefined): string {
  if (value == null) return '未披露'
  return new Intl.NumberFormat('zh-CN').format(value)
}

function stageLabel(stage: string): string {
  const map: Record<string, string> = {
    idea: 'Idea',
    pre_seed: 'Pre-Seed',
    seed: 'Seed',
    pre_a: 'Pre-A',
    a: 'A',
  }
  return map[stage] ?? stage
}

function App() {
  const [keyword, setKeyword] = useState('')
  const [loading, setLoading] = useState(false)
  const [errorText, setErrorText] = useState('')
  const [founders, setFounders] = useState<Founder[]>([])
  const [hasSearched, setHasSearched] = useState(false)

  const resultCountText = useMemo(() => {
    if (!hasSearched) return '输入关键词开始搜索创业者'
    return `共 ${founders.length} 条结果`
  }, [founders.length, hasSearched])

  async function onSearch(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setLoading(true)
    setErrorText('')

    try {
      const response = await FounderService.searchFounders({
        keyword: keyword.trim() || null,
        limit: 20,
      })

      setFounders(response.data?.items ?? [])
      setHasSearched(true)
    } catch (error) {
      if (error instanceof ApiError) {
        handleApiErrorStatus(error.status)
        setErrorText(`请求失败（${error.status}）：${error.message}`)
      } else if (error instanceof Error) {
        setErrorText(error.message)
      } else {
        setErrorText('未知错误，请稍后重试')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="page">
      <section className="hero">
        <p className="eyebrow">Venture OS</p>
        <h1>Founder Search Console</h1>
        <p className="subtitle">统一 API Client 已接入，直接搜索 founder 数据并渲染结果。</p>
      </section>

      <section className="panel">
        <form className="searchRow" onSubmit={onSearch}>
          <input
            className="searchInput"
            placeholder="例如：AI 医疗 / Seed / 上海"
            value={keyword}
            onChange={(event) => setKeyword(event.target.value)}
          />
          <button className="searchButton" type="submit" disabled={loading}>
            {loading ? '搜索中...' : '搜索'}
          </button>
        </form>

        <p className="hint">{resultCountText}</p>
        {errorText ? <p className="errorText">{errorText}</p> : null}
      </section>

      <section className="grid">
        {founders.map((founder) => (
          <article className="card" key={founder.id}>
            <header className="cardHeader">
              <h2>{founder.name}</h2>
              <span className="badge">{stageLabel(founder.stage)}</span>
            </header>
            <p className="meta">赛道：{founder.sectors.join(' / ') || '未设置'}</p>
            <p className="meta">地区：{founder.location ?? '未设置'}</p>
            <p className="meta">融资目标：¥ {formatMoney(founder.funding_target_cny)}</p>
          </article>
        ))}
      </section>
    </main>
  )
}

export default App
