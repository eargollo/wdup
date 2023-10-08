import {useState} from 'react';
import { useRef } from 'react';
import './App.css';
import {Greet} from "../wailsjs/go/main/App";
import {PathSelect} from "../wailsjs/go/main/App";
import {DuplicateSearch} from "../wailsjs/go/main/App";

function App() {
    const [folderList, setFolderList] = useState([]);
    const [resultText, setResultText] = useState([]);
    const updateResultText = (result) => setResultText(result);

    function greet() {
        Greet("").then(updateResultText);
    }

    function pathSelect(e) { 
        PathSelect("a").then(addFile);
        e.preventDefault();
    }

    function addFile(file) {
        console.log(file)
        setFolderList([...folderList, file]);
    }

    function deleteFolderFromList(e) {
        console.log(e);
        console.log(e.target.id);
        const index = e.target.id;
        const newList = [...folderList];
        newList.splice(index, 1);
        setFolderList(newList);
    }

    function deleteFile(e) {
        console.log(e)
        e.preventDefault();
    }

    function textHighlight(e) {
        console.log(e)
        console.log(e.target.id)
        e.target.classList.toggle("file-selected");
        parent.document.getElementById(`filename_${e.target.id}`).classList.toggle("file-selected");
    }

    function runSearch(e) {
        console.log(e)
        DuplicateSearch(folderList).then(updateResultText);
    }

    return (
        <div id="App">
            <div id="input" className="input-box">
                <button className="btn" onClick={pathSelect}>Add Folter</button>
                <button className="btn" onClick={greet}>Emulate</button>
                <button className="btn" onClick={runSearch}>Search Duplicates</button>
                <div>
                    {folderList.map((item, index) => (
                        <div key={index}>
                            <span>{item}</span>
                            <button id={index} className="btn" onClick={deleteFolderFromList}>Delete</button>
                        </div>
                    ))}  
                </div>
            </div>
            <div id="result">
                <button className='btn' onClick={deleteFile}>Delete selected</button>
                { resultText.map((item, index) => ( 
                    <div id={index} key={index} className={index%2?"odd":"even"}>
                    {item.map((file, fid) => (
                        <div className="file" key={`${index}_${fid}`}>
                            <input type="checkbox" id={`${index}_${fid}`} key={`${index}_${fid}`} onClick={textHighlight} />
                            <span id={`filename_${index}_${fid}`}>{file.path}</span>
                        </div>
                    ))}
                    </div>
                ))}
            </div>
        </div>
    )
}

export default App
