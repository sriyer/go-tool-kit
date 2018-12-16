package tabwriter

import (
	"fmt"
	"math/rand"
	"text/tabwriter"
	"os"
	"encoding/json"
	"log"
	"strings"
	"reflect"
	"strconv"
)

func main() {
	reservations := getReservations()
	fmt.Println("reservations : ", reservations)
	formatPrint(reservations)
}


type Reservation struct {
	Host string `table:"Host name,0"`
	CPU  float64 `table:"Cores,1"`
	Memory float64 `table:"Memory in MB"`
	DockerDisk float64 `table:"Docker Disk"`
	DataDisk float64 `table:"Data Disk, 4"`
	Group string `table:"Node Group, 5"`
}

func getReservations() []*Reservation {
	var reservations []*Reservation

	for i := 0 ; i < 10 ; i++ {
		reservations = append(reservations, &Reservation{
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

func formatPrint(reservations []*Reservation) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 2, ' ', 0)
	headers := OrderHeaders(*reservations[0])
	fmt.Fprintln(w, FormatHeaders(headers, *reservations[0]))
	fmt.Fprintln(w, FormatHeaderUnderscore(headers, *reservations[0]))
	for _, reservation := range reservations {
		fmt.Fprintln(w, FormatData(reservation, headers))
	}
	w.Flush()
}

func OrderHeaders(reservation Reservation) []string {
	m := toMap(reservation)
	t := reflect.TypeOf(reservation)
	h := make([]string, t.NumField(),t.NumField())
	var unindexed []string
	var headers []string

	for i := 0 ; i < t.NumField() ; i++ {
		f := t.Field(i)
		k := f.Name
		if _, ok := m[k]; !ok {continue}
		if isTagged(f) {

			label, pos := GetLabelAndPosition(f)
			if pos > -1 {
				if h[pos] != "" {
					unindexed = append(unindexed, k)
				} else {
					h[pos] = k
				}
			} else if label != "" {
				unindexed = append(unindexed, k)
			}
		}
	}

	for _, v := range h {
		if v == "" {
			continue
		}
		headers = append(headers, v)
	}
	headers = append(headers, unindexed...)
	return headers
}

func isTagged(f reflect.StructField) bool {
	_, present := f.Tag.Lookup("table")
	return present
}


func GetLabelAndPosition(f reflect.StructField) (string, int) {
	pos := -1
	var err error
	label := ""
	if tag := f.Tag.Get("table"); tag != "" {
		values := strings.Split(tag, ",")
		if len(values) > 0 {
			label = values[0]
			if len(values) > 1 {
				pos, err = strconv.Atoi(strings.TrimSpace(values[1]))
				if err != nil {
				//	log.Println("[ERROR] ", err)
					pos = -1
				}
			}
		}
	}
	return label, pos
}

func FormatData(reservation *Reservation, headers []string) string {
	m := toMap(reservation)
	s := ""
	for _, k := range headers {
		s = s + fmt.Sprint(m[k]) + "\t"
	}
	return s
}

func FormatHeaders(headers []string, reservation Reservation) string {
	s := ""
	t := reflect.TypeOf(reservation)
	for _, h := range headers {
		l := h
		if f, ok := t.FieldByName(h); ok {
			l, _ = GetLabelAndPosition(f)
			if l == "" {
				l = h
			}
		}
		s = s + l + "\t"
	}
	return s
}

func FormatHeaderUnderscore(headers []string, reservation Reservation) string {
	s := ""
	t := reflect.TypeOf(reservation)
	for _, h := range headers {
		l := h
		if f, ok := t.FieldByName(h); ok {
			l, _ = GetLabelAndPosition(f)
			if l == "" {
				l = h
			}
		}
		s = s + fmt.Sprintf("%0*s", len(l) + 1, "\t" )
	}
	s = strings.Replace(s, "0", "-", -1)
	return s
}

func GetHeaders(reservation *Reservation) []string {
	m := toMap(reservation)
	var s []string
	for k := range m {
		s = append(s, k)
	}
	return s
}

func toMap(i interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	b, err := json.Marshal(i)
	if err != nil {
		log.Println("[ERROR] ", err)
		return m
	}
	err = json.Unmarshal(b, &m)
	if err != nil {
		log.Println("[ERROR] ", err)
		return make(map[string]interface{})
	}
	return m
}


