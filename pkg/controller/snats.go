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
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	"github.com/Sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kubernetes/pkg/controller"
	snatclientset "github.com/noironetworks/aci-containers/pkg/snatpolicy/clientset/versioned"
	snatpolicy "github.com/noironetworks/aci-containers/pkg/snatpolicy/apis/aci.snat/v1"
)

type MyLabel struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type MyPodSelector struct {
	Labels     []MyLabel
	Deployment string
	Namespace  string
}

type MyPortRange struct {
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}

type MySnatPolicy struct {
	SnatIp    []string
	Selector  MyPodSelector
	PortRange []MyPortRange
	Protocols []string
}

func SnatPolicyLogger(log *logrus.Logger, snat *snatpolicy.SnatPolicy) *logrus.Entry {
	return log.WithFields(logrus.Fields{
		"namespace": snat.ObjectMeta.Namespace,
		"name":      snat.ObjectMeta.Name,
		"spec":      snat.Spec,
	})
}

func (cont *AciController) initSnatInformerFromClient(
	snatClient *snatclientset.Clientset) {
	cont.initSnatInformerBase(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
//				options.FieldSelector =
//					fields.Set{"metadata.name": cont.config.KubeConfig}.String()
//					fields.Set{"metadata.name": cont.config.AciVmmDomainType}.String()
				return snatClient.AciV1().SnatPolicies(metav1.NamespaceAll).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
//				options.FieldSelector =
//					fields.Set{"metadata.name": cont.config.KubeConfig}.String()
				return snatClient.AciV1().SnatPolicies(metav1.NamespaceAll).Watch(options)
			},
		})
}

func (cont *AciController) initSnatInformerBase(listWatch *cache.ListWatch) {
	cont.snatInformer = cache.NewSharedIndexInformer(
		listWatch,
		&snatpolicy.SnatPolicy{},
		controller.NoResyncPeriodFunc(),
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	)
	cont.snatInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			cont.log.Debug("POLICY ADDED")
			cont.snatPolicyUpdate(obj)
		},
		UpdateFunc: func(_ interface{}, obj interface{}) {
			cont.snatPolicyUpdate(obj)
		},
		DeleteFunc: func(obj interface{}) {
			cont.snatPolicyDelete(obj)
		},
	})
	cont.log.Debug("Initializing Snat Local info Informers")
}

func (cont *AciController) snatPolicyUpdate(obj interface{}) {
	cont.indexMutex.Lock()
	snat := obj.(*snatpolicy.SnatPolicy)
	key, err := cache.MetaNamespaceKeyFunc(snat)
	if err != nil {
		SnatPolicyLogger(cont.log, snat).
			Error("Could not create key:" + err.Error())
		return
	}
	//cont.log.Info("Snat Policy Object added/Updated ", snat)
	//cont.log.Info("Key produced ", key)
	cont.indexMutex.Unlock()
	cont.doUpdateSnatPolicy(key)
}

func (cont *AciController) doUpdateSnatPolicy(key string) {
	snatobj, exists, err :=
		cont.snatInformer.GetStore().GetByKey(key)
	if err != nil {
		cont.log.Error("Could not lookup snat for " +
			key + ": " + err.Error())
		return
	}
	if !exists || snatobj == nil {
		return
	}
	snat := snatobj.(*snatpolicy.SnatPolicy)
	logger := SnatPolicyLogger(cont.log, snat)
	cont.snatPolicyChanged(snatobj, logger)
}

var Create bool
func (cont *AciController) snatPolicyChanged(snatobj interface{}, logger *logrus.Entry) {
	snatpolicy := snatobj.(*snatpolicy.SnatPolicy)
	cont.indexMutex.Lock()
//	var mypolicy MySnatPolicy
	if len(cont.snatPolicyCache) == 0 {
		cont.log.Debug("EMPTY SNATPOLICYCACHE")
		cont.updateSnatPolicyCache(snatpolicy.ObjectMeta.Name, snatobj)
//		mypolicy = cont.updateSnatPolicyCache(snatpolicy.ObjectMeta.Name, snatobj)
		go cont.updateServiceDeviceInstanceSnat("MYSNAT")
	} else {
		cont.log.Debug("APPENDING TO SNATPOLICYCACHE")
		policyName := snatpolicy.ObjectMeta.Name
		if _, ok := cont.snatPolicyCache[policyName]; ok {
			cont.log.Debug("KEY ALREADY THERE")
//			mypolicy = cont.updateSnatPolicyCache(snatpolicy.ObjectMeta.Name, snatobj)
			cont.updateSnatPolicyCache(snatpolicy.ObjectMeta.Name, snatobj)
			go cont.updateServiceDeviceInstanceSnat("MYSNAT")
		} else {
//			mypolicy = cont.updateSnatPolicyCache(snatpolicy.ObjectMeta.Name, snatobj)
			cont.updateSnatPolicyCache(snatpolicy.ObjectMeta.Name, snatobj)
			go cont.updateServiceDeviceInstanceSnat("MYSNAT")
			cont.log.Debug("KEY NOT THERE")
		}
	}
	cont.log.Debug("map issssssssss ", cont.snatPolicyCache)
	cont.indexMutex.Unlock()
}

func (cont *AciController) updateSnatPolicyCache(key string, snatobj interface{}) MySnatPolicy {
	snatpolicy := snatobj.(*snatpolicy.SnatPolicy)
	var mypolicy MySnatPolicy
	mypolicy.SnatIp = snatpolicy.Spec.SnatIp
	cont.log.Debug("MYPOLICY ISSSSSSSS ....", mypolicy.SnatIp)
	cont.log.Debug("incoing policy isss....", snatpolicy.Spec.SnatIp)
	cont.log.Debug("snat policy cache is ", cont.snatPolicyCache)
	snatLabels := snatpolicy.Spec.Selector.Labels
	snatDeploy := snatpolicy.Spec.Selector.Deployment
	snatNS := snatpolicy.Spec.Selector.Namespace
	var myLabels []MyLabel
	for _, val := range snatLabels {
		lab := MyLabel{Key: val.Key, Value: val.Value}
		myLabels = append(myLabels, lab)
	}
	mypolicy.Selector = MyPodSelector{Labels: myLabels, Deployment: snatDeploy, Namespace: snatNS}
	cont.snatPolicyCache[key] = mypolicy
	return mypolicy
}

func (cont *AciController) snatPolicyDelete(snatobj interface{}) {
	cont.log.Debug("CONTROLLER======= DELETE==== UPDATE")
	cont.indexMutex.Lock()
        snatpolicy := snatobj.(*snatpolicy.SnatPolicy)
//	mypolicy := cont.snatPolicyCache[snatpolicy.ObjectMeta.Name]
	delete(cont.snatPolicyCache, snatpolicy.ObjectMeta.Name)

	cont.log.Debug("cache after deleting is ", cont.snatPolicyCache)
        if len(cont.snatPolicyCache) == 0 {
                cont.log.Debug("Cache is empty now....")
		graphName := cont.aciNameForKey("snat", "MYSNAT")
		//vk8s_1_snat_SNAT
		go cont.apicConn.ClearApicObjects(graphName)
//		go cont.apicConn.ClearApicObjects(cont.aciNameForKey("snat-vmm", "snat"))
        } else {
		go cont.updateServiceDeviceInstanceSnat("MYSNAT")
	}
	cont.indexMutex.Unlock()
}
