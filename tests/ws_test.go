package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

var usernames = [100]string{"Wade", "Dave", "Seth", "Ivan", "Riley", "Gilbert", "Jorge", "Dan", "Brian", "Roberto", "Ramon", "Miles", "Liam", "Nathaniel", "Ethan", "Lewis", "Milton", "Claude", "Joshua", "Glen", "Harvey", "Blake", "Antonio", "Connor", "Julian", "Aidan", "Harold", "Conner", "Peter", "Hunter", "Eli", "Alberto", "Carlos", "Shane", "Aaron", "Marlin", "Paul", "Ricardo", "Hector", "Alexis", "Adrian", "Kingston", "Douglas", "Gerald", "Joey", "Johnny", "Charlie", "Scott", "Martin", "Tristin", "Troy", "Tommy", "Rick", "Victor", "Jessie", "Neil", "Ted", "Nick", "Wiley", "Morris", "Clark", "Stuart", "Orlando", "Keith", "Marion", "Marshall", "Noel", "Everett", "Romeo", "Sebastian", "Stefan", "Robin", "Clarence", "Sandy", "Ernest", "Samuel", "Benjamin", "Luka", "Fred", "Albert", "Greyson", "Terry", "Cedric", "Joe", "Paul", "George", "Bruce", "Christopher", "Mark", "Ron", "Craig", "Philip", "Jimmy", "Arthur", "Jaime", "Perry", "Harold", "Jerry", "Shawn", "Walter"}

var generalIndex int = 0

func connectWS(t *testing.T, server *httptest.Server) *websocket.Conn {
	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect to the server
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	return conn
}

func readWS(t *testing.T, c *websocket.Conn) response {
	_, message, err := c.ReadMessage()
	if err != nil {
		t.Errorf("read: " + err.Error())
	}

	res := byteArray2Response(t, message)
	return res
}

func sendWS(t *testing.T, c *websocket.Conn, req request) {
	r := req2ByteArray(t, req)
	err := c.WriteMessage(websocket.TextMessage, r)
	if err != nil {
		t.Errorf("write error: " + err.Error())
	}
}

func connectMember(t *testing.T, server *httptest.Server, username string, memberID int, expectingErrorResponse bool, expectedRes response) (member, *websocket.Conn) {
	conn := connectWS(t, server)

	req := request{
		Action: "connect",
		Member: member{
			Username: username,
		},
	}

	sendWS(t, conn, req)
	res := readWS(t, conn)
	res.Room = room{}

	if !expectingErrorResponse {
		expectedRes = response{
			Code:    3,
			Message: username + " has joined the lobby",
			Room:    room{},
			Member: member{
				Username: username,
				ID:       memberID,
			},
		}
	}

	compare(t, res, expectedRes)
	return res.Member, conn
}

func compare(t *testing.T, res response, expected response) {

	resStr := string(res2ByteArray(t, res))
	expectedResString := string(res2ByteArray(t, expected))

	if resStr != expectedResString {
		t.Errorf("handler returned unexpected body: got %v want %v", resStr, expectedResString)
	}
}

var addr = flag.String("addr", "localhost:8080", "http service address")

func byteArray2Response(t *testing.T, b []byte) response {
	var res response
	err := json.Unmarshal(b, &res)
	if err != nil {
		t.Errorf("json unmarshal error: " + err.Error())
	}
	return res
}

func req2ByteArray(t *testing.T, r request) []byte {
	b, err := json.Marshal(r)
	if err != nil {
		t.Errorf("json marsha error: " + err.Error())
	}
	return b
}

func res2ByteArray(t *testing.T, r response) []byte {
	b, err := json.Marshal(r)
	if err != nil {
		t.Errorf("json marsha error: " + err.Error())
	}
	return b
}

func sendAndCompare(t *testing.T, c *websocket.Conn, req request, expectedRes response) {
	sendWS(t, c, req)
	res := readWS(t, c)
	compare(t, res, expectedRes)
}

// User join successfully to the system
func TestConnect1(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(wsEndpoint))
	defer server.Close()

	username := usernames[generalIndex]
	_, conn := connectMember(t, server, username, generalIndex+1, false, response{})
	defer conn.Close()

	generalIndex++
}

// User join successfully to the system
func TestConnect2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(wsEndpoint))
	defer server.Close()

	username := ""

	expectedResponse := response{
		Code:    1,
		Message: "You need to specify an Username",
	}

	_, conn := connectMember(t, server, username, -1, true, expectedResponse)
	defer conn.Close()
}

