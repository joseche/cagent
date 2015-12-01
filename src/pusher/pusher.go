package pusher

import (
	"bytes"
	"io/ioutil"
	"misc"
	"net/http"
	"time"
	"encoding/json"
	"strings"
	"fmt"
)

const (
	PUSH_URL        string        = "http://localhost:3000/api/host/loadavg"
	PUSHER_INTERVAL time.Duration = time.Second * 5
	API_TOKEN       string        = "c804190ce588faefae423bbda1ae3926"
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

func push_loadavg() {
	conn := misc.OpenConn()
	if conn == nil {
		misc.Err("push_loadavg: Local db not found")
	} else {
		//-=-=-=-=-=-=- Start of collector cycle
		
		// extract data, create json request
		loadavg := ""
		query := "SELECT id,hostid,unixtime,load1,load5,load15 FROM " + misc.LOADAVG_TB + " LIMIT 100"
		stmt, err := conn.Query(query)
		if err != nil {
			misc.Err("Selecting data from  " + misc.LOADAVG_TB + " table: " + err.Error())
		} else {
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
					if len(loadavg) > 1 {
						loadavg += ","
					}
					loadavg += "{\"loadid\":" + string(id) + ",\"host_sig\":\"" + host_sig +
						"\",\"dt\":\"" + dt + "\",\"load1\":" + load1 +
						",\"load5\":" + load5 + ",\"load15\":" + load15 + "}"
				}
			}
		}
		stmt.Close()
		conn.Close()
		
		if len(loadavg)<10{
			misc.Debug("No loadavg records to send")
			return
		}
		
		
		var jsonstr string
		jsonstr = "{\"signature\": \"" + misc.Host + "\"," +
			"\"data\": {" +
			"\"loadavg\": [" + loadavg +
			"]}}"
			
		// json request created	
		// send http post request	
		req, err := http.NewRequest("POST", PUSH_URL, bytes.NewReader([]byte(jsonstr)))
		if err != nil {
			misc.Err("http request failed: " + err.Error())
		}

		req.Header.Set("Authentication", API_TOKEN)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			misc.Err("http request failed: " + err.Error())
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		misc.Debug("Response Status:" + resp.Status)
		
		if resp.StatusCode != 200 {
			misc.Err("http post request failed, status code: "+resp.Status)
			return
		}
		
		// http post successful, check response fields
		
		// convert successful ids to string for query
		var data map[string][]interface{}
		err = json.Unmarshal(body, &data)	
		if err != nil {
			misc.Err("Can not unmarshal json: "+err.Error())
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
		misc.Debug("Deleting submitted ids: "+ids_str)
		querystr := "DELETE FROM "+misc.LOADAVG_TB+" WHERE id IN ("+ids_str+")"
		conn = misc.OpenConn()
		err = conn.Exec(querystr)
		conn.Close()
		if err != nil {
			misc.Err("Deleting submitted ids: "+err.Error())
			return
		}
		
		// print errors from server
		errs := data["error_msgs"]
		for e := range errs {
			str := fmt.Sprintf("%v", errs[e] )
			misc.Err(str)
		}
		
		// last errored records
		ids = data["error_ids"]
		ids_str = ""
		for i := range ids {
			ids_str += fmt.Sprintf("%v,", ids[i])
		}
		ids_str = strings.TrimRight(ids_str, ",")
		misc.Debug("Deleting errored ids: "+ids_str)
		querystr = "DELETE FROM "+misc.LOADAVG_TB+" WHERE id IN ("+ids_str+")"
		conn = misc.OpenConn()
		err = conn.Exec(querystr)
		conn.Close()
		if err != nil {
			misc.Err("Deleting errored ids: "+err.Error())
			return
		}
		
		//-=-=-=-=-=-=- End of collector cycle
	}
}

func pusher_task() {
	push_loadavg()
}
