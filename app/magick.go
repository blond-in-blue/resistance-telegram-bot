package main

import (
	"fmt"
	"log"
	"os/exec"
)

// StichPicturesTogether stitches pictures together
func StichPicturesTogether(frames string) {
	// convert -delay 10 F_*.png -loop 0 movie.gif
	// cmd := exec.Command("convert", "-delay", "10", fmt.Sprintf("%s/F_*.png", frames), "-loop", "0", "staging.gif")
	cmd := exec.Command("ffmpeg", "-r", "10", "-i", fmt.Sprintf("%s/F_%%03d.png", frames), "-c:v", "libx264", "-vf", "fps=25", "-pix_fmt", "yuv420p", "movie.mp4")
	err := cmd.Run()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	} else {
		log.Printf("Successful Stiching")
	}

	// cmd = exec.Command("convert", "staging.gif", "-fuzz", "10%", "-layers", "Optimize", "movie.gif")
	// err = cmd.Run()
	// if err != nil {
	// 	log.Printf("Command finished with error: %v", err)
	// } else {
	// 	log.Printf("Successful Compression")
	// }

}
