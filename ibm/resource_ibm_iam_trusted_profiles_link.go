// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func resourceIBMIamTrustedProfilesLink() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIamTrustedProfilesLinkCreate,
		ReadContext:   resourceIBMIamTrustedProfilesLinkRead,
		DeleteContext: resourceIBMIamTrustedProfilesLinkDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"profile_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the trusted profile.",
			},
			"cr_type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The compute resource type. Valid values are VSI, IKS_SA, ROKS_SA.",
			},
			"link": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				ForceNew:    true,
				Description: "Link details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The CRN of the compute resource.",
						},
						"namespace": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The compute resource namespace, only required if cr_type is IKS_SA or ROKS_SA.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the compute resource, only required if cr_type is IKS_SA or ROKS_SA.",
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional name of the Link.",
			},
			"entity_tag": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "version of the claim rule.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If set contains a date time string of the creation date in ISO format.",
			},
			"modified_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If set contains a date time string of the last modification date in ISO format.",
			},
		},
	}
}

func resourceIBMIamTrustedProfilesLinkCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	createLinkOptions := &iamidentityv1.CreateLinkOptions{}

	createLinkOptions.SetProfileID(d.Get("profile_id").(string))
	createLinkOptions.SetCrType(d.Get("cr_type").(string))
	link := resourceIBMIamTrustedProfilesLinkMapToCreateProfileLinkRequestLink(d.Get("link.0").(map[string]interface{}))
	createLinkOptions.SetLink(&link)
	if _, ok := d.GetOk("name"); ok {
		createLinkOptions.SetName(d.Get("name").(string))
	}

	profileLink, response, err := iamIdentityClient.CreateLink(createLinkOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateLink failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateLink failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createLinkOptions.ProfileID, *profileLink.ID))

	return resourceIBMIamTrustedProfilesLinkRead(context, d, meta)
}

func resourceIBMIamTrustedProfilesLinkMapToCreateProfileLinkRequestLink(createProfileLinkRequestLinkMap map[string]interface{}) iamidentityv1.CreateProfileLinkRequestLink {
	createProfileLinkRequestLink := iamidentityv1.CreateProfileLinkRequestLink{}

	createProfileLinkRequestLink.CRN = core.StringPtr(createProfileLinkRequestLinkMap["crn"].(string))
	createProfileLinkRequestLink.Namespace = core.StringPtr(createProfileLinkRequestLinkMap["namespace"].(string))
	if createProfileLinkRequestLinkMap["name"] != nil {
		createProfileLinkRequestLink.Name = core.StringPtr(createProfileLinkRequestLinkMap["name"].(string))
	}

	return createProfileLinkRequestLink
}

func resourceIBMIamTrustedProfilesLinkRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getLinkOptions := &iamidentityv1.GetLinkOptions{}

	getLinkOptions.SetProfileID(d.Get("profile_id").(string))
	getLinkOptions.SetLinkID(d.Get("link_id").(string))

	profileLink, response, err := iamIdentityClient.GetLink(getLinkOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetLink failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetLink failed %s\n%s", err, response))
	}

	if err = d.Set("profile_id", getLinkOptions.ProfileID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting profile_id: %s", err))
	}
	if err = d.Set("cr_type", profileLink.CrType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting cr_type: %s", err))
	}
	linkMap := resourceIBMIamTrustedProfilesLinkCreateProfileLinkRequestLinkToMap(*profileLink.Link)
	if err = d.Set("link", []map[string]interface{}{linkMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting link: %s", err))
	}
	if err = d.Set("name", profileLink.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("entity_tag", profileLink.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting entity_tag: %s", err))
	}
	if err = d.Set("created_at", dateTimeToString(profileLink.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("modified_at", dateTimeToString(profileLink.ModifiedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting modified_at: %s", err))
	}

	return nil
}

func resourceIBMIamTrustedProfilesLinkCreateProfileLinkRequestLinkToMap(createProfileLinkRequestLink iamidentityv1.ProfileLinkLink) map[string]interface{} {
	createProfileLinkRequestLinkMap := map[string]interface{}{}

	createProfileLinkRequestLinkMap["crn"] = createProfileLinkRequestLink.CRN
	createProfileLinkRequestLinkMap["namespace"] = createProfileLinkRequestLink.Namespace
	if createProfileLinkRequestLink.Name != nil {
		createProfileLinkRequestLinkMap["name"] = createProfileLinkRequestLink.Name
	}

	return createProfileLinkRequestLinkMap
}

func resourceIBMIamTrustedProfilesLinkDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteLinkOptions := &iamidentityv1.DeleteLinkOptions{}

	deleteLinkOptions.SetProfileID(d.Get("profile_id").(string))
	deleteLinkOptions.SetLinkID(d.Get("link_id").(string))

	response, err := iamIdentityClient.DeleteLink(deleteLinkOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteLink failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteLink failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
