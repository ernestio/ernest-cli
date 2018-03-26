/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package manager

import (
	h "github.com/ernestio/ernest-cli/helper"
	eclient "github.com/ernestio/ernest-go-sdk/client"
	emodels "github.com/ernestio/ernest-go-sdk/models"
)

// Policy : ernest-go-sdk Policy wrapper
type Policy struct {
	cli *eclient.Client
}

// Get : Gets a policy by name
func (c *Policy) Get(id string) *emodels.Policy {
	policy, err := c.cli.Policies.Get(id)
	if err != nil {
		h.PrintError(err.Error())
	}
	return policy
}

// Update : Updates a policy
func (c *Policy) Update(policy *emodels.Policy) {
	if err := c.cli.Policies.Update(policy); err != nil {
		h.PrintError(err.Error())
	}
}

// Create : Creates a new policy
func (c *Policy) Create(policy *emodels.Policy) {
	if err := c.cli.Policies.Create(policy); err != nil {
		h.PrintError(err.Error())
	}
}

// List : Lists all policies on the system
func (c *Policy) List() []*emodels.Policy {
	policies, err := c.cli.Policies.List()
	if err != nil {
		h.PrintError(err.Error())
	}
	return policies
}

// Delete : Deletes a policy and all its relations
func (c *Policy) Delete(policy string) {
	if err := c.cli.Policies.Delete(policy); err != nil {
		h.PrintError(err.Error())
	}
}

// GetDocument : Gets a policy document by revision
func (c *Policy) GetDocument(policy, revision string) *emodels.PolicyDocument {
	document, err := c.cli.Policies.GetDocument(policy, revision)
	if err != nil {
		h.PrintError(err.Error())
	}

	return document
}

// ListDocuments : Lists all policy documents by policy name
func (c *Policy) ListDocuments(policy string) []*emodels.PolicyDocument {
	documents, err := c.cli.Policies.ListDocuments(policy)
	if err != nil {
		h.PrintError(err.Error())
	}

	return documents
}

// CreateDocument : Creates a policy document and all its relations
func (c *Policy) CreateDocument(policy, document string) {
	if _, err := c.cli.Policies.CreateDocument(policy, document); err != nil {
		h.PrintError(err.Error())
	}
}
