package setting

import (
	"edetector_API/config"
	"edetector_API/pkg/errhandler"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/gin-gonic/gin"
)

type SettingData struct {
	ServerAndEmail serverAndEmail `json:"serverAndEmail"`
	Function       function       `json:"function"`
	Agent          agent          `json:"agent"`
	Interface      interface_     `json:"interface"`
}

type serverAndEmail struct {
	WorkerPort    int    `json:"workerPort"`
	DetectPort    int    `json:"detectPort"`
	DetectDefault bool   `json:"detectDefault"`
	UpdatePort    int    `json:"updatePort"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	IP            string `json:"ip"`
	Port          int    `json:"port"`
	Encryption    string `json:"encryption"`
}

type function struct {
	EmailNotification emailNotification `json:"emailNotification"`
	Report            report            `json:"report"`
	AutoDefense       autoDefense       `json:"autoDefense"`
}

type emailNotification struct {
	On             bool   `json:"on"`
	RiskLevel      string `json:"riskLevel"`
	RecepientEmail string `json:"recepientEmail"`
}

type report struct {
	On             bool   `json:"on"`
	GenerateTime   string `json:"generateTime"`
	RiskLevel      string `json:"riskLevel"`
	RecepientEmail string `json:"recepientEmail"`
}

type autoDefense struct {
	On                             bool   `json:"on"`
	RiskLevel                      string `json:"riskLevel"`
	Action                         string `json:"action"`
	MemoryDump                     bool   `json:"memoryDump"`
	DumpUploadVirusTotalAnalysis   bool   `json:"dumpUploadVirusTotalAnalysis"`
	DumpUploadHybridAnalysis       bool   `json:"dumpUploadHybridAnalysis"`
	SampleDownload                 bool   `json:"sampleDownload"`
	SampleUploadVirusTotalAnalysis bool   `json:"sampleUploadVirusTotalAnalysis"`
	SampleUploadHybridAnalysis     bool   `json:"sampleUploadHybridAnalysis"`
}

type agent struct {
	Mission       mission       `json:"mission"`
	GenerateAgent generateAgent `json:"generateAgent"`
	CPUPriority   string        `json:"cpuPriority"`
}

type mission struct {
	FileAnalysis         bool   `json:"fileAnalysis"`
	FileAnalysisMainDisk string `json:"fileAnalysisMainDisk"`
	MemoryScan           bool   `json:"memoryScan"`
	Collection           bool   `json:"collection"`
	GenerateImage        bool   `json:"generateImage"`
}

type generateAgent struct {
	ServerIP      string `json:"serverIP"`
	MemoryScan    bool   `json:"memoryScan"`
	Collection    bool   `json:"collection"`
	GenerateImage bool   `json:"generateImage"`
}

type interface_ struct {
	CEFLog      bool   `json:"cefLog"`
	RiskLevel   string `json:"riskLevel"`
	IP          string `json:"ip"`
	Port        int    `json:"port"`
	SignatureID string `json:"signatureID"`
}

func GetSetting() (SettingData, error) {
	var settings SettingData
	file, err := os.Open(config.Viper.GetString("SETTING_PATH"))
	if err != nil {
		return settings, err
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal([]byte(bytes), &settings)
	if err != nil {
		return settings, err
	}
	return settings, nil
}

func GetSettingField(c *gin.Context) {
	fieldName := c.Param("field")
	setting, err := GetSetting()
	if err != nil {
		errhandler.Error(c, err, "Error loading setting data")
		return
	}
	r := reflect.ValueOf(setting)
	field := r.FieldByName(fieldName)
	if !field.IsValid() {
		errhandler.Error(c, fmt.Errorf("setting field not found"), "Error checking field")
		return
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
		"data":      field.Interface(),
	})
}

func UpdateSetting(setting SettingData) error {
	bytes, err := json.MarshalIndent(setting, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(config.Viper.GetString("SETTING_PATH"), bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func UpdateSettingField(c *gin.Context) {
	fieldName := c.Param("field")
	setting, err := GetSetting()
	if err != nil {
		errhandler.Error(c, err, "Error loading setting data")
		return
	}
	r := reflect.ValueOf(&setting)
	field := r.Elem().FieldByName(fieldName)
	if !field.IsValid() {
		errhandler.Error(c, fmt.Errorf("setting field not found"), "Error checking field")
		return
	}
	req := reflect.New(field.Type()).Interface()
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Info(c, err, "Invalid request format")
		return
	}
	field.Set(reflect.ValueOf(req).Elem())
	if err := UpdateSetting(setting); err != nil {
		errhandler.Error(c, err, "Error updating setting data")
		return
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
		"message":   "update success",
	})
}