// User join successfully to the system and then other user try to join with the same username
func TestConnect3(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(wsEndpoint))
	defer server.Close()

	username := usernames[generalIndex]
	_, conn1 := connectMember(t, server, username, generalIndex+1, false, response{})
	defer conn1.Close()

	expectedResponse := response{
		Code:    2,
		Message: "Username " + username + " has already been taken",
	}
	_, conn2 := connectMember(t, server, username, -1, true, expectedResponse)
	defer conn2.Close()

	generalIndex++
}

// A bunch of users join successfully
func TestConnect4(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(wsEndpoint))
	defer server.Close()

	target := generalIndex + 6
	for ; generalIndex < target; generalIndex++ {
		username := usernames[generalIndex]
		_, conn := connectMember(t, server, username, generalIndex+1, false, response{})
		defer conn.Close()
	}
}

func TestConnect5(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(wsEndpoint))
	defer server.Close()

	firstTarget := generalIndex + 11
	secondTarget := firstTarget + 11

	for ; generalIndex < secondTarget; generalIndex++ {
		username := usernames[generalIndex]
		_, conn := connectMember(t, server, username, generalIndex+1, false, response{})
		defer conn.Close()
	}

	usernameRepeated := usernames[firstTarget+5]

	expectedResponse := response{
		Code:    2,
		Message: "Username " + usernameRepeated + " has already been taken",
	}

	_, conn := connectMember(t, server, usernameRepeated, -1, true, expectedResponse)
	defer conn.Close()

	for ; generalIndex < secondTarget; generalIndex++ {
		username := usernames[generalIndex]
		_, conn := connectMember(t, server, username, generalIndex+1, false, response{})
		defer conn.Close()
	}

}

// It will create a member and then make a request with an invalid action
func TestWrongAction(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(wsEndpoint))
	defer server.Close()

	username := usernames[generalIndex]
	m, conn := connectMember(t, server, username, generalIndex+1, false, response{})
	defer conn.Close()
	generalIndex++

	req := request{
		Action: "invalid",
		Member: m,
	}
	sendWS(t, conn, req)
	res := readWS(t, conn)

	expectedRes := response{
		Code:    1001,
		Message: "Action not found. You can use: create | join | lock | leave",
	}
	compare(t, res, expectedRes)
}

// It will create a member and then make a request with an invalid action
func TestCreateAction(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(wsEndpoint))
	defer server.Close()

	username := usernames[generalIndex]
	m, conn := connectMember(t, server, username, generalIndex+1, false, response{})
	defer conn.Close()
	generalIndex++

	req := request{
		Action: "create",
		Member: m,
	}

	sendWS(t, conn, req)
	res := readWS(t, conn)

	expectedRes := response{
		Code:    51,
		Message: username + " has joined the room",
		Room: room{
			ID:       2,
			IsLocked: false,
			Members:  []member{0: m},
			Master:   m.ID,
		},
		Member: m,
	}

	compare(t, res, expectedRes)

}

// Member A will create a room. Member B will join the room. Member B will create his own room
func TestAction1(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(wsEndpoint))
	defer server.Close()

	username := usernames[generalIndex]
	m1, conn1 := connectMember(t, server, username, generalIndex+1, false, response{})
	defer conn1.Close()
	generalIndex++

	req1 := request{
		Action: "create",
		Member: m1,
	}

	sendWS(t, conn1, req1)

	res1 := readWS(t, conn1)

	expectedRes1 := response{
		Code:    51,
		Message: username + " has joined the room",
		Room: room{
			ID:       2,
			IsLocked: false,
			Members:  []member{0: m1},
			Master:   m1.ID,
		},
		Member: m1,
	}

	compare(t, res1, expectedRes1)

	username = usernames[generalIndex]
	m2, conn2 := connectMember(t, server, username, generalIndex+1, false, response{})
	defer conn2.Close()
	generalIndex++

	req2 := request{
		Action: "join",
		Member: m2,
		Room: room{
			ID: 2,
		},
	}

	sendWS(t, conn2, req2)

	resFirstMember := readWS(t, conn1)
	res2 := readWS(t, conn2)

	expectedRes2 := response{
		Code:    51,
		Message: username + " has joined the room",
		Room: room{
			ID:       2,
			IsLocked: false,
			Members:  []member{0: m1, 1: m2},
			Master:   m1.ID,
		},
		Member: m2,
	}

	// Check if the member that was already in the room receive the correct data
	compare(t, resFirstMember, expectedRes2)

	// Check if the second member receive the correct data
	compare(t, res2, expectedRes2)

	req3 := request{
		Action: "create",
		Member: m2,
	}

	sendWS(t, conn2, req3)
	res3 := readWS(t, conn1)
	res4 := readWS(t, conn2)

	expectedRes3 := response{
		Code:    52,
		Message: username + " has left the room",
		Room: room{
			ID:       2,
			IsLocked: false,
			Members:  []member{0: m1},
			Master:   m1.ID,
		},
		Member: m2,
	}

	expectedRes4 := response{
		Code:    51,
		Message: username + " has joined the room",
		Room: room{
			ID:       3,
			IsLocked: false,
			Members:  []member{0: m2},
			Master:   m2.ID,
		},
		Member: m2,
	}

	compare(t, res4, expectedRes4)
	compare(t, res3, expectedRes3)

}

