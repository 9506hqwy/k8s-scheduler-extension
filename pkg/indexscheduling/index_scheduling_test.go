package indexscheduling_test

import (
	"context"
	"testing"

	"github.com/9506hqwy/k8s-scheduler-extension/pkg/indexscheduling"
	v1 "k8s.io/api/core/v1"
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
