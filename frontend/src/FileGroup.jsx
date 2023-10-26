import {OpenPath, OpenFile, DeleteFile} from "../wailsjs/go/main/App";
import { useState } from 'react';
import React from 'react';
import prettyBytes from 'pretty-bytes';

function FileGroup({id, files}) {
    const [duplicatedFiles, setDuplicatedFiles] = useState([])
    const [alert, setAlert] = useState({open: false, type: "", message: ""});
    // const [alert, setAlert] = useState({open: true, type: "danger", message: "Attention"});
 
    React.useEffect(() => {
        var result = []
        const unique = new Map()
        console.log("Files: ", files)
        files.forEach((file) => {
            var dFile = Object.assign({}, file)
            dFile.status = "duplicated"
            if (unique.get(file.path) !== true) {
                // Add only if unique
                result.push(dFile);
            }
            unique.set(file.path, true)
        })
        console.log("Result: ", result)
        setDuplicatedFiles(result)
    }, []);
    
    
    function deleteFile(e, id, file) {
        console.log("Delete file: ", file)
        DeleteFile(file).then(({success, error}) => {
            console.log(success)
            if (success != true) {
                //Show error
                setAlert({open: true, type: "danger", message: error});
                console.log(error)
            } else {
                const result = [...duplicatedFiles];
                result.splice(id, 1);
                setDuplicatedFiles(result);
            }
        })
        e.preventDefault();
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
            }
        })
    }

    function dismissAlert() {
        setAlert({open: false, type: "", message: ""});
    }
    
    if (duplicatedFiles.length <=1) {
        return <></>
    }
    return (
        <ul className="list-group">
            <li className="list-group-item active">
                Files {duplicatedFiles.length} File size: {prettyBytes(duplicatedFiles[0].size)} Duplicated size: {prettyBytes(duplicatedFiles[0].size * (duplicatedFiles.length - 1))}  
            </li>
            { alert.open ? <li> <div className={`alert alert-${alert.type} alert-dismissible`}> {alert.message} <button type="button" class="btn-close" onClick={dismissAlert} aria-label="Close"></button></div></li> : <></>}
            { duplicatedFiles.map((file, index) => (
                <li id={`${id}_${index}`} className="list-group-item">
                    <button type="button" className="btn btn-danger" onClick={(e) => {deleteFile(e, index, file.path)}}>
                        <i class="bi bi-trash3-fill"></i>
                    </button>
                    <button type="button" className="btn btn-secondary" onClick={() => openPath(file.path)}>
                        <i class="bi bi-folder2-open"></i>
                    </button>
                    <button type="button" className="btn btn-secondary" onClick={() => openFile(file.path)}>
                        <i class="bi bi-file-earmark-image-fill"></i>
                    </button>
                    <span>{file.path}</span>
                </li>
            ))}
        </ul>
    )
}

export default FileGroup;