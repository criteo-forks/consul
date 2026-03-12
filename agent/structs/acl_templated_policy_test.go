// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package structs

import (
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/require"
)

func TestDeduplicate(t *testing.T) {
	type testCase struct {
		templatedPolicies ACLTemplatedPolicies
		expectedCount     int
	}
	tcases := map[string]testCase{
		"multiple-of-the-same-template": {
			templatedPolicies: ACLTemplatedPolicies{
				&ACLTemplatedPolicy{
					TemplateName: api.ACLTemplatedPolicyServiceName,
					TemplateVariables: &ACLTemplatedPolicyVariables{
						Name: "api",
					},
				},
				&ACLTemplatedPolicy{
					TemplateName: api.ACLTemplatedPolicyServiceName,
					TemplateVariables: &ACLTemplatedPolicyVariables{
						Name: "api",
					},
				},
			},
			expectedCount: 1,
		},
		"separate-templates-with-matching-variables": {
			templatedPolicies: ACLTemplatedPolicies{
				&ACLTemplatedPolicy{
					TemplateName: api.ACLTemplatedPolicyNodeName,
					TemplateVariables: &ACLTemplatedPolicyVariables{
						Name: "api",
					},
				},
				&ACLTemplatedPolicy{
					TemplateName: api.ACLTemplatedPolicyServiceName,
					TemplateVariables: &ACLTemplatedPolicyVariables{
						Name: "api",
					},
				},
			},
			expectedCount: 2,
		},
		"separate-templates-with-multiple-matching-variables": {
			templatedPolicies: ACLTemplatedPolicies{
				&ACLTemplatedPolicy{
					TemplateName: api.ACLTemplatedPolicyServiceName,
					TemplateVariables: &ACLTemplatedPolicyVariables{
						Name: "api",
					},
				},
				&ACLTemplatedPolicy{
					TemplateName: api.ACLTemplatedPolicyNodeName,
					TemplateVariables: &ACLTemplatedPolicyVariables{
						Name: "api",
					},
				},
				&ACLTemplatedPolicy{
					TemplateName: api.ACLTemplatedPolicyNodeName,
					TemplateVariables: &ACLTemplatedPolicyVariables{
						Name: "web",
					},
				},
				&ACLTemplatedPolicy{
					TemplateName: api.ACLTemplatedPolicyServiceName,
					TemplateVariables: &ACLTemplatedPolicyVariables{
						Name: "api",
					},
				},
				&ACLTemplatedPolicy{
					TemplateName: api.ACLTemplatedPolicyDNSName,
				},
				&ACLTemplatedPolicy{
					TemplateName: api.ACLTemplatedPolicyServiceName,
					TemplateVariables: &ACLTemplatedPolicyVariables{
						Name: "web",
					},
				},
				&ACLTemplatedPolicy{
					TemplateName: api.ACLTemplatedPolicyDNSName,
				},
			},
			expectedCount: 5,
		},
		"same-template-with-different-exact-name": {
			templatedPolicies: ACLTemplatedPolicies{
				&ACLTemplatedPolicy{
					TemplateName: api.ACLTemplatedPolicyAllowServiceName,
					TemplateVariables: &ACLTemplatedPolicyVariables{
						Name:      "api-",
						ExactName: "api",
					},
				},
				&ACLTemplatedPolicy{
					TemplateName: api.ACLTemplatedPolicyAllowServiceName,
					TemplateVariables: &ACLTemplatedPolicyVariables{
						Name:      "api-",
						ExactName: "web",
					},
				},
			},
			expectedCount: 2,
		},
	}

	for name, tcase := range tcases {
		t.Run(name, func(t *testing.T) {
			policies := tcase.templatedPolicies.Deduplicate()

			require.Equal(t, tcase.expectedCount, len(policies))
		})
	}
}

func TestValidateTemplatedPolicyAllowService(t *testing.T) {
	tcases := map[string]struct {
		templatedPolicy *ACLTemplatedPolicy
		expectedErr     string
	}{
		"legacy-prefix-only": {
			templatedPolicy: &ACLTemplatedPolicy{
				TemplateName: api.ACLTemplatedPolicyAllowServiceName,
				TemplateVariables: &ACLTemplatedPolicyVariables{
					Name: "api-",
				},
			},
		},
		"exact-name-valid": {
			templatedPolicy: &ACLTemplatedPolicy{
				TemplateName: api.ACLTemplatedPolicyAllowServiceName,
				TemplateVariables: &ACLTemplatedPolicyVariables{
					Name:      "api-",
					ExactName: "api",
				},
			},
		},
		"name-does-not-match-exact-name": {
			templatedPolicy: &ACLTemplatedPolicy{
				TemplateName: api.ACLTemplatedPolicyAllowServiceName,
				TemplateVariables: &ACLTemplatedPolicyVariables{
					Name:      "api-",
					ExactName: "web",
				},
			},
			expectedErr: `requires name "api-" to match exact_name "web" plus a trailing '-'`,
		},
	}

	for name, tcase := range tcases {
		t.Run(name, func(t *testing.T) {
			err := tcase.templatedPolicy.ValidateTemplatedPolicy(ACLTemplatedPolicyServiceSchema)
			if tcase.expectedErr == "" {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			require.Contains(t, err.Error(), tcase.expectedErr)
		})
	}
}
