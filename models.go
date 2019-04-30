package main

import (
	"strings"
	"sync"
	"time"
)

type mapProjects map[string]*project

type CacheProjects struct {
	Cycle    int64
	Projects mapProjects
	Mu       sync.RWMutex
}

func (c *CacheProjects) Set(k string, f *project) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.Projects[k] = f
}

func (c *CacheProjects) Get(k string) (*project, bool) {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	found, ok := c.Projects[k]
	return found, ok
}

func (c *CacheProjects) Keys() []string {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	keys := []string{}
	for k := range c.Projects {
		keys = append(keys, k)
	}
	return keys
}

func (c *CacheProjects) Update(newItems mapProjects) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.Cycle = time.Now().UnixNano()
	c.Projects = newItems
}

func (c *CacheProjects) Length() int {
	return len(c.Projects)
}

func (c *CacheProjects) LastCycle() int64 {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	return c.Cycle
}

type filemMetaDatas []*fileMetaData

type project struct {
	Name           string
	fileMetaDatass filemMetaDatas
}

func (p *project) Add(f *fileMetaData) {
	p.fileMetaDatass = append(p.fileMetaDatass, f)
}

func (p project) getLatest() *fileMetaData {
	pMData := p.fileMetaDatass
	amount := len(pMData)
	latestProject := pMData[0]
	for i := 1; i < amount; i++ {
		if pMData[i].TimeStamp.After(latestProject.TimeStamp) {
			latestProject = pMData[i]
		}
	}
	return latestProject
}

type fileMetaData struct {
	ProjectName string
	Filename    string
	TimeStamp   time.Time
	Path        string
}

func parsePathString(s string) (string, string, time.Time, error) {
	items := strings.Split(s, "/")
	filename := items[len(items)-1]
	projectDate := strings.Split(filename, "_")

	strDate := strings.Split(projectDate[len(projectDate)-1], ".")[0]
	t, err := time.Parse("2006-01-02", strDate)

	projectName := strings.Join(projectDate[:len(projectDate)-1], "_")
	return projectName, filename, t, err
}

func parseContainerName(s string) *fileMetaData {
	projectName, filename, timeStamp, err := parsePathString(s)
	if err != nil {
		panic(err)
	}
	return &fileMetaData{
		ProjectName: projectName,
		Filename:    filename,
		TimeStamp:   timeStamp,
		Path:        s,
	}
}

type projectBlackList map[string]bool
const projectBlackListDelimiter = ","

func createProjectBlackList(s string) projectBlackList {
	m := make(projectBlackList)
	if len(s) == 0 {
		return m
	}
	for _, item := range strings.Split(s, projectBlackListDelimiter) {
		m[item] = true
	}
	return m
}

func (b projectBlackList) IsBlackListed(s string) bool {
	return b[s]
}
