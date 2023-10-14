import {useState} from 'react';
import './App.css';
import {DuplicateSearch, Emulate} from "../wailsjs/go/main/App";
import SelectPaths from './SelectPaths';
import DuplicatedFiles from './DuplicatedFiles';
import {EventsOn} from "../wailsjs/runtime/runtime";

function App() {
    const [duplicatedFiles, setDuplicatedFiles] = useState([]);
    const [loading, setLoading] = useState(false);
    let result;

    function emulate() {
        Emulate().then(setDuplicatedFiles);
    }

    function runSearch(folderList) {
        setLoading(true);
        DuplicateSearch(folderList).then(setDuplicatedFiles).finally(() => setLoading(false));
    }

    EventsOn("Progress", (data) => {
        console.log(data)
        parent.document.getElementById("event").innerHTML = data.description;
    })

    return (
        <div id="App" className="">
            {/* <button className="btn" onClick={emulate}>Emulate</button> */}
            <SelectPaths runSearch={runSearch}/>
            <div><span id="event"></span></div>
            <DuplicatedFiles files={duplicatedFiles} loading={loading}/>
        </div>
    )
}

export default App
