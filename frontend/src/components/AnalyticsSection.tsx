'use client'

import { useState } from 'react'
import { useAnalyticsWithCorrelation, useAnalyticsSummary } from '@/hooks/useAnalytics'
import { useASCCredentials } from '@/hooks/useASCCredentials'
import { ASCCredentialsForm } from './ASCCredentialsForm'
import { AnalyticsChart } from './AnalyticsChart'
import { triggerAnalyticsFetch } from '@/lib/api'

type Props = {
  appId: string
  platform: 'ios' | 'android'
}

export function AnalyticsSection({ appId, platform }: Props) {
  const { credentials, isLoading: credLoading, refetch: refetchCreds } = useASCCredentials(appId)
  const { analytics, isLoading: analyticsLoading, refetch: refetchAnalytics } = useAnalyticsWithCorrelation(appId, 30)
  const { summary, refetch: refetchSummary } = useAnalyticsSummary(appId, 30)
  const [showCredentialsForm, setShowCredentialsForm] = useState(false)
  const [isFetching, setIsFetching] = useState(false)

  // Only show for iOS apps
  if (platform !== 'ios') {
    return null
  }

  const handleFetchAnalytics = async () => {
    setIsFetching(true)
    try {
      await triggerAnalyticsFetch(appId)
      refetchAnalytics()
      refetchSummary()
    } catch {
      alert('Failed to fetch analytics data')
    } finally {
      setIsFetching(false)
    }
  }

  if (credLoading) {
    return (
      <div className="bg-white rounded-lg shadow p-6">
        <div className="animate-pulse h-8 bg-gray-200 rounded w-1/3 mb-4"></div>
        <div className="animate-pulse h-32 bg-gray-200 rounded"></div>
      </div>
    )
  }

  return (
    <div className="bg-white rounded-lg shadow">
      <div className="p-4 border-b flex items-center justify-between">
        <h3 className="text-lg font-semibold">App Store Analytics</h3>
        <div className="flex gap-2">
          <button
            onClick={() => setShowCredentialsForm(!showCredentialsForm)}
            className="px-3 py-1 text-sm border rounded hover:bg-gray-50"
          >
            {credentials ? 'API Settings' : 'Configure API'}
          </button>
          {credentials?.is_valid && (
            <button
              onClick={handleFetchAnalytics}
              disabled={isFetching}
              className="px-3 py-1 text-sm bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50"
            >
              {isFetching ? 'Fetching...' : 'Refresh Data'}
            </button>
          )}
        </div>
      </div>

      {showCredentialsForm && (
        <div className="p-4 border-b bg-gray-50">
          <ASCCredentialsForm
            appId={appId}
            existingCredentials={credentials}
            onSuccess={() => {
              refetchCreds()
              setShowCredentialsForm(false)
            }}
          />
        </div>
      )}

      {!credentials?.is_valid && !showCredentialsForm && (
        <div className="p-4 bg-yellow-50 border-l-4 border-yellow-400">
          <p className="text-sm text-yellow-700">
            App Store Connect API credentials need to be configured to view analytics.
          </p>
        </div>
      )}

      {credentials?.is_valid && (
        <>
          {/* Summary cards */}
          {summary && (
            <div className="p-4 grid grid-cols-2 md:grid-cols-4 gap-4">
              <div className="bg-blue-50 p-4 rounded">
                <p className="text-sm text-gray-600">Total Impressions</p>
                <p className="text-2xl font-bold">{summary.total_impressions.toLocaleString()}</p>
                <p className={`text-sm ${summary.impressions_change_percent >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                  {summary.impressions_change_percent >= 0 ? '+' : ''}{summary.impressions_change_percent.toFixed(1)}%
                </p>
              </div>
              <div className="bg-green-50 p-4 rounded">
                <p className="text-sm text-gray-600">Total Downloads</p>
                <p className="text-2xl font-bold">{summary.total_downloads.toLocaleString()}</p>
                <p className={`text-sm ${summary.downloads_change_percent >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                  {summary.downloads_change_percent >= 0 ? '+' : ''}{summary.downloads_change_percent.toFixed(1)}%
                </p>
              </div>
              <div className="bg-purple-50 p-4 rounded">
                <p className="text-sm text-gray-600">Conversion Rate</p>
                <p className="text-2xl font-bold">{(summary.average_conversion * 100).toFixed(2)}%</p>
              </div>
              <div className="bg-gray-50 p-4 rounded">
                <p className="text-sm text-gray-600">Period</p>
                <p className="text-2xl font-bold">30 days</p>
              </div>
            </div>
          )}

          {/* Chart */}
          {analyticsLoading ? (
            <div className="p-4">
              <div className="animate-pulse h-64 bg-gray-200 rounded"></div>
            </div>
          ) : analytics.length > 0 ? (
            <div className="p-4">
              <AnalyticsChart analytics={analytics} />
            </div>
          ) : (
            <div className="p-4 text-center text-gray-500">
              No analytics data available. Click &quot;Refresh Data&quot; to fetch.
            </div>
          )}

          {/* Correlation insights */}
          {analytics.length > 0 && (
            <div className="p-4 border-t">
              <h4 className="font-semibold mb-2">Correlation Insights</h4>
              <ul className="space-y-1 text-sm">
                {analytics
                  .filter(a => a.new_version || a.ranking_change)
                  .slice(0, 5)
                  .map((a, i) => (
                    <li key={i} className="flex items-center gap-2">
                      <span className="text-gray-500">{new Date(a.date).toLocaleDateString()}</span>
                      {a.new_version && (
                        <span className="bg-blue-100 text-blue-800 px-2 py-0.5 rounded text-xs">
                          v{a.new_version} released
                        </span>
                      )}
                      {a.ranking_change && a.ranking_change > 0 && (
                        <span className="bg-green-100 text-green-800 px-2 py-0.5 rounded text-xs">
                          Ranking improved by {a.ranking_change}
                        </span>
                      )}
                      {a.ranking_change && a.ranking_change < 0 && (
                        <span className="bg-red-100 text-red-800 px-2 py-0.5 rounded text-xs">
                          Ranking dropped by {Math.abs(a.ranking_change)}
                        </span>
                      )}
                    </li>
                  ))}
                {analytics.filter(a => a.new_version || a.ranking_change).length === 0 && (
                  <li className="text-gray-500">No significant events in this period</li>
                )}
              </ul>
            </div>
          )}
        </>
      )}
    </div>
  )
}
