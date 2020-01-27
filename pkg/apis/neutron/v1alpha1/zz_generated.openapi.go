// +build !ignore_autogenerated

// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/neutron-operator/pkg/apis/neutron/v1alpha1.OvsAgent":       schema_pkg_apis_neutron_v1alpha1_OvsAgent(ref),
		"github.com/neutron-operator/pkg/apis/neutron/v1alpha1.OvsAgentSpec":   schema_pkg_apis_neutron_v1alpha1_OvsAgentSpec(ref),
		"github.com/neutron-operator/pkg/apis/neutron/v1alpha1.OvsAgentStatus": schema_pkg_apis_neutron_v1alpha1_OvsAgentStatus(ref),
	}
}

func schema_pkg_apis_neutron_v1alpha1_OvsAgent(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "OvsAgent is the Schema for the ovsagents API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/neutron-operator/pkg/apis/neutron/v1alpha1.OvsAgentSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/neutron-operator/pkg/apis/neutron/v1alpha1.OvsAgentStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/neutron-operator/pkg/apis/neutron/v1alpha1.OvsAgentSpec", "github.com/neutron-operator/pkg/apis/neutron/v1alpha1.OvsAgentStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_neutron_v1alpha1_OvsAgentSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "OvsAgentSpec defines the desired state of OvsAgent",
				Properties: map[string]spec.Schema{
					"label": {
						SchemaProps: spec.SchemaProps{
							Description: "Label is the value of the 'daemon=' label to set on a node that should run the daemon",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"openvswitchImage": {
						SchemaProps: spec.SchemaProps{
							Description: "Image is the Docker image to run for the daemon",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"label", "openvswitchImage"},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_neutron_v1alpha1_OvsAgentStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "OvsAgentStatus defines the observed state of OvsAgent",
				Properties: map[string]spec.Schema{
					"count": {
						SchemaProps: spec.SchemaProps{
							Description: "Count is the number of nodes the daemon is deployed to",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
				},
				Required: []string{"count"},
			},
		},
		Dependencies: []string{},
	}
}
