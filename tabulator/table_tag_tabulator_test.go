package tabulator

import (
	"fmt"
	"math/rand"
	"testing"
	"os"
)

func TestDefaultTabWriter(t *testing.T) {
	tab := DefaultTableTagTabulator(os.Stdout)
	tab.Tabulate(getReservations())
}

type Reservation struct {
	Host string `table:"Host name,0"`
	CPU  float64 `table:"Cores,1"`
	Memory float64 `table:"Memory in MB"`
	DockerDisk float64 `table:"Docker Disk"`
	DataDisk float64 `table:"Data Disk, 4"`
	Group string `table:"Node Group, 5"`
}

func getReservations() []interface{} {
	var reservations []interface{}

	for i := 0 ; i < 10 ; i++ {
		reservations = append(reservations, Reservation{
			Host: fmt.Sprint("Host", i),
			CPU: rand.Float64(),
			Memory: rand.Float64(),
			DockerDisk: rand.Float64(),
			DataDisk: rand.Float64(),
			Group: fmt.Sprint("Group", i),
		})
	}

	return reservations
}
