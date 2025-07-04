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
	var freqHz C.double = 107000000
	var actualFreq C.double
	C.fobos_rx_set_frequency(dev, freqHz, &actualFreq)
	fmt.Println("Actual frequency:", float64(actualFreq))

	C.fobos_rx_set_direct_sampling(dev, 0)
	C.fobos_rx_set_lna_gain(dev, 1)
	C.fobos_rx_set_vga_gain(dev, 10)

	var srHz C.double = 5e6
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
		fmt.Printf("Received %d sampls\n", actual)
	}
}
