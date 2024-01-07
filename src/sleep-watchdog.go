package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	ticker := time.NewTicker(1 * time.Minute) // Check every minute
	defer ticker.Stop()

	conn, err := net.Dial("unix", "/var/run/acpid.socket")
	if err != nil {
		log.Println("Error connecting to acpid socket:", err)
		os.Exit(1)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	go func() {
		ticker := time.NewTicker(1 * time.Minute) // Check every minute
		defer ticker.Stop()
		onTick()
		for {
			<-ticker.C
			onTick()
		}
	}()

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading line:", err)
			continue
		}

		// Remove trailing newline character
		line = strings.TrimRight(line, "\n")

		// Analyze the line
		analyzeLine(line)
	}
}

func onTick() {
	maxIdleTime, err := getMaxIdleTime()
	if err != nil {
		panic(err)
	}
	fmt.Println("Max idle time:", maxIdleTime)
}
func analyzeLine(line string) {
	fields := strings.Fields(line)
	if len(fields) >= 3 {
		switch fields[0] {
		case "button/lid":
			switch fields[2] {
			case "close":
				log.Println("Lid is closed.")
				// Add your logic here to handle the lid being closed
			case "open":
				log.Println("Lid is open.")
				// Add your logic here to handle the lid being open
			}
		case "ac_adapter":
			switch fields[len(fields)-1] {
			case "00000000":
				log.Println("AC disconnected.")
				// Add your logic here to handle the AC being disconnected
			case "00000001":
				log.Println("AC connected.")
				// Add your logic here to handle the AC being connected
			}
		case "button/power":
			log.Println("Power button pressed.")
			// Add your logic here to handle the power button press
		}
	}
}

func checkIdleTime(user, display string, sudoUid int) (int, error) {
	// Your implementation of checkIdleTime goes here
	// For example:
	cmd := exec.Command("bash", "-c", fmt.Sprintf(`sudo -u %s %s DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/%d/bus dbus-send --print-reply --dest=org.gnome.Mutter.IdleMonitor /org/gnome/Mutter/IdleMonitor/Core org.gnome.Mutter.IdleMonitor.GetIdletime | awk '{print $NF}' | tail -n 1`, user, display, sudoUid))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	idleTimeStr := strings.TrimSpace(string(out))
	idleTime, err := strconv.Atoi(idleTimeStr)
	if err != nil {
		return 0, err
	}

	return idleTime, nil
}

func getMaxIdleTime() (int, error) {
	// Get list of all logged-in users
	cmd := exec.Command("bash", "-c", "who | cut -d' ' -f1 | sort -u")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, err
	}
	loggedInUsers := strings.Split(out.String(), "\n")
	sort.Strings(loggedInUsers)

	maxIdleTime := 0
	for _, user := range loggedInUsers {
		user = strings.TrimSpace(user)
		if user == "" {
			continue
		}

		// Get sudo uid
		cmd = exec.Command("id", "-u", user)
		var sudoUidOut bytes.Buffer
		cmd.Stdout = &sudoUidOut
		err = cmd.Run()
		if err != nil {
			return 0, err
		}
		sudoUid, err := strconv.Atoi(strings.TrimSpace(sudoUidOut.String()))
		if err != nil {
			return 0, err
		}

		// Get display info
		cmd := exec.Command("bash", "-c", fmt.Sprintf("w -h %s | awk '$3 ~ /:[0-9.]*/{print $3}'", user))
		var displayInfoOut bytes.Buffer
		cmd.Stdout = &displayInfoOut
		err = cmd.Run()
		if err != nil {
			return 0, err
		}
		displayInfo := strings.TrimSpace(displayInfoOut.String())
		display := ""
		if strings.Contains(displayInfo, ":") {
			// This is a X11 display, export DISPLAY
			display = "DISPLAY=" + displayInfo
		} else if strings.HasSuffix(displayInfo, ".0") {
			// This is a Wayland display, export WAYLAND_DISPLAY
			display = "WAYLAND_DISPLAY=" + displayInfo
		}

		// Call checkIdleTime and update maxIdleTime if necessary
		idleTime, err := checkIdleTime(user, display, sudoUid)
		if err != nil {
			return 0, err
		}
		if idleTime > maxIdleTime {
			maxIdleTime = idleTime
		}
	}

	return maxIdleTime, nil
}
