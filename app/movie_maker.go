package main

import (
	"log"
	"os/exec"
)

// StichPicturesTogether stitches pictures together
func StichPicturesTogether(frames string) {
	cmd := exec.Command("ffmpeg", "-y", "-r", "10", "-i", frames+"/F_%03d.png", "-c:v", "libx264", "-vf", "fps=25", "-pix_fmt", "yuv420p", "movie.mp4")
	err := cmd.Run()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	} else {
		log.Printf("Successful Stiching")
	}
}
