// Copyright 2019 Cisco Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	nodeinfo "github.com/noironetworks/aci-containers/pkg/nodeinfo/apis/aci.snat/v1"
	snatglobalinfo "github.com/noironetworks/aci-containers/pkg/snatglobalinfo/apis/aci.snat/v1"
	snatpolicy "github.com/noironetworks/aci-containers/pkg/snatpolicy/apis/aci.snat/v1"
	tu "github.com/noironetworks/aci-containers/pkg/testutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

var dummy struct{}

func snatpolicydata(name string, namespace string,
	snatIp []string, labels map[string]string) *snatpolicy.SnatPolicy {
	policy := &snatpolicy.SnatPolicy{
		Spec: snatpolicy.SnatPolicySpec{
			SnatIp: snatIp,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	var podSelector snatpolicy.PodSelector
	podSelector.Namespace = namespace
	podSelector.Labels = labels
	policy.Spec.Selector = podSelector
	return policy
}

func Nodeinfodata(name string, namespace string, macaddr string, snatplcy map[string]struct{}) *nodeinfo.NodeInfo {
	info := &nodeinfo.NodeInfo{
		Spec: nodeinfo.NodeInfoSpec{
			SnatPolicyNames: snatplcy,
			Macaddress:      macaddr,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	return info
}

type policy struct {
	namespace string
	name      string
	snatip    []string
	labels    map[string]string
}

type nodedata struct {
	namespace       string
	name            string
	snatpolicynames map[string]struct{}
	macaddr         string
}

var snatTests = []policy{
	{
		"testns",
		"policy1",
		[]string{"10.1.1.8"},
		map[string]string{"key": "value"},
	},
	{
		"testns",
		"policy2",
		[]string{"10.1.1.9"},
		map[string]string{"key": "value"},
	},
}

var nodeTests = []nodedata{
	{
		"testns",
		"node-1",
		map[string]struct{}{"policy1": dummy},
		"01:02:03:04",
	},
	{
		"testns",
		"node-2",
		map[string]struct{}{"policy1": dummy},
		"01:02:03:05",
	},
	{
		"testns",
		"node-1",
		map[string]struct{}{"policy2": dummy, "policy1": dummy},
		"01:02:03:04",
	},
}

func snatWait(t *testing.T, desc string, expected map[string]ContSnatGlobalInfo, actual map[string]*ContSnatGlobalInfo) {
	tu.WaitFor(t, desc, 500*time.Millisecond, func(last bool) (bool, error) {
		for key, v := range expected {
			val, ok := actual[key]
			if ok {
				if v.SnatIp != val.SnatIp && v.SnatPortRange != val.SnatPortRange {
					return false, nil
				}
			} else {
				return false, nil
			}
		}
		return true, nil
	})

}

func TestSnatnodeInfo(t *testing.T) {
	cont := testController()
	cont.run()
	for _, pt := range snatTests {
		snatObj := snatpolicydata(pt.name, pt.namespace, pt.snatip, pt.labels)
		cont.fakeSnatPolicySource.Add(snatObj)
	}
	nodinfo := make(map[string]bool)
	for _, pt := range nodeTests {
		nodeobj := Nodeinfodata(pt.name, pt.namespace, pt.macaddr, pt.snatpolicynames)
		if _, ok := nodinfo[pt.name]; !ok {
			cont.fakeNodeInfoSource.Add(nodeobj)
			cont.log.Debug("NodeInfo Added: ", nodeobj)
			nodinfo[pt.name] = true
		} else {
			cont.log.Debug("NodeInfo Modified: ", nodeobj)
			cont.fakeNodeInfoSource.Modify(nodeobj)
		}
	}
	time.Sleep(2 * time.Second)
	cont.log.Debug("snatGlobalInfoCache: ", cont.AciController.snatGlobalInfoCache)
	expected := map[string]ContSnatGlobalInfo{
		"node-1": ContSnatGlobalInfo{SnatIp: "10.1.1.8", SnatPortRange: snatglobalinfo.PortRange{Start: 5000, End: 7999}},
		"node-2": ContSnatGlobalInfo{SnatIp: "10.1.1.8", SnatPortRange: snatglobalinfo.PortRange{Start: 8000, End: 10999}},
	}
	expected1 := map[string]ContSnatGlobalInfo{
		"node-1": ContSnatGlobalInfo{SnatIp: "10.1.1.9", SnatPortRange: snatglobalinfo.PortRange{Start: 5000, End: 7999}},
	}
	snatWait(t, "snat test", expected, cont.AciController.snatGlobalInfoCache["10.1.1.8"])
	snatWait(t, "snat test", expected1, cont.AciController.snatGlobalInfoCache["10.1.1.9"])
	cont.stop()
}
