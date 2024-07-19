package main

import (
	"encoding/json"
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-awsx/sdk/v2/go/awsx/ec2"
	"github.com/pulumi/pulumi-eks/sdk/go/eks"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"strings"
)

func createCluster(ctx *pulumi.Context, vpc *ec2.Vpc) (*eks.Cluster, error) {

	// Create an EKS cluster inside of the VPC.
	cluster, err := eks.NewCluster(ctx, getProjectName(ctx), &eks.ClusterArgs{
		Name:                         pulumi.String(getProjectName(ctx)),
		VpcId:                        vpc.VpcId,
		PublicSubnetIds:              vpc.PublicSubnetIds,
		PrivateSubnetIds:             vpc.PrivateSubnetIds,
		NodeAssociatePublicIpAddress: pulumi.BoolRef(false),
		UseDefaultVpcCni:             pulumi.BoolRef(true),
		CreateOidcProvider:           pulumi.BoolPtr(true),
		DesiredCapacity:              pulumi.Int(2),
		MinSize:                      pulumi.Int(1),
		MaxSize:                      pulumi.Int(3),
		InstanceType:                 pulumi.String("t3.medium"),
		EnabledClusterLogTypes: pulumi.StringArray{
			pulumi.String("api"),
			pulumi.String("audit"),
			pulumi.String("authenticator"),
		},
	})
	if err != nil {
		return nil, err
	}

	err = createCertManagerIam(ctx, cluster)
	if err != nil {
		return nil, err
	}

	err = createFileShareIam(ctx, cluster)
	if err != nil {
		return nil, err
	}

	// Export the cluster's kubeconfig
	ctx.Export("kubeconfig", cluster.Kubeconfig)

	return cluster, nil
}

// ClusterConfig define the structure to match the JSON structure
type ClusterConfig struct {
	Clusters []struct {
		Cluster struct {
			Server string `json:"server"`
		} `json:"cluster"`
	} `json:"clusters"`
}

// getEksHash processes the kubeConfig JSON and extracts the EKS hash
func getEksHash(kubeConfig pulumi.StringOutput) pulumi.StringOutput {
	return kubeConfig.ApplyT(func(jsonStr string) (string, error) {
		var clusterConfig ClusterConfig

		// Parse the JSON string
		err := json.Unmarshal([]byte(jsonStr), &clusterConfig)
		if err != nil {
			return "", err
		}

		// Extract the server URL
		if len(clusterConfig.Clusters) > 0 {
			server := clusterConfig.Clusters[0].Cluster.Server
			// Remove the https:// prefix and extract the value from the server URL
			if strings.HasPrefix(server, "https://") {
				server = server[len("https://"):]
			}
			parts := strings.Split(server, ".")
			if len(parts) > 0 {
				return parts[0], nil
			}
			return "", fmt.Errorf("server URL format is not as expected")
		}
		return "", fmt.Errorf("no clusters found in the configuration")
	}).(pulumi.StringOutput)
}

