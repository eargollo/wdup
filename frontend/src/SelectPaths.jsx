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
        <div id="input" className="container-flex border border-2">
            <div className="row">
                <div className="col">
                    <div className="row">
                        {/* <div className="col">
                            <label className="form-label">Folders to be scanned:</label>
                        </div> */}
                        <div className="col-2">
                            <button className="btn btn-secondary" onClick={selectAndAddPath}>
                            <i className="bi-folder-plus"></i> Add
                            </button>
                        </div>
                    </div>
                    {folderList.map((item, index) => (
                        <div key={index} className="container row mb-0 input-group">
                            <span className="col text-start badge input-group-text">{item}</span>
                            <button id={index} className="col col-1 badge btn btn-outline-danger" onClick={() => deleteFolderFromList(index)}>
                                <i className="bi-folder-minus"></i>
                            </button>
                        </div>
                    ))}  

                </div>
                <div className="col-2">
                    <button className="btn btn-primary" onClick={run}>
                        <i className="bi bi-search"></i>
                        <div>Find Duplicates</div>
                    </button>
                </div>
            </div>                
        </div>

    )
}

export default SelectPaths;