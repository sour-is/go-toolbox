/* modules:
MQTTMessage:
	model: sour.is/x/toolbox/mqtt.GraphMessage
*/
package mqtt

import (
	"context"
)

// GraphMqtt implements graphql resolver
type GraphMqtt struct{}

// GraphMessage returned to client
type GraphMessage struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

// Mqtt subscribes to a given topic
func (GraphMqtt) Mqtt(ctx context.Context, topic string, qos *int) (out <-chan GraphMessage, err error) {
	var unsub func()
	var in <-chan Message
	if qos == nil {
		v := 0
		qos = &v
	}

	in, unsub, err = Subscribe(topic, int8(*qos))
	if err != nil {
		return
	}

	ch := make(chan GraphMessage)
	out = ch

	go func() {
		defer close(ch)
		defer unsub()

		for {
			select {
			case <-ctx.Done():
				return

			case m := <-in:
				pm := GraphMessage{
					Topic:   m.Topic(),
					Message: string(m.Message),
				}

				select {
				case <-ctx.Done():

					return
				case ch <- pm:
				}
			}
		}
	}()

	return
}

// MqttPublish publishes a message onto the queue
func (GraphMqtt) MqttPublish(
	ctx context.Context,
	topic string,
	message string,
) (ok bool, err error) {
	var m Message

	m, err = NewMessage(topic, message)
	if err != nil {
		return
	}

	err = Publish(m)
	ok = err == nil

	return
}
