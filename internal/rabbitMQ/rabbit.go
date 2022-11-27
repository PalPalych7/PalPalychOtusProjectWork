package rabbitmq

import (
	"context"

	ms "github.com/PalPalych7/OtusProjectWork/internal/mainstructs"
	"github.com/streadway/amqp"
)

type RabbitQueue struct {
	rabCfg  ms.RabbitCFG
	conn    *amqp.Connection
	channel *amqp.Channel
	done    chan error
	ctx     context.Context
}

func New(ctx context.Context, q ms.RabbitCFG) (*RabbitQueue, error) {
	c := &RabbitQueue{
		rabCfg:  q,
		conn:    nil,
		channel: nil,
		done:    make(chan error),
		ctx:     ctx,
	}
	return c, nil
}

func (r *RabbitQueue) Start() error {
	var err error
	r.conn, err = amqp.Dial(r.rabCfg.URI)
	if err != nil {
		return err
	}
	r.channel, err = r.conn.Channel()
	if err != nil {
		return err
	}

	if err = r.channel.ExchangeDeclare(
		r.rabCfg.Exchange,
		r.rabCfg.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	_, err = r.channel.QueueDeclare(
		r.rabCfg.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	if err = r.channel.QueueBind(
		r.rabCfg.Queue,
		r.rabCfg.BindingKey,
		r.rabCfg.Exchange,
		false,
		nil,
	); err != nil {
		return err
	}
	return err
}

func (r *RabbitQueue) SendMess(myMes []byte) error {
	err := r.channel.Publish(
		r.rabCfg.Exchange,
		r.rabCfg.BindingKey,
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            myMes,
			DeliveryMode:    amqp.Transient,
			Priority:        0,
		},
	)
	return err
}

func (r *RabbitQueue) Shutdown() error {
	err := r.channel.Cancel(r.rabCfg.ConsumerTag, true)
	if err != nil {
		return err
	}
	return r.conn.Close()
}
