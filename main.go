package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/net/ipv4"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"sync"
	"time"
)

// --- YOUR CORE LOGIC (VARS) ---
var (
	LocalPort      = "8053"
	RemoteDNS      = "8.8.8.8"
	RemoteAddr     = "8.8.8.8:53"
	CheckInterval  = 5 * time.Minute
	TOS_MAP        = map[string]int{
		"Default": 0, "CS1": 32, "CS2": 64, "CS3": 96, "CS4": 128, "CS5": 160, "CS6": 192, "CS7": 224,
		"AF11": 40, "AF12": 48, "AF13": 56, "AF21": 72, "AF22": 80, "AF23": 88,
		"AF31": 104, "AF32": 112, "AF33": 120, "AF41": 136, "AF42": 144, "AF43": 152, "EF": 184,
	}
	currentBestTOS  = 136
	mu              sync.Mutex
)

// --- THE ENGINE FUNCTION ---
func startEngine(statusLabel *widget.Label) {
	addr, _ := net.ResolveUDPAddr("udp4", "0.0.0.0:"+LocalPort)
	localConn, _ := net.ListenUDP("udp4", addr)
	localConn.SetReadBuffer(4 * 1024 * 1024)
	localConn.SetWriteBuffer(4 * 1024 * 1024)

	// Scanner Loop
	go func() {
		for {
			var bestLat float64 = 9999
			var bestName string
			var bestValue int
			for name, value := range TOS_MAP {
				cmd := exec.Command("ping", "-c", "1", "-W", "1", "-Q", strconv.Itoa(value), RemoteDNS)
				out, _ := cmd.Output()
				re := regexp.MustCompile(`time=([\d.]+)`)
				match := re.FindStringSubmatch(string(out))
				if len(match) > 1 {
					lat, _ := strconv.ParseFloat(match[1], 64)
					if lat < bestLat {
						bestLat, bestName, bestValue = lat, name, value
					}
				}
			}
			if bestValue != 0 {
				mu.Lock()
				currentBestTOS = bestValue
				mu.Unlock()
				statusLabel.SetText(fmt.Sprintf("ONLINE | TOS: %s | LAT: %.2fms", bestName, bestLat))
			}
			time.Sleep(CheckInterval)
		}
	}()

	// Forwarding Loop
	for {
		buf := make([]byte, 2048)
		n, clientAddr, _ := localConn.ReadFromUDP(buf)
		go func(data []byte, addr *net.UDPAddr) {
			if len(data) > 12 { data = append(data, 0x00) } // DPI Bypass
			remoteAddr, _ := net.ResolveUDPAddr("udp4", RemoteAddr)
			remoteConn, _ := net.DialUDP("udp4", nil, remoteAddr)
			defer remoteConn.Close()
			mu.Lock()
			p := ipv4.NewConn(remoteConn)
			_ = p.SetTOS(currentBestTOS)
			mu.Unlock()
			remoteConn.Write(data)
			resp := make([]byte, 2048)
			remoteConn.SetReadDeadline(time.Now().Add(2 * time.Second))
			rn, _, _ := remoteConn.ReadFromUDP(resp)
			localConn.WriteToUDP(resp[:rn], addr)
		}(buf[:n], clientAddr)
	}
}

func main() {
	myApp := app.New()
	window := myApp.NewWindow("Sigma Engine")
	window.Resize(fyne.NewSize(300, 200))

	status := widget.NewLabel("Status: Ready to Blast")
	btn := widget.NewButton("ACTIVATE", func() {
		status.SetText("Engine Starting...")
		go startEngine(status)
	})

	window.SetContent(container.NewVBox(status, btn))
	window.ShowAndRun()
}
