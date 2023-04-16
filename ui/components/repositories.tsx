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
        fetch("http://localhost:8090/api/repositories", {mode:"cors"})
        .then(r => r.json())
        .then((r: string[])=>{
            const repos: Repository[] = r.map((item) => {
                return {name: item, selected: false};
            });
            setRepositories(repos);
        });
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