// Members 1 to 10 will create its own room
// Members 11 to 20 will join
// Member 13 will leave his room
func TestAction2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(wsEndpoint))
	defer server.Close()

	expectedRes := response{
		Code: 51,
	}

	n := 10
	members := make([]member, n*2)
	roomsTest := make([]room, n)

	for i := 0; i < n; i++ {
		username := usernames[generalIndex]
		m, conn := connectMember(t, server, username, generalIndex+1, false, response{})
		generalIndex++

		req := request{
			Action: "create",
			Member: m,
		}

		sendWS(t, conn, req)
		res := readWS(t, conn)

		expectedRes.Message = username + " has joined the room"
		expectedRes.Room.ID = i + 2
		expectedRes.Room.Master = m.ID
		expectedRes.Room.Members = []member{0: m}
		expectedRes.Member = m
		compare(t, res, expectedRes)

		members[i] = member{
			Conn:     conn,
			Username: username,
			ID:       m.ID,
		}
		roomsTest[i] = expectedRes.Room
	}

	for i := n; i < n*2; i++ {
		username := usernames[generalIndex]
		m, conn := connectMember(t, server, username, generalIndex+1, false, response{})
		generalIndex++

		req := request{
			Action: "join",
			Member: m,
			Room:   roomsTest[i-n],
		}

		sendWS(t, conn, req)

		expectedRes.Message = username + " has joined the room"
		expectedRes.Room = roomsTest[i-n]
		expectedRes.Room.Members = []member{0: members[i-n], 1: m}
		expectedRes.Member = m

		res1 := readWS(t, members[i-n].Conn)
		res2 := readWS(t, conn)

		compare(t, res1, expectedRes)
		compare(t, res2, expectedRes)

		members[i] = member{
			Conn:     conn,
			Username: username,
			ID:       m.ID,
		}
		roomsTest[i-n] = expectedRes.Room
	}

	req := request{
		Action: "leave",
		Member: members[12],
		Room: room{
			ID: roomsTest[12-n].ID,
		},
	}
	sendWS(t, members[12].Conn, req)

	roomsTest[12%len(roomsTest)] = room{
		ID:       roomsTest[12-n].ID,
		IsLocked: false,
		Master:   roomsTest[12-n].Master,
		Members:  []member{0: members[12-n]},
	}
	expectedRes1 := response{
		Code:    52,
		Message: members[12].Username + " has left the room",
		Room:    roomsTest[12-n],
		Member:  members[12],
	}

	expectedRes2 := response{
		Code:    51,
		Message: members[12].Username + " has joined the room",
		Room: room{
			ID:       1,
			IsLocked: false,
			Members:  []member{0: members[12]},
			Master:   members[12].ID,
		},
		Member: members[12],
	}
	res1 := readWS(t, members[12-n].Conn)
	res2 := readWS(t, members[12].Conn)

	compare(t, res1, expectedRes1)
	compare(t, res2, expectedRes2)

	for _, m := range members {
		m.Conn.Close()
	}
}

