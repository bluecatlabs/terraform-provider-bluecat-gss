// Copyright 2022 BlueCat Networks. All rights reserved

package utils

import (
	"encoding/json"
	"fmt"
	"terraform-provider-bluecat/bluecat/entities"
	"terraform-provider-bluecat/bluecat/models"
)

// ObjectManager The BlueCat object manager
type ObjectManager struct {
	Connector BCConnector
}

// GSS Application

// CreateGSSApplication Create the GSS Application
func (objMgr *ObjectManager) CreateGSSApplication(configuration string, view string, zone string, absoluteName string, fallback []interface{}, ttl int, properties string, healthCheck []interface{}, searchOrder []interface{}) (*entities.GSSApplication, error) {
	gssApplication := models.NewGSSApplication(entities.GSSApplication{
		Configuration: configuration,
		View:          view,
		Zone:          zone,
		Fallback:      fallback,
		AbsoluteName:  absoluteName,
		TTL:           ttl,
		Properties:    properties,
		HealthCheck:   healthCheck,
		SearchOrder:   searchOrder,
	})

	_, err := objMgr.Connector.CreateObject(gssApplication)
	return gssApplication, err
}

// GetGSSApplication the Host record
func (objMgr *ObjectManager) GetGSSApplication(configuration string, view string, absoluteName string) (*entities.GSSApplication, error) {

	gssApplication := models.GSSApplication(entities.GSSApplication{
		Configuration: configuration,
		View:          view,
		AbsoluteName:  absoluteName,
	})

	err := objMgr.Connector.GetObject(gssApplication, &gssApplication)
	return gssApplication, err
}

// UpdateGSSApplication Update the GSS Application
func (objMgr *ObjectManager) UpdateGSSApplication(configuration string, view string, zone string, absoluteName string, fallback []interface{}, ttl int, properties string, healthCheck []interface{}, searchOrder []interface{}) (*entities.GSSApplication, error) {

	gssApplication := models.GSSApplication(entities.GSSApplication{
		Configuration: configuration,
		View:          view,
		Zone:          zone,
		Fallback:      fallback,
		AbsoluteName:  absoluteName,
		TTL:           ttl,
		Properties:    properties,
		HealthCheck:   healthCheck,
		SearchOrder:   searchOrder,
	})

	err := objMgr.Connector.UpdateObject(gssApplication, &gssApplication)
	return gssApplication, err
}

// DeleteGSSApplication Delete the GSS Application
func (objMgr *ObjectManager) DeleteGSSApplication(configuration string, view string, absoluteName string) (string, error) {

	gssApplication := models.GSSApplication(entities.GSSApplication{
		Configuration: configuration,
		View:          view,
		AbsoluteName:  absoluteName,
	})

	return objMgr.Connector.DeleteObject(gssApplication)
}

// GSS Answers

// CreateGSSAnswer Create the GSS Answers
func (objMgr *ObjectManager) CreateGSSAnswer(applicationId int, addresses []interface{}, region string, name string, answerType string) (*entities.GSSAnswer, error) {

	gssAnswer := models.NewGSSAnswer(entities.GSSAnswer{
		ApplicationId: applicationId,
		Addresses:     addresses,
		Region:        region,
		Name:          name,
		Type:          answerType,
	})

	response, err := objMgr.Connector.CreateObject(gssAnswer)
	newGssAnswer := models.NewGSSAnswer(entities.GSSAnswer{})

	if err != nil {
		msg := fmt.Sprintf("Failure to convert respose created answer: %s", err)
		log.Debug(msg)
		return gssAnswer, err
	}
	json.Unmarshal([]byte(response), &newGssAnswer)

	return newGssAnswer, err
}

// GetGSSAnswer the GSS Answers
func (objMgr *ObjectManager) GetGSSAnswer(applicationId int, answerId int) (*entities.GSSAnswer, error) {

	gssAnswer := models.GSSAnswer(entities.GSSAnswer{
		ApplicationId: applicationId,
		AnswerId:      answerId,
	})

	newGssAnswer := models.GSSAnswer(entities.GSSAnswer{})
	err := objMgr.Connector.GetObject(gssAnswer, &newGssAnswer)
	return newGssAnswer, err
}

// UpdateGSSAnswer Update the GSS Answers
func (objMgr *ObjectManager) UpdateGSSAnswer(applicationId int, answerId int, addresses []interface{}, region string, name string, answerType string) (*entities.GSSAnswer, error) {

	gssAnswer := models.GSSAnswer(entities.GSSAnswer{
		ApplicationId: applicationId,
		Addresses:     addresses,
		Region:        region,
		Name:          name,
		AnswerId:      answerId,
		Type:          answerType,
	})

	err := objMgr.Connector.UpdateObject(gssAnswer, &gssAnswer)
	return gssAnswer, err
}

// DeleteGSSAnswer Delete the GSS Answers
func (objMgr *ObjectManager) DeleteGSSAnswer(applicationId int, answerId int) (string, error) {

	gssAnswer := models.GSSAnswer(entities.GSSAnswer{
		ApplicationId: applicationId,
		AnswerId:      answerId,
	})

	return objMgr.Connector.DeleteObject(gssAnswer)
}

// GSS Search Order

// CreateGSSSearchOrder Create the GSS Search Order
func (objMgr *ObjectManager) CreateGSSSearchOrder(name string, nodes []interface{}, links []interface{}) (*entities.GSSSearchOrder, error) {
	gssSearchOrder := models.NewSearchOrder(entities.GSSSearchOrder{
		Nodes: nodes,
		Links: links,
		Name:  name,
	})

	_, err := objMgr.Connector.CreateObject(gssSearchOrder)
	return gssSearchOrder, err
}

// GetGSSSearchOrder Get the GSS Search Order
func (objMgr *ObjectManager) GetGSSSearchOrder(searchOrderID int, name string, links []interface{}) (*entities.GSSSearchOrder, error) {

	gssSearchOrder := models.GSSSearchOrder(entities.GSSSearchOrder{
		SearchOrderId: searchOrderID,
		Name:          name,
		Links:         links,
	})

	newGssSearchOrder := models.GSSSearchOrder(entities.GSSSearchOrder{})
	err := objMgr.Connector.GetObject(gssSearchOrder, &newGssSearchOrder)
	return newGssSearchOrder, err
}

// UpdateGSSSearchOrder Update the GSS Search Order
func (objMgr *ObjectManager) UpdateGSSSearchOrder(name string, searchOrderId int, nodes []interface{}, links []interface{}) (*entities.GSSSearchOrder, error) {

	gssSearchOrder := models.GSSSearchOrder(entities.GSSSearchOrder{
		Name:          name,
		Nodes:         nodes,
		Links:         links,
		SearchOrderId: searchOrderId,
	})

	err := objMgr.Connector.UpdateObject(gssSearchOrder, &gssSearchOrder)
	return gssSearchOrder, err
}

// DeleteGSSSearchOrder Delete the GSS Search Order
func (objMgr *ObjectManager) DeleteGSSSearchOrder(name string) (string, error) {

	gssSearchOrder := models.GSSSearchOrder(entities.GSSSearchOrder{
		Name: name,
	})

	return objMgr.Connector.DeleteObject(gssSearchOrder)
}
