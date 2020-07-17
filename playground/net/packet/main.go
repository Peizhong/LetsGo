package main

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// sudo tcpdump -i lo dst port 3000

var (
	device       string = "lo"
	snapshot_len int32  = 512000
	promiscuous  bool   = false
	err          error
	timeout      time.Duration = 1 * time.Second
	handle       *pcap.Handle
)

var (
	targetPort = 3000
)

// sudo apt-get install libpcap-dev
func main() {
	// Find all devices
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Devices found:")
	for _, device := range devices {
		fmt.Println("\nName: ", device.Name)
		fmt.Println("Description: ", device.Description)
		fmt.Println("Devices addresses: ", device.Description)
		for _, address := range device.Addresses {
			fmt.Println("- IP address: ", address.IP)
			fmt.Println("- Subnet mask: ", address.Netmask)
		}
	}

	// Open device
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	var bytes int64
	var count int64
	go func() {
		for {
			<-time.After(time.Second)
			log.Printf("tcp[%d] recv: %dtimes@%dbytes\n", targetPort, atomic.LoadInt64(&count), atomic.LoadInt64(&bytes))
			atomic.StoreInt64(&count, 0)
			atomic.StoreInt64(&bytes, 0)
		}
	}()

	// Use the handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			if tcp, ok := tcpLayer.(*layers.TCP); ok && tcp.DstPort == layers.TCPPort(targetPort) {
				atomic.AddInt64(&bytes, int64(len(tcpLayer.LayerContents())))
				atomic.AddInt64(&count, 1)
				// spew.Dump(tcp)
			}
		}
	}
	// fmt.Println(packet)
}
