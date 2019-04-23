package mqtt

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
	"sour.is/x/toolbox/log"
)

// Message is a value received from mqtt
type Message struct {
	TopicName []byte
	Message   []byte
	QoS       uint8
	Retain    bool
}

// Topic message was sent on
func (m Message) Topic() string {
	return string(m.TopicName)
}

func (m Message) String() string {
	return fmt.Sprintf("%s: %s", m.Topic(), string(m.Message))
}

// JSON unmarshalls content into passed struct.
func (m Message) JSON(s interface{}) error {
	return json.Unmarshal(m.Message, s)
}

// NewMessage creates new message
func NewMessage(topic string, s interface{}) (m Message, err error) {
	m.TopicName = []byte(topic)
	switch v := s.(type) {
	case fmt.Stringer:
		m.Message = []byte(v.String())
	case []byte:
		m.Message = v
	case string:
		m.Message = []byte(v)
	default:
		m.Message, err = json.Marshal(s)
	}

	return
}

// Client holds a threadsafe client
type Client struct {
	*client.Client
	Topics map[string][]chan<- Message
	mutex  *sync.RWMutex
}

// ConnectOptions defines options for connecting to mqtt
type ConnectOptions = client.ConnectOptions

// Default shared connection to mqtt
var Default Client

// Dial starts a shared connection to mqtt
func Dial(c *ConnectOptions) error {
	var err error
	Default, err = New(c)
	return err
}

// Terminate disconnects shared connection
func Terminate() {
	Default.Terminate()
}

// Terminate closes all subsciptions and terminates connection.
func (c Client) Terminate() {
	log.Debug("cleaning up subscriptions")
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for name, topics := range c.Topics {
		for _, ch := range topics {
			close(ch)
		}
		delete(c.Topics, name)
	}

	c.Client.Terminate()
}

// Publish sends message on shared connection
func Publish(m Message) error {
	return Default.Publish(m)
}

// Subscribe starts a watcher for shared connecton
func Subscribe(topic string, qos int8) (<-chan Message, func(), error) {
	return Default.Subscribe(topic, qos)
}

// New connects to MQTT
func New(c *ConnectOptions) (Client, error) {
	client := Client{
		Client: client.New(
			&client.Options{
				// Define the processing of the error handler.
				ErrorHandler: func(err error) {
					log.Error(err)
				},
			}),
		Topics: make(map[string][]chan<- Message),
		mutex:  &sync.RWMutex{},
	}
	// Connect to the MQTT Server.
	err := client.Connect(c)

	return client, err
}

// Publish a JSON mesg to the network.
func (c Client) Publish(m Message) error {
	return c.Client.Publish(
		&client.PublishOptions{
			QoS:       m.QoS,
			TopicName: m.TopicName,
			Message:   m.Message,
			Retain:    m.Retain,
		},
	)
}

// Subscribe to a topic and get messages on channel
func (c Client) Subscribe(
	topic string,
	qos int8,
) (
	out <-chan Message,
	unsubscribe func(),
	err error,
) {
	if !mqtt.ValidQoS(byte(qos)) {
		err = fmt.Errorf("invalid qos = %v", qos)
		return
	}

	ch := make(chan Message)
	out = ch
	var in chan<- Message = ch

	unsubscribe = func() {
		err = c.unsubscribe(topic, in)
		if err != nil {
			log.Error(err)
		}
	}

	fn := func(topicName, message []byte) {
		c.mutex.RLock()
		defer c.mutex.RUnlock()

		m := Message{
			TopicName: topicName,
			Message:   message,
		}
		for _, ch := range c.Topics[topic] {
			ch <- m
		}
	}

	c.Client.Subscribe(
		&client.SubscribeOptions{
			SubReqs: []*client.SubReq{
				&client.SubReq{
					TopicFilter: []byte(topic),
					QoS:         byte(qos),
					Handler:     fn,
				},
			},
		},
	)

	log.Debug("subscribe to ", topic)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Topics[topic] = append(c.Topics[topic], in)

	return
}

func (c Client) unsubscribe(topic string, in chan<- Message) (err error) {
	log.Debug("unsubscribe from ", topic)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i, ch := range c.Topics[topic] {
		if ch == in {
			close(ch)

			if len(c.Topics[topic]) == 1 {
				delete(c.Topics, topic)
				err = c.Client.Unsubscribe(
					&client.UnsubscribeOptions{
						TopicFilters: [][]byte{[]byte(topic)},
					},
				)
				return
			}

			c.Topics[topic] = append(c.Topics[topic][:i], c.Topics[topic][i+1:]...)

			return
		}
	}

	return
}
