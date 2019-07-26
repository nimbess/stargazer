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
	CRDPlural   string = "unpconfigs"
	CRDGroup    string = "nimbess.com"
	CRDVersion  string = "v1"
	FullCRDName string = CRDPlural + "." + CRDGroup
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type UnpConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              UnpConfigSpec   `json:"spec"`
	Status            UnpConfigStatus `json:"status,omitempty"`
}

type UnpConfigSpec struct {
	Type       string `json:"type"`
	PodLabel   string `json:"podLabel"`
	Attributes string `json:"attributes"`
}

type UnpConfigStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type UnpConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []UnpConfig `json:"items"`
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
				Kind:   reflect.TypeOf(UnpConfig{}).Name(),
			},
		},
	}

	_, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil && apierrors.IsAlreadyExists(err) {
		log.Info("UNP CRD successfully registered")
		return nil
	}
	if err == nil {
		log.Info("UNP CRD successfully registered")
	}
	return err
}
