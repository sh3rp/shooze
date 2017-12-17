package shooze

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var API_VERSION = 1

type Webservice struct {
	db     *gorm.DB
	engine *gin.Engine
}

func NewWebservice() Webservice {
	r := gin.Default()

	db, err := gorm.Open("sqlite3", "/tmp/shooze.db")
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Config{}, &ConfigParameter{}, &Schedule{}, &Probe{}, &Deploy{}, &Agent{})
	db.Model(&Config{}).Related(&ConfigParameter{})
	db.Model(&Probe{}).Related(&Config{})
	db.Model(&Probe{}).Related(&Schedule{})

	return Webservice{db, r}
}

func (w Webservice) Init(port int) Webservice {
	w.engine.GET(endpoint("config"), w.getConfigs)
	w.engine.GET(endpoint("config/:id"), w.getConfig)
	w.engine.POST(endpoint("config"), w.postConfig)
	w.engine.DELETE(endpoint("config/:id"), w.deleteConfig)

	w.engine.GET(endpoint("schedule"), w.getSchedules)
	w.engine.GET(endpoint("schedule/:id"), w.getSchedule)
	w.engine.POST(endpoint("schedule"), w.postSchedule)
	w.engine.DELETE(endpoint("schedule/:id"), w.deleteSchedule)

	w.engine.GET(endpoint("probe"), w.getProbes)
	w.engine.GET(endpoint("probe/:id"), w.getProbe)
	w.engine.POST(endpoint("probe"), w.postProbe)
	w.engine.DELETE(endpoint("probe/:id"), w.deleteProbe)

	w.engine.GET(endpoint("agent"), w.getAgents)
	w.engine.GET(endpoint("agent/:id"), w.getAgent)
	w.engine.POST(endpoint("agent"), w.postAgent)
	w.engine.DELETE(endpoint("agent/:id"), w.deleteAgent)

	w.engine.GET(endpoint("deploy"), w.getDeploys)
	w.engine.GET(endpoint("deploy/:id"), w.getDeploy)
	w.engine.POST(endpoint("deploy"), w.postDeploy)
	w.engine.DELETE(endpoint("deploy/:id"), w.deleteDeploy)

	go w.engine.Run(fmt.Sprintf(":%d", port))
	return w
}

func (w Webservice) getConfigs(c *gin.Context) {
	var configs []Config
	w.db.Preload("Parameters").Find(&configs)
	c.JSON(200, WSObject{0, "OK", configs})
}
func (w Webservice) getConfig(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(200, WSObject{Status: 1, Message: fmt.Sprintf("%v", err)})
	} else {
		var config Config
		w.db.Preload("Parameters").First(&config, id)
		c.JSON(200, WSObject{Status: 0, Message: "OK", Data: config})
	}
}
func (w Webservice) postConfig(c *gin.Context) {
	actionInt := c.PostForm("_action")
	action, err := strconv.Atoi(actionInt)

	if err == nil {
		postVars := c.Request.PostForm

		var parameters []ConfigParameter

		for key, val := range postVars {
			if key != "_action" {
				parameter := ConfigParameter{Key: key, Value: val[0]}
				if w.db.NewRecord(parameter) {
					w.db.Create(&parameter)
					parameters = append(parameters, parameter)
				}
				fmt.Printf("%+v\n", parameters)
			}
		}
		config := Config{Action: ActionType(action), Parameters: parameters}
		w.db.Create(&config)

		c.JSON(200, WSObject{Status: 0, Message: "OK", Data: config})
	} else {
		c.JSON(200, WSObject{Status: 1, Message: fmt.Sprintf("%v", err)})
	}
}
func (w Webservice) deleteConfig(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(200, WSObject{Status: 1, Message: fmt.Sprintf("%v", err)})
	} else {
		var config Config
		w.db.Delete(&config, id)
		c.JSON(200, WSObject{Status: 0, Message: "OK"})
	}
}

