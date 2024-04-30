package main

import (
	"fmt"
	"stage2024/pkg/settings"

	"gopkg.in/yaml.v3"
)

func main() {
	s := settings.Settings{}
	s.SetDefault()
	y, _ := yaml.Marshal(&s)
	fmt.Print(string(y))
}
