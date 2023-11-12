package main

import (
	"encoding/json"
	"os"

	"github.com/gowool/grabber"
)

const target = "https://www.cbsnews.com/news/chef-gordon-ramsay-and-wife-tana-welcome-6th-child-jesse-james/"

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
	enc.Encode(page)
}
