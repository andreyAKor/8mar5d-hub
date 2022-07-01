package devices

// StatusType Статус ответа.
type StatusType string

const (
	StatusTypeOk    StatusType = "ok"
	StatusTypeError StatusType = "error"
)

// DeviceType Тип устройства.
type DeviceType string

const (
	DeviceType1 DeviceType = "AVR"
)

// SensorType Тип сенсора.
type SensorType string

const (
	SensorTypeDHT11 SensorType = "DHT-11"
	SensorTypeMQ2   SensorType = "MQ-2"
)

// Status Статус ответа.
type Status struct {
	Status StatusType `json:"status"`
}

// Error Вывод ошибки.
type Error struct {
	Code    *int    `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
}

// Info информации о контроллере.
type Info struct {
	Status
	Error
	Data *struct {
		Type       DeviceType `json:"type"`
		Controller string     `json:"controller"`
		Types      []string   `json:"types"`
		Version    int        `json:"version"`
	} `json:"data,omitempty"`
}

// Sensors информации по сенсорам.
type Sensors struct {
	Status
	Error
	Data *[]struct {
		Type    SensorType `json:"type"`
		Version int        `json:"version"`
	} `json:"data,omitempty"`
}

// SensorDHT11 Метрики сенсора DHT-11.
type SensorDHT11 struct {
	Status
	Error
	Data *struct {
		Temperature struct {
			Value   float64 `json:"value"`
			Measure string  `json:"measure"`
		} `json:"temperature"`
		Humidity struct {
			Value   float64 `json:"value"`
			Measure string  `json:"measure"`
		} `json:"humidity"`
	} `json:"data,omitempty"`
}

// SensorMQ2 Метрики сенсора MQ-2.
type SensorMQ2 struct {
	Status
	Error
	Data *struct {
		Ratio struct {
			Value   float64 `json:"value"`
			Measure string  `json:"measure"`
		} `json:"ratio"`
		Lpg struct {
			Value   int    `json:"value"`
			Measure string `json:"measure"`
		} `json:"lpg"`
		Methane struct {
			Value   int    `json:"value"`
			Measure string `json:"measure"`
		} `json:"methane"`
		Smoke struct {
			Value   int    `json:"value"`
			Measure string `json:"measure"`
		} `json:"smoke"`
		Hydrogen struct {
			Value   int    `json:"value"`
			Measure string `json:"measure"`
		} `json:"hydrogen"`
	} `json:"data,omitempty"`
}
