import Alert from "@/components/alert";
import Layout from "@/components/layout";
import { Inter } from "next/font/google";
import { ReactEventHandler, useEffect, useState } from "react";

const inter = Inter({ subsets: ["latin"] });

interface Repository {
  name: string;
  selected: boolean;
}

interface RegistrationRequest {
  clientId: string;
  repositories: string[];
}

interface RegistrationResponse {
  clientId: string;
  clientSecret: string;
}

export default function Register() {
  const [repositories, setRepositories] = useState<Repository[]>([]);
  const [clientId, setClientId] = useState<string>("");
  const [clientSecret, setClientSecret] = useState<string>("");
  const [fetchError, setFetchError] = useState<string>("");
  const [postError, setPostError] = useState<string>("");

  useEffect(() => {
    const fetchRepositories = async () => {
      setFetchError("");
      const response = await fetch("http://localhost:8090/api/repositories", {
        mode: "cors",
      });
      const names: string[] = await response.json();
      const repos: Repository[] = names.map((i) => {
        return { name: i, selected: false };
      });
      setRepositories(repos);
    };

    fetchRepositories().catch((reason) => {
      setFetchError(reason.message);
    });
  }, []);

  function select(repo: Repository) {
    repo.selected = !repo.selected;
  }

  const isPostError = () => postError.length > 0;

  function generateClientSecret() {
    const postRegistrationRequest = async (request: RegistrationRequest) => {
      const response = await fetch("http://localhost:8090/api/register", {
        mode: "cors",
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(request),
      });
      const registration: RegistrationResponse = await response.json();
      setClientSecret(registration.clientSecret);
      if (clientId.length == 0) {
        setClientId(registration.clientId);
      }
    };

    postRegistrationRequest({
      clientId: clientId,
      repositories: repositories.map((r) => r.name),
    }).catch((reason) => setPostError(reason));
  }

  if (fetchError.length > 0) {
    return (
      <Layout selected="register">
        <Alert title="Woops" visible={true}>
          Dang! Something wrong happened while loading the list of configured
          repositories. Please make sure you configured at least one repository
          and that your instance is running properly.
          <span>{fetchError}</span>
        </Alert>
      </Layout>
    );
  }

  return (
    <Layout selected="register">
      <div className="mx-auto max-w-2xl text-center">
        <h2 className="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
          Register Client
        </h2>
        <p className="mt-2 text-lg leading-8 text-gray-600">
          Use this form to register a new Client ID/Client Secret pair.
          <br />
          Don&apos;t forget to register the new Client ID in your configuration.
          <br />
          If Client ID is ommitted, a new random identifier will be genreated.
        </p>
      </div>
      <div className="mx-auto mt-16 max-w-xl sm:mt-20">
        <div className="grid grid-cols-1 gap-x-8 gap-y-6 sm:grid-cols-2">
          <div className="sm:col-span-2">
            <div>
              <label
                htmlFor="clientid"
                className="block text-sm font-semibold leading-6 text-gray-900"
              >
                Client Id
              </label>
              <div className="mt-2.5">
                <input
                  type="text"
                  name="clientid"
                  id="clientid"
                  autoComplete="off"
                  value={clientId}
                  onChange={(e) => setClientId(e.target.value)}
                  className="block w-full rounded-md border-0 px-3.5 py-2 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                />
              </div>
            </div>
          </div>
          <div className="sm:col-span-2">
            <div>
              <label
                htmlFor="repository"
                className="block text-sm font-semibold leading-6 text-gray-900"
              >
                Repository
              </label>
              <div className="mt-2.5">
                {repositories.map((item, index) => (
                  <div onClick={() => select(item)} key={index}>
                    <input
                      defaultChecked={item.selected}
                      id="checked-checkbox"
                      type="checkbox"
                      value=""
                      className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
                    />
                    <label
                      htmlFor="checked-checkbox"
                      className="ml-2 text-lg text-gray-600 dark:text-gray-300"
                    >
                      {item.name}
                    </label>
                  </div>
                ))}
              </div>
            </div>
          </div>
          <div className="sm:col-span-2">
            <div>
              <label
                htmlFor="clientid"
                className="block text-sm font-semibold leading-6 text-gray-900"
              >
                Client Secret
              </label>
              <div className="mt-2.5">
                <Alert title="Woops" visible={isPostError()}>
                  Woopsie! Something unexpected happend while trying to generate
                  a new Client ID and Secret pair
                  <span>{postError}</span>
                </Alert>
                <span className="text-lg text-gray-600 dark:text-gray-300">
                  {clientSecret}
                </span>
              </div>
            </div>
          </div>
        </div>
        <div className="mt-10">
          <button
            onClick={generateClientSecret}
            className="block w-full rounded-md bg-indigo-600 px-3.5 py-2.5 text-center text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
          >
            Generate
          </button>
        </div>
      </div>
    </Layout>
  );
}
