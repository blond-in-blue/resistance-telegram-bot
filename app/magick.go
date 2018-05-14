package main

import (
	"fmt"
	"log"
	"os/exec"
)

func StichPicturesTogether(frames string) {
	// convert -delay 10 F_*.png -loop 0 movie.gif
	cmd := exec.Command("convert", "-delay", "10", fmt.Sprintf("%s/F_*.png", frames), "-loop", "0", "staging.gif")
	err := cmd.Run()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	} else {
		log.Printf("Successful Stiching")
	}

	cmd = exec.Command("convert", "staging.gif", "-fuzz", "10%", "-layers", "Optimize", "movie.gif")
	err = cmd.Run()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	} else {
		log.Printf("Successful Compression")
	}

}