// Member A will create a room and lock it
// Member B will try to join the room
func TestAction3(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(wsEndpoint))
	defer server.Close()

	username := usernames[generalIndex]
	m1, conn1 := connectMember(t, server, username, generalIndex+1, false, response{})
	defer conn1.Close()
	generalIndex++

	req1 := request{
		Action: "create",
		Member: m1,
	}

	sendWS(t, conn1, req1)

	res1 := readWS(t, conn1)

	r := room{
		ID:       2,
		IsLocked: false,
		Members:  []member{0: m1},
		Master:   m1.ID,
	}

	expectedRes1 := response{
		Code:    51,
		Message: username + " has joined the room",
		Room:    r,
		Member:  m1,
	}

	compare(t, res1, expectedRes1)

	req2 := request{
		Action: "lock",
		Member: m1,
		Room:   r,
	}

	r.IsLocked = true

	sendWS(t, conn1, req2)
	res2 := readWS(t, conn1)

	expectedRes2 := response{
		Code:    404,
		Message: "The room has been locked",
		Room:    r,
	}

	compare(t, res2, expectedRes2)

	username = usernames[generalIndex]
	m2, conn2 := connectMember(t, server, username, generalIndex+1, false, response{})
	defer conn2.Close()
	generalIndex++

	req3 := request{
		Action: "join",
		Member: m2,
		Room: room{
			ID: r.ID,
		},
	}

	sendWS(t, conn2, req3)

	res3 := readWS(t, conn2)

	expectedRes3 := response{
		Code:    303,
		Message: "The room is locked",
	}

	compare(t, res3, expectedRes3)

}

// It will made a create, join request without member info
// It will made a lock request without room info
// A member will create a room. Other member who is not in the room will try to lock the room
// Other member will join the room and it will try to lock the room
func TestAction4(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(wsEndpoint))
	defer server.Close()

	conn := connectWS(t, server)
	defer conn.Close()

	req1 := request{
		Action: "create",
		Member: member{},
	}

	sendWS(t, conn, req1)

	res1 := readWS(t, conn)

	expectedRes1 := response{
		Code:    201,
		Message: "Member not found",
	}

	compare(t, res1, expectedRes1)

	req2 := request{
		Action: "join",
		Member: member{},
	}

	sendWS(t, conn, req2)

	res2 := readWS(t, conn)

	expectedRes2 := response{
		Code:    301,
		Message: "Member not found",
	}

	compare(t, res2, expectedRes2)

	req3 := request{
		Action: "lock",
		Member: member{},
		Room:   room{},
	}

	sendWS(t, conn, req3)

	res3 := readWS(t, conn)

	expectedRes3 := response{
		Code:    401,
		Message: "Room not found",
	}

	compare(t, res3, expectedRes3)

	username := usernames[generalIndex]
	m1, conn1 := connectMember(t, server, username, generalIndex+1, false, response{})
	defer conn1.Close()
	generalIndex++

	req4 := request{
		Action: "create",
		Member: m1,
	}

	sendWS(t, conn, req4)

	res4 := readWS(t, conn1)

	r := room{
		ID:       2,
		IsLocked: false,
		Members:  []member{0: m1},
		Master:   m1.ID,
	}

	expectedRes4 := response{
		Code:    51,
		Message: username + " has joined the room",
		Room:    r,
		Member:  m1,
	}

	compare(t, res4, expectedRes4)

	req5 := request{
		Action: "lock",
		Member: member{},
		Room:   r,
	}

	sendWS(t, conn, req5)

	res5 := readWS(t, conn)

	expectedRes5 := response{
		Code:    402,
		Message: "Member not found",
	}

	compare(t, res5, expectedRes5)

	username = usernames[generalIndex]
	m2, conn2 := connectMember(t, server, username, generalIndex+1, false, response{})
	defer conn2.Close()
	generalIndex++

	req6 := request{
		Action: "join",
		Member: m2,
		Room:   r,
	}

	sendWS(t, conn, req6)

	res6 := readWS(t, conn2)

	r = room{
		ID:       2,
		IsLocked: false,
		Members:  []member{0: m1, 1: m2},
		Master:   m1.ID,
	}

	expectedRes6 := response{
		Code:    51,
		Message: username + " has joined the room",
		Room:    r,
		Member:  m2,
	}

	compare(t, res6, expectedRes6)

	req7 := request{
		Action: "lock",
		Member: m2,
		Room:   r,
	}

	sendWS(t, conn2, req7)

	res7 := readWS(t, conn2)

	r = room{
		ID:       2,
		IsLocked: false,
		Members:  []member{0: m1, 1: m2},
		Master:   m1.ID,
	}

	expectedRes7 := response{
		Code:    403,
		Message: "You are not themaster of the room",
		Room:    room{},
		Member:  member{},
	}

	compare(t, res7, expectedRes7)

}

