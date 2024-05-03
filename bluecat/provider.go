// Copyright 2022 BlueCat Networks. All rights reserved

package bluecat

import (
	log "github.com/sirupsen/logrus"
	"terraform-provider-bluecat-gss/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider BlueCat provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "BlueCat Gateway IP address.",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User to authenticate with BlueCat Gateway server.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Password to authenticate with BlueCat Gateway server.",
			},
			"api_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "API Version of REST_API workflow server",
			},
			"port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Port number used for connection for BlueCat Gateway Server.",
			},
			"transport": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Transport type (HTTP or HTTPS).",
			},
			"encrypt_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Default is false, to indicate if the password is encrypted",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"bluecatgss_application":  ResourceGSSApplication(),
			"bluecatgss_answer":       ResourceGSSAnswer(),
			"bluecatgss_search_order": ResourceGSSSearchOrder(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			// TODO: Support read data from GSS in the future (don't need right now)
			// 			"bluecatgss_application":  DataSourceGSSApplication(),
			// 			"bluecatgss_answer":       DataSourceGSSAnswer(),
			//			"bluecatgss_search_order": DataSourceGSSSearchOrder(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	hostConfig := utils.HostConfig{
		Host:            d.Get("server").(string),
		Port:            d.Get("port").(string),
		Transport:       d.Get("transport").(string),
		Username:        d.Get("username").(string),
		Password:        d.Get("password").(string),
		Version:         d.Get("api_version").(string),
		EncryptPassword: d.Get("encrypt_password").(bool),
	}

	requestBuilder := &utils.APIRequestBuilder{}
	requester := &utils.APIHttpRequester{}

	conn, err := utils.NewConnector(hostConfig, requestBuilder, requester)
	if err != nil {
		log.Debugf("Failed to initialize the provider: %s", err)
		return nil, err
	}
	return conn, err
}
