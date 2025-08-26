package Utils

import (
	"fmt"
	"strings"
	"text/template"
)

func BuildTemplate(data interface{}, templateString string) string {

	var tempBuff strings.Builder
	tmpl, _ := template.New("buildTemplators").Parse(templateString)

	if err := tmpl.Execute(&tempBuff, data); err != nil {
		fmt.Println(err)
		return ""
	}
	return tempBuff.String()
}
