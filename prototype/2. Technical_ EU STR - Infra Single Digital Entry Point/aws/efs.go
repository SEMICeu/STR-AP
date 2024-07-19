package main

import (
	awsec2 "github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/efs"
	awsvpc "github.com/pulumi/pulumi-aws/sdk/v6/go/aws/vpc"
	"github.com/pulumi/pulumi-awsx/sdk/v2/go/awsx/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createFileShare(ctx *pulumi.Context, vpc *ec2.Vpc) (*efs.FileSystem, error) {
	projectName := getProjectName(ctx)

	fileShare, _ := efs.NewFileSystem(ctx, projectName, &efs.FileSystemArgs{
		CreationToken: pulumi.String(projectName),
		Tags: pulumi.StringMap{
			"Name": pulumi.String(projectName),
		},
	})

	// Create SG for fileShare
	allowFileShare, _ := awsec2.NewSecurityGroup(ctx, "allowFileShare", &awsec2.SecurityGroupArgs{
		Name:        pulumi.String("allow_efs_file_share"),
		Description: pulumi.String("Allow efs file share"),
		VpcId:       vpc.VpcId,
		Tags: pulumi.StringMap{
			"Name": pulumi.String("allow_efs_file_share"),
		},
	})

	vpc.PrivateSubnetIds.ApplyT(func(subnetIds []string) error {
		for _, subnetId := range subnetIds {
			subnet, err := awsec2.LookupSubnet(ctx, &awsec2.LookupSubnetArgs{
				Id: pulumi.StringRef(subnetId),
			}, nil)
			if err != nil {
				return err
			}

			_, err = awsvpc.NewSecurityGroupIngressRule(ctx, "allow_efs_ingress-"+subnetId, &awsvpc.SecurityGroupIngressRuleArgs{
				SecurityGroupId: allowFileShare.ID(),
				CidrIpv4:        pulumi.String(subnet.CidrBlock),
				FromPort:        pulumi.Int(2049),
				IpProtocol:      pulumi.String("tcp"),
				ToPort:          pulumi.Int(2049),
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	_, _ = awsvpc.NewSecurityGroupEgressRule(ctx, "allow_efs_egress", &awsvpc.SecurityGroupEgressRuleArgs{
		SecurityGroupId: allowFileShare.ID(),
		CidrIpv4:        pulumi.String("0.0.0.0/0"),
		FromPort:        pulumi.Int(0),
		IpProtocol:      pulumi.String("-1"),
		ToPort:          pulumi.Int(0),
	})

	// Create EFS mount targets for each private subnet
	vpc.PrivateSubnetIds.ApplyT(func(subnetIds []string) error {
		for _, subnetId := range subnetIds {
			_, err := efs.NewMountTarget(ctx, projectName+"-"+subnetId, &efs.MountTargetArgs{
				FileSystemId: fileShare.ID(),
				SubnetId:     pulumi.String(subnetId),
				SecurityGroups: pulumi.StringArray{
					allowFileShare.ID(),
				},
			}, pulumi.DependsOn([]pulumi.Resource{fileShare}))
			if err != nil {
				return err
			}
		}
		return nil
	})

	return fileShare, nil
}
