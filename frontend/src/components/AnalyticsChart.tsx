'use client'

import { AnalyticsWithCorrelation } from '@/lib/api'

type Props = {
  analytics: AnalyticsWithCorrelation[]
}

export function AnalyticsChart({ analytics }: Props) {
  // Reverse to show oldest first
  const data = [...analytics].reverse()

  if (data.length === 0) {
    return (
      <div className="text-center text-gray-500 py-8">
        No analytics data available
      </div>
    )
  }

  // Calculate max values for scaling
  const maxImpressions = Math.max(...data.map(d => d.impressions), 1)
  const maxDownloads = Math.max(...data.map(d => d.downloads), 1)

  const chartWidth = 800
  const chartHeight = 200
  const padding = { top: 20, right: 60, bottom: 30, left: 60 }
  const plotWidth = chartWidth - padding.left - padding.right
  const plotHeight = chartHeight - padding.top - padding.bottom

  // Generate path for line chart
  const generatePath = (values: number[], maxValue: number) => {
    const points = values.map((v, i) => {
      const x = padding.left + (i / (values.length - 1)) * plotWidth
      const y = padding.top + plotHeight - (v / maxValue) * plotHeight
      return `${x},${y}`
    })
    return `M ${points.join(' L ')}`
  }

  const impressionsPath = generatePath(data.map(d => d.impressions), maxImpressions)
  const downloadsPath = generatePath(data.map(d => d.downloads), maxDownloads)

  // Find version release points
  const versionPoints = data
    .map((d, i) => d.new_version ? { index: i, version: d.new_version } : null)
    .filter((v): v is { index: number; version: string } => v !== null)

  return (
    <div className="overflow-x-auto">
      <svg width={chartWidth} height={chartHeight} className="mx-auto">
        {/* Grid lines */}
        {[0, 0.25, 0.5, 0.75, 1].map(ratio => (
          <line
            key={ratio}
            x1={padding.left}
            y1={padding.top + plotHeight * (1 - ratio)}
            x2={chartWidth - padding.right}
            y2={padding.top + plotHeight * (1 - ratio)}
            stroke="#e5e7eb"
            strokeDasharray="4"
          />
        ))}

        {/* Impressions line */}
        <path
          d={impressionsPath}
          fill="none"
          stroke="#3b82f6"
          strokeWidth="2"
        />

        {/* Downloads line */}
        <path
          d={downloadsPath}
          fill="none"
          stroke="#10b981"
          strokeWidth="2"
        />

        {/* Version release markers */}
        {versionPoints.map(({ index, version }) => {
          const x = padding.left + (index / (data.length - 1)) * plotWidth
          return (
            <g key={index}>
              <line
                x1={x}
                y1={padding.top}
                x2={x}
                y2={padding.top + plotHeight}
                stroke="#f59e0b"
                strokeWidth="2"
                strokeDasharray="4"
              />
              <text
                x={x}
                y={padding.top - 5}
                textAnchor="middle"
                fontSize="10"
                fill="#f59e0b"
              >
                v{version}
              </text>
            </g>
          )
        })}

        {/* Y-axis labels - Impressions */}
        <text
          x={padding.left - 10}
          y={padding.top}
          textAnchor="end"
          fontSize="10"
          fill="#3b82f6"
        >
          {maxImpressions.toLocaleString()}
        </text>
        <text
          x={padding.left - 10}
          y={padding.top + plotHeight}
          textAnchor="end"
          fontSize="10"
          fill="#3b82f6"
        >
          0
        </text>

        {/* Y-axis labels - Downloads */}
        <text
          x={chartWidth - padding.right + 10}
          y={padding.top}
          textAnchor="start"
          fontSize="10"
          fill="#10b981"
        >
          {maxDownloads.toLocaleString()}
        </text>
        <text
          x={chartWidth - padding.right + 10}
          y={padding.top + plotHeight}
          textAnchor="start"
          fontSize="10"
          fill="#10b981"
        >
          0
        </text>

        {/* X-axis labels */}
        {data.length > 0 && (
          <>
            <text
              x={padding.left}
              y={chartHeight - 5}
              textAnchor="start"
              fontSize="10"
              fill="#6b7280"
            >
              {new Date(data[0].date).toLocaleDateString()}
            </text>
            <text
              x={chartWidth - padding.right}
              y={chartHeight - 5}
              textAnchor="end"
              fontSize="10"
              fill="#6b7280"
            >
              {new Date(data[data.length - 1].date).toLocaleDateString()}
            </text>
          </>
        )}
      </svg>

      {/* Legend */}
      <div className="flex justify-center gap-6 mt-2 text-sm">
        <div className="flex items-center gap-1">
          <div className="w-3 h-3 bg-blue-500 rounded-full"></div>
          <span>Impressions</span>
        </div>
        <div className="flex items-center gap-1">
          <div className="w-3 h-3 bg-green-500 rounded-full"></div>
          <span>Downloads</span>
        </div>
        <div className="flex items-center gap-1">
          <div className="w-3 h-0.5 bg-yellow-500"></div>
          <span>Version Release</span>
        </div>
      </div>
    </div>
  )
}
