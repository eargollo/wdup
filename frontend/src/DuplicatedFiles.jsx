import {OpenPath, OpenFile} from "../wailsjs/go/main/App";

function DuplicatedFiles({files, loading}) {
    let result;

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

    function openPath(path) {
        console.log(path)
        OpenPath(path).then((err) => {
            if (err !== "") {
                console.log(err)
                alert(err)
            }
        })
    }

    function openFile(path) {
        console.log(path)
        OpenFile(path).then((err) => {
            if (err !== "") {
                console.log(err)
                alert(err)
            }
        })
    }

    if (loading) {
        result = <div id="result">
                <span>Loading...</span>
            </div>;
    } else {
        if (files.length == 0) {
            result = <div id="result"/>
        } else {
            result = <div id="result">
                <button className='btn' onClick={deleteFile}>Delete selected</button>
                { files.map((item, index) => ( 
                    <div id={index} key={index} className={index%2?"odd":"even"}>
                    {item.map((file, fid) => (
                        <div className="file" key={`${index}_${fid}`}>
                            <input type="checkbox" id={`${index}_${fid}`} key={`${index}_${fid}`} onClick={textHighlight} />
                            <span id={`filename_${index}_${fid}`}>{file.path}</span>
                            <button className="btn" id={`bt_${index}_${fid}`} key={`bt_${index}_${fid}`} onClick={() => openPath(file.path)}>Open Folder</button>
                            <button className="btn" id={`of_${index}_${fid}`} key={`of_${index}_${fid}`} onClick={() => openFile(file.path)}>Open File</button>
                        </div>
                    ))}
                    </div>
                ))}
            </div>;
        }
    }

    return result;
}

export default DuplicatedFiles;