func createCertManagerIam(ctx *pulumi.Context, cluster *eks.Cluster) error {
	policyName := fmt.Sprintf("%s-%s", getProjectName(ctx), certManagerPolicyNameSuffix)
	roleName := policyName

	policyDocument, err := getCertManagerPolicy(ctx)
	if err != nil {
		return err
	}

	policy, err := iam.NewPolicy(ctx, policyName, &iam.PolicyArgs{
		Name:   pulumi.String(policyName),
		Path:   pulumi.String("/"),
		Policy: pulumi.String(policyDocument.Json),
	})
	if err != nil {
		return err
	}

	eksHash := getEksHash(cluster.KubeconfigJson)
	eksHash.ApplyT(func(eksHashString string) error {
		sa := fmt.Sprintf("system:serviceaccount:cert-manager:cert-manager")
		assumePolicy, err := getAssumePolicy(ctx, eksHashString, sa)
		if err != nil {
			return err
		}

		role, err := iam.NewRole(ctx, policyName, &iam.RoleArgs{
			Name:             pulumi.String(roleName),
			AssumeRolePolicy: pulumi.String(assumePolicy.Json),
		})
		if err != nil {
			return err
		}

		_, err = iam.NewRolePolicyAttachment(ctx, policyName, &iam.RolePolicyAttachmentArgs{
			Role:      role.Name,
			PolicyArn: policy.Arn,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return nil
}

func createFileShareIam(ctx *pulumi.Context, cluster *eks.Cluster) error {
	policyName := fmt.Sprintf("%s-%s", getProjectName(ctx), fileSharePolicyNameSuffix)
	roleName := policyName

	policyDocument, err := getFileSharePolicy(ctx)
	if err != nil {
		return err
	}

	policy, err := iam.NewPolicy(ctx, policyName, &iam.PolicyArgs{
		Name:   pulumi.String(policyName),
		Path:   pulumi.String("/"),
		Policy: pulumi.String(policyDocument.Json),
	})
	if err != nil {
		return err
	}

	eksHash := getEksHash(cluster.KubeconfigJson)
	eksHash.ApplyT(func(eksHashString string) error {
		sa := fmt.Sprintf("system:serviceaccount:kube-system:efs-csi-controller-sa")
		assumePolicy, err := getAssumePolicy(ctx, eksHashString, sa)
		if err != nil {
			return err
		}

		role, err := iam.NewRole(ctx, policyName, &iam.RoleArgs{
			Name:             pulumi.String(roleName),
			AssumeRolePolicy: pulumi.String(assumePolicy.Json),
		})
		if err != nil {
			return err
		}

		_, err = iam.NewRolePolicyAttachment(ctx, policyName, &iam.RolePolicyAttachmentArgs{
			Role:      role.Name,
			PolicyArn: policy.Arn,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return nil
}

// getAssumePolicy Construct an assume policy for EKS with OIDC
func getAssumePolicy(ctx *pulumi.Context, eksHashString string, sa string) (*iam.GetPolicyDocumentResult, error) {
	// Get current identity
	current, err := aws.GetCallerIdentity(ctx, nil, nil)
	if err != nil {
		return nil, err
	}

	region := config.Get(ctx, "aws:region")

	federated := fmt.Sprintf("arn:aws:iam::%s:oidc-provider/oidc.eks.%s.amazonaws.com/id/%s", current.AccountId, region, eksHashString)
	variable := fmt.Sprintf("oidc.eks.eu-west-1.amazonaws.com/id/%s:sub", eksHashString)

	assumePolicy, err := iam.GetPolicyDocument(ctx, &iam.GetPolicyDocumentArgs{
		Statements: []iam.GetPolicyDocumentStatement{
			{
				Effect: pulumi.StringRef("Allow"),
				Principals: []iam.GetPolicyDocumentStatementPrincipal{
					{
						Type: "Federated",
						Identifiers: []string{
							federated,
						},
					},
				},
				Conditions: []iam.GetPolicyDocumentStatementCondition{
					{
						Test:     "StringEquals",
						Variable: variable,
						Values:   []string{sa},
					},
				},
				Actions: []string{
					"sts:AssumeRoleWithWebIdentity",
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return assumePolicy, nil
}

func getCertManagerPolicy(ctx *pulumi.Context) (*iam.GetPolicyDocumentResult, error) {
	policyDocument, err := iam.GetPolicyDocument(ctx, &iam.GetPolicyDocumentArgs{
		Statements: []iam.GetPolicyDocumentStatement{
			{
				Effect:    pulumi.StringRef("Allow"),
				Resources: []string{"arn:aws:route53:::change/*"},
				Actions: []string{
					"route53:GetChange",
				},
			},
			{
				Effect:    pulumi.StringRef("Allow"),
				Resources: []string{"arn:aws:route53:::hostedzone/*"},
				Actions: []string{
					"route53:ChangeResourceRecordSets",
					"route53:ListResourceRecordSets",
				},
			},
			{
				Effect:    pulumi.StringRef("Allow"),
				Resources: []string{"*"},
				Actions: []string{
					"route53:ListHostedZonesByName",
				},
			},
		},
	}, nil)
	if err != nil {
		return nil, err
	}

	return policyDocument, nil
}

func getFileSharePolicy(ctx *pulumi.Context) (*iam.GetPolicyDocumentResult, error) {
	policyDocument, err := iam.GetPolicyDocument(ctx, &iam.GetPolicyDocumentArgs{
		Statements: []iam.GetPolicyDocumentStatement{
			{
				Effect:    pulumi.StringRef("Allow"),
				Resources: []string{"*"},
				Actions: []string{
					"elasticfilesystem:DescribeAccessPoints",
					"elasticfilesystem:DescribeFileSystems",
					"elasticfilesystem:DescribeMountTargets",
					"ec2:DescribeAvailabilityZones",
				},
			},
			{
				Effect:    pulumi.StringRef("Allow"),
				Actions:   []string{"elasticfilesystem:CreateAccessPoint"},
				Resources: []string{"*"},
				Conditions: []iam.GetPolicyDocumentStatementCondition{
					{
						Test:     "StringLike",
						Variable: "aws:RequestTag/efs.csi.aws.com/cluster",
						Values:   []string{"true"},
					},
				},
			},
			{
				Effect:    pulumi.StringRef("Allow"),
				Actions:   []string{"elasticfilesystem:TagResource"},
				Resources: []string{"*"},
				Conditions: []iam.GetPolicyDocumentStatementCondition{
					{
						Test:     "StringLike",
						Variable: "aws:RequestTag/efs.csi.aws.com/cluster",
						Values:   []string{"true"},
					},
				},
			},
			{
				Effect:    pulumi.StringRef("Allow"),
				Actions:   []string{"elasticfilesystem:DeleteAccessPoint"},
				Resources: []string{"*"},
				Conditions: []iam.GetPolicyDocumentStatementCondition{
					{
						Test:     "StringLike",
						Variable: "aws:RequestTag/efs.csi.aws.com/cluster",
						Values:   []string{"true"},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return nil, err
	}

	return policyDocument, nil
}
