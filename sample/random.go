package sample

import (
	"math/rand"
	"time"

	"github.com/ak-karimzai/go-grpc/pb"
	"github.com/google/uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomKeyboardLayout() pb.Keyboard_Layout {
	switch rand.Intn(3) {
	case 1:
		return pb.Keyboard_QWERTY
	case 2:
		return pb.Keyboard_QWERTZ
	default:
		return pb.Keyboard_AZERTY
	}
}

func randomCPUBrand() string {
	return randomStringFromSet("Intel", "AMD")
}

func randomBool() bool {
	return rand.Intn(2) == 1
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func randomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randomFloat32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func randomScreenPanel() pb.Screen_Panel {
	if rand.Intn(2) == 1 {
		return pb.Screen_IPS
	}
	return pb.Screen_OLE
}

func randomScreenResoultion() *pb.Screen_Resolution {
	height := randomInt(1080, 4320)
	width := height * 16 / 9

	return &pb.Screen_Resolution{
		Height: uint32(height),
		Width:  uint32(width),
	}
}

func randomStringFromSet(arr ...string) string {
	n := len(arr)

	if n == 0 {
		return ""
	}
	return arr[rand.Intn(n)]
}

func randomCPUName(brand string) string {
	if brand == "Intel" {
		return randomStringFromSet("CORE I7",
			"CORE I5",
			"CORE I3",
			"CORE 2DUO")
	}

	return randomStringFromSet("A5",
		"A6",
		"A7",
		"A9",
		"A10",
		"RYZEN 3",
		"RYZEN 5")
}

func randomGPUBrand() string {
	return randomStringFromSet("NVIDIA", "AMD")
}

func randomGPUName(gpuBrand string) string {
	if gpuBrand == "NVIDIA" {
		return randomStringFromSet("RTX 2060",
			"RTX 2070",
			"GTX 1660-Ti",
			"GTX-1070")
	}
	return randomStringFromSet("RX 590",
		"RX 580",
		"RX 5700-XT",
		"RX VEGA-56")
}

func randomID() string {
	return uuid.New().String()
}

func randomLaptopBrand() string {
	return randomStringFromSet(
		"DELL",
		"HP",
		"COMPAQ",
		"APPLE",
		"AOC")
}

func randomLaptopName(brand string) string {
	if brand == "DELL" {
		return randomStringFromSet("XPS-15",
			"XPS-13",
			"DELL latitude")
	}
	return "UNKOWN"
}
