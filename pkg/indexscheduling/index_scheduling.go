package indexscheduling

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	extenderv1 "k8s.io/kube-scheduler/extender/v1"
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
	if err := compareName(pod.GetName(), nodeInfo.Node().GetName()); err != nil {
		return framework.NewStatus(framework.Unschedulable, "Failed to schedule.", err.Error())
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

func Filter(args *extenderv1.ExtenderArgs) extenderv1.ExtenderFilterResult {
	ret := extenderv1.ExtenderFilterResult{}

	podName := args.Pod.GetName()

	if args.Nodes != nil && args.Nodes.Items != nil {
		retNodes := make([]v1.Node, 0, len(args.Nodes.Items))
		failedNodes := make(map[string]string)
		for _, node := range args.Nodes.Items {
			if err := compareName(podName, node.GetName()); err != nil {
				failedNodes[node.GetName()] = err.Error()
			} else {
				retNodes = append(retNodes, node)
			}
		}

		ret.Nodes = &v1.NodeList{Items: retNodes}
		ret.FailedNodes = failedNodes
	}

	if args.NodeNames != nil {
		retNodeNames := make([]string, 0, len(*args.NodeNames))
		failedNodes := make(map[string]string)
		for _, nodeName := range *args.NodeNames {
			if err := compareName(podName, nodeName); err != nil {
				failedNodes[nodeName] = err.Error()
			} else {
				retNodeNames = append(retNodeNames, nodeName)
			}
		}

		ret.NodeNames = &retNodeNames
		ret.FailedNodes = failedNodes
	}

	return ret
}

func compareName(podName string, nodeName string) error {
	re, _ := regexp.Compile("([0-9]+)$")

	podNum := re.FindString(podName)
	if podNum == "" {
		return fmt.Errorf("pod's name does not end with number")
	}

	hostName := strings.SplitN(nodeName, ".", 2)
	nodeNum := re.FindString(hostName[0])
	if nodeNum == "" {
		return fmt.Errorf("node's name does not end with number")
	}

	podNumber, _ := strconv.ParseInt(podNum, 10, 32)
	nodeNumer, _ := strconv.ParseInt(nodeNum, 10, 32)

	if podNumber != nodeNumer {
		return fmt.Errorf("numner does not match")
	}

	return nil
}
