// Copyright 2022 BlueCat Networks. All rights reserved

package models

import (
	"terraform-provider-bluecat-gss/bluecat/entities"
)

// Application
// Application Initialize the new Application to be added
func NewGSSApplication(application entities.GSSApplication) *entities.GSSApplication {
	res := application
	res.SetObjectType("bluecat_gss_application")
	res.SetSubPath("")
	return &res
}

// Application Initialize the Application to be loaded, updated or deleted
func GSSApplication(application entities.GSSApplication) *entities.GSSApplication {
	res := application
	res.SetObjectType("")
	res.SetSubPath("/bluecat_gss_application")
	return &res
}

// Answer
// Answer Initialize the new Answer to be added
func NewGSSAnswer(answer entities.GSSAnswer) *entities.GSSAnswer {
	res := answer
	res.SetObjectType("bluecat_gss_answer")
	res.SetSubPath("")
	return &res
}

// Application Initialize the Application to be loaded, updated or deleted
func GSSAnswer(application entities.GSSAnswer) *entities.GSSAnswer {
	res := application
	res.SetObjectType("")
	res.SetSubPath("/bluecat_gss_answer")
	return &res
}

// SearchOrder
// SearchOrder Initialize the new SearchOrder to be added
func NewSearchOrder(answer entities.GSSSearchOrder) *entities.GSSSearchOrder {
	res := answer
	res.SetObjectType("bluecat_gss_search_order")
	res.SetSubPath("")
	return &res
}

// SearchOrder Initialize the SearchOrder to be loaded, updated or deleted
func GSSSearchOrder(application entities.GSSSearchOrder) *entities.GSSSearchOrder {
	res := application
	res.SetObjectType("")
	res.SetSubPath("/bluecat_gss_search_order")
	return &res
}
