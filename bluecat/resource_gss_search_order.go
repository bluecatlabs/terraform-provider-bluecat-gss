// Copyright 2022 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
	"strconv"
	"terraform-provider-bluecat/bluecat/utils"
)

// ResourceGSSSearchOrder The GSS Search Order
func ResourceGSSSearchOrder() *schema.Resource {

	return &schema.Resource{
		Create: createGSSSearchOrder,
		Read:   getGSSSearchOrder,
		Update: updateGSSSearchOrder,
		Delete: deleteGSSSearchOrder,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The Search Order name",
				ValidateFunc: StringIsNotWhiteSpace,
			},
			"nodes": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The list of nodes will be config in Search Order",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of a node",
						},
					},
				},
			},

			"links": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Defined the Search Order Configuration to each Client Region",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Source node",
							ValidateFunc: StringIsNotWhiteSpace,
						},
						"target": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Destination node",
						},
						"cost": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Distance between source and destination node",
						},
						"enable_link": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Enable link",
						},
					},
				},
			},
		},
	}
}

// createGSSSearchOrder Create the new GSS Search Order
func createGSSSearchOrder(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to create GSS Search Order %s", d.Get("name"))
	name := d.Get("name").(string)
	nodes := d.Get("nodes").([]interface{})
	links := d.Get("links").([]interface{})

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	searchOrderObj, err := objMgr.CreateGSSSearchOrder(name, nodes, links)
	if err != nil {
		msg := fmt.Sprintf("Error creating GSS Search Order %s: %s", name, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.SetId(strconv.Itoa(searchOrderObj.SearchOrderId))
	log.Debugf("Completed to create Search Order %s", d.Get("name"))
	return getGSSSearchOrder(d, m)
}

// getGSSSearchOrder Get the GSS Search Order
func getGSSSearchOrder(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to get GSS Search Order %s", d.Get("name"))
	searchOrderID, _ := strconv.Atoi(d.Id())
	name := d.Get("name").(string)
	links := d.Get("links").([]interface{})

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	gssSearchOrder, err := objMgr.GetGSSSearchOrder(searchOrderID, name, links)
	if err != nil {
		msg := fmt.Sprintf("Getting GSS Search Order %s failed: %s", gssSearchOrder.Name, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}

	d.SetId(strconv.Itoa(gssSearchOrder.SearchOrderId))
	d.Set("nodes", gssSearchOrder.Nodes)
	d.Set("links", gssSearchOrder.Links)
	d.Set("name", gssSearchOrder.Name)
	log.Debugf("Completed reading GSS Search Order %s", gssSearchOrder.Name)
	return nil
}

// updateGSSSearchOrder Update the existing GSS Search Order
func updateGSSSearchOrder(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to update GSS Search Order %s", d.Get("name"))
	nodes := d.Get("nodes").([]interface{})
	links := d.Get("links").([]interface{})
	name := d.Get("name").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	searchOrderID, _ := strconv.Atoi(d.Id())
	gssSearchOrder, err := objMgr.UpdateGSSSearchOrder(name, searchOrderID, nodes, links)
	if err != nil {
		msg := fmt.Sprintf("Error updating GSS Search Order %s: %s", name, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.SetId(strconv.Itoa(gssSearchOrder.SearchOrderId))
	log.Debugf("Completed to update GSS Search Order %s", d.Get("name"))
	return getGSSSearchOrder(d, m)
}

// deleteGSSSearchOrder Delete the GSS Search Order
func deleteGSSSearchOrder(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to delete GSS Search Order %s", d.Get("name"))
	name := d.Get("name").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err := objMgr.DeleteGSSSearchOrder(name)
	if err != nil {
		msg := fmt.Sprintf("Delete GSS Search Order %s failed: %s", name, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.SetId("")
	log.Debugf("Completed to delete GSS Search Order: %s", name)
	return nil
}
