// dispositivo.go
package packge

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"time"
)

type DeviceActive struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	IDEdge       int
	IDNodo       int
	Name         string
	IsActive     bool
	Notification bool
	mu           sync.Mutex
	UpdatedAt    time.Time
}

func NewDeviceActive(idEdge int, idNodo int, name string) *DeviceActive {
	return &DeviceActive{
		ID:       primitive.NewObjectID(),
		IDEdge:   idEdge,
		IDNodo:   idNodo,
		Name:     name,
		IsActive: false,
	}
}

func (d *DeviceActive) SetActive(active bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.IsActive = active
	d.UpdatedAt = time.Now()
	d.Notification = false
	if active {
		fmt.Printf("Dispositivo %s (IDNodo: %d) está ativo.\n", d.Name, d.IDNodo)
	} else {
		fmt.Printf("Dispositivo %s (IDNodo: %d) está offline.\n", d.Name, d.IDNodo)
	}

}

func (d *DeviceActive) SetNotif(notif bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.Notification = notif
	if notif {
		fmt.Printf("Dispositivo %s (IDNodo: %d) já foi notificacao.\n", d.Name, d.IDNodo)
	} else {
		fmt.Printf("Dispositivo %s (IDNodo: %d) ainda não foi notificacao.\n", d.Name, d.IDNodo)
	}
}

func (d *DeviceActive) SetActiveAndNotif(active bool, notif bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.IsActive = active
	d.Notification = notif
}

func (d *DeviceActive) GetStatus() bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.IsActive
}

func (d *DeviceActive) GetLastUpdateTime() time.Time {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.UpdatedAt
}

func SetDeviceStatus(deviceID int, active bool) error {
	// Find the device in the list
	for _, device := range listDeviceActive {
		if device.IDNodo == deviceID {
			// Set the status of the device
			device.SetActive(active)
			return nil
		}
	}

	return fmt.Errorf("dispositivo com ID %d não encontrado", deviceID)
}

func SetDeviceNotif(deviceID int, notif bool) error {
	// Find the device in the list
	for _, device := range listDeviceActive {
		if device.IDNodo == deviceID {
			// Set the status of the device
			device.SetNotif(notif)
			return nil
		}
	}

	return fmt.Errorf("dispositivo com ID %d não encontrado", deviceID)
}

// Helper function to check inactive devices every minute
func checkInactiveDevices() {

	for range time.Tick(time.Minute) {
		checkTime := time.Now().Add(-3 * time.Minute)
		//fmt.Printf("checkInactiveDevices\n")
		fmt.Println("----------------------------------------------------------------------------------------")
		for _, d := range listDeviceActive {
			d.mu.Lock()
			if d.IsActive {
				if !d.Notification && d.UpdatedAt.Before(checkTime) {
					fmt.Printf("Dispositivo %s (IDNodo: %d) está inativo por mais de 3 minutos.\n", d.Name, d.IDNodo)
					d.IsActive = false
				} else {
					fmt.Printf("Dispositivo ativo %s (IDNodo: %d) (Data: %s)\n", d.Name, d.IDNodo, d.UpdatedAt.Local().String())
				}
			} else {
				fmt.Printf("Dispositivo inativo %s (IDNodo: %d) (Data: %s)\n", d.Name, d.IDNodo, d.UpdatedAt.Local().String())
			}

			d.mu.Unlock()
		}
	}
}

var listDeviceActive []*DeviceActive

func SetDispositivos(d []*DeviceActive) {
	listDeviceActive = d

	// Start a background goroutine to check inactive devices
	go checkInactiveDevices()
}

func GetDispositivos() []*DeviceActive {
	return listDeviceActive
}
