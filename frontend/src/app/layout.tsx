import type { Metadata } from 'next'
import './globals.css'
import { Providers } from '@/components/Providers'
import { Header } from '@/components/Header'

export const metadata: Metadata = {
  title: 'ASO Compass',
  description: 'App Store Optimization Tool',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="ja" suppressHydrationWarning>
      <head>
        <script dangerouslySetInnerHTML={{ __html: `
          (function() {
            var stored = localStorage.getItem('theme');
            var preferred = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
            if ((stored || preferred) === 'dark') document.documentElement.classList.add('dark');
          })();
        ` }} />
      </head>
      <body className="bg-gray-50 dark:bg-[#0d0f14] min-h-screen transition-colors">
        <Providers>
          <Header />
          <main className="max-w-7xl mx-auto px-4 py-6">
            {children}
          </main>
        </Providers>
      </body>
    </html>
  )
}
