package samsung

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	logger "log"
	"net"
	"time"
)

// EnableLogging allows you to turn on debug logging
var EnableLogging = true

var port = 55000

var errDenied = errors.New("Command denied. TV User has not allowed access")
var errWaitingForUser = errors.New("Waiting for TV user to grant or deny access")
var errTimeoutUser = errors.New("Timeout or cancelled by TV user")

// TV is a single television, identified by it's hostname/IP address.
type TV struct {
	Host            string // The hostname or IP address of the TV
	ApplicationID   string // This ID is used by the television to store the user's Allow/Deny response. If this is changed, it will ask again.
	ApplicationName string // ApplicationName is displayed on the screen the first time the ApplicationID is used.
}

// OnPowerChange allows you to monitor the on/off state. The provided callback will be fired whenever this TV goes online/offline.
func (tv *TV) OnPowerChange(interval time.Duration, callback func(bool)) {
	log("Monitoring power state of TV %s", tv.Host)

	address := fmt.Sprintf("%s:%d", tv.Host, port)

	var err error
	var conn net.Conn

	var lastState *bool

	for {

		// TODO: This should just be pinging, not connecting
		conn, err = net.DialTimeout("tcp", address, interval)

		online := err == nil

		if lastState == nil || *lastState != online {
			callback(online)
			lastState = &online
		}

		if err == nil {
			conn.Close()
			time.Sleep(interval)
		}

	}
}

// SendCommand sends a single named command to the TV.
// Even if a connection is made, an error can be returned indicating
// A list of possible commands (not all will be available on all models) is available at http://wiki.samygo.tv/index.php5/Key_codes
func (tv *TV) SendCommand(cmd string) error {

	log("Sending command %s to TV %s", cmd, tv.Host)

	appString := "iphone..iapp.samsung"
	tvAppString := "iphone.UN60D6000.iapp.samsung"

	header := header(tv.Host, tv.ApplicationID, tv.ApplicationName, appString)
	command := command(cmd, tvAppString)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", tv.Host, port))
	if err != nil {
		log("Failed to connect to tv: %s", err)
		return err
	}

	conn.Write(header)
	conn.Write(command)

	var response = make([]byte, 64)
	length, err := conn.Read(response)
	if err != nil {
		return err
	}

	err = readResponse(response[0:length])
	conn.Close()
	return err
}

func header(ip, mac, name, app string) []byte {
	msg := bytes.Buffer{}

	msg.WriteByte(0x64)
	msg.WriteByte(0x0)
	addB64(ip, &msg)
	addB64(mac, &msg)
	addB64(name, &msg)

	wrapped := bytes.Buffer{}

	wrapped.WriteByte(0x0)
	wrapped.WriteByte(uint8(len(app)))
	wrapped.WriteByte(0x0)
	wrapped.Write([]byte(app))
	wrapped.WriteByte(uint8(msg.Len()))
	wrapped.WriteByte(0x0)
	wrapped.Write(msg.Bytes())

	return wrapped.Bytes()
}

func command(command, app string) []byte {
	msg := bytes.Buffer{}

	msg.Write([]byte{0, 0, 0})
	addB64(command, &msg)

	wrapped := bytes.Buffer{}

	wrapped.WriteByte(0x0)
	wrapped.WriteByte(uint8(len(app)))
	wrapped.WriteByte(0x0)
	wrapped.Write([]byte(app))
	wrapped.WriteByte(uint8(msg.Len()))
	wrapped.WriteByte(0x0)
	wrapped.Write(msg.Bytes())

	return wrapped.Bytes()
}

func readResponse(bytes []byte) error {

	//unknown := bytes[0]
	nameLength := bytes[1] // It is actually little-endian uint16. Whatever.
	name := string(bytes[3 : nameLength+3])

	log("Response from name: %s", name)
	payloadLength := bytes[nameLength+3] // Also little-endian uint16.
	payload := bytes[nameLength+5 : nameLength+5+payloadLength]

	switch string(payload) {
	case string([]byte{0x64, 0x00, 0x01, 0x00}):
		print("Command sent successfully")
		return nil
	case string([]byte{0x64, 0x00, 0x00, 0x00}):
		return errDenied
	case string([]byte{0x0A, 0x00, 0x01, 0x00, 0x00, 0x00}):
		fallthrough
	case string([]byte{0x0A, 0x00, 0x02, 0x00, 0x00, 0x00}):
		return errWaitingForUser
	case string([]byte{0x65, 0x00}):
		return errTimeoutUser
	default:
		log("Failed to read response... assuming ok. %+x", payload)
		return nil
	}
}

func addB64(str string, msg *bytes.Buffer) {
	enc := []byte(base64.StdEncoding.EncodeToString([]byte(str)))

	msg.WriteByte(uint8(len(enc)))
	msg.WriteByte(0)
	msg.Write(enc)
}

func log(msg string, args ...interface{}) {
	if EnableLogging {
		logger.Printf("samsung-tv: "+msg, args...)
	}
}
