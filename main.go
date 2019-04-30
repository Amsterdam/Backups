package main

// Backups

import (
	"fmt"
	"github.com/ncw/swift"
	"log"
	"net/http"
	"time"
)

// Globals
var Cache *CacheProjects
var connObj swift.Connection

func init() {
	// Take input from default, env or cmd line, in that order.
	SETTINGS.Set("OBJUSERNAME", "backup_objectstore", "set username for object store")
	SETTINGS.Set("OBJAPIKEY", "redacted", "Set api key for object store")
	SETTINGS.Set("OBJAUTHURL", "redacted", "Set auth url for objectstore")

	SETTINGS.Set("OBJCONTAINER", "dbbackup", "Set name of the container from object store")
	SETTINGS.Set("OBJPATHNAME", "postgres/weekly", "Set Path filter")
	SETTINGS.Set("OBJTENANT", "redacted", "Set Tenant name")
	SETTINGS.SetInt("OBJSYNC", 60, "Pauze between syncs with object store in seconds")

	SETTINGS.Set("BLACKLIST", "", "Projects that are blacklisted from results. Add projects with ',' as delimiter")

	SETTINGS.Set("HOST", "0.0.0.0:8000", "enter host with port")

	// create Global Cache
	Cache = &CacheProjects{
		Projects: make(mapProjects),
	}

	//TODO reconnect/authenticate after connection is lost
	connObj = swift.Connection{
		UserName: SETTINGS.Get("OBJUSERNAME"),
		ApiKey:   SETTINGS.Get("OBJAPIKEY"),
		AuthUrl:  SETTINGS.Get("OBJAUTHURL"),
		Tenant:   SETTINGS.Get("OBJTENANT"),
	}
	// Authenticate
	err := connObj.Authenticate()
	if err != nil {
		panic(err)
	}
}

func cacheWorker() {
	userName := SETTINGS.Get("OBJUSERNAME")
	apiKey := SETTINGS.Get("OBJAPIKEY")
	authURL := SETTINGS.Get("OBJAUTHURL")
	tenant := SETTINGS.Get("OBJTENANT")
	containerName := SETTINGS.Get("OBJCONTAINER")
	pathName := SETTINGS.Get("OBJPATHNAME")
	blacklist := createProjectBlackList(SETTINGS.Get("BLACKLIST"))
	fmt.Println("blacklisted items", blacklist)
	pauzeSync := SETTINGS.GetInt("OBJSYNC")
	for {
		syncCache(userName, apiKey, authURL, tenant, containerName, pathName, blacklist)
		time.Sleep(time.Duration(pauzeSync) * time.Second)
	}
}

func main() {
	host := SETTINGS.Get("HOST")
	// Start cache worker to sync object store projects with the cache
	go cacheWorker()
	fmt.Println("started sync worker with seconds", SETTINGS.GetInt("OBJSYNC"), "pauze")
	fmt.Println("start server under host", host)

	// Register funcs and start server
	http.HandleFunc("/", listRest)
	http.HandleFunc("/item/", getRest)
	if err := http.ListenAndServe(host, nil); err != nil {
		log.Fatal(err)
	}
}
