package etcdv3

import (
	"go.uber.org/zap"
	"sync"
	"time"

	"common.bojiu.com/log"
)

const minHeartBeatTime = 500 * time.Millisecond

// Registrar registers service instance liveness information to etcd.
type Registrar struct {
	client  Client
	service Service
	logger  *zap.Logger

	quitmtx sync.Mutex
	quit    chan struct{}
}

// Service holds the instance identifying data you want to publish to etcd. Key
// must be unique, and value is the string returned to subscribers, typically
// called the "instance" string in other parts of package sd.
type Service struct {
	Key   string // unique key, e.g. "/service/foobar/1.2.3.4:8080"
	Value string // returned to subscribers, e.g. "http://1.2.3.4:8080"
	TTL   *TTLOption
}

// TTLOption allow setting a key with a TTL. This option will be used by a loop
// goroutine which regularly refreshes the lease of the key.
type TTLOption struct {
	heartbeat time.Duration // e.g. time.Second * 3
	ttl       time.Duration // e.g. time.Second * 10
}

// NewTTLOption returns a TTLOption that contains proper TTL settings. Heartbeat
// is used to refresh the lease of the key periodically; its value should be at
// least 500ms. TTL defines the lease of the key; its value should be
// significantly greater than heartbeat.
//
// Good default values might be 3s heartbeat, 10s TTL.
func NewTTLOption(heartbeat, ttl time.Duration) *TTLOption {
	if heartbeat <= minHeartBeatTime {
		heartbeat = minHeartBeatTime
	}
	if ttl <= heartbeat {
		ttl = 3 * heartbeat
	}
	return &TTLOption{
		heartbeat: heartbeat,
		ttl:       ttl,
	}
}

// NewRegistrar returns a etcd Registrar acting on the provided catalog
// registration (service).
func NewRegistrar(client Client, service Service, logger *zap.Logger) *Registrar {
	if logger == nil {
		logger = log.NewLogger(log.SetDevelopment(true), log.SetLevel(zap.DebugLevel))
	}

	return &Registrar{
		client:  client,
		service: service,
		logger:  logger,
	}
}

// Register implements the sd.Registrar interface. Call it when you want your
// service to be registered in etcd, typically at startup.
func (r *Registrar) Register() {
	if err := r.client.Register(r.service); err != nil {
		r.logger.With(zap.Error(err), zap.Stack("trace")).Info("show error")
		//r.logger.With(zap.Error(err), zap.Any("instance", instance), zap.Stack("trace")).Info("show error")
		return
	}
	if r.service.TTL != nil {
		r.logger.With(zap.Any("LeaseID", r.client.LeaseID())).Info("action,register,lease")
	} else {
		r.logger.With().Info("action,register")
	}
}

// Deregister implements the sd.Registrar interface. Call it when you want your
// service to be deregistered from etcd, typically just prior to shutdown.
func (r *Registrar) Deregister() {
	if err := r.client.Deregister(r.service); err != nil {
		r.logger.With(zap.Error(err), zap.Stack("trace")).Info("show error")
		//r.logger.Log("err", err)
	} else {
		r.logger.With().Info("action,deregister")
		//r.logger.Log("action", "deregister")
	}

	r.quitmtx.Lock()
	defer r.quitmtx.Unlock()
	if r.quit != nil {
		close(r.quit)
		r.quit = nil
	}
}
