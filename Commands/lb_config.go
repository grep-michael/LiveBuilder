package commands

import (
	filesystem "LiveBuilder/Filesystem"
	"bytes"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type LBConfigData struct {
	ISOVolume      string
	ISOPublisher   string
	ISOApplication string
	ISOImageName   string
}

func NewlbConfig() *LBConfigData {
	return &LBConfigData{
		ISOVolume:      "DefaultVolume",
		ISOPublisher:   "DefaultPublisher",
		ISOApplication: "DefaultApplication",
		ISOImageName:   "DefaultName",
	}
}

func (cfg *LBConfigData) GetConfigTemplate() *template.Template {

	lb_confg_template := filepath.Join(filesystem.GetFileManager().GetAppDataDir(), filesystem.LBCONFIG_TEMPLATE_ID)
	content, err := os.ReadFile(lb_confg_template)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	tmpl, err := template.New("lbconfig").Parse(string(content))
	if err != nil {
		panic(err)
	}
	return tmpl
}
func (cfg *LBConfigData) BuildTemplate() string {
	template := cfg.GetConfigTemplate()
	var buf bytes.Buffer
	if err := template.Execute(&buf, cfg); err != nil {
		panic(err)
	}
	str := buf.String()
	log.Printf("Build lb config from template: %s\n", str)
	return str
}
