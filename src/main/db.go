package main

import (
  "github.com/mxk/go-sqlite/sqlite3"
  "os"
)

type Table struct {
	tname string
	tcreate string
}

const (
  DB_FILE string = MASTER_DIR + "/pending"
  
  CPUTIMES_TB string = "CpuTimes"
  LOADAVG_TB string = "LoadAVG"
  MEMORY_TB string = "Memory"
  
  CPUTIMES_CREATE string = `CREATE TABLE `+CPUTIMES_TB+`(
    id INTEGER PRIMARY KEY ASC,
    cpuid TEXT,
    unixtime INTEGER,
    user REAL,
    sys REAL,
    idle REAL
    )`
  
  LOADAVG_CREATE string = `CREATE TABLE `+LOADAVG_TB+`(
    id INTEGER PRIMARY KEY ASC,
    unixtime INTEGER,
    load1 REAL,
    load5 REAL,
    load15 REAL
    )`

  MEMORY_CREATE string = `CREATE TABLE `+MEMORY_TB+`(
    id INTEGER PRIMARY KEY ASC,
    unixtime INTEGER,
    total_ram INTEGER,
    ram_free INTEGER,
    ram_used_percent REAL,
    total_swap INTEGER,
    swap_free INTEGER,
    swap_used_percent REAL
    )`
   	
)

func create_table(tbname string, tbcreate string, conn *sqlite3.Conn) (created bool){ 
  tbcount := -1	
  created = false
  query := "SELECT count(*) as count FROM sqlite_master "+
           "WHERE type='table' AND name='"+tbname+"'"
  ret,err := conn.Query(query)
  if (err!=nil){
    Err("create_table error: '"+tbname+"': "+err.Error())
  }else{
    ret.Scan(&tbcount)
    ret.Close()
    if tbcount <= 0 {
      err := conn.Exec(tbcreate)
      if err != nil {
    	panic("Can't create table: "+tbname+", "+err.Error())
      }else{
    	Debug(tbname+" created")
    	created = true
      }
    }
  }
  return created
}

func init_db(){
	
  dir_exists,_ := File_exists(MASTER_DIR)
  if ! dir_exists {
  	err := os.MkdirAll(MASTER_DIR,0755)
  	if err != nil {
  		panic("Can't create "+MASTER_DIR+", "+err.Error())
  	}
  	Info(MASTER_DIR +" created")
  }
	
  file_exists,_ := File_exists(DB_FILE)
  if ! file_exists {
  	Info("Local DB doesn't exist, creating new")
  }
  
  conn, err := sqlite3.Open( DB_FILE )
  if err != nil {
    panic(err)
  }
  
  create_table(CPUTIMES_TB, CPUTIMES_CREATE, conn)
  create_table(LOADAVG_TB, LOADAVG_CREATE, conn)
  create_table(MEMORY_TB, MEMORY_CREATE, conn)
  
  err = conn.Close()
  if err != nil {
  	Err("Closing local DB: "+err.Error())
  }
}

func openConn() (*sqlite3.Conn) {
	file_exists,_ := File_exists(DB_FILE)
  	if ! file_exists {
  		init_db()
  	}
  	conn, err := sqlite3.Open( DB_FILE )
    if err !=nil {
    	Err("Unable to access local db: "+err.Error())
    	return nil
  	}
	return conn
}