package hostprovider

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/internal/freezer"
	"time"
)

type HostProvider interface {
	Provider() (string, error)
	Freeze(host string, cause error, duration time.Duration) error
}

func NewWithHosts(hosts []string) HostProvider {
	return &arrayHostProvider{
		hosts:   hosts,
		freezer: freezer.New(),
	}
}

type arrayHostProvider struct {
	hosts         []string
	freezer       freezer.Freezer
	lastFreezeErr error
}

func (a *arrayHostProvider) Provider() (string, error) {
	if len(a.hosts) == 0 {
		return "", errors.New("no host found")
	}

	for _, host := range a.hosts {
		if a.freezer.Available(host) {
			return host, nil
		}
	}

	if a.lastFreezeErr != nil {
		return "", a.lastFreezeErr
	} else {
		return "", errors.New("all hosts are frozen")
	}
}

func (a *arrayHostProvider) Freeze(host string, cause error, duration time.Duration) error {
	if duration <= 0 {
		return nil
	}

	a.lastFreezeErr = cause
	return a.freezer.Freeze(host, duration)
}
