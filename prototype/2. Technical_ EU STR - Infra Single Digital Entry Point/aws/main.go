package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

const deployCluster = true
const chartsDeployed = true
const certManagerPolicyNameSuffix = "eks-certmanager-route53"
const fileSharePolicyNameSuffix = "eks-efs"

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		vpc, err := createVPC(ctx)
		if err != nil {
			return err
		}

		_, err = createFileShare(ctx, vpc)
		if err != nil {
			return err
		}

		if deployCluster {
			_, err = createCluster(ctx, vpc)
			if err != nil {
				return err
			}
		}

		if chartsDeployed {
			err = setDnsAlias(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func getProjectName(ctx *pulumi.Context) string {
	projectName := config.Get(ctx, "project:name")
	return projectName
}
