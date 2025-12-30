'use client'

import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'
import { Ranking } from '@/lib/api'
import { formatDate } from '@/lib/date'

type Props = {
  rankings: Ranking[]
}

export function RankingChart({ rankings }: Props) {
  if (rankings.length === 0) {
    return <p className="text-gray-500 text-sm">データがありません</p>
  }

  const data = rankings
    .filter((r) => r.rank !== null)
    .map((r) => ({
      date: formatDate(r.recorded_at),
      rank: r.rank,
    }))
    .reverse()

  if (data.length === 0) {
    return <p className="text-gray-500 text-sm">順位データがありません</p>
  }

  return (
    <ResponsiveContainer width="100%" height={200}>
      <LineChart data={data} margin={{ top: 5, right: 20, bottom: 5, left: 0 }}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="date" tick={{ fontSize: 12 }} />
        <YAxis reversed domain={[1, 100]} tick={{ fontSize: 12 }} allowDecimals={false} />
        <Tooltip
          formatter={(value: number) => [`${value}位`, '順位']}
          labelFormatter={(label) => `日付: ${label}`}
        />
        <Line
          type="monotone"
          dataKey="rank"
          stroke="#3b82f6"
          strokeWidth={2}
          dot={{ fill: '#3b82f6', strokeWidth: 2, r: 3 }}
        />
      </LineChart>
    </ResponsiveContainer>
  )
}
