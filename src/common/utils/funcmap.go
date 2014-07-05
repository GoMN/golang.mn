package utils

import (
	"html/template"
)

var (
	TemplateFuncMap = template.FuncMap{
	"addint": addInt,
}
)

func addInt(num int, num2 int) int {
	return num + num2
}
