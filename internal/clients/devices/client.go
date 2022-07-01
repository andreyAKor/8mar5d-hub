package devices

import (
	"encoding/json"
	"io"
	"time"

	"github.com/pkg/errors"
	"github.com/tarm/serial"
)

const (
	serialBaud        = 9600
	serialReadTimeout = time.Second * 3
	bufSize           = 256

	// Команды контроллера.
	serialCommandInfo        = "/\n"              // Получение информации о контроллере
	serialCommandSensors     = "/sensors\n"       // Получение информации по сенсорам
	serialCommandSensorDHT11 = "/sensor/DHT-11\n" // Получение метрик сенсора DHT-11
	serialCommandSensorMQ2   = "/sensor/MQ-2\n"   // Получение метрик сенсора MQ-2
)

var (
	ErrSerialNotConn = errors.New("serial not connected")

	connections map[string]*serial.Port

	_ io.Closer = (*Client)(nil)
)

type Client struct {
	host string
}

func init() {
	connections = make(map[string]*serial.Port)
}

// New Клиент для работы с контроллерами через seral-порт.
func New(host string) (*Client, error) {
	if _, ok := connections[host]; !ok {
		s, err := serial.OpenPort(&serial.Config{
			Name:        host,
			Baud:        serialBaud,
			ReadTimeout: serialReadTimeout,
		})
		if err != nil {
			return nil, errors.Wrapf(err, `serial open port fail on host "%s"`, host)
		}

		connections[host] = s
	}

	return &Client{
		host: host,
	}, nil
}

// Host Хост устройства.
func (c *Client) Host() string {
	return c.host
}

// Info Получение информации о контроллере.
func (c *Client) Info() (*Info, error) {
	data, err := c.send(serialCommandInfo)
	if err != nil {
		return nil, errors.Wrap(err, "send fail")
	}

	res := &Info{}
	if err := json.Unmarshal(data, res); err != nil {
		return nil, errors.Wrap(err, "json unmarshal fail")
	}
	if res.Status.Status == StatusTypeError {
		if res.Message != nil {
			return nil, errors.New(*res.Message)
		}
		return nil, errors.New("controller unknown error")
	}

	return res, nil
}

// Sensors Получение информации по сенсорам.
func (c *Client) Sensors() (*Sensors, error) {
	data, err := c.send(serialCommandSensors)
	if err != nil {
		return nil, errors.Wrap(err, "send fail")
	}

	res := &Sensors{}
	if err := json.Unmarshal(data, res); err != nil {
		return nil, errors.Wrap(err, "json unmarshal fail")
	}
	if res.Status.Status == StatusTypeError {
		if res.Message != nil {
			return nil, errors.New(*res.Message)
		}
		return nil, errors.New("controller unknown error")
	}

	return res, nil
}

// SensorDHT11 Получение метрик сенсора DHT-11.
func (c *Client) SensorDHT11() (*SensorDHT11, error) {
	data, err := c.send(serialCommandSensorDHT11)
	if err != nil {
		return nil, errors.Wrap(err, "send fail")
	}

	res := &SensorDHT11{}
	if err := json.Unmarshal(data, res); err != nil {
		return nil, errors.Wrap(err, "json unmarshal fail")
	}
	if res.Status.Status == StatusTypeError {
		if res.Message != nil {
			return nil, errors.New(*res.Message)
		}
		return nil, errors.New("controller unknown error")
	}

	return res, nil
}

// SensorMQ2 Получение метрик сенсора MQ-2.
func (c *Client) SensorMQ2() (*SensorMQ2, error) {
	data, err := c.send(serialCommandSensorMQ2)
	if err != nil {
		return nil, errors.Wrap(err, "send fail")
	}

	res := &SensorMQ2{}
	if err := json.Unmarshal(data, res); err != nil {
		return nil, errors.Wrap(err, "json unmarshal fail")
	}
	if res.Status.Status == StatusTypeError {
		if res.Message != nil {
			return nil, errors.New(*res.Message)
		}
		return nil, errors.New("controller unknown error")
	}

	return res, nil
}

// Close Закрывает соединение с serial-портом.
func (c *Client) Close() error {
	s, ok := connections[c.host]
	if !ok {
		return ErrSerialNotConn
	}
	if s == nil {
		delete(connections, c.host)
		return ErrSerialNotConn
	}

	if err := s.Close(); err != nil {
		return errors.Wrap(err, "serial close fail")
	}

	return nil
}

// send Посылает команду на контроллер.
func (c *Client) send(cmd string) ([]byte, error) {
	s, ok := connections[c.host]
	if !ok {
		return nil, ErrSerialNotConn
	}
	if s == nil {
		delete(connections, c.host)
		return nil, ErrSerialNotConn
	}

	if s == nil {
		return nil, ErrSerialNotConn
	}

	n, err := s.Write([]byte(cmd))
	if err != nil {
		return nil, errors.Wrap(err, "serial write fail")
	}

	var res []byte
	for {
		buf := make([]byte, bufSize)
		n, err = s.Read(buf)
		if err != nil && err != io.EOF {
			return nil, errors.Wrap(err, "serial read fail")
		}

		res = append(res, buf[:n]...)

		if err == io.EOF || (buf[n-1] == '\n') {
			break
		}
	}

	return res, nil
}
