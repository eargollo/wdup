import {useState} from 'react';
import { useRef } from 'react';
import './App.css';
import {Greet} from "../wailsjs/go/main/App";
import {PathSelect} from "../wailsjs/go/main/App";

function App() {
    const [resultText, setResultText] = useState([]);
    const [folderText, setFolderText] = useState("/");
    const [name, setName] = useState('');
    const updateName = (e) => setName(e.target.value);
    const updateResultText = (result) => setResultText(result);

    const inputFile = useRef(null);

    function greet() {
        Greet(name).then(updateResultText);
    }

    function pathSelect() { 
        PathSelect("a").then(setFolderText);
        e.preventDefault();
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

    function onChangeFile(e) {
        // e,stopPropagation();
        e.preventDefault();
        console.log(e);
        var files = e.target.files;
        e.currentTarget.attributes[2].ownerElement.files
        console.log(files)
    }

    return (
        <div id="App">
            <div>
                {/* <input onChange={onChangeFile} type="file" id="filepicker" name="fileList" webkitdirectory={""} directory={""} multiple={true} /> */}
                {/* <input onChange={onChangeFile} type='file' id='file' ref={inputFile} webkitdirectory multiple/> */}
            </div>
            <div id="input" className="input-box">
                <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text" value={folderText}/>
                <button className="btn" onClick={pathSelect}>Select Path</button>
                <button className="btn" onClick={greet}>Run</button>
                {/* <input type='file' id='file' ref={inputFile} style={{display: 'none'}}/> */}
            </div>
            <div id="result">
                <button className='btn' onClick={deleteFile}>Delete selected</button>
                { resultText.map((item, index) => ( 
                    <div id={index} className={index%2?"odd":"even"}>
                    {item.map((file, fid) => (
                        <div className="file">
                            <input type="checkbox" id={`${index}_${fid}`} key={file.name} onClick={textHighlight} />
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
