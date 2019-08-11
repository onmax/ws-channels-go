function connectWS(isDefault) {
    const port = isDefault ? 8080 : document.getElementById("port_input").value


    s = new WebSocket(`ws://localhost:${port}/ws`)
    errorBox("none")

    s.onopen = () => {
        console.log("Connected", new Date())
        createMemberForm("inherit")
    }

    s.onerror = () => {
        const msg = `<p>Couldn't connect to WebSocket server on port ${port}. Type the port and try again: </p><input id="port_input" type="number"><button onclick="connectWS(false)">Try again</button>`
        errorBox("inherit", msg)
    }

    s.onmessage = msg => {
        const res = JSON.parse(msg.data)
        console.log(res)
        errorBox("none")

        switch (res.code) {
            case 3:
                youData = res.member
                currentRoom = res.room
                console.table(youData)
                console.log("Current room", currentRoom)
                updateRoomView(res)
                break;
            case 51:
                currentRoom = res.room
                console.log("Current room", currentRoom)
                updateRoomView(res)
                break;
            case 52:
                currentRoom = res.room
                console.log("Current room", currentRoom)
                updateRoomView(res)
                break;
            case 404:
                currentRoom = res.room
                console.log("Current room", currentRoom)
                updateRoomView(res)
                break;
            default:
                errorBox("inherit", res.message)
        }
    }

    return s
}
