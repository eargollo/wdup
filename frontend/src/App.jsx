import {useState} from 'react';
import './App.css';
import {DuplicateSearch, Emulate} from "../wailsjs/go/main/App";
import SelectPaths from './SelectPaths';
import DuplicatedFiles from './DuplicatedFiles';

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

    return (
        <div id="App">
            {/* <button className="btn" onClick={emulate}>Emulate</button> */}
            <SelectPaths runSearch={runSearch}/>
            <DuplicatedFiles files={duplicatedFiles} loading={loading}/>
        </div>
    )
}

export default App
