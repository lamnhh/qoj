function handleSocketData(json) {
    console.log(json);
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

let accessToken = "";

fetch("/api/login", {
    method: "POST",
    headers: {
        "Content-Type": "application/json"
    },
    body: JSON.stringify({
        username: "lamnhh",
        password: "123456"
    })
}).then(function(res) {
    if (res.ok) {
        return res.json();
    }
    throw res.json();
}).then(function(res) {
    accessToken = res.accessToken;
});

function handleSubmit(form) {
    let problemId = form.problemId.value;
    let file = form.file.files[0];

    let body = new FormData();
    body.append("problemId", problemId);
    body.append("file", file);
    fetch("/api/submission", {
        method: "POST",
        headers: {
            "Authorization": "Bearer " + accessToken
        },
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

function subscribe() {
    let submissionId = parseInt(document.getElementById("submissionId").value);
    socket.send(JSON.stringify({
        "type": "subscribe",
        "message": String(submissionId)
    }));
}

function unsubscribe() {
    let submissionId = parseInt(document.getElementById("submissionId").value);
    socket.send(JSON.stringify({
        "type": "unsubscribe",
        "message": String(submissionId)
    }));
}