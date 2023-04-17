import Image from 'next/image'
import { Inter } from 'next/font/google'
import Repositories from '@/components/registration'

const inter = Inter({ subsets: ['latin'] })

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
      <div className="mb-32 grid text-center lg:mb-0 lg:grid-cols-4 lg:text-left">
        <Repositories></Repositories>
      </div>
    </main>
  )
}
