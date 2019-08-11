function connectMember() {
    const username = document.getElementById("username-connect").value

    const req = {
        "action": "connect",
        "member": { username }
    }
    server.send(JSON.stringify(req))
}

function createRoom() {
    const req = {
        "action": "create",
        "member": youData
    }
    server.send(JSON.stringify(req))
}

function joinRoom() {
    const id = ~~(document.getElementById("id-room-input").value)
    const req = {
        "action": "join",
        "room": { id },
        "member": youData
    }
    server.send(JSON.stringify(req))
}

function leaveRoom() {
    const req = {
        "action": "leave",
        "room": currentRoom,
        "member": youData
    }
    server.send(JSON.stringify(req))
}

function lockRoom() {
    const req = {
        "action": "lock",
        "room": currentRoom,
        "member": youData
    }
    server.send(JSON.stringify(req))
}