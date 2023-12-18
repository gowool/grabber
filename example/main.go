package main

import (
	"encoding/json"
	"os"

	"github.com/gowool/grabber"
)

const target = "https://www.rainews.it/articoli/2023/12/massiccio-attacco-di-hacker-russi-a-enti-pubblici-italiani-chiesto-un-riscatto-6e384558-28ab-4aa4-9c64-9994eaf83ed8.html"

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
