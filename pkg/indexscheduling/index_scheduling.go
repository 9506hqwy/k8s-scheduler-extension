package indexscheduling

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

const (
	Name = "IndexScheduling"
)

type IndexScheduling struct {
}

func (c *IndexScheduling) Name() string {
	return Name
}

func (c *IndexScheduling) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	re, err := regexp.Compile("([0-9]+)$")
	if err != nil {
		return framework.NewStatus(framework.Error, "Failed to compile regex.")
	}

	podNum := re.FindString(pod.GetName())
	if podNum == "" {
		return framework.NewStatus(framework.Unschedulable, "Failed to schedule.", "Pod's name does not end with number.")
	}

	hostName := strings.SplitN(nodeInfo.Node().GetName(), ".", 2)
	nodeNum := re.FindString(hostName[0])
	if nodeNum == "" {
		return framework.NewStatus(framework.Unschedulable, "Failed to schedule.", "Node's name does not end with number.")
	}

	podNumber, _ := strconv.ParseInt(podNum, 10, 32)
	nodeNumer, _ := strconv.ParseInt(nodeNum, 10, 32)

	if podNumber != nodeNumer {
		return framework.NewStatus(framework.Unschedulable, "Failed to schedule.", "Numner does not match.")
	}

	s := framework.NewStatus(framework.Success, "")
	s.SetPlugin(Name)
	return s
}

func New(ctx context.Context, obj runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	s := &IndexScheduling{}
	logger := klog.FromContext(ctx)

	logger.Info("IndexScheduling start")
	return s, nil
}
