package etcdv3

import (
	"common.bojiu.com/discover/kit/sd"
	"common.bojiu.com/discover/kit/sd/internal/instance"
	"common.bojiu.com/log"
	"go.uber.org/zap"
)

// Instancer yields instances stored in a certain etcd keyspace. Any kind of
// change in that keyspace is watched and will update the Instancer's Instancers.
type Instancer struct {
	cache  *instance.Cache
	client Client
	prefix string
	logger *zap.Logger
	quitc  chan struct{}
}

// NewInstancer returns an etcd instancer. It will start watching the given
// prefix for changes, and update the subscribers.
func NewInstancer(c Client, prefix string, logger *zap.Logger) (*Instancer, error) {
	if logger == nil {
		logger = log.NewLogger(log.SetDevelopment(true), log.SetLevel(zap.DebugLevel))
	}

	s := &Instancer{
		client: c,
		prefix: prefix,
		cache:  instance.NewCache(),
		logger: logger,
		quitc:  make(chan struct{}),
	}

	instances, err := s.client.GetEntries(s.prefix)
	if err == nil {
		s.logger.With(zap.String("prefix", s.prefix), zap.Any("instances", len(instances))).Info("提示")
		//logger.Log("prefix", s.prefix, "instances", len(instances))
	} else {
		s.logger.With(zap.String("prefix", s.prefix), zap.Error(err), zap.Stack("trace")).Info("Error")
		//logger.Log("prefix", s.prefix, "err", err)
	}
	s.cache.Update(sd.Event{Instances: instances, Err: err})

	go s.loop()
	return s, nil
}

func (s *Instancer) loop() {
	ch := make(chan struct{})
	go s.client.WatchPrefix(s.prefix, ch)

	for {
		select {
		case <-ch:
			instances, err := s.client.GetEntries(s.prefix)
			if err != nil {
				s.logger.With(zap.String("msg", "failed to retrieve entries"), zap.Error(err), zap.Stack("trace")).Info("Error")
				s.cache.Update(sd.Event{Err: err})
				continue
			}
			s.cache.Update(sd.Event{Instances: instances})

		case <-s.quitc:
			return
		}
	}
}

// Stop terminates the Instancer.
func (s *Instancer) Stop() {
	close(s.quitc)
}

// Register implements Instancer.
func (s *Instancer) Register(ch chan<- sd.Event) {
	s.cache.Register(ch)
}

// Deregister implements Instancer.
func (s *Instancer) Deregister(ch chan<- sd.Event) {
	s.cache.Deregister(ch)
}
