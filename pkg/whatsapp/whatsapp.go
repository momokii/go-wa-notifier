package whatsapp

import (
	"context"
	"fmt"
	"os"
	"sync"

	_ "github.com/lib/pq" // Import PostgreSQL driver
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
	_ "modernc.org/sqlite" // Import SQLite driver and with this can use SQLite without cgo enabled = 1
)

type whatsApp struct {
	client  *whatsmeow.Client
	qrCode  string
	qrReady bool
	mutex   sync.RWMutex
}

// Singleton instance
var (
	instance *whatsApp
	once     sync.Once
	initErr  error
	mutex    sync.Mutex // Mutex to protect instance reset
)

// ResetInstance clears the singleton instance
func ResetInstance() {
	mutex.Lock()
	defer mutex.Unlock()
	instance = nil
	once = sync.Once{}
}

// NewWhatsApp returns a singleton WhatsApp instance
func NewWhatsApp() (*whatsApp, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if instance == nil {
		instance, initErr = initWhatsApp()
	}
	return instance, initErr
}

// initWhatsApp creates a new WhatsApp client instance
func initWhatsApp() (*whatsApp, error) {
	pgHost := os.Getenv("HOST_POSTGRES")
	pgPort := os.Getenv("PORT_POSTGRES")
	pgUser := os.Getenv("USER_POSTGRES")
	pgPassword := os.Getenv("PASSWORD_POSTGRES")
	pgDatabase := os.Getenv("DATABASE_POSTGRES")

	pgConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pgHost, pgPort, pgUser, pgPassword, pgDatabase)

	// Initialize with PostgreSQL instead of SQLite
	container, err := sqlstore.New("postgres", pgConnString, waLog.Noop)

	// using SQLite instead of PostgreSQL
	// container, err := sqlstore.New("sqlite", "file:wapp.db?_pragma=foreign_keys(1)", waLog.Noop)

	if err != nil {
		return nil, err
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return nil, err
	}

	wa := &whatsApp{
		client:  whatsmeow.NewClient(deviceStore, waLog.Noop),
		qrCode:  "",
		qrReady: false,
		mutex:   sync.RWMutex{},
	}

	if wa.client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := wa.client.GetQRChannel(context.Background())
		err = wa.client.Connect()
		if err != nil {
			return nil, err
		}

		go func() {
			for evt := range qrChan {
				if evt.Event == "code" {
					// log.Println("QR Code received:", evt.Code)

					// Thread-safe update of QR code
					wa.mutex.Lock()
					wa.qrCode = evt.Code
					wa.qrReady = true
					wa.mutex.Unlock()

					// Log this important transition
					// log.Println("QR code is now ready for scanning")
				} else {
					if evt.Event == "success" {
						wa.mutex.Lock()
						wa.qrReady = false
						wa.mutex.Unlock()
						// log.Println("Login successful, QR code no longer needed")
					}
				}
			}
		}()
	} else {
		err := wa.client.Connect()
		if err != nil {
			return nil, err
		}
		// log.Println("WhatsApp client connected with existing session")
	}

	return wa, nil
}

func (w *whatsApp) GetClient() *whatsmeow.Client {
	return w.client
}

func (w *whatsApp) SendMessage(to, message string, with_disconnect bool) error {

	if w.client == nil {
		return fmt.Errorf("client is nil")
	}

	if w.client.Store.ID == nil {
		return fmt.Errorf("client is not connected")
	}

	ctx := context.Background()
	_, err := w.client.SendMessage(ctx, types.JID{
		User:   to, // Replace with the recipient's phone number
		Server: types.DefaultUserServer,
	}, &waE2E.Message{
		Conversation: proto.String(message),
	})
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	// if with_disconnect, so disconnect the client after sending the message
	if with_disconnect {
		defer w.client.Disconnect()
	}

	return nil
}

func (w *whatsApp) Disconnect() error {
	if w.client == nil {
		return fmt.Errorf("client is nil")
	}

	w.client.Disconnect()

	return nil
}

func (w *whatsApp) Logout() error {
	if w.client == nil {
		return fmt.Errorf("client is nil")
	}

	if w.client.Store.ID == nil {
		return fmt.Errorf("client is not connected")
	}

	// Check if client is connected
	if !w.client.IsConnected() {
		return fmt.Errorf("client is not connected")
	}

	// Perform logout
	err := w.client.Logout()
	if err != nil {
		return fmt.Errorf("error during logout: %v", err)
	}

	// Disconnect after logout
	w.client.Disconnect()

	// Reset the singleton instance to force reinitialization on next call
	ResetInstance()

	return nil
}

// GetQRCode returns the current QR code and a boolean indicating if it's ready
func (w *whatsApp) GetQRCode() (string, bool) {
	// Thread-safe read of QR code
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	return w.qrCode, w.qrReady
}

// IsConnected checks if the client is currently connected
func (w *whatsApp) IsConnected() bool {
	return w.client != nil && w.client.IsConnected() && w.client.Store.ID != nil
}
