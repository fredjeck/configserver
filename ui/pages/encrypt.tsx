import Alert from "@/components/alert";
import Layout from "@/components/layout";
import { Inter } from "next/font/google";
import { ReactEventHandler, useEffect, useState } from "react";

const inter = Inter({ subsets: ["latin"] });

interface Repository {
  name: string;
  selected: boolean;
}



export default function Register() {
  const [postError, setPostError] = useState<string>("");

  const [value, setValue] = useState<string>("");
  const [token, setToken] = useState<string>("");

  const isPostError = () => postError.length > 0;

  function generateToken() {
    const postValue = async () => {
      const response = await fetch("http://localhost:8090/api/encrypt", {
        mode: "cors",
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: token,
      });
      const registration = await response.text();
      setToken(registration);
    };

    postValue().catch((reason) => setPostError(reason));
  }

  return (
    <Layout>
      <div className="mx-auto max-w-2xl text-center">
        <h2 className="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
          Encrypt Secret
        </h2>
        <p className="mt-2 text-lg leading-8 text-gray-600">
          Use this form to generate an encrypted token.
          <br />
          Simply the value in your configuration file by the generated token in your repository
        </p>
      </div>
      <div className="mx-auto mt-16 max-w-xl sm:mt-20">
        <div className="grid grid-cols-1 gap-x-8 gap-y-6 sm:grid-cols-2">
          <div className="sm:col-span-2">
            <div>
              <label
                htmlFor="value"
                className="block text-sm font-semibold leading-6 text-gray-900"
              >
                Value
              </label>
              <div className="mt-2.5">
                <textarea
                  name="value"
                  rows={3}
                  id="value"
                  value={value}
                  onChange={(e) => setValue(e.target.value)}
                  className="block w-full rounded-md border-0 px-3.5 py-2 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                />
              </div>
            </div>
          </div>
          <div className="sm:col-span-2">
            <div>
              <label
                htmlFor="token"
                className="block text-sm font-semibold leading-6 text-gray-900"
              >
                Token
              </label>
              <div className="mt-2.5">
                <Alert title="Woops" visible={isPostError()}>
                  Woopsie! Something unexpected happend while trying to generate
                  a new Client ID and Secret pair
                  <span>{postError}</span>
                </Alert>
                <span className="text-lg text-gray-600 dark:text-gray-300">
                  {token}
                </span>
              </div>
            </div>
          </div>
        </div>
        <div className="mt-10">
          <button
            onClick={generateToken}
            className="block w-full rounded-md bg-indigo-600 px-3.5 py-2.5 text-center text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
          >
            Generate
          </button>
        </div>
      </div>
    </Layout>
  );
}
