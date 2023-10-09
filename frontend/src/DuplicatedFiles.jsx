import {OpenPath, OpenFile, DeleteFile} from "../wailsjs/go/main/App";

function DuplicatedFiles({files, loading}) {
    // Big todo, put all changes into state and use react to update the DOM
    let result;

    function deleteFile(e, id, file) {
        console.log("Delete file: ", file)
        DeleteFile(file).then(({success, error}) => {
            console.log(success)
            if (success != true) {
                console.log(error)
            } else {
                console.log("ID: ", id)
                parent.document.getElementById(id).classList.add("file-deleted");
                var children = parent.document.getElementById(id).children;
                for (var i = 0; i < children.length; i++) {
                    children[i].disabled = true;
                }
            }
        })
        e.preventDefault();
    }

    function checkBoxUpdate(e, id) {
        parent.document.getElementById(`filename_${id}`).classList.toggle("file-selected");
        // if (e.target.checked) {
        //     parent.document.getElementById(`delete_${id}`).disabled=false;
        // } else {
        //     parent.document.getElementById(`delete_${id}`).disabled=true;
        // }
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
                {/* <button className='btn' onClick={deleteAll}>Delete selected</button> */}
                { files.map((item, index) => ( 
                    <div id={index} key={index} className={index%2?"odd":"even"}>
                    {item.map((file, fid) => (
                        <div className="file" id={`${index}_${fid}`} key={`${index}_${fid}`}>
                            <input type="checkbox" id={`check_${index}_${fid}`} key={`check_${index}_${fid}`} onChange={(e) => checkBoxUpdate(e, `${index}_${fid}`)} />
                            <span id={`filename_${index}_${fid}`}>{file.path}</span>
                            <button className="btn" id={`bt_${index}_${fid}`} key={`bt_${index}_${fid}`} onClick={() => openPath(file.path)}>Open Folder</button>
                            <button className="btn" id={`of_${index}_${fid}`} key={`of_${index}_${fid}`} onClick={() => openFile(file.path)}>Open File</button>
                            <button className="btn btn-warn" id={`delete_${index}_${fid}`} key={`delete_${index}_${fid}`} onClick={(e) => {deleteFile(e, `${index}_${fid}`, file.path)}}>Delete File</button>
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