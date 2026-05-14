'use client'

import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts'
import { Ranking } from '@/lib/api'
import { formatDate } from '@/lib/date'

const COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899']

export type KeywordRankingData = {
  id: string
  keyword: string
  rankings: Ranking[]
}

type Props = {
  keywords: KeywordRankingData[]
}

export function RankingChart({ keywords }: Props) {
  const active = keywords.filter((k) => k.rankings.some((r) => r.rank !== null))

  if (active.length === 0) {
    return <p className="text-gray-500 text-sm">データがありません</p>
  }

  // Merge all dates across keywords, then build one row per date
  const dateMap = new Map<string, Record<string, number | null>>()

  for (const kw of active) {
    for (const r of kw.rankings) {
      const date = formatDate(r.recorded_at)
      if (!dateMap.has(date)) dateMap.set(date, {})
      dateMap.get(date)![kw.id] = r.rank
    }
  }

  const data = Array.from(dateMap.entries())
    .sort(([a], [b]) => a.localeCompare(b))
    .map(([date, ranks]) => ({ date, ...ranks }))

  return (
    <ResponsiveContainer width="100%" height={220}>
      <LineChart data={data} margin={{ top: 5, right: 20, bottom: 5, left: 0 }}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="date" tick={{ fontSize: 11 }} />
        <YAxis reversed domain={[1, 'auto']} tick={{ fontSize: 11 }} allowDecimals={false} />
        <Tooltip
          formatter={(value: number, name: string) => {
            const kw = active.find((k) => k.id === name)
            return [`${value}位`, kw?.keyword ?? name]
          }}
          labelFormatter={(label) => `日付: ${label}`}
        />
        {active.length > 1 && (
          <Legend
            formatter={(value) => active.find((k) => k.id === value)?.keyword ?? value}
          />
        )}
        {active.map((kw, i) => (
          <Line
            key={kw.id}
            type="monotone"
            dataKey={kw.id}
            stroke={COLORS[i % COLORS.length]}
            strokeWidth={2}
            dot={{ fill: COLORS[i % COLORS.length], r: 3 }}
            connectNulls
          />
        ))}
      </LineChart>
    </ResponsiveContainer>
  )
}
