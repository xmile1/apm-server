// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package otel_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/model/pdata"

	"github.com/elastic/apm-server/model"
)

func TestResourceConventions(t *testing.T) {
	defaultAgent := model.Agent{Name: "otlp", Version: "unknown"}
	defaultService := model.Service{
		Name:     "unknown",
		Language: model.Language{Name: "unknown"},
	}

	for name, test := range map[string]struct {
		attrs    map[string]pdata.AttributeValue
		expected model.APMEvent
	}{
		"empty": {
			attrs:    nil,
			expected: model.APMEvent{Agent: defaultAgent, Service: defaultService},
		},
		"service": {
			attrs: map[string]pdata.AttributeValue{
				"service.name":           pdata.NewAttributeValueString("service_name"),
				"service.version":        pdata.NewAttributeValueString("service_version"),
				"deployment.environment": pdata.NewAttributeValueString("service_environment"),
			},
			expected: model.APMEvent{
				Agent: model.Agent{Name: "otlp", Version: "unknown"},
				Service: model.Service{
					Name:        "service_name",
					Version:     "service_version",
					Environment: "service_environment",
					Language:    model.Language{Name: "unknown"},
				},
			},
		},
		"agent": {
			attrs: map[string]pdata.AttributeValue{
				"telemetry.sdk.name":     pdata.NewAttributeValueString("sdk_name"),
				"telemetry.sdk.version":  pdata.NewAttributeValueString("sdk_version"),
				"telemetry.sdk.language": pdata.NewAttributeValueString("language_name"),
			},
			expected: model.APMEvent{
				Agent: model.Agent{Name: "sdk_name/language_name", Version: "sdk_version"},
				Service: model.Service{
					Name:     "unknown",
					Language: model.Language{Name: "language_name"},
				},
			},
		},
		"runtime": {
			attrs: map[string]pdata.AttributeValue{
				"process.runtime.name":    pdata.NewAttributeValueString("runtime_name"),
				"process.runtime.version": pdata.NewAttributeValueString("runtime_version"),
			},
			expected: model.APMEvent{
				Agent: model.Agent{Name: "otlp", Version: "unknown"},
				Service: model.Service{
					Name:     "unknown",
					Language: model.Language{Name: "unknown"},
					Runtime: model.Runtime{
						Name:    "runtime_name",
						Version: "runtime_version",
					},
				},
			},
		},
		"cloud": {
			attrs: map[string]pdata.AttributeValue{
				"cloud.provider":          pdata.NewAttributeValueString("provider_name"),
				"cloud.region":            pdata.NewAttributeValueString("region_name"),
				"cloud.account.id":        pdata.NewAttributeValueString("account_id"),
				"cloud.availability_zone": pdata.NewAttributeValueString("availability_zone"),
				"cloud.platform":          pdata.NewAttributeValueString("platform_name"),
			},
			expected: model.APMEvent{
				Agent:   defaultAgent,
				Service: defaultService,
				Cloud: model.Cloud{
					Provider:         "provider_name",
					Region:           "region_name",
					AccountID:        "account_id",
					AvailabilityZone: "availability_zone",
					ServiceName:      "platform_name",
				},
			},
		},
		"container": {
			attrs: map[string]pdata.AttributeValue{
				"container.name":       pdata.NewAttributeValueString("container_name"),
				"container.id":         pdata.NewAttributeValueString("container_id"),
				"container.image.name": pdata.NewAttributeValueString("container_image_name"),
				"container.image.tag":  pdata.NewAttributeValueString("container_image_tag"),
				"container.runtime":    pdata.NewAttributeValueString("container_runtime"),
			},
			expected: model.APMEvent{
				Agent:   defaultAgent,
				Service: defaultService,
				Container: model.Container{
					Name:      "container_name",
					ID:        "container_id",
					Runtime:   "container_runtime",
					ImageName: "container_image_name",
					ImageTag:  "container_image_tag",
				},
			},
		},
		"kubernetes": {
			attrs: map[string]pdata.AttributeValue{
				"k8s.namespace.name": pdata.NewAttributeValueString("kubernetes_namespace"),
				"k8s.node.name":      pdata.NewAttributeValueString("kubernetes_node_name"),
				"k8s.pod.name":       pdata.NewAttributeValueString("kubernetes_pod_name"),
				"k8s.pod.uid":        pdata.NewAttributeValueString("kubernetes_pod_uid"),
			},
			expected: model.APMEvent{
				Agent:   defaultAgent,
				Service: defaultService,
				Kubernetes: model.Kubernetes{
					Namespace: "kubernetes_namespace",
					NodeName:  "kubernetes_node_name",
					PodName:   "kubernetes_pod_name",
					PodUID:    "kubernetes_pod_uid",
				},
			},
		},
		"host": {
			attrs: map[string]pdata.AttributeValue{
				"host.name": pdata.NewAttributeValueString("host_name"),
				"host.id":   pdata.NewAttributeValueString("host_id"),
				"host.type": pdata.NewAttributeValueString("host_type"),
				"host.arch": pdata.NewAttributeValueString("host_arch"),
			},
			expected: model.APMEvent{
				Agent:   defaultAgent,
				Service: defaultService,
				Host: model.Host{
					Hostname:     "host_name",
					ID:           "host_id",
					Type:         "host_type",
					Architecture: "host_arch",
				},
			},
		},
		"process": {
			attrs: map[string]pdata.AttributeValue{
				"process.pid":             pdata.NewAttributeValueInt(123),
				"process.command_line":    pdata.NewAttributeValueString("command_line"),
				"process.executable.path": pdata.NewAttributeValueString("executable_path"),
			},
			expected: model.APMEvent{
				Agent:   defaultAgent,
				Service: defaultService,
				Process: model.Process{
					Pid:         123,
					CommandLine: "command_line",
					Executable:  "executable_path",
				},
			},
		},
		"os": {
			attrs: map[string]pdata.AttributeValue{
				"os.type":        pdata.NewAttributeValueString("DARWIN"),
				"os.description": pdata.NewAttributeValueString("Mac OS Mojave"),
			},
			expected: model.APMEvent{
				Agent:   defaultAgent,
				Service: defaultService,
				Host: model.Host{
					OS: model.OS{
						Platform: "darwin",
						Type:     "macos",
						Full:     "Mac OS Mojave",
					},
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			meta := transformResourceMetadata(t, test.attrs)
			assert.Equal(t, test.expected, meta)
		})
	}
}

func TestResourceLabels(t *testing.T) {
	stringArray := pdata.NewAttributeValueArray()
	stringArray.ArrayVal().AppendEmpty().SetStringVal("abc")
	stringArray.ArrayVal().AppendEmpty().SetStringVal("def")

	intArray := pdata.NewAttributeValueArray()
	intArray.ArrayVal().AppendEmpty().SetIntVal(123)
	intArray.ArrayVal().AppendEmpty().SetIntVal(456)

	metadata := transformResourceMetadata(t, map[string]pdata.AttributeValue{
		"string_array": stringArray,
		"int_array":    intArray,
	})
	assert.Equal(t, model.Labels{
		"string_array": {Values: []string{"abc", "def"}},
	}, metadata.Labels)
	assert.Equal(t, model.NumericLabels{
		"int_array": {Values: []float64{123, 456}},
	}, metadata.NumericLabels)
}

func transformResourceMetadata(t *testing.T, resourceAttrs map[string]pdata.AttributeValue) model.APMEvent {
	traces, spans := newTracesSpans()
	traces.ResourceSpans().At(0).Resource().Attributes().InitFromMap(resourceAttrs)
	otelSpan := spans.Spans().AppendEmpty()
	otelSpan.SetTraceID(pdata.NewTraceID([16]byte{1}))
	otelSpan.SetSpanID(pdata.NewSpanID([8]byte{2}))
	events := transformTraces(t, traces)
	events[0].Transaction = nil
	events[0].Trace = model.Trace{}
	events[0].Event.Outcome = ""
	events[0].Timestamp = time.Time{}
	events[0].Processor = model.Processor{}
	return events[0]
}
