package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type request struct {
	Action string
	Member member
	Room   room
}

type member struct {
	ID       int             `json:"id,omitempty"`
	Username string          `json:"username,omitempty"`
	Role     string          `json:"role,omitempty"`
	Conn     *websocket.Conn `json:"-"`
}

type room struct {
	// TODO: Add date of creation and name
	ID       int      `json:"id,"`
	Members  []member `json:"members,"`
	IsLocked bool     `json:"is_locked,"`
	Master   int      `json:"master,"`
}

type response struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Room    room   `json:"room,omitempty"`
	Member  member `json:"member,omitempty"`
}

type appData struct {
	roomID   int
	memberID int
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// This array will contain rooms data
// First room with ID = 1 will be the lobby room where users that no belong to a room will stay
var rooms = []room{
	room{
		ID:       1,
		Members:  []member{},
		IsLocked: false,
	},
}

var ad = appData{
	roomID:   2,
	memberID: 1,
}

func isInLobby(memberID int) (member, error) {
	for _, m := range rooms[0].Members {
		if m.ID == memberID {
			return m, nil
		}
	}
	return member{}, errors.New("User with ID " + string(memberID) + " is not in the lobby")
}

func getIndexRoom(roomID int) (int, error) {
	for i, r := range rooms {
		if r.ID == roomID {
			return i, nil
		}
	}
	return -1, errors.New("Room not found")
}

func getMemberLocationByID(memberID int) (member, room, int, int, error) {
	for i, r := range rooms {
		for j, m := range r.Members {
			if m.ID == memberID {
				return m, r, i, j, nil
			}
		}
	}
	return member{}, room{}, -1, -1, errors.New("Member not found")
}

func getMemberLocationByUsername(username string) (member, *room, int, int, error) {
	for i, r := range rooms {
		for j, m := range r.Members {
			if m.Username == username {
				return m, &r, i, j, nil
			}
		}
	}
	return member{}, nil, -1, -1, errors.New("Member not found")
}

func getRoomAndIsMember(memberID int, roomID int) (int, string, int) {
	indexDestinationRoom, err := getIndexRoom(roomID)

	if err != nil {
		return -1, err.Error(), 401
	}

	_, _, _, _, err = getMemberLocationByID(memberID)

	if err != nil {
		return -1, err.Error(), 402
	}

	return indexDestinationRoom, "", -1
}

func responseToByteArray(res interface{}) []byte {
	jsonToSend, errJSON := json.Marshal(res)
	if errJSON != nil {
		panic(errJSON)
	}
	return jsonToSend
}

func moveMember(m member, indexOriginRoom int, roomID int) (error, int) {
	indexMember := -1
	for i, mm := range rooms[indexOriginRoom].Members {
		if mm.ID == m.ID {
			indexMember = i
		}
	}

	if indexMember < 0 {
		return errors.New("Member not found"), 75
	}

	indexDestinationRoom, err := getIndexRoom(roomID)

	if err != nil {
		return err, 76
	}

	// Nothing to do. Member is already in his/her destination
	if indexOriginRoom == indexDestinationRoom {
		return errors.New("Member is already in the room"), 77
	}

	//fmt.Printf("Room destination before: %+v\n", rooms[indexDestinationRoom])
	// Move the member
	rooms[indexOriginRoom].Members = append(rooms[indexOriginRoom].Members[:indexMember], rooms[indexOriginRoom].Members[indexMember+1:]...)
	rooms[indexDestinationRoom].Members = append(rooms[indexDestinationRoom].Members, m)
	//fmt.Printf("Room destination after: %+v\n", rooms[indexDestinationRoom])

	// If the destination room is empty, the new member will be the master
	if len(rooms[indexDestinationRoom].Members) == 1 {
		rooms[indexDestinationRoom].Master = m.ID
	}

	// If the member was the last member in the room
	if len(rooms[indexOriginRoom].Members) == 0 {
		// We remove that room if it is not the lobby
		if indexOriginRoom != 1 {
			// TODO: Remove the room
		}
	} else if rooms[indexOriginRoom].Master == m.ID {
		// If the member was a master in his/her room, a new master will be chosen
		rand.Seed(time.Now().Unix())
		index := rand.Int() % len(rooms[indexOriginRoom].Members)
		rooms[indexOriginRoom].Master = rooms[indexOriginRoom].Members[index].ID
	}

	// Multicast members of the origin and the destination room
	destinationRes := response{
		Code:    51,
		Message: m.Username + " has joined the room",
		Room:    rooms[indexDestinationRoom],
		Member:  m,
	}

	originRes := response{
		Code:    52,
		Message: m.Username + " has left the room",
		Room:    rooms[indexOriginRoom],
		Member:  m,
	}

	multicast(indexDestinationRoom, destinationRes)
	multicast(indexOriginRoom, originRes)

	return nil, -1
}

func setRoles(index int) {
	rand.Seed(time.Now().Unix())
	spyIndex := rand.Int() % len(rooms[index].Members)

	for i := 1; i < len(rooms[index].Members); i++ {
		if spyIndex == i {
			rooms[index].Members[i].Role = "spy"
		} else {
			rooms[index].Members[i].Role = "citizen"
		}
	}
}

func connect(req request) {
	if req.Member.Username == "" {
		unicastError(errors.New("You need to specify an Username"), req.Member.Conn, 1)
		return
	}

	// 2. Check that username has not been taken
	_, _, _, _, err := getMemberLocationByUsername(req.Member.Username)

	if err == nil {
		unicastError(errors.New("Username "+req.Member.Username+" has already been taken"), req.Member.Conn, 2)
		return
	}

	// 3. Add new member to the lobby (room with ID = 1)
	var newMember = member{
		Username: req.Member.Username,
		Role:     "",
		ID:       ad.memberID,
		Conn:     req.Member.Conn,
	}

	ad.memberID++

	rooms[0].Members = append(rooms[0].Members, newMember)

	res := response{
		Code:    3,
		Message: newMember.Username + " has joined the lobby",
		Room:    rooms[0],
		Member:  newMember,
	}

	multicast(0, res)

}

func createRoom(req request) {
	// 1. Find member
	m, _, roomIndexOrigin, _, err := getMemberLocationByID(req.Member.ID)

	if err != nil {
		unicastError(err, req.Member.Conn, 201)
		return
	}

	// 2. Create room
	var r = room{
		ID:       ad.roomID,
		IsLocked: false,
	}
	ad.roomID++
	rooms = append(rooms, r)

	// 3. Move user
	if err, errorCode := moveMember(m, roomIndexOrigin, r.ID); err != nil {
		unicastError(err, m.Conn, errorCode)
		return
	}
}

func joinRoom(req request) {
	// 1. Find member
	m, _, roomIndexOrigin, _, err := getMemberLocationByID(req.Member.ID)
	if err != nil {
		unicastError(err, req.Member.Conn, 301)
		return
	}

	// 2. Find room
	indexDestination, err := getIndexRoom(req.Room.ID)

	if err != nil {
		unicastError(err, req.Member.Conn, 302)
		return
	}

	// 3. Check if the room is unlocked
	if rooms[indexDestination].IsLocked {
		err = errors.New("The room is locked")
		unicastError(err, req.Member.Conn, 303)
		return
	}

	// 4. Move user
	err, errorCode := moveMember(m, roomIndexOrigin, req.Room.ID)
	if err != nil {
		unicastError(err, req.Member.Conn, errorCode)
	}
}

func lockRoom(req request) {
	// 1. Check if the user who requested is a member of the room and if the room exists
	indexDestinationRoom, errorMessage, errorCode := getRoomAndIsMember(req.Member.ID, req.Room.ID)

	if errorCode != -1 {
		unicastError(errors.New(errorMessage), req.Member.Conn, errorCode)
		return
	}

	// 3. Check if the user who requested is the master of the room
	if rooms[indexDestinationRoom].Master != req.Member.ID {
		err := errors.New("You are not the master of the room")
		unicastError(err, req.Member.Conn, 403)
		return
	}

	// 4. Lock the room
	rooms[indexDestinationRoom].IsLocked = true

	// 5. Set roles
	setRoles(indexDestinationRoom)

	// 6. Multicast all members of the room
	res := response{
		Code:    404,
		Message: "The room has been locked",
		Room:    rooms[indexDestinationRoom],
	}
	multicast(indexDestinationRoom, res)
}

func leaveRoom(req request) {
	// 1. Check if the user who requested is a member of the room and if the room exists
	indexOriginRoom, errorMessage, errorCode := getRoomAndIsMember(req.Member.ID, req.Room.ID)

	if errorCode != -1 {
		unicastError(errors.New(errorMessage), req.Member.Conn, errorCode)
		return
	}

	// 3. Move member to lobby
	err, errorCode := moveMember(req.Member, indexOriginRoom, 1)
	if err != nil {
		unicastError(err, req.Member.Conn, errorCode)
		return
	}
}

func multicast(indexRoom int, res response) {
	jsonToSend, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Printf("Multicast to room #%d the following JSON:\n\t%s\n", rooms[indexRoom].ID, string(jsonToSend))
	for _, m := range rooms[indexRoom].Members {
		if m.Conn == nil {
			fmt.Printf("Conn nil.\nRoom:%+v\n", rooms[indexRoom])
		}
		if err = m.Conn.WriteMessage(1, jsonToSend); err != nil {
			fmt.Println(err)
		}
	}

}

func unicastError(err error, conn *websocket.Conn, code int) {
	res := response{
		Code:    code,
		Message: err.Error(),
	}
	unicast(res, conn)
}

func unicast(res response, conn *websocket.Conn) {
	jsonToSend := responseToByteArray(res)

	if err := conn.WriteMessage(1, jsonToSend); err != nil {
		fmt.Println(err)
	}

	//fmt.Println("JSON sent:\n\t", string(jsonToSend))
}

func reader(req request) {
	//fmt.Printf("Request: %+v\n", req)
	switch req.Action {
	case "connect":
		connect(req)

	case "create":
		createRoom(req)

	case "join":
		joinRoom(req)

	case "lock":
		lockRoom(req)

	case "leave":
		leaveRoom(req)

	default:
		res := response{
			Code:    1001,
			Message: "Action not found. You can use: create | join | lock | leave",
		}
		unicast(res, req.Member.Conn)
	}

	return

}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}
	log.Println("Client connected successfully")

	for {

		req := request{}

		err = conn.ReadJSON(&req)

		if err != nil {
			res := response{
				Code:    1002,
				Message: "Unable to read JSON",
			}
			unicast(res, conn)
			return
		}

		req.Member.Conn = conn

		reader(req)
	}
}

func setupRoutes() {
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
