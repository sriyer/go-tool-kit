package tabulator

import (
	"text/tabwriter"
	"io"
	"fmt"
	"reflect"
	"strings"
	"strconv"
	"encoding/json"
)

const tableTag = "table"

type tableTagTabulator struct {
	*tabwriter.Writer
}

func DefaultTableTagTabulator(out io.Writer) Tabulator {
	w := new(tabwriter.Writer)
	w.Init(out, 0, 0, 2, ' ', 0)
	return &tableTagTabulator{w}
}

func (w tableTagTabulator) Tabulate(s []interface{}) error {
	return w.write(s)	
}

func (w tableTagTabulator) write(s []interface{}) error {
	if len(s) < 1 {
		return nil
	}
	headers := orderHeaders(s[0])
	fmt.Fprintln(w, formatHeaders(headers, s[0]))
	fmt.Fprintln(w, formatHeaderUnderscore(headers, s[0]))
	for _, r := range s {
		fmt.Fprintln(w, formatData(r, headers))
	}
	w.Flush()
	return nil
}

func orderHeaders(s interface{}) []string {
	m := toMap(s)
	t := reflect.TypeOf(s)
	h := make([]string, t.NumField(),t.NumField())
	var unindexed []string
	var headers []string

	for i := 0 ; i < t.NumField() ; i++ {
		f := t.Field(i)
		k := f.Name
		if _, ok := m[k]; !ok {continue}
		if isTagged(f) {

			label, pos := getLabelAndPosition(f)
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

func formatHeaders(headers []string, r interface{}) string {
	s := ""
	t := reflect.TypeOf(r)
	for _, h := range headers {
		l := h
		if f, ok := t.FieldByName(h); ok {
			l, _ = getLabelAndPosition(f)
			if l == "" {
				l = h
			}
		}
		s = s + l + "\t"
	}
	return s
}

func formatHeaderUnderscore(headers []string, r interface{}) string {
	s := ""
	t := reflect.TypeOf(r)
	for _, h := range headers {
		l := h
		if f, ok := t.FieldByName(h); ok {
			l, _ = getLabelAndPosition(f)
			if l == "" {
				l = h
			}
		}
		s = s + fmt.Sprintf("%0*s", len(l) + 1, "\t" )
	}
	s = strings.Replace(s, "0", "-", -1)
	return s
}

func formatData(r interface{}, headers []string) string {
	m := toMap(r)
	s := ""
	for _, k := range headers {
		s = s + fmt.Sprint(m[k]) + "\t"
	}
	return s
}

func isTagged(f reflect.StructField) bool {
	_, present := f.Tag.Lookup(tableTag)
	return present
}


func getLabelAndPosition(f reflect.StructField) (string, int) {
	pos := -1
	var err error
	label := ""
	if tag := f.Tag.Get(tableTag); tag != "" {
		values := strings.Split(tag, ",")
		if len(values) > 0 {
			label = values[0]
			if len(values) > 1 {
				pos, err = strconv.Atoi(strings.TrimSpace(values[1]))
				if err != nil {
					pos = -1
				}
			}
		}
	}
	return label, pos
}

func toMap(i interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	b, err := json.Marshal(i)
	if err != nil {
		fmt.Println("[ERROR] ", err)
		return m
	}
	err = json.Unmarshal(b, &m)
	if err != nil {
		fmt.Println("[ERROR] ", err)
		return make(map[string]interface{})
	}
	return m
}