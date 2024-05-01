// Copyright 2022 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
	"strconv"
	"terraform-provider-bluecat-gss/bluecat/utils"
)

// ResourceGSSAnswer The GSS Answer
func ResourceGSSAnswer() *schema.Resource {
	return &schema.Resource{
		Create: createGSSAnswer,
		Read:   getGSSAnswer,
		Update: updateGSSAnswer,
		Delete: deleteGSSAnswer,

		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The Application ID",
			},
			"addresses": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The list of IP addresses that will be linked to the GSS Answer",
				Elem: &schema.Schema{
					Type: schema.TypeString},
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region name",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "GSS Answer name",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Type Answer",
				ValidateFunc: StringInSlice([]string{"ip_address", "fqdn"}, false),
			},
		},
	}
}

// createGSSAnswer Create the new GSS Answer
func createGSSAnswer(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to create GSS Answer with Application ID %v", d.Get("application_id"))
	application_id := d.Get("application_id").(int)
	addresses := d.Get("addresses").([]interface{})
	region := d.Get("region").(string)
	name := d.Get("name").(string)
	answer_type := d.Get("type").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	gssAnswer, err := objMgr.CreateGSSAnswer(application_id, addresses, region, name, answer_type)
	if err != nil {
		msg := fmt.Sprintf("Error creating GSS Answer %v: %s", gssAnswer.AnswerId, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.SetId(strconv.Itoa(gssAnswer.AnswerId))
	d.Set("application_id", application_id)
	log.Debugf("Completed to create GSS Answers %v", gssAnswer.AnswerId)
	return getGSSAnswer(d, m)
}

// getGSSAnswer Get the GSS Answer
func getGSSAnswer(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to get GSS Answer with Application ID: %v", d.Get("application_id"))
	application_id := d.Get("application_id").(int)
	answerID, _ := strconv.Atoi(d.Id())

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	gssAnswer, err := objMgr.GetGSSAnswer(application_id, answerID)
	if err != nil {
		msg := fmt.Sprintf("Getting GSS Answer %v failed: %s", gssAnswer.AnswerId, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}

	d.SetId(strconv.Itoa(gssAnswer.AnswerId))
	d.Set("application_id", application_id)
	d.Set("addresses", gssAnswer.Addresses)
	d.Set("region", gssAnswer.Region)
	d.Set("name", gssAnswer.Name)
	d.Set("type", gssAnswer.Type)
	log.Debugf("Completed reading GSS Answer %v", gssAnswer.AnswerId)
	return nil
}

// updateGSSAnswer Update the existing GSS Answer
func updateGSSAnswer(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to update GSS Answer with Application Id %v", d.Get("application_id"))
	application_id := d.Get("application_id").(int)
	addresses := d.Get("addresses").([]interface{})
	region := d.Get("region").(string)
	name := d.Get("name").(string)
	answer_type := d.Get("type").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	answerID, _ := strconv.Atoi(d.Id())
	gssAnswer, err := objMgr.UpdateGSSAnswer(application_id, answerID, addresses, region, name, answer_type)
	if err != nil {
		msg := fmt.Sprintf("Error updating GSS Answer %v: %s", application_id, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.SetId(strconv.Itoa(gssAnswer.AnswerId))
	d.Set("application_id", application_id)
	log.Debugf("Completed to update GSS Answer %s", d.Get("application_id"))
	return getGSSAnswer(d, m)
}

// deleteGSSAnswer Delete the GSS Answer
func deleteGSSAnswer(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to delete GSS Answer %s", d.Get("application_id"))
	application_id := d.Get("application_id").(int)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	answerID, _ := strconv.Atoi(d.Id())
	_, err := objMgr.DeleteGSSAnswer(application_id, answerID)
	if err != nil {
		msg := fmt.Sprintf("Delete GSS Answer %v failed: %s", application_id, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.SetId("")
	log.Debugf("Completed to delete GSS Answer: %s", answerID)
	return nil
}
