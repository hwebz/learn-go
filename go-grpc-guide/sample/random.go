package sample

import (
	"github.com/google/uuid"
	pb "github.com/hwebz/go-grpc-guide/pb"
	"math/rand"
	"time"
)

// init() is automatically called before any other logic
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

func randomBool() bool {
	return rand.Intn(2) == 1
}

func randomCPUBrand() string {
	return randomStringFromSet("Intel", "AMD")
}

func randomStringFromSet(a ...string) string {
	n := len(a)
	if n == 0 {
		return ""
	}
	return a[rand.Intn(n)]
}

func randomCPUName(brand string) string {
	if brand == "Intel" {
		return randomStringFromSet("Xeon E-1234v", "Core i9-9900K", "Core i7-7700K", "Core i5-5500K", "Core i3 3300K")
	}

	return randomStringFromSet("Ryzen 7 Pro 2700U", "Ryzen 5 Pro 5100U", "Ryzen 3 Pro 3200U")
}

func randomInt(min int, max int) int {
	return min + rand.Intn(max-min+1)
}

func randomFloat32(min float32, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func randomFloat64(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randomGPUBrand() string {
	return randomStringFromSet("NVIDIA", "AMD")
}

func randomGPUName(brand string) string {
	if brand == "NVIDIA" {
		return randomStringFromSet("RTX 3090", "RTX 3080", "RTX 3070", "RTX 3060", "RTX 3050")
	}

	return randomStringFromSet("RX 590", "RX 580", "RX 5700XT", "RX Vega-56")
}

func randomScreenResolution() *pb.Screen_Resolution {
	height := randomInt(1080, 4320)
	width := height * 16 / 9

	resolution := &pb.Screen_Resolution{
		Width:  uint32(width),
		Height: uint32(height),
	}

	return resolution
}

func randomScreenPanel() pb.Screen_Panel {
	if rand.Intn(2) == 1 {
		return pb.Screen_OLED
	}

	return pb.Screen_IPS
}

func randomUUID() string {
	return uuid.New().String()
}

func randomLaptopBrand() string {
	return randomStringFromSet("Asus", "Dell", "Acer", "HP", "Lenovo", "Apple")
}

func randomLaptopName(brand string) string {
	switch brand {
	case "Apple":
		return randomStringFromSet("MacBook Air", "MacBook Pro")
	case "Asus":
		return randomStringFromSet("ZenBook", "VivoBook", "ROG Zephyrus", "TUF Gaming")
	case "Dell":
		return randomStringFromSet("XPS", "Inspiron", "Alienware", "Latitude", "Vostro")
	case "Acer":
		return randomStringFromSet("Aspire", "Swift", "Predator", "Nitro", "Spin")
	case "HP":
		return randomStringFromSet("Spectre", "Pavilion", "Envy", "Omen", "EliteBook")
	case "Lenovo":
		return randomStringFromSet("ThinkPad", "IdeaPad", "Legion", "Yoga")
	default:
		return "Unknown Model"
	}
}
