package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

func main() {
	Key("key.txt")
	Req(os.Args[1])
}

func Key(filename string) {
	re := regexp.MustCompile("key=([A-z0-9_.]+);")
	b, _ := ioutil.ReadFile(filename)
	k := re.FindAllStringSubmatch(string(b), 1)[0][1]
	fmt.Println(k)
}

func Req(url string) {
	res, _ := http.Get(url)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	switch path.Ext(os.Args[1]) {
	case ".csv":
		r := csv.NewReader(strings.NewReader(string(body)))

		result, _ := r.ReadAll()

		data := []map[string]interface{}{}
		cols := result[0]

		for _, row := range result {
			m := map[string]interface{}{}
			for i, cell := range row {
				m[cols[i]] = cell
			}
			data = append(data, m)
		}

		j, _ := json.MarshalIndent(data, "", "\t")
		fmt.Println(string(j))

	case ".json":
		var v interface{}
		json.Unmarshal(body, &v)

		data := []map[string]interface{}{}

		for _, vv := range v.(map[string]interface{}) {
			meta := vv.(map[string]interface{})
			rows := meta["data"].([]interface{})
			cols := meta["column_names"].([]interface{})

			for _, row := range rows {
				m := map[string]interface{}{}
				for i, cell := range row.([]interface{}) {
					m[cols[i].(string)] = cell
				}
				data = append(data, m)
			}
		}

		j, _ := json.MarshalIndent(data, "", "\t")
		fmt.Println(string(j))

	default:
		fmt.Println(string(body))
	}

}
