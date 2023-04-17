import { ReactEventHandler, useEffect, useState } from "react";

interface Repository{
    name:string;
    selected:boolean;
}

interface RegistrationRequest{
    clientId:string;
    repositories:string[];
}

interface RegistrationResponse{
    ClientId:string;
    ClientSecret:string;
}


export default function Repositories(){

    const [repositories, setRepositories] = useState<Repository[]>([]);
    const [clientId, setClientId] = useState<string>("");
    const [clientSecret, setClientSecret] = useState<string>("");

    useEffect(()=>{
        const fetchRepositories = async () =>{
            const response = await fetch("http://localhost:8090/api/repositories", {mode:"cors"});
            const names:string[] = await response.json();
            const repos: Repository[] = names.map((i) => {
                return {name: i, selected: false};
            });
            setRepositories(repos);
        }

        fetchRepositories().catch(reason =>{});
    },[]);

    function select(repo:Repository){
        repo.selected=!repo.selected;
    }

    function generateClientSecret(){
        const request: RegistrationRequest = {
            clientId: clientId,
            repositories: repositories.map(r => r.name)
        };

        fetch("http://localhost:8090/api/register", {mode:"cors", method:"POST", headers: {
            "Content-Type": "application/json"
          }, body: JSON.stringify(request)})
         .then(r=>r.json())
         .then((r:RegistrationResponse)=>setClientSecret(r.ClientSecret));

    }

    return (
        <div>
<div className="bg-teal-100 border-t-4 border-teal-500 rounded-b text-teal-900 px-4 py-3 shadow-md" role="alert">
  <div className="flex">
    <div className="py-1"><svg className="fill-current h-6 w-6 text-teal-500 mr-4" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"><path d="M2.93 17.07A10 10 0 1 1 17.07 2.93 10 10 0 0 1 2.93 17.07zm12.73-1.41A8 8 0 1 0 4.34 4.34a8 8 0 0 0 11.32 11.32zM9 11V9h2v6H9v-4zm0-6h2v2H9V5z"/></svg></div>
    <div>
      <p className="font-bold">Our privacy policy has changed</p>
      <p className="text-sm">Make sure you know how these changes affect you.</p>
    </div>
  </div>
</div>
        <input value={clientId} onChange={e => setClientId(e.target.value)}></input>
        <ul>
            {repositories.map((item, index)=>(<li key={index} onClick={()=>select(item)} defaultChecked={item.selected}><input type="checkbox"></input>{item.name}</li>))}
        </ul>
        <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 border border-blue-700 rounded" onClick={generateClientSecret}>
            Generate Client Secret
            </button>
            <h3>{clientSecret}</h3>
        </div>
        
    );
}