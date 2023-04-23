import Alert from "@/components/alert";
import Layout from "@/components/layout";
import { Inter } from "next/font/google";
import { ReactEventHandler, useEffect, useState } from "react";

const inter = Inter({ subsets: ["latin"] });

interface Repository {
  name: string;
  lastUpdate: string;
  nextUpdate: string;
  hitCount: number;
  lastError: string;
}

export default function Register() {
  const [repositories, setRepositories] = useState<Repository[]>([]);
  const [fetchError, setFetchError] = useState<string>("");

  useEffect(() => {
    const fetchRepositories = async () => {
      setFetchError("");
      const response = await fetch("http://localhost:8090/api/stats", {
        mode: "cors",
      });
      const repos: Repository[] = await response.json();
      setRepositories(repos);
    };

    fetchRepositories().catch((reason) => {
      setFetchError(reason.message);
    });
  }, []);

  if (fetchError.length > 0) {
    return (
      <Layout>
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
    <Layout>
      <div className="mx-auto max-w-2xl text-center">
        <h2 className="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
          Statistics
        </h2>
        <p className="mt-2 text-lg leading-8 text-gray-600">
          I'm serving, they're hating
        </p>
      </div>
      <div className="mx-auto mt-16 max-w-xl sm:mt-20">
        <div className="grid grid-cols-1 gap-x-8 gap-y-6 sm:grid-cols-2">
          <table className="table-auto">
            <thead>
              <tr>
                <th>Repository</th>
                <th>Hit Count</th>
                <th>Last Update</th>
                <th>Next Update</th>
                <th>Last Error</th>
              </tr>
            </thead>
            <tbody>
              {repositories.map((repo, index) => (
                <tr key={index}>
                  <th>{repo.name}</th>
                  <td>{repo.hitCount}</td>
                  <td>{new Date(repo.lastUpdate).toLocaleString(undefined,{dateStyle:"medium", timeStyle:"medium"})}</td>
                  <td>{new Date(repo.nextUpdate).toLocaleString(undefined,{dateStyle:"medium", timeStyle:"medium"})}</td>
                  <td>{repo.lastError}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </Layout>
  );
}
