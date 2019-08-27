package v1

import (
	log "github.com/sirupsen/logrus"
	apiextensionv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
)

const (
	CRDPlural   string = "unifiednetworkpolicies"
	CRDGroup    string = "nimbess.com"
	CRDVersion  string = "v1"
	FullCRDName string = CRDPlural + "." + CRDGroup
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type UnifiedNetworkPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              UnifiedNetworkPolicySpec   `json:"spec"`
	Status            UnifiedNetworkPolicyStatus `json:"status,omitempty"`
}

type DefaultPolicy struct {
	Action string `json:"action,omitempty"`
}

type URLFilter struct {
	Urls        []string             `json:"urls,omitempty"`
	Action      string               `json:"action,omitempty"`
	PodSelector metav1.LabelSelector `json:"podSelector,omitempty"`
	Network     string               `json:"network,omitempty"`
}

type L7Policy struct {
	Default   DefaultPolicy `json:"default,omitempty"`
	UrlFilter URLFilter     `json:"urlFilter,omitempty"`
}

type UnifiedNetworkPolicySpec struct {
	L7Policies  []L7Policy           `json:"l7Policies"`
	PodSelector metav1.LabelSelector `json:"podSelector"`
	Network     string               `json:"network"`
	Attributes  string               `json:"attributes"`
}

type UnifiedNetworkPolicyStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type UnifiedNetworkPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []UnifiedNetworkPolicy `json:"items"`
}

func CreateCRD(clientset *clientset.Clientset) error {
	ver := apiextensionv1beta1.CustomResourceDefinitionVersion{Name: CRDVersion, Served: true, Storage: true}
	crd := &apiextensionv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{Name: FullCRDName},
		Spec: apiextensionv1beta1.CustomResourceDefinitionSpec{
			Group:    CRDGroup,
			Versions: []apiextensionv1beta1.CustomResourceDefinitionVersion{ver},
			Scope:    apiextensionv1beta1.NamespaceScoped,
			Names: apiextensionv1beta1.CustomResourceDefinitionNames{
				Plural: CRDPlural,
				Kind:   reflect.TypeOf(UnifiedNetworkPolicy{}).Name(),
			},
		},
	}
	_, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil && apierrors.IsAlreadyExists(err) {
		log.Info("UNP CRD already registered")
		err = clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(FullCRDName, &metav1.DeleteOptions{})
		if err != nil {
			log.Fatalf("Error deleting existing CRD: %v", err)
		}
		_, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
		if err != nil {
			log.Fatalf("Failed to create UNP CRD: %v", err)
		}
	}
	if err == nil {
		log.Info("UNP CRD successfully registered")
	}
	return err
}