func (w Webservice) getSchedules(c *gin.Context) {
	var schedules []Schedule
	w.db.Find(&schedules)
	c.JSON(200, WSObject{0, "OK", schedules})
}
func (w Webservice) getSchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(200, WSObject{Status: 1, Message: fmt.Sprintf("%v", err)})
	} else {
		var schedule Schedule
		w.db.First(&schedule, id)
		c.JSON(200, WSObject{Status: 0, Message: "OK", Data: schedule})
	}
}
func (w Webservice) postSchedule(c *gin.Context) {
	label := c.PostForm("label")

	if label == "" {
		c.JSON(200, WSObject{Status: 2, Message: "Required parameter 'label' not supplied"})
		return
	}

	cron := c.PostForm("crontab")

	if cron == "" {
		c.JSON(200, WSObject{Status: 2, Message: "Required parameter 'cron' not supplied"})
		return
	}

	schedule := Schedule{Label: label, Cron: cron}
	w.db.Create(&schedule)

	c.JSON(200, WSObject{Status: 0, Message: "OK", Data: schedule})
}
func (w Webservice) deleteSchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(200, WSObject{Status: 1, Message: fmt.Sprintf("%v", err)})
	} else {
		var schedule Schedule
		w.db.Delete(&schedule, id)
		c.JSON(200, WSObject{Status: 0, Message: "OK"})
	}
}

func (w Webservice) getProbes(c *gin.Context) {
	var probes []Probe
	w.db.Preload("Config").Preload("Config.Parameters").Preload("Schedule").Find(&probes)
	c.JSON(200, WSObject{0, "OK", probes})
}
func (w Webservice) getProbe(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(200, WSObject{Status: 1, Message: fmt.Sprintf("%v", err)})
	} else {
		var probe Probe
		w.db.First(&probe, id)
		c.JSON(200, WSObject{Status: 0, Message: "OK", Data: probe})
	}
}
func (w Webservice) postProbe(c *gin.Context) {
	configId, err := strconv.Atoi(c.PostForm("config_id"))

	if err != nil {
		c.JSON(200, WSObject{Status: 2, Message: "Required parameter 'config_id' not supplied"})
		return
	}

	var config Config

	w.db.First(&config, configId)

	if &config == nil {
		c.JSON(200, WSObject{Status: 3, Message: "No config found with id"})
		return
	}

	scheduleId, err := strconv.Atoi(c.PostForm("schedule_id"))

	if err != nil {
		c.JSON(200, WSObject{Status: 2, Message: "Required parameter 'schedule_id' not supplied"})
		return
	}

	var schedule Schedule

	w.db.First(&schedule, scheduleId)

	if &schedule == nil {
		c.JSON(200, WSObject{Status: 3, Message: "No schedule found with id"})
		return
	}

	probe := Probe{
		Config:   config,
		Schedule: schedule,
	}

	w.db.Create(&probe)

	c.JSON(200, WSObject{Status: 0, Message: "OK", Data: probe})
}
func (w Webservice) deleteProbe(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(200, WSObject{Status: 1, Message: fmt.Sprintf("%v", err)})
	} else {
		var probe Probe
		w.db.Delete(&probe, id)
		c.JSON(200, WSObject{Status: 0, Message: "OK"})
	}
}

func (w Webservice) getAgents(c *gin.Context) {
	var agents []Agent
	w.db.Find(&agents)
	c.JSON(200, WSObject{0, "OK", agents})
}
func (w Webservice) getAgent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(200, WSObject{Status: 1, Message: fmt.Sprintf("%v", err)})
	} else {
		var agent Agent
		w.db.First(&agent, id)
		c.JSON(200, WSObject{Status: 0, Message: "OK", Data: agent})
	}
}
func (w Webservice) postAgent(c *gin.Context) {
	label := c.PostForm("label")

	if label == "" {
		c.JSON(200, WSObject{Status: 2, Message: "Required parameter 'label' not supplied"})
		return
	}

	ip := c.PostForm("ip")

	if ip == "" {
		c.JSON(200, WSObject{Status: 2, Message: "Required parameter 'ip' not supplied"})
		return
	}

	agent := Agent{Label: label, IPAddress: ip}
	w.db.Create(&agent)

	c.JSON(200, WSObject{Status: 0, Message: "OK", Data: agent})
}
func (w Webservice) deleteAgent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(200, WSObject{Status: 1, Message: fmt.Sprintf("%v", err)})
	} else {
		var agent Agent
		w.db.Delete(&agent, id)
		c.JSON(200, WSObject{Status: 0, Message: "OK"})
	}
}

func (w Webservice) getDeploys(c *gin.Context)   {}
func (w Webservice) getDeploy(c *gin.Context)    {}
func (w Webservice) postDeploy(c *gin.Context)   {}
func (w Webservice) deleteDeploy(c *gin.Context) {}

func endpoint(name string) string {
	return fmt.Sprintf("v%d/%s", API_VERSION, name)
}
