'use client'

import { useReviews } from '@/hooks/useReviews'
import { Review } from '@/lib/api'
import { formatDateTime } from '@/lib/date'

function StarRating({ rating }: { rating: number }) {
  return (
    <span className="text-yellow-500">
      {'★'.repeat(rating)}{'☆'.repeat(5 - rating)}
    </span>
  )
}

function ReviewCard({ review }: { review: Review }) {
  return (
    <div className="border-b py-4 last:border-b-0">
      <div className="flex items-center justify-between mb-2">
        <div className="flex items-center gap-2">
          <StarRating rating={review.rating} />
          <span className="font-medium">{review.author}</span>
        </div>
        <span className="text-sm text-gray-500 dark:text-gray-400">{formatDateTime(review.posted_at)}</span>
      </div>
      {review.title && <p className="font-medium mb-1">{review.title}</p>}
      <p className="text-gray-700 dark:text-gray-300 text-sm">{review.content}</p>
      {review.version && (
        <p className="text-xs text-gray-400 dark:text-gray-500 mt-2">v{review.version}</p>
      )}
    </div>
  )
}

export function ReviewsSection({ appId }: { appId: string }) {
  const { reviews, stats, isLoading } = useReviews(appId)

  if (isLoading) {
    return (
      <div className="bg-white dark:bg-[#12161e] rounded-lg shadow dark:shadow-black/20 p-4">
        <h3 className="text-lg font-semibold mb-4">レビュー</h3>
        <p className="text-gray-500 dark:text-gray-400">読み込み中...</p>
      </div>
    )
  }

  return (
    <div className="bg-white dark:bg-[#12161e] rounded-lg shadow dark:shadow-black/20">
      <div className="p-4 border-b">
        <h3 className="text-lg font-semibold">レビュー ({reviews.length}件)</h3>
        {stats && stats.rating_counts && (
          <div className="flex items-center gap-4 mt-2">
            <div className="flex items-center gap-1">
              <StarRating rating={Math.round(stats.average_rating || 0)} />
              <span className="font-bold">{(stats.average_rating || 0).toFixed(1)}</span>
            </div>
            <div className="text-sm text-gray-500 dark:text-gray-400">
              {[5, 4, 3, 2, 1].map((r) => (
                <span key={r} className="mr-2">
                  {r}★: {stats.rating_counts?.[r] || 0}
                </span>
              ))}
            </div>
          </div>
        )}
      </div>

      <div className="p-4">
        {reviews.length === 0 ? (
          <p className="text-gray-500 dark:text-gray-400">レビューがありません</p>
        ) : (
          <div className="max-h-96 overflow-y-auto">
            {reviews.slice(0, 20).map((review) => (
              <ReviewCard key={review.id} review={review} />
            ))}
            {reviews.length > 20 && (
              <p className="text-center text-gray-500 dark:text-gray-400 py-2">
                他 {reviews.length - 20} 件のレビュー
              </p>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
