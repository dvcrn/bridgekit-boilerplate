package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"maunium.net/go/mautrix/bridge/bridgeconfig"
	"maunium.net/go/mautrix/bridge/commands"
	"maunium.net/go/mautrix/format"

	"github.com/dvcrn/bridgekit/pkg"
	"github.com/dvcrn/bridgekit/pkg/domain"
	"maunium.net/go/mautrix/bridge"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

var _ pkg.BridgeConnector = (*MyBridgeConnector)(nil)
var _ pkg.MatrixRoomEventHandler = (*MyBridgeConnector)(nil)

type MemDB struct {
	Users map[id.UserID]*domain.User
	Rooms map[id.RoomID]*domain.Room
}

func (mdb *MemDB) Store() {
	data, err := json.Marshal(mdb)
	if err != nil {
		fmt.Println("Error marshaling database:", err)
		return
	}

	err = os.WriteFile("db.json", data, 0644)
	if err != nil {
		fmt.Println("Error writing database to file:", err)
		return
	}
}

func (mdb *MemDB) Load() {
	data, err := os.ReadFile("db.json")
	if err != nil {
		fmt.Println("Error reading database file:", err)
		return
	}

	err = json.Unmarshal(data, mdb)
	if err != nil {
		fmt.Println("Error unmarshaling database:", err)
		return
	}
}

func NewMemDB() *MemDB {
	return &MemDB{
		Users: map[id.UserID]*domain.User{},
		Rooms: map[id.RoomID]*domain.Room{},
	}
}

type MyBridgeConnector struct {
	kit *pkg.BridgeKit

	memDb *MemDB
}

func (m MyBridgeConnector) Init(ctx context.Context) error {
	fmt.Println("Initializing MyBridgeConnector")
	m.memDb.Load()

	for _, room := range m.memDb.Rooms {
		m.kit.RoomManager.LoadRoomIntent(room)
	}

	m.kit.RegisterCommand(&commands.FullHandler{
		Func: func(e *commands.Event) {
			// dummy function function
			user := e.User.(*domain.User)
			var room *domain.Room
			if e.Portal != nil {
				room = e.Portal.(*domain.Room)
			}
			fmt.Println("Login called", user.DisplayName, room.MXID.String())

			e.Reply("Okay, you logged in!!")

			// authenticate here. room is usually the management room
			// the user is already in the mem-DB because GetUser is called first
			// Let's create some rooms for them
			m.createDummyRooms(ctx, user)

		},
		Name: "login",
		Help: commands.HelpMeta{
			Section:     commands.HelpSectionAuth,
			Description: "Authenticate with the bridge",
		},
	})

	return nil
}

func (m MyBridgeConnector) createDummyRooms(ctx context.Context, user *domain.User) {
	ghost := m.kit.GhostMaster.NewGhost(
		"SomeUserID",
		"Test User",
		"user_name",
		id.ContentURI{},
	)

	room := m.kit.RoomManager.NewRoom("Test Room", "Some Topic", ghost)
	createdRoom, _, err := m.kit.CreateRoom(ctx, room, user)
	if err != nil {
		fmt.Println("err : ", err.Error())
		return
	}

	content := format.RenderMarkdown("Hello, I'm a bot", true, false)
	m.kit.SendBotMessageInRoom(ctx, createdRoom, &content)

	content = format.RenderMarkdown("Hello, I'm a ghost", true, false)
	m.kit.SendMessageInRoom(ctx, createdRoom, createdRoom.MainIntent(), &content)
}

func (m MyBridgeConnector) Start(ctx context.Context) {
	fmt.Println("Starting MyBridgeConnector")

	// Bind remote connection, eg websocket
	// Start interval polling
	// Create ghosts

	bridgeConfig := m.kit.Config.Bridge.(MyBridgeConfig)
	fmt.Printf("%v+\n", bridgeConfig.SomeKey)

	// TODO: for debug. remove me.
	func(v interface{}) {
		j, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		buf := bytes.NewBuffer(j)
		fmt.Printf("%v\n", buf.String())
	}(m.kit.Config)

	fmt.Println("--------------------")
	fmt.Println("Started bridge!! Go ahead and message", m.kit.Bot.UserID.String())
	fmt.Println("Use the login command :)")
	fmt.Println("--------------------")

}

func (m MyBridgeConnector) Stop() {
	//TODO implement me
	fmt.Println("Stopping MyBridgeConnector")
	m.memDb.Store()
}

func (m MyBridgeConnector) GetRoom(ctx context.Context, roomID id.RoomID) *domain.Room {
	// check if in DB
	if r, ok := m.memDb.Rooms[roomID]; ok {
		return r
	}

	// not in db
	return &domain.Room{
		MXID: roomID,
	}
}

func (m MyBridgeConnector) GetAllRooms(ctx context.Context) []bridge.Portal {
	//TODO implement me
	panic("implement me")
}

func (m MyBridgeConnector) IsGhost(ctx context.Context, userID id.UserID) bool {
	//TODO implement me
	return false
}

func (m MyBridgeConnector) GetGhost(ctx context.Context, userID id.UserID) *domain.Ghost {
	fmt.Println("GetGhost unimplemented")
	//TODO implement me
	return nil
}

func (m MyBridgeConnector) GetUser(ctx context.Context, uid id.UserID, create bool) *domain.User {
	if u, ok := m.memDb.Users[uid]; ok {
		return u
	}

	u := &domain.User{
		MXID:             uid,
		RemoteID:         "whatsapp_id",
		RemoteName:       "WhatsApp Name",
		DisplayName:      "whatsapp_user",
		PermissionLevel:  bridgeconfig.PermissionLevelAdmin,
		ManagementRoomID: "", // dont have a management room id yet, otherwise return it here
	}

	m.memDb.Users[uid] = u

	return u
}

func (m MyBridgeConnector) SetManagementRoom(ctx context.Context, user *domain.User, roomID id.RoomID) {
	//TODO implement me
	fmt.Println("SetSetManagementRoom for ", user.DisplayName)
	if _, ok := m.memDb.Users[user.MXID]; ok {
		m.memDb.Users[user.MXID].ManagementRoomID = roomID
	}
}

func (m MyBridgeConnector) HandleMatrixRoomEvent(ctx context.Context, room *domain.Room, user bridge.User, evt *event.Event) error {
	switch evt.Type {
	case event.EventMessage:
		fmt.Println("got message event")
		// TODO: for debug. remove me.
		func(v interface{}) {
			j, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			buf := bytes.NewBuffer(j)
			fmt.Printf("%v\n", buf.String())
		}(evt)
	default:
		fmt.Println("got unhandled event type: ", evt.Type.String())
	}

	// marking the message as read to indicate to the user that the bridge has processed it
	if err := m.kit.MarkRead(ctx, evt, room); err != nil {
		fmt.Println("error marking as read: ", err)
	}
	return nil
}

func NewBridgeConnector(bk *pkg.BridgeKit) *MyBridgeConnector {
	br := &MyBridgeConnector{
		kit:   bk,
		memDb: NewMemDB(),
	}
	return br
}
