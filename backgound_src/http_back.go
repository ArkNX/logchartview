package main

import (
		"net/http"
		"encoding/json"
	  "fmt"
		"os"
		"strings"
)

//go run http_back.go

type logconfig struct {
	Loginfo [] struct {
		Logpath      string
		Filename     string
		Readcount    int
	}
	ListenInfo struct {
		Urlinterface string
		Address      string
		Port         string
	}
}

type json_data_info struct {
	Name string
	X []string
	Data []string
}

type handle_info struct {
		Logpath      string
		Filename     string	
		Readcount    int
}

type read_chart_handle struct{
	curr_handle_info []handle_info 
}

func (h *read_chart_handle)Do_Chart(file_name string, Readcount int, json_data *json_data_info) bool {
	text_line, _ := tail(file_name, Readcount)
	
	line_count := len(text_line)
	i := 0
	json_data.X = make([]string, line_count)
	json_data.Data = make([]string, line_count)

	for _, line := range text_line {
		//format json
		fields := strings.Split(line, "|")

		if(len(fields) < 3) {
			//err data
			continue
		} else {
			json_data.Name = fields[0]
			json_data.X[i] = fields[1]
			json_data.Data[i] = strings.Replace(fields[2], "\r", "", -1)
			i++
		}
	}

	return true;
}

func (h *read_chart_handle) Do_Json_String(json_data_list []json_data_info) string {
	response_data, err := json.Marshal(json_data_list)
	 if err != nil { 
		 fmt.Println(err);
		 return ""
	 }

	 return string(response_data)
}

func (h *read_chart_handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var json_data_list []json_data_info
	json_data_list = make([]json_data_info, len(h.curr_handle_info))

	for i,x:= range h.curr_handle_info {
		fmt.Println("[ServeHTTP]", h.curr_handle_info[i].Logpath)
		file_name := x.Logpath + x.Filename

		h.Do_Chart(file_name, x.Readcount, &json_data_list[i])
	}

	w.Write([]byte(h.Do_Json_String(json_data_list)))
}

func read_config(currlogconfig *logconfig) bool {
	file, _ := os.Open("logconfig.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(currlogconfig)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	return true
}

func main() {
		currlogconfig := logconfig {}

		read_ret := read_config(&currlogconfig)

		if(read_ret == false) {
			fmt.Println("read config fail.")
			return
		}

		listenaddr := currlogconfig.ListenInfo.Address + ":" + currlogconfig.ListenInfo.Port

		fmt.Println("listenaddr=", listenaddr)
		fmt.Println("Urlinterface=", currlogconfig.ListenInfo.Urlinterface)

		chart_handler := read_chart_handle{}

		line_count := len(currlogconfig.Loginfo)
		chart_handler.curr_handle_info = make([]handle_info, line_count)

		for i := 0; i < line_count; i++ {
			chart_handler.curr_handle_info[i].Logpath = currlogconfig.Loginfo[i].Logpath
			chart_handler.curr_handle_info[i].Filename = currlogconfig.Loginfo[i].Filename
			chart_handler.curr_handle_info[i].Readcount = currlogconfig.Loginfo[i].Readcount
		}

    http.Handle(currlogconfig.ListenInfo.Urlinterface, &chart_handler)
		http.ListenAndServe(listenaddr, nil)
}
