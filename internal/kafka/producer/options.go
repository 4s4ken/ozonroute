package producer

import (
	"github.com/IBM/sarama"
)

type Option interface {
	Apply(*sarama.Config) error
}
