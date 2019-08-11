function errorBox(status, msg) {
    const errorDiv = document.querySelector(".error-msg")
    errorDiv.style.display = status
    errorDiv.innerHTML = msg ? msg : ""
}


function createMemberForm(status) {
    document.querySelector(".create-member").style.display = status
}

function updateRoomView(json) {
    document.querySelector(".create-member").style.display = "none"

    const room = document.querySelector(".room")

    room.style.display = "inherit"

    const usernameMember = room.querySelector("span#my-username")
    usernameMember.innerText = youData.username

    const idMember = room.querySelector("span#my-id")
    idMember.innerText = youData.id

    const idSpan = room.querySelector("span#id-room")
    idSpan.innerText = json.room.id

    const lockedSpan = room.querySelector("span#is-locked")
    lockedSpan.innerText = json.room.is_locked ? "locked" : "unlocked"


    const membersDiv = room.querySelector(".members")
    membersDiv.innerHTML = "<h6>Members</h6>"

    json.room.members.map(member => {
        const mDiv = document.createElement("div")
        mDiv.innerText = `${member.id} - ${member.username}`
        mDiv.innerText += json.room.master === member.id ? " (master)" : ""
        if (json.member.id === member.id) {
            mDiv.classList.add("shine")
        }
        membersDiv.appendChild(mDiv)
    })

    document.querySelector(".message").innerText = `${json.code} - ${json.message}`
}

let server = connectWS(true)
let youData = {}
let currentRoom = {}

