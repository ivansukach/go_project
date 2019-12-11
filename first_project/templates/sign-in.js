document.addEventListener("DOMContentLoaded", setTimeout);
setTimeout(ready, 5000);

function ready() {
    let data=document.getElementById('setValue');
    let value=document.getElementById('tkn').innerText;
    let button=document.getElementById('putToken');
    data.value=value;
    button.click();
}
