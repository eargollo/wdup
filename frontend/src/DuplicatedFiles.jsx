import FileGroup from "./FileGroup";

function DuplicatedFiles({files, loading}) {
    // Big todo, put all changes into state and use react to update the DOM
    let result;


    function checkBoxUpdate(e, id) {
        parent.document.getElementById(`filename_${id}`).classList.toggle("file-selected");
        // if (e.target.checked) {
        //     parent.document.getElementById(`delete_${id}`).disabled=false;
        // } else {
        //     parent.document.getElementById(`delete_${id}`).disabled=true;
        // }
    }

    // Take loading out of this component
    if (loading) {
        result = <div id="result">
                <span>Loading...</span>
            </div>;
    } else {
        if (files.length == 0) {
            result = <div id="result"/>
        } else {
            result = <div id="result">
                { files.map((item, index) => ( 
                    <FileGroup id={index} files={item} key={index}/>
                ))}
            </div>;
        }
    }

    return result;
}

export default DuplicatedFiles;