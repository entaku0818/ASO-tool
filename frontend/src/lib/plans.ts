// Single source of truth for plan limits and displayed prices.
// Actual billing amounts are defined in Stripe (see STRIPE_PRICE_MONTHLY/
// YEARLY/LICENSE on the backend) — the values below are for UI display only
// and must be kept in sync with the corresponding Stripe Price objects.

export const PLAN_LIMITS = {
  freeAppLimit: 1,
  freeKeywordLimit: 10,
} as const

export const PRICES = {
  webProMonthly: 1980,
  webProYearlyMonthlyEquivalent: 1650,
  macosLicenseYearly: 9800,
} as const

export function formatYen(amount: number): string {
  return `¥${amount.toLocaleString('ja-JP')}`
}
