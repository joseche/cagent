package pusher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"misc"
	"net/http"
	"strings"
	"time"
)

const (
	PUSH_SERVER      string        = "http://localhost:3000/"
	PUSH_LOADAVG_URL string        = PUSH_SERVER + "api/host/loadavg"
	PUSH_CPUTIME_URL string        = PUSH_SERVER + "api/host/cputime"
	PUSHER_INTERVAL  time.Duration = time.Second * 10
	API_TOKEN        string        = "c804190ce588faefae423bbda1ae3926"
)

func Pusher_routine(done chan bool) {
	ticker := misc.UpdateTicker(PUSHER_INTERVAL, "pusher")
	for misc.Running {
		<-ticker.C
		pusher_task()
		ticker = misc.UpdateTicker(PUSHER_INTERVAL, "pusher")
	}
	done <- true
}

func process_push_request(tablename string, jsonstr string, url string) {
	req, err := http.NewRequest("POST", url, bytes.NewReader([]byte(jsonstr)))
	if err != nil {
		misc.Err("http request failed: " + err.Error())
		return
	}

	req.Header.Set("Authentication", API_TOKEN)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		misc.Err("http request failed: " + err.Error())
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	misc.Debug("Response Status:" + resp.Status)

	if resp.StatusCode != 200 {
		misc.Err("http post request failed, status code: " + resp.Status)
		return
	}

	// http post successful, check response fields

	// convert successful ids to string for query
	var data map[string][]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		misc.Err("Can not unmarshal json: " + err.Error())
		misc.Err("Response Body:" + string(body))
		return
	}

	ids := data["successful_ids"]

	// couldn't find a join in golang...
	ids_str := ""
	for i := range ids {
		ids_str += fmt.Sprintf("%v,", ids[i])
	}
	ids_str = strings.TrimRight(ids_str, ",")
	misc.Debug("Deleting submitted ids: " + ids_str)
	querystr := "DELETE FROM " + tablename + " WHERE id IN (" + ids_str + ")"
	conn := misc.OpenConn()
	err = conn.Exec(querystr)
	conn.Close()
	if err != nil {
		misc.Err("Deleting submitted ids: " + err.Error())
		return
	}

	// print errors from server
	errs := data["error_msgs"]
	for e := range errs {
		str := fmt.Sprintf("%v", errs[e])
		misc.Err(str)
	}

	// last errored records
	ids = data["error_ids"]
	ids_str = ""
	for i := range ids {
		ids_str += fmt.Sprintf("%v,", ids[i])
	}

	if len(ids_str) > 2 {
		ids_str = strings.TrimRight(ids_str, ",")
		misc.Debug("Deleting errored ids: " + ids_str)
		querystr = "DELETE FROM " + tablename + " WHERE id IN (" + ids_str + ")"
		conn = misc.OpenConn()
		err = conn.Exec(querystr)
		conn.Close()
		if err != nil {
			misc.Err("Deleting errored ids: " + err.Error())
		}
	}
}

func push_loadavg() {
	conn := misc.OpenConn()
	if conn == nil {
		misc.Err("push_loadavg: Local db not found")
		return
	}

	//-=-=-=-=-=-=- Start of collector cycle
	// extract data, create json request
	loadavg := ""
	query := "SELECT id,hostid,unixtime,load1,load5,load15 FROM " + misc.LOADAVG_TB + " LIMIT 100"
	stmt, err := conn.Query(query)
	if err != nil && err.Error() != "EOF" {
		misc.Err("Selecting data from  " + misc.LOADAVG_TB + " table: " + err.Error())
		conn.Close()
		return
	}

	for ; err == nil; err = stmt.Next() {
		var id string
		var host_sig string
		var dt string
		var load1 string
		var load5 string
		var load15 string
		err := stmt.Scan(&id, &host_sig, &dt, &load1, &load5, &load15)
		if err != nil {
			misc.Err("Invalid data in table: " + err.Error())
		} else {
			loadavg += "{\"loadid\":" + id + ",\"host_sig\":\"" + host_sig +
				"\",\"dt\":\"" + dt + "\",\"load1\":" + load1 +
				",\"load5\":" + load5 + ",\"load15\":" + load15 + "},"
		}
	}
	stmt.Close()
	conn.Close()

	if len(loadavg) < 10 {
		misc.Debug("No loadavg records to send")
		return
	}
	var jsonstr string
	jsonstr = "{\"signature\": \"" + misc.Host + "\"," +
		"\"data\": {" +
		"\"loadavg\": [" + strings.TrimRight(loadavg, ",") +
		"]}}"
	process_push_request(misc.LOADAVG_TB, jsonstr, PUSH_LOADAVG_URL)
	//-=-=-=-=-=-=- End of collector cycle
}

func push_cputimes() {
	conn := misc.OpenConn()
	if conn == nil {
		misc.Err("push_cputimes: Local db not found")
		return
	}

	//-=-=-=-=-=-=- Start of collector cycle
	// extract data, create json request
	cputimes := ""
	query := "SELECT id,hostid,cpuid,unixtime,user,sys,idle,nice,iowait,irq,softirq,steal,guest,guest_nice,stolen FROM " + misc.CPUTIMES_TB + " LIMIT 100"
	stmt, err := conn.Query(query)
	if err != nil && err.Error() != "EOF" {
		misc.Err("Selecting data from  " + misc.CPUTIMES_TB + " table: " + err.Error())
		conn.Close()
		return
	}

	for ; err == nil; err = stmt.Next() {
		var id string
		var host_sig string
		var cpuid string
		var dt string
		var user string
		var sys string
		var idle string
		var nice string
		var iowait string
		var irq string
		var softirq string
		var steal string
		var guest string
		var guest_nice string
		var stolen string
		err := stmt.Scan(&id, &host_sig, &cpuid, &dt, &user, &sys, &idle, &nice, &iowait, &irq, &softirq, &steal, &guest, &guest_nice, &stolen)
		if err != nil {
			misc.Err("Invalid data in table: " + err.Error())
		} else {
			cputimes += "{\"cputid\":" + id + "," +
				"\"host_sig\":\"" + host_sig + "\"," +
				"\"cpuname\":\"" + cpuid + "\"," +
				"\"dt\":\"" + dt + "\"," +
				"\"user\":" + user + "," +
				"\"sys\":" + sys + "," +
				"\"idle\":" + idle + "," +
				"\"nice\":" + nice + "," +
				"\"iowait\":" + iowait + "," +
				"\"irq\":" + irq + "," +
				"\"softirq\":" + softirq + "," +
				"\"steal\":" + steal + "," +
				"\"guest\":" + guest + "," +
				"\"guest_nice\":" + guest_nice + "," +
				"\"stolen\":" + stolen + "},"

		}
	}
	stmt.Close()
	conn.Close()

	if len(cputimes) < 10 {
		misc.Debug("No cputimes records to send")
		return
	}
	var jsonstr string
	jsonstr = "{\"signature\": \"" + misc.Host + "\"," +
		"\"data\": {" +
		"\"cputime\": [" + strings.TrimRight(cputimes, ",") +
		"]}}"
	process_push_request(misc.CPUTIMES_TB, jsonstr, PUSH_CPUTIME_URL)
	//-=-=-=-=-=-=- End of collector cycle
}

func pusher_task() {
	push_loadavg()
	push_cputimes()
}
