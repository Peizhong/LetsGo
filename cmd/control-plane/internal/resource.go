// Copyright 2020 Envoyproxy Authors
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
package internal

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	apiv2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	apiv2core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	listener "github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	apiv2route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v2"
	resourcev2 "github.com/envoyproxy/go-control-plane/pkg/resource/v2"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
)

const (
	ClusterName  = "example_proxy_cluster"
	RouteName    = "local_route"
	ListenerName = "listener_0"
	ListenerPort = 10000
	UpstreamHost = "localhost"
	UpstreamPort = 8000
)

func makeCluster(clusterName string) *apiv2.Cluster {
	return &apiv2.Cluster{
		Name:                 clusterName,
		ConnectTimeout:       ptypes.DurationProto(5 * time.Second),
		ClusterDiscoveryType: &apiv2.Cluster_Type{Type: apiv2.Cluster_LOGICAL_DNS},
		LbPolicy:             apiv2.Cluster_ROUND_ROBIN,
		LoadAssignment:       makeEndpoint(clusterName),
		DnsLookupFamily:      apiv2.Cluster_V4_ONLY,
	}
}

func makeEndpoint(clusterName string) *apiv2.ClusterLoadAssignment {
	return &apiv2.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []*endpoint.LocalityLbEndpoints{{
			LbEndpoints: []*endpoint.LbEndpoint{{
				HostIdentifier: &endpoint.LbEndpoint_Endpoint{
					Endpoint: &endpoint.Endpoint{
						Address: &apiv2core.Address{
							Address: &apiv2core.Address_SocketAddress{
								SocketAddress: &apiv2core.SocketAddress{
									Protocol: apiv2core.SocketAddress_TCP,
									Address:  UpstreamHost,
									PortSpecifier: &apiv2core.SocketAddress_PortValue{
										PortValue: UpstreamPort,
									},
								},
							},
						},
					},
				},
			}},
		}},
	}
}

func makeRoute(routeName string, clusterName string) *apiv2.RouteConfiguration {
	return &apiv2.RouteConfiguration{
		Name: routeName,
		VirtualHosts: []*apiv2route.VirtualHost{{
			Name:    "local_service",
			Domains: []string{"*"},
			Routes: []*apiv2route.Route{{
				Match: &apiv2route.RouteMatch{
					PathSpecifier: &apiv2route.RouteMatch_Prefix{
						Prefix: "/",
					},
				},
				Action: &apiv2route.Route_Route{
					Route: &apiv2route.RouteAction{
						ClusterSpecifier: &apiv2route.RouteAction_Cluster{
							Cluster: clusterName,
						},
						HostRewriteSpecifier: &apiv2route.RouteAction_HostRewrite{
							HostRewrite: UpstreamHost,
						},
					},
				},
			}},
		}},
	}
}

func makeHTTPListener(listenerName string, route string) *apiv2.Listener {
	// HTTP filter configuration
	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.HttpConnectionManager_AUTO,
		StatPrefix: "http",
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				ConfigSource:    makeConfigSource(),
				RouteConfigName: route,
			},
		},
		HttpFilters: []*hcm.HttpFilter{{
			Name: wellknown.Router,
		}},
	}
	pbst, err := ptypes.MarshalAny(manager)
	if err != nil {
		panic(err)
	}

	return &apiv2.Listener{
		Name: listenerName,
		Address: &apiv2core.Address{
			Address: &apiv2core.Address_SocketAddress{
				SocketAddress: &apiv2core.SocketAddress{
					Protocol: apiv2core.SocketAddress_TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &apiv2core.SocketAddress_PortValue{
						PortValue: ListenerPort,
					},
				},
			},
		},
		FilterChains: []*listener.FilterChain{{
			Filters: []*listener.Filter{{
				Name: wellknown.HTTPConnectionManager,
				ConfigType: &listener.Filter_TypedConfig{
					TypedConfig: pbst,
				},
			}},
		}},
	}
}

func makeConfigSource() *apiv2core.ConfigSource {
	source := &apiv2core.ConfigSource{}
	source.ResourceApiVersion = resourcev2.DefaultAPIVersion
	source.ConfigSourceSpecifier = &apiv2core.ConfigSource_ApiConfigSource{
		ApiConfigSource: &apiv2core.ApiConfigSource{
			TransportApiVersion:       resourcev2.DefaultAPIVersion,
			ApiType:                   apiv2core.ApiConfigSource_GRPC,
			SetNodeOnFirstMessageOnly: true,
			GrpcServices: []*apiv2core.GrpcService{{
				TargetSpecifier: &apiv2core.GrpcService_EnvoyGrpc_{
					EnvoyGrpc: &apiv2core.GrpcService_EnvoyGrpc{ClusterName: "xds_cluster"},
				},
			}},
		},
	}
	return source
}

func GenerateSnapshot() cache.Snapshot {
	return cache.NewSnapshot(
		"1",
		[]types.Resource{}, // endpoints
		[]types.Resource{makeCluster(ClusterName)},
		[]types.Resource{makeRoute(RouteName, ClusterName)},
		[]types.Resource{makeHTTPListener(ListenerName, RouteName)},
		[]types.Resource{}, // runtimes
	)
}
