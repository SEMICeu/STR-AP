package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/route53"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func setDnsAlias(ctx *pulumi.Context) error {

	// Get config values
	hostname := config.Get(ctx, "project:hostname")
	domain := config.Get(ctx, "project:domain")
	domainHostedZoneId := config.Get(ctx, "project:domainHostedZoneId")
	nlbHostedZoneId := config.Get(ctx, "project:nlbHostedZoneId")
	nlbDnsName := config.Get(ctx, "project:nlbDnsName")

	fqdn := hostname + "." + domain

	_, err := route53.NewRecord(ctx, hostname, &route53.RecordArgs{
		ZoneId: pulumi.String(domainHostedZoneId),
		Type:   pulumi.String(route53.RecordTypeA),
		Name:   pulumi.String(fqdn),
		Aliases: route53.RecordAliasArray{
			&route53.RecordAliasArgs{
				Name:                 pulumi.String(nlbDnsName),
				ZoneId:               pulumi.String(nlbHostedZoneId),
				EvaluateTargetHealth: pulumi.Bool(true),
			},
		},
	})
	if err != nil {
		return err
	}

	// Export the application route53 hostedzoneId and fqdn
	ctx.Export("domainHostedZoneId", pulumi.String(domainHostedZoneId))
	ctx.Export("fqdn", pulumi.String(fqdn))

	return nil
}
