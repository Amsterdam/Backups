package main

import (
	"fmt"
	"github.com/ncw/swift"
)

func reconnectObjStore(userName, apiKey, authURL, tenant string) error {
	// Create a connection
	connObj = swift.Connection{
		UserName: userName,
		ApiKey:   apiKey,
		AuthUrl:  authURL,
		Tenant:   tenant, // Name of the tenant (v2 auth only)
	}
	// Authenticate
	err := connObj.Authenticate()
	if err != nil {
		return err
	}
	return nil
}

func syncCache(userName, apiKey, authURL, tenant, containerName, pathName string, blacklist projectBlackList) {
	objects := make([]string, 0)
	opt := &swift.ObjectsOpts{Path: pathName}
	err := connObj.ObjectsWalk(containerName, opt, func(opts *swift.ObjectsOpts) (interface{}, error) {
		newObjects, err := connObj.ObjectNames(containerName, opts)
		if err == nil {
			objects = append(objects, newObjects...)
		}
		return newObjects, err
	})
	if err != nil {
		err = reconnectObjStore(userName, apiKey, authURL, tenant)
		if err != nil {
			// Unable to reconnect to Obj Store fail hard.
			panic(err)
		}
	}
	newCache := make(mapProjects)
	//var objLast *fileMetaData
	for _, obj := range objects {

		objMeta := parseContainerName(obj)
		for k, v := range blacklist {
			fmt.Sprintln(k, v)
		}
		if blacklist.IsBlackListed(objMeta.ProjectName) {
			continue
		}

		if newCache[objMeta.ProjectName] == nil {
			newCache[objMeta.ProjectName] = &project{Name: objMeta.ProjectName}
		}

		newCache[objMeta.ProjectName].Add(objMeta)

		fmt.Println(objMeta.ProjectName, objMeta.Filename)
		/*
			a, err := connObj.ObjectGetBytes(containerName, objLast.Path)
			if err != nil {
				fmt.Println(err)
			}
			err = ioutil.WriteFile("blabla" + strconv.Itoa(i), a, 0644)
			if err != nil {
				fmt.Println(err)
			}
		*/
	}
	Cache.Update(newCache)
}
