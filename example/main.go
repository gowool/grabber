package main

import (
	"encoding/json"
	"os"

	"github.com/gowool/grabber"
)

const target = "https://www.cnbc.com/2025/06/05/auto-groups-sound-the-alarm-as-chinas-rare-earth-curbs-start-to-bite.html"

func main() {
	req, err := grabber.NewRequest(target)
	if err != nil {
		panic(err)
	}

	page, err := grabber.Do(req)
	if err != nil {
		panic(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(page)
}
