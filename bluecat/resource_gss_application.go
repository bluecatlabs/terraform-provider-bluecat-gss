// Copyright 2022 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"terraform-provider-bluecat/bluecat/utils"
)

// ResourceGSSApplication The GSS Application
func ResourceGSSApplication() *schema.Resource {

	return &schema.Resource{
		Create: createGSSApplication,
		Read:   getGSSApplication,
		Update: updateGSSApplication,
		Delete: deleteGSSApplication,

		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The Configuration name",
				ValidateFunc: StringIsNotWhiteSpace,
			},
			"view": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The View name",
				ValidateFunc: StringIsNotWhiteSpace,
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Zone in which you want to update a GSS Application. If not provided, the absolute name must be FQDN ones",
			},
			"absolute_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the GSS Application. Must be FQDN if the Zone is not provided",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					zone := d.Get("zone").(string)
					return checkDiffName(old, new, zone)
				},
				ValidateFunc: StringIsNotWhiteSpace,
			},
			"fallback": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "The list of IP addresses that will be linked to the GSS Application",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The TTL value",
				Default:     -1,
			},
			"properties": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Host record's properties. Example: attribute=value|",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return checkDiffProperties(old, new)
				},
			},
			"health_check": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Health Check configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Required attribute
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "NO_HEALTH_CHECK",
							Description:  "Health Check support these Type: TCP, HTTP_HEAD, CUSTOMIZE",
							ValidateFunc: StringInSlice([]string{"TCP", "HTTP_HEAD", "CUSTOMIZE", "NO_HEALTH_CHECK"}, false),
						},
						"check_every": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     30,
							Description: "Define the interval health check",
						},
						// For HTTP HEAD Type
						"secure_connection": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Only for HTTP HEAD Type",
						},
						"appended_url_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "/",
							Description: "Only for HTTP HEAD Type",
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								// Suppress diff if the attribute is absent in the new configuration
								if new == "" {
									return true
								}
								return false
							},
						},
						"optional_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Only for HTTP HEAD Type",
						},
						"header_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Only for HTTP HEAD Type",
						},

						// For TCP Type
						"port": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "22",
							Description: "Only for TCP Type",
						},

						// For customize
						"custom_data": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Specifies the Custom data",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: StringIsNotWhiteSpace,
							},
						},
					},
				},
			},

			"search_order": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Defined the Search Order Configuration to each Client Region",
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

// StringInSlice returns a SchemaValidateFunc which tests if the provided value
// is of type string and matches the value of an element in the valid slice
// will test with in lower case if ignoreCase is true
func StringInSlice(valid []string, ignoreCase bool) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
			return warnings, errors
		}

		for _, str := range valid {
			if v == str || (ignoreCase && strings.EqualFold(v, str)) {
				return warnings, errors
			}
		}

		errors = append(errors, fmt.Errorf("expected %s to be one of %v, got %s", k, valid, v))
		return warnings, errors
	}
}

// StringIsNotWhiteSpace is a ValidateFunc that ensures a string is not empty or consisting entirely of whitespace characters
func StringIsNotWhiteSpace(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if strings.TrimSpace(v) == "" {
		return nil, []error{fmt.Errorf("expected %q to not be an empty string or whitespace", k)}
	}

	return nil, nil
}

func getFQDN(rrName, zone string) string {
	if !strings.HasSuffix(rrName, ".") && len(zone) > 0 && !strings.HasSuffix(rrName, zone) {
		return fmt.Sprintf("%s.%s", rrName, zone)
	}
	return rrName
}

func getZoneFromRRName(rrName string) (zoneFQDN string) {
	zoneFQDN = ""
	index := strings.Index(rrName, ".")
	if index > 0 {
		zoneFQDN = rrName[index+1:]
	}
	return
}

