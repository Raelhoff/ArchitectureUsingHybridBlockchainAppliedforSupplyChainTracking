// dispositivo.go
package packge

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"time"
)

type EdgeActive struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	IDEdge       int
	Name         string
	IsActive     bool
	Notification bool
	mu           sync.Mutex
	UpdatedAt    time.Time
}

func NewEdgeActive(idEdge int, name string) *EdgeActive {
	return &EdgeActive{
		ID:       primitive.NewObjectID(),
		IDEdge:   idEdge,
		Name:     name,
		IsActive: false,
	}
}

func (d *EdgeActive) SetActive(active bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.IsActive = active
	d.UpdatedAt = time.Now()
	d.Notification = false
	if active {
		fmt.Printf("Dispositivo (Edge) %s (ID: %d) está ativo.\n", d.Name, d.IDEdge)
	} else {
		fmt.Printf("Dispositivo (Edge) %s (ID: %d) está offline.\n", d.Name, d.IDEdge)
	}

}

func (d *EdgeActive) SetNotif(notif bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.Notification = notif
	if notif {
		fmt.Printf("Dispositivo (Edge) %s (ID: %d) já foi notificacao.\n", d.Name, d.IDEdge)
	} else {
		fmt.Printf("Dispositivo (Edge) %s (ID: %d) ainda não foi notificacao.\n", d.Name, d.IDEdge)
	}
}

func (d *EdgeActive) SetActiveAndNotif(active bool, notif bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.IsActive = active
	d.Notification = notif
}

func (d *EdgeActive) GetStatus() bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.IsActive
}

func (d *EdgeActive) GetLastUpdateTime() time.Time {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.UpdatedAt
}

func SetEdgeStatus(edgeID int, active bool) error {
	// Find the device in the list
	for _, device := range listEdgeActive {
		if device.IDEdge == edgeID {
			// Set the status of the device
			device.SetActive(active)
			return nil
		}
	}

	return fmt.Errorf("dispositivo edge com ID %d não encontrado", edgeID)
}

func SetEdgeNotif(deviceID int, notif bool) error {
	// Find the device in the list
	for _, device := range listEdgeActive {
		if device.IDEdge == deviceID {
			// Set the status of the device
			device.SetNotif(notif)
			return nil
		}
	}

	return fmt.Errorf("dispositivo com ID %d não encontrado", deviceID)
}

// Helper function to check inactive devices every minute
func checkInactiveEdges() {

	for range time.Tick(time.Minute) {
		checkTime := time.Now().Add(-5 * time.Minute)
		//fmt.Printf("checkInactiveDevices\n")
		fmt.Println("----------------------------------------------------------------------------------------")
		for _, d := range listEdgeActive {
			d.mu.Lock()
			if d.IsActive {
				if !d.Notification && d.UpdatedAt.Before(checkTime) {
					fmt.Printf("Dispositivo (Edge) %s (ID: %d) está inativo por mais de 5 minutos.\n", d.Name, d.IDEdge)
					d.IsActive = false
				} else {
					fmt.Printf("Dispositivo (Edge) ativo %s (ID: %d) (Data: %s)\n", d.Name, d.IDEdge, d.UpdatedAt.Local().String())
				}
			} else {
				fmt.Printf("Dispositivo inativo %s (ID: %d) (Data: %s)\n", d.Name, d.IDEdge, d.UpdatedAt.Local().String())
			}

			d.mu.Unlock()
		}
	}
}

var listEdgeActive []*EdgeActive

func SetEdges(d []*EdgeActive) {
	listEdgeActive = d

	// Start a background goroutine to check inactive devices
	go checkInactiveEdges()
}

func GetEdges() []*EdgeActive {
	return listEdgeActive
}
