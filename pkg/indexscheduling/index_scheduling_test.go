package indexscheduling_test

import (
	"context"
	"testing"

	"github.com/9506hqwy/k8s-scheduler-extension/pkg/indexscheduling"
	v1 "k8s.io/api/core/v1"
	extenderv1 "k8s.io/kube-scheduler/extender/v1"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

func Test_Name(t *testing.T) {
	name := indexscheduling.Name
	if name != "IndexScheduling" {
		t.Error("Failed.", name)
	}
}

func Test_FilterScheduling(t *testing.T) {
	i := indexscheduling.IndexScheduling{}
	p := v1.Pod{}
	p.SetName("Pod01")
	n := v1.Node{}
	n.SetName("Node01")
	ni := framework.NodeInfo{}
	ni.SetNode(&n)

	s := i.Filter(context.Background(), framework.NewCycleState(), &p, &ni)
	if s.Code() != framework.Success {
		t.Error("Failed.", s)
	}
}

func Test_FilterSchedulingWithDomain(t *testing.T) {
	i := indexscheduling.IndexScheduling{}
	p := v1.Pod{}
	p.SetName("Pod01")
	n := v1.Node{}
	n.SetName("Node01.domain")
	ni := framework.NodeInfo{}
	ni.SetNode(&n)

	s := i.Filter(context.Background(), framework.NewCycleState(), &p, &ni)
	if s.Code() != framework.Success {
		t.Error("Failed.", s)
	}
}

func Test_FilterNoMatchName(t *testing.T) {
	i := indexscheduling.IndexScheduling{}
	p := v1.Pod{}
	p.SetName("Pod01")
	n := v1.Node{}
	n.SetName("Node02")
	ni := framework.NodeInfo{}
	ni.SetNode(&n)

	s := i.Filter(context.Background(), framework.NewCycleState(), &p, &ni)
	if s.Code() != framework.Unschedulable {
		t.Error("Failed.", s)
	}
}

func Test_FilterPodNameNotNumber(t *testing.T) {
	i := indexscheduling.IndexScheduling{}
	p := v1.Pod{}
	p.SetName("PodAA")
	n := v1.Node{}
	n.SetName("Node01")
	ni := framework.NodeInfo{}
	ni.SetNode(&n)

	s := i.Filter(context.Background(), framework.NewCycleState(), &p, &ni)
	if s.Code() != framework.Unschedulable {
		t.Error("Failed.", s)
	}
}

func Test_FilterNodeNameNotNumber(t *testing.T) {
	i := indexscheduling.IndexScheduling{}
	p := v1.Pod{}
	p.SetName("Pod01")
	n := v1.Node{}
	n.SetName("NodeAA")
	ni := framework.NodeInfo{}
	ni.SetNode(&n)

	s := i.Filter(context.Background(), framework.NewCycleState(), &p, &ni)
	if s.Code() != framework.Unschedulable {
		t.Error("Failed.", s)
	}
}

func Test_FilterFuncNoCache(t *testing.T) {
	args := extenderv1.ExtenderArgs{
		Pod:   &v1.Pod{},
		Nodes: &v1.NodeList{},
	}

	args.Pod.SetName("Pod01")

	args.Nodes.Items = []v1.Node{{}, {}}
	args.Nodes.Items[0].SetName("Node01")
	args.Nodes.Items[1].SetName("Node02")

	ret := indexscheduling.Filter(&args)

	if ret.Nodes.Items[0].GetName() != "Node01" {
		t.Error("Failed.", ret.Nodes.Items[0])
	}

	if _, ok := ret.FailedNodes["Node02"]; !ok {
		t.Error("Failed.", ret.FailedNodes)
	}
}

func Test_FilterFuncCache(t *testing.T) {
	args := extenderv1.ExtenderArgs{
		Pod: &v1.Pod{},
	}

	args.Pod.SetName("Pod01")

	nodeNames := make([]string, 0, 2)
	nodeNames = append(nodeNames, "Node01")
	nodeNames = append(nodeNames, "Node02")
	args.NodeNames = &nodeNames

	ret := indexscheduling.Filter(&args)

	if (*ret.NodeNames)[0] != "Node01" {
		t.Error("Failed.", ret.NodeNames)
	}

	if _, ok := ret.FailedNodes["Node02"]; !ok {
		t.Error("Failed.", ret.FailedNodes)
	}
}
