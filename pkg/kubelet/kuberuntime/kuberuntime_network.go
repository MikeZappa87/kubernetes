package kuberuntime

import (
	"context"
	"errors"

	"github.com/MikeZappa87/kni-api/pkg/apis/runtime/beta"
	v1 "k8s.io/api/core/v1"
	ref "k8s.io/client-go/tools/reference"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	"k8s.io/kubernetes/pkg/kubelet/events"
)

func (m *kubeGenericRuntimeManager) AttachInterface(ctx context.Context, result *kubecontainer.PodSyncResult,
	pod *v1.Pod, podSandboxID string) (*runtimeapi.PodSandboxStatus, error) {

	resp, err := m.runtimeService.PodSandboxStatus(ctx, podSandboxID, false)
	if err != nil {
		ref, referr := ref.GetReference(legacyscheme.Scheme, pod)
		if referr != nil {
			klog.ErrorS(referr, "Couldn't make a ref to pod", "pod", klog.KObj(pod))
		}
		m.recorder.Eventf(ref, v1.EventTypeWarning, events.FailedStatusPodSandBox, "Unable to get pod sandbox status: %v", err)
		klog.ErrorS(err, "Failed to get pod sandbox status; Skipping pod", "pod", klog.KObj(pod))
		result.Fail(err)
		return nil, err
	}
	if resp.GetStatus() == nil {
		result.Fail(errors.New("pod sandbox status is nil"))
		return nil, err
	}

	dnsConfig, err := m.runtimeHelper.GetPodDNS(pod)
	if err != nil {
		return nil, err
	}

	portMappings := []*runtimeapi.PortMapping{}
	for _, c := range pod.Spec.Containers {
		containerPortMappings := kubecontainer.MakePortMappings(&c)

		for idx := range containerPortMappings {
			port := containerPortMappings[idx]
			hostPort := int32(port.HostPort)
			containerPort := int32(port.ContainerPort)
			protocol := toRuntimeProtocol(port.Protocol)
			portMappings = append(portMappings, &runtimeapi.PortMapping{
				HostIp:        port.HostIP,
				HostPort:      hostPort,
				ContainerPort: containerPort,
				Protocol:      protocol,
			})
		}
	}

	_, err = m.networkService.AttachInterface(ctx, &beta.AttachInterfaceRequest{
		Id:           podSandboxID,
		Labels:       resp.Status.GetLabels(),
		Annotations:  resp.Status.GetAnnotations(),
		Extradata:    resp.Info,
		DnsConfig:    toKNIDnsConfig(dnsConfig),
		PortMappings: toKNIPortMapping(portMappings),
		Name: 		  pod.Name,
		Namespace:    pod.Namespace,
	})

	if err != nil {
		ref, referr := ref.GetReference(legacyscheme.Scheme, pod)
		if referr != nil {
			klog.ErrorS(referr, "Couldn't make a ref to pod", "pod", klog.KObj(pod))
		}

		m.recorder.Eventf(ref, v1.EventTypeWarning, events.FailedNetworkAttach, "Unable to attach network to pod: %v", err)
		klog.ErrorS(err, "Failed to attach pod network; Skipping pod", "pod", klog.KObj(pod))
		result.Fail(err)

		return nil, err
	}

	netResp, err := m.networkService.QueryPodNetwork(ctx, podSandboxID)

	if err != nil {
		ref, referr := ref.GetReference(legacyscheme.Scheme, pod)
		if referr != nil {
			klog.ErrorS(referr, "Couldn't make a ref to pod", "pod", klog.KObj(pod))
		}

		m.recorder.Eventf(ref, v1.EventTypeWarning, events.FailedNetworkPodStatus, "Unable to query pod network status: %v", err)
		klog.ErrorS(err, "Failed to query pod network status; Skipping pod", "pod", klog.KObj(pod))
		result.Fail(err)

		return nil, err
	}

	resp.Status.Network = m.toCRINetworkStatus(netResp, "eth0")

	return resp.Status, nil
}

func (m *kubeGenericRuntimeManager) DetachInterface(ctx context.Context, podSandboxId string) error {
	err := m.networkService.DetachInterface(ctx, podSandboxId)

	if err != nil {
		return err
	}

	return nil
}

func (m *kubeGenericRuntimeManager) toCRINetworkStatus(net *beta.QueryPodNetworkResponse, ifName string) *runtimeapi.PodSandboxNetworkStatus {
	var podIps []*runtimeapi.PodIP

	if net.Ipconfigs == nil {
		return &runtimeapi.PodSandboxNetworkStatus{
			Ip: "",
		}
	}

	if val, ok := net.Ipconfigs[ifName]; ok {
		for _, v := range val.Ip {
			podIps = append(podIps, &runtimeapi.PodIP{
				Ip: v,
			})
		}

		if len(podIps) == 0 {
			return &runtimeapi.PodSandboxNetworkStatus{
				Ip: "",
			}
		}

		return &runtimeapi.PodSandboxNetworkStatus{
			Ip:            podIps[0].Ip,
			AdditionalIps: podIps[1:],
		}
	}

	return &runtimeapi.PodSandboxNetworkStatus{
		Ip: "",
	}
}

func toKNIDnsConfig(dns *runtimeapi.DNSConfig) *beta.DNSConfig {
	return &beta.DNSConfig{
		Servers:  dns.Servers,
		Options:  dns.Options,
		Searches: dns.Searches,
	}
}

func toKNIPortMapping(mapping []*runtimeapi.PortMapping) []*beta.PortMapping {
	var ports []*beta.PortMapping

	for _, v := range mapping {
		ports = append(ports, &beta.PortMapping{
			Protocol:      beta.Protocol(v.Protocol),
			ContainerPort: v.ContainerPort,
			HostPort:      v.HostPort,
			HostIp:        v.HostIp,
		})
	}

	return ports
}