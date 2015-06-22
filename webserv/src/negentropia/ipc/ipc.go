package ipc

// IPC code shared among server and client

const (
	CM_CODE_FATAL           = 0
	CM_CODE_INFO            = 1
	CM_CODE_AUTH            = 2  // client->server: let me in
	CM_CODE_ECHO            = 3  // client->server: please echo this
	CM_CODE_KILL            = 4  // server->client: do not attempt reconnect on same session
	CM_CODE_REQZ            = 5  // client->server: please send current zone
	CM_CODE_ZONE            = 6  // server->client: reset client zone info
	CM_CODE_SKYBOX          = 7  // server->client: set full skybox
	CM_CODE_PROGRAM         = 8  // server->client: set shader program
	CM_CODE_INSTANCE        = 9  // server->client: set instance
	CM_CODE_INSTANCE_UPDATE = 10 // server->client: update instance
	CM_CODE_MESSAGE         = 11 // server->client: message for user
	CM_CODE_MISSION_NEXT    = 12 // client->server: switch mission
	CM_CODE_SWITCH_ZONE     = 13 // client->server: switch zone
)

type ClientMsg struct {
	Code int
	Data string
	Tab  map[string]string
}
