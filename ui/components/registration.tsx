import { ReactEventHandler, useEffect, useState } from "react";
import Alert from "./alert";

interface Repository {
  name: string;
  selected: boolean;
}

interface RegistrationRequest {
  clientId: string;
  repositories: string[];
}

interface RegistrationResponse {
  ClientId: string;
  ClientSecret: string;
}

export default function Repositories() {
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
      setClientSecret(registration.ClientSecret);
    };

    postRegistrationRequest({
      clientId: clientId,
      repositories: repositories.map((r) => r.name),
    }).catch((reason) => setPostError(reason));
  }

  if (fetchError.length > 0) {
    return (
      <Alert title="Woops" message="">
        Dang! Something wrong happened while loading the list of configured
        repositories. Please make sure you configured at least one repository
        and that your instance is running properly.
        <div>{fetchError}</div>
      </Alert>
    );
  }

  return (
    <div>
      <input
        value={clientId}
        onChange={(e) => setClientId(e.target.value)}
      ></input>
      <ul>
        {repositories.map((item, index) => (
          <li
            key={index}
            onClick={() => select(item)}
            defaultChecked={item.selected}
          >
            <input type="checkbox"></input>
            {item.name}
          </li>
        ))}
      </ul>
      <button
        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 border border-blue-700 rounded"
        onClick={generateClientSecret}
      >
        Generate Client Secret
      </button>
      <h3>{clientSecret}</h3>
      if(isPostError)
      {
        <Alert title="Woops" message="">
          Woopsie! Something unexpected happend while trying to generate a new Client ID and Secret pair
          <div>{postError}</div>
        </Alert>
      }
    </div>
  );
}
