import {useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Greet} from "../wailsjs/go/main/App";

function App() {
    const [resultText, setResultText] = useState([]);
    const [name, setName] = useState('');
    const updateName = (e) => setName(e.target.value);
    const updateResultText = (result) => setResultText(result);

    function greet() {
        Greet(name).then(updateResultText);
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

    return (
        <div id="App">
            <div id="input" className="input-box">
                <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text"/>
                <button className="btn" onClick={greet}>Greet</button>
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
