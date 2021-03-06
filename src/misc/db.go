package misc

import (
	"github.com/mxk/go-sqlite/sqlite3"
)

type Table struct {
	tname string
	tcreate string
}

const (
  DEBUG      bool   = true
  MASTER_DIR string = "/opt/clarity/"
	
  DB_FILE string = MASTER_DIR + "/pending"
  
  LOADAVG_TB string = "LoadAVG"
  CPUTIMES_TB string = "CpuTimes"
  VIRTUAL_MEMORY_TB string = "VirtualMemory"
  SWAP_MEMORY_TB string = "SwapMemory"
  
  LOADAVG_CREATE string = `CREATE TABLE `+LOADAVG_TB+`(
    id INTEGER PRIMARY KEY ASC,
	hostid TEXT,
    unixtime INTEGER,
    load1 REAL,
    load5 REAL,
    load15 REAL
    )`
  
  CPUTIMES_CREATE string = `CREATE TABLE `+CPUTIMES_TB+`(
    id INTEGER PRIMARY KEY ASC,
    hostid TEXT,
    cpuid TEXT,
    unixtime INTEGER,
    user REAL,
    sys REAL,
    idle REAL,
    nice REAL, 
    iowait REAL,
    irq REAL,
    softirq REAL,
    steal REAL,
    guest REAL,
    guest_nice REAL,
    stolen REAL
    )`
  
  VIRTUAL_MEMORY_CREATE string = `CREATE TABLE `+VIRTUAL_MEMORY_TB+`(
    id INTEGER PRIMARY KEY ASC,
    hostid TEXT,
    unixtime INTEGER,
    total INTEGER,
    available INTEGER,
    used INTEGER,
    usedpercent REAL,
    free INTEGER,
    active INTEGER,
    inactive INTEGER,
    buffers INTEGER,
    cached INTEGER,
    wired INTEGER,
    shared INTEGER
    )`
  
  SWAP_MEMORY_CREATE string = `CREATE TABLE `+SWAP_MEMORY_TB+`(
    id INTEGER PRIMARY KEY ASC,
    hostid TEXT,
    unixtime INTEGER,
    total INTEGER,
    used INTEGER,
    free INTEGER,
    usedpercent REAL,
    sin INTEGER,
    sout INTEGER
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


func OpenConn() (*sqlite3.Conn) {
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