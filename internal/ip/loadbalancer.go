// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ip

import (
	"errors"
	"fmt"
	"net"

	"github.com/siderolabs/go-loadbalancer/loadbalancer"
	"github.com/siderolabs/go-loadbalancer/upstream"
	"go.uber.org/zap"
)

// LoadBalancer is an interface for loadbalancer instances.
type LoadBalancer interface {
	AddRoute(ipPort string, upstreamAddrs []string, options ...upstream.ListOption) error
	Start() error
	Close() error
}

// LoadBalancerProvider is a factory for LoadBalancer instances.
type LoadBalancerProvider interface {
	New(logger *zap.Logger) (LoadBalancer, error)
}

// TCPLoadBalancerProvider is a LoadBalancerProvider that creates and returns loadbalancer.TCP instances.
type TCPLoadBalancerProvider struct{}

// New returns a new loadbalancer.TCP instance.
func (t *TCPLoadBalancerProvider) New(logger *zap.Logger) (LoadBalancer, error) {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &tcpLoadBalancer{
		TCP: loadbalancer.TCP{
			Logger: logger,
		},
	}, nil
}

// tcpLoadBalancer is a wrapper around loadbalancer.TCP.
type tcpLoadBalancer struct {
	loadbalancer.TCP
}

// Close closes the TCP load balancer and waits for it to stop.
func (lb *tcpLoadBalancer) Close() error {
	if err := lb.TCP.Close(); err != nil {
		return fmt.Errorf("failed to close TCP load balancer: %w", err)
	}

	if err := lb.TCP.Wait(); err != nil && !errors.Is(err, net.ErrClosed) {
		return fmt.Errorf("failed to wait for TCP load balancer: %w", err)
	}

	return nil
}
