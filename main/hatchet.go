// Copyright 2022-present Kuei-chun Chen. All rights reserved.

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/simagix/hatchet"
)

var repo = "simagix/hatchet"
var version = "devel-xxxxxx"

func main() {
	if version == "devel-xxxxxx" {
		version = "devel-" + time.Now().Format("20060102")
	}
	fullVersion := fmt.Sprintf(`%v %v`, repo, version)
	hatchet.Run(fullVersion)
}

func main1() {
	a := []string{"find", " -count", "update"}
	print(strings.Join(a, "\",\""))
}
