import { PathSelect} from "../wailsjs/go/main/App";
import { useState } from 'react';

// TODO: avoid duplicate paths in list
function SelectPaths({runSearch}) {
    const [folderList, setFolderList] = useState([]);

    function selectAndAddPath(e) { 
        PathSelect().then((file) => setFolderList([...folderList, file]))
        e.preventDefault();
    }

    function run() {
        runSearch(folderList);
    }

    function deleteFolderFromList(index) {
        const newList = [...folderList];
        newList.splice(index, 1);
        setFolderList(newList);
    }

    return (
        <div id="input" className="input-box">
            <button className="btn" onClick={selectAndAddPath}>Add Folter</button>
            <button className="btn" onClick={run}>Search Duplicates</button>
            <div>
                {folderList.map((item, index) => (
                    <div key={index}>
                        <span>{item}</span>
                        <button id={index} className="btn" onClick={() => deleteFolderFromList(index)}>Delete</button>
                    </div>
                ))}  
            </div>
        </div>

    )
}

export default SelectPaths;