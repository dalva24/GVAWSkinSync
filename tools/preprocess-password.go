package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	in, err := os.Open("res/id_top200.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	out, err := os.Create("res/passwords.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	badwords := []string{
		"penis",
		"vagina",
		"kontol",
		"kemaluan",
		"kelamin",
		"payudara",
		"puting",
		"pantat",
		"sex",
		"perkosa",
		"bunuh",
		"fff",
		"rokok",
		"judi",
	}

	scanner := bufio.NewScanner(in)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.Split(text, " ")[0]
		if len(text) >= 5 && len(text) <= 9 && !strings.ContainsAny(text, "0123456789!@#$%^&*()`~-=_+[]{};':\",./<>?\\|") {
			bad := false
			for _, word := range badwords {
				if strings.Contains(text, word) {
					bad = true
					break
				}
			}
			if !bad {
				fmt.Println(text)
				out.WriteString(text + "\n")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
