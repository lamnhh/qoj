function handleSocketData(json) {
    switch (json.type) {
        case "result":
            console.log(json.message);
            break;
        default:
            console.log("invalid type");
            break;
    }
}

function createSocket() {
    let socket = new WebSocket("ws://localhost:3000/ws");

    socket.addEventListener("open", function(event) {
        socket.send(JSON.stringify({
            "type": "hello",
            "message": "Hello there"
        }));
    });

    socket.addEventListener("message", function(event) {
        let json = JSON.parse(event.data);
        handleSocketData(json);
    });

    return socket;
}

let socket = createSocket();

function handleSubmit(form) {
    let problemId = form.problemId.value;
    let file = form.file.files[0];

    let body = new FormData();
    body.append("problemId", problemId);
    body.append("file", file);
    fetch("/api/submission", {
        method: "POST",
        body
    }).then(function (res) {
        if (res.ok) {
            return res.json();
        }
        throw res.json();
    }).then(function({submissionId}) {
        socket.send(JSON.stringify({
            "type": "subscribe",
            "message": String(submissionId)
        }));
    }).catch(console.log);
}