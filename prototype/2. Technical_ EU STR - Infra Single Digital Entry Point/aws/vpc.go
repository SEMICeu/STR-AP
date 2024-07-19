package main

import (
	"github.com/pulumi/pulumi-awsx/sdk/v2/go/awsx/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func createVPC(ctx *pulumi.Context) (*ec2.Vpc, error) {

	vpcCidr := config.Get(ctx, "project:vpcCidr")

	vpc, err := ec2.NewVpc(ctx, getProjectName(ctx), &ec2.VpcArgs{
		CidrBlock: pulumi.StringRef(vpcCidr),
		NatGateways: &ec2.NatGatewayConfigurationArgs{
			Strategy: ec2.NatGatewayStrategySingle,
		},
	})
	if err != nil {
		return nil, err
	}

	ctx.Export("vpcId", vpc.VpcId)
	ctx.Export("vpcCidr", pulumi.String(vpcCidr))
	ctx.Export("vpcPrivateSubnetIds", vpc.PrivateSubnetIds)
	ctx.Export("vpcPublicSubnetIds", vpc.PublicSubnetIds)
	return vpc, nil
}
