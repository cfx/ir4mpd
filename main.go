package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

var codes = map[int]string{
	//	0x45: "CH-",
	//	0x46: "CH",
	//	0x47: "CH+",
	0x44: "PREV",
	0x40: "NEXT",
	0x43: "PLAY/PAUSE",
	0x7:  "VOL-",
	0x15: "VOL+",
	//	0x9:  "EQ",
	//	0x16: "0",
	//	0x19: "100+",
	//	0xD:  "200+",
	0xC:  "1",
	0x18: "2",
	0x5E: "3",
	0x8:  "4",
	0x1C: "5",
	0x5A: "6",
	0x42: "7",
	0x52: "8",
	0x4a: "9",
}

const event_delay_ms = 150

func main() {
	var ts time.Time
	buff := make([]byte, 16)

	f, err := os.Open("/dev/input/event0")

	if err != nil {
		panic(err)
	}

	defer f.Close()

	for {
		f.Read(buff)
		code := buff[12]
		code_str := codes[int(code)]
		now := time.Now()
		time_diff := now.Sub(ts)

		if code != 0 && code_str != "" && time_diff.Milliseconds() > event_delay_ms {
			ts = time.Now()
			go func(code_str string) {
				switch code_str {
				case "VOL+":
					execMpc("volume", "+10")
					return

				case "VOL-":
					execMpc("volume", "-10")
					return

				case "PLAY/PAUSE":
					execMpc("toggle")
					return

				case "NEXT":
					execMpc("next")
					return

				case "PREV":
					execMpc("prev")
					return

				}

				execMpc("play", code_str)
				return

			}(code_str)
		}
	}

}

func execMpc(name string, params ...string) {
	args := append([]string{"-q", name}, params...)
	_, err := exec.Command("mpc", args...).Output()

	if err != nil {
		fmt.Errorf("mpc error: %v\n", err)
	}
}
