package main

/*
#cgo CFLAGS: -Ifobos
#cgo LDFLAGS: -Lfobos -lfobos -lusb-1.0
#include "fobos.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	count := C.fobos_rx_get_device_count()
	if count <= 0 {
		fmt.Println("Cannot find the Fobos devices")
		return
	}
	fmt.Printf("Found devices: %d\n", int(count))

	var dev *C.struct_fobos_dev_t
	res := C.fobos_rx_open(&dev, 0)
	if res != 0 {
		fmt.Printf("Cannot open device: %d\n", int(res))
		return
	}
	defer C.fobos_rx_close(dev)

	// Device information
	hw := make([]C.char, 32)
	fw := make([]C.char, 32)
	man := make([]C.char, 32)
	prod := make([]C.char, 32)
	ser := make([]C.char, 32)

	C.fobos_rx_get_board_info(dev, &hw[0], &fw[0], &man[0], &prod[0], &ser[0])
	fmt.Printf("HW: %s, FW: %s, SN: %s\n",
		C.GoString(&hw[0]), C.GoString(&fw[0]), C.GoString(&ser[0]))

	// Set up the device
	var freqHz C.double = 106000000 // 106 MHz, for example
	var actualFreq C.double
	C.fobos_rx_set_frequency(dev, freqHz, &actualFreq)
	fmt.Println("Actual frequency:", float64(actualFreq))

	C.fobos_rx_set_direct_sampling(dev, 0)
	C.fobos_rx_set_lna_gain(dev, 1)
	C.fobos_rx_set_vga_gain(dev, 10)

	//Set the sampling rate to 8 MHz, it is the minimum supported by Fobos
	var srHz C.double = 8e6

	var actualSr C.double
	C.fobos_rx_set_samplerate(dev, srHz, &actualSr)
	fmt.Println("Actual sampling frequency:", float64(actualSr))

	C.fobos_rx_set_clk_source(dev, 0)

	// Synchronous receiving
	bufLen := 1024
	C.fobos_rx_start_sync(dev, C.uint32_t(bufLen))
	defer C.fobos_rx_stop_sync(dev)

	iqBuf := make([]float32, bufLen*2)

	for {
		var actual C.uint32_t
		C.fobos_rx_read_sync(dev, (*C.float)(unsafe.Pointer(&iqBuf[0])), &actual)
		if actual != 0 {
			//fmt.Printf("Received %d samples\n", int(actual))
			// Process the received samples in iqBuf
			// For demonstration, we will just print the first few samples
			//for i := 0; i < 10 && i < int(actual); i++ {
			//fmt.Printf("Sample %d: I = %f, Q = %f\n", i, iqBuf[i*2], iqBuf[i*2+1])
			//}
		} else {
			fmt.Println("No samples received")
		}
	}
}
