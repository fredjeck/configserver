import Image from 'next/image'
import { Inter } from 'next/font/google'

const inter = Inter({ subsets: ['latin'] })

export default function Home() {
  return (
    <div className="isolate bg-white px-6 py-24 sm:py-32 lg:px-8">
      <div
        className="absolute inset-x-0 -z-10 transform-gpu overflow-hidden blur-3xl sm:top-[-10rem]"
        aria-hidden="true"
      >
        <div
          className="relative left-1/2 -z-10 aspect-[1155/678] w-[25rem] max-w-none bg-gradient-to-r from-green-300 via-blue-500 to-purple-600 opacity-30"
        />
      </div>
      <div className="mx-auto max-w-2xl text-center">
        <h2 className="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">Register Client</h2>
        <p className="mt-2 text-lg leading-8 text-gray-600">
          Use this form to register a new Client ID/Client Secret pair.<br/>
          Don't forget to register the new Client ID in your configuration.<br/>
          If Client ID is ommitted, a new random identifier will be genreated.
        </p>
      </div>
      <form action="#" method="POST" className="mx-auto mt-16 max-w-xl sm:mt-20">
        <div className="sm:col-span-2">
          <div>
            <label htmlFor="clientid" className="block text-sm font-semibold leading-6 text-gray-900">
              Client Id
            </label>
            <div className="mt-2.5">
              <input
                type="text"
                name="clientid"
                id="clientid"
                autoComplete="off"
                className="block w-full rounded-md border-0 px-3.5 py-2 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
              />
            </div>
          </div>
        </div>
        <div className="mt-10">
          <button
            type="submit"
            className="block w-full rounded-md bg-indigo-600 px-3.5 py-2.5 text-center text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
          >
            Generate
          </button>
        </div>
      </form>
    </div>
  )
}