func checkDiffProperties(old string, new string) bool {
	newProperties := strings.Split(new, "|")
	for i := 0; i < len(newProperties); i++ {
		if newProperties[i] != "" && !strings.Contains(fmt.Sprintf("|%s|", old), fmt.Sprintf("|%s|", newProperties[i])) {
			return false
		}
	}
	return true
}

func checkDiffName(old string, new string, zone string) bool {
	if old == getFQDN(new, zone) {
		return true
	}
	return false
}

// createGSSApplication Create the new GSS Application
func createGSSApplication(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to createA GSS Application %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)
	absoluteName := d.Get("absolute_name").(string)
	fallback := d.Get("fallback").([]interface{})
	ttl := d.Get("ttl").(int)
	properties := d.Get("properties").(string)
	healthCheck := d.Get("health_check").([]interface{})
	searchOrder := d.Get("search_order").([]interface{})

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	fqdnName := absoluteName

	if len(zone) > 0 {
		fqdnName = getFQDN(absoluteName, zone)
	} else {
		zone = getZoneFromRRName(fqdnName)
	}

	applicationObj, err := objMgr.CreateGSSApplication(configuration, view, zone, fqdnName, fallback, ttl, properties, healthCheck, searchOrder)
	if err != nil {
		msg := fmt.Sprintf("Error creating GSS Application %s: %s", fqdnName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.Set("absolute_name", fqdnName)
	d.SetId(strconv.Itoa(applicationObj.ApplicationId))
	log.Debugf("Completed to create Application %s", d.Get("absolute_name"))
	return getGSSApplication(d, m)
}

// getGSSApplication Get the GSS Application
func getGSSApplication(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to get GSS Application: %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	absoluteName := d.Get("absolute_name").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	gssApplication, err := objMgr.GetGSSApplication(configuration, view, absoluteName)
	if err != nil {
		msg := fmt.Sprintf("Getting GSS Application %s failed: %s", absoluteName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.SetId(strconv.Itoa(gssApplication.ApplicationId))
	d.Set("absolute_name", gssApplication.AbsoluteName)
	d.Set("properties", gssApplication.Properties)
	log.Debugf("Completed reading GSS Application %s", d.Get("absolute_name"))
	return nil
}

// updateGSSApplication Update the existing GSS Application
func updateGSSApplication(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to update GSS Application %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)
	absoluteName := d.Get("absolute_name").(string)
	fallback := d.Get("fallback").([]interface{})
	ttl := d.Get("ttl").(int)
	properties := d.Get("properties").(string)
	healthCheck := d.Get("health_check").([]interface{})
	searchOrder := d.Get("search_order").([]interface{})

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	fqdnName := absoluteName

	if len(zone) > 0 {
		fqdnName = getFQDN(absoluteName, zone)
	} else {
		zone = getZoneFromRRName(fqdnName)
	}
	gssApplication, err := objMgr.UpdateGSSApplication(configuration, view, zone, fqdnName, fallback, ttl, properties, healthCheck, searchOrder)
	if err != nil {
		msg := fmt.Sprintf("Error updating GSS Application %s: %s", fqdnName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.SetId(strconv.Itoa(gssApplication.ApplicationId))
	log.Debugf("Completed to update GSS Application %s", d.Get("absolute_name"))
	return getGSSApplication(d, m)
}

// deleteGSSApplication Delete the GSS Application
func deleteGSSApplication(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to delete GSS Application %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	absoluteName := d.Get("absolute_name").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	// Check the host exist or not
	_, err := objMgr.GetGSSApplication(configuration, view, absoluteName)
	if err != nil {
		log.Debugf("The Application %s not found", absoluteName)
	} else {
		_, err := objMgr.DeleteGSSApplication(configuration, view, absoluteName)
		if err != nil {
			msg := fmt.Sprintf("Delete GSS Application %s failed: %s", absoluteName, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
	}
	d.SetId("")
	log.Debugf("Completed to delete GSS Application %s", d.Get("absolute_name"))
	return nil
}
