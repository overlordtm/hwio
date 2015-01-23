package hwio

type AriaPin struct {
	names       []string
	modules     []string
	gpioLogical int
}

type AriaG25Driver struct {
	ariaPins []*AriaPin
	modules  map[string]Module
}

func (d *AriaG25Driver) Init() error {
	d.createPinData()
	return d.initializeModules()
}

func (d *AriaG25Driver) makePin(names []string, modules []string, gpioLogical int) *AriaPin {
	return &AriaPin{names, modules, gpioLogical}
}

func (d *AriaG25Driver) createPinData() {
	d.ariaPins = []*AriaPin{
		d.makePin([]string{"PA22"}, []string{"gpio"}, 22),
		d.makePin([]string{"PA28"}, []string{"gpio"}, 28),
	}
}

func (d *AriaG25Driver) initializeModules() error {
	d.modules = make(map[string]Module)

	gpio := NewAriaGPIOModule("gpio")
	e := gpio.SetOptions(d.getGPIOOptions())
	if e != nil {
		return e
	}

	d.modules["gpio"] = gpio

	return nil

}

func (d *AriaG25Driver) getGPIOOptions() map[string]interface{} {
	result := make(map[string]interface{})

	pins := make(AriaGPIOModulePinDefMap)

	for i, hw := range d.ariaPins {
		if d.usedBy(hw, "gpio") {
			pins[Pin(i)] = &AriaGPIOModulePinDef{pin: Pin(i), gpioLogical: hw.gpioLogical}
		}
	}
	result["pins"] = pins

	return result
}

func (d *AriaG25Driver) GetModules() map[string]Module {
	return d.modules
}

func (d *AriaG25Driver) Close() {
	for _, module := range d.modules {
		module.Disable()
	}
}

func (d *AriaG25Driver) PinMap() (pinMap HardwarePinMap) {
	pinMap = make(HardwarePinMap)

	for i, hw := range d.ariaPins {
		pinMap.add(Pin(i), hw.names, hw.modules)
	}
	return
}

func (d *AriaG25Driver) usedBy(pinDef *AriaPin, module string) bool {
	for _, n := range pinDef.modules {
		if n == module {
			return true
		}
	}
	return false
}
