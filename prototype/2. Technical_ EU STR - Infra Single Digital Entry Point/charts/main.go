package main

import (
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	helmv4 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v4"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	yamlv2 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

const deployStrChart = true
const deployNginxIngressChart = true
const requestCertificate = true
const deployCertManagerChart = true

const certManagerRoleNameSuffix = "eks-certmanager-route53"
const fileShareRoleNameSuffix = "eks-efs"

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		awsStackRef := config.Get(ctx, "project:awsStackRef")
		awsStack, err := pulumi.NewStackReference(ctx, awsStackRef, nil)
		if err != nil {
			return err
		}

		err = deployEfsChart(ctx, awsStack)

		if deployCertManagerChart {
			err := deployCertManager(ctx, awsStack)
			if err != nil {
				return err
			}
		}

		if requestCertificate {
			err := requestLetsencryptCertificate(ctx, awsStack)
			if err != nil {
				return err
			}
		}

		if deployNginxIngressChart {
			err := deployNginxIngress(ctx, awsStack)
			if err != nil {
				return err
			}
		}

		// Always create the namespace, so that cert-manager can inject the secret
		namespace := config.Get(ctx, "project:namespace")

		_, err = corev1.NewNamespace(ctx, getProjectName(ctx), &corev1.NamespaceArgs{
			Metadata: &metav1.ObjectMetaArgs{Name: pulumi.String(namespace)},
		})
		if err != nil {
			return err
		}

		if deployStrChart {
			err = deployApp(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func deployApp(ctx *pulumi.Context) error {
	// deploy the str application
	namespace := config.Get(ctx, "project:namespace")

	// Get config values
	imagePassword := config.Get(ctx, "image:password")
	chart := config.Get(ctx, "chart:file")
	valuesFile := config.Get(ctx, "values:file")
	kafkaBoostrapServers := config.Get(ctx, "kafka:boostrapServers")
	username := config.Get(ctx, "kafka:username")
	password := config.Get(ctx, "kafka:password")
	fileSystemId := config.Get(ctx, "project:efs")

	crt := config.Get(ctx, "tls:crt")
	key := config.Get(ctx, "tls:key")

	_, err := helmv4.NewChart(ctx, namespace, &helmv4.ChartArgs{
		Namespace: pulumi.String(namespace),
		Chart:     pulumi.String(chart),
		ValueYamlFiles: pulumi.AssetOrArchiveArray{
			pulumi.NewFileAsset(valuesFile),
		},
		Values: pulumi.Map{
			"image": pulumi.Map{
				"password": pulumi.String(imagePassword),
			},
			"kafka": pulumi.Map{
				"boostrapServers": pulumi.String(kafkaBoostrapServers),
				"username":        pulumi.String(username),
				"password":        pulumi.String(password),
			},
			"tls": pulumi.Map{
				"crt": pulumi.String(crt),
				"key": pulumi.String(key),
			},
			"efsStorage": pulumi.Map{
				"fileSystemId": pulumi.String(fileSystemId),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func deployNginxIngress(ctx *pulumi.Context, awsStack *pulumi.StackReference) error {
	vpcCidr := awsStack.GetOutput(pulumi.String("vpcCidr"))

	namespace := "ingress-nginx"
	chart := fmt.Sprintf("./%s", namespace)
	valuesFile := fmt.Sprintf("%s/values.yaml", chart)

	// Namespace
	ns, err := corev1.NewNamespace(ctx, namespace, &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{Name: pulumi.String(namespace)},
	})
	if err != nil {
		return err
	}

	// Deploy chart
	_, err = helmv4.NewChart(ctx, namespace, &helmv4.ChartArgs{
		Namespace: ns.Metadata.Name(),
		Chart:     pulumi.String(chart),
		ValueYamlFiles: pulumi.AssetOrArchiveArray{
			pulumi.NewFileAsset(valuesFile),
		},
		Values: pulumi.Map{
			"configmap": pulumi.Map{
				"proxyRealIpCidr": vpcCidr,
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func deployCertManager(ctx *pulumi.Context, awsStack *pulumi.StackReference) error {

	current, err := aws.GetCallerIdentity(ctx, nil, nil)
	if err != nil {
		return err
	}
	roleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", current.AccountId, getProjectName(ctx)+"-"+certManagerRoleNameSuffix)

	// Namespace
	namespace := "cert-manager"
	ns, err := corev1.NewNamespace(ctx, namespace, &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{Name: pulumi.String(namespace)},
	})
	if err != nil {
		return err
	}

	// Deploy chart
	_, err = helmv4.NewChart(ctx, namespace, &helmv4.ChartArgs{
		Namespace: ns.Metadata.Name(),
		RepositoryOpts: &helmv4.RepositoryOptsArgs{
			Repo: pulumi.String("https://charts.jetstack.io"),
		},
		Chart:   pulumi.String("cert-manager"),
		Version: pulumi.String("v1.15.0"),
		Values: pulumi.Map{
			"crds": pulumi.Map{
				"enabled": pulumi.Bool(true),
			},
			"webhook": pulumi.Map{
				"timeoutSeconds": pulumi.Int(5),
			},
			"serviceAccount": pulumi.Map{
				"annotations": pulumi.Map{
					"eks.amazonaws.com/role-arn": pulumi.String(roleArn),
				},
			},
			"securityContext": pulumi.Map{
				"fsGroup": pulumi.Int(1001),
			},
			"extraArgs": pulumi.StringArray{
				pulumi.String("--issuer-ambient-credentials"),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func requestLetsencryptCertificate(ctx *pulumi.Context, awsStack *pulumi.StackReference) error {
	// get config values
	email := config.Get(ctx, "cert:emailContact")
	server := config.Get(ctx, "acme:serverUrl")
	region := config.Get(ctx, "aws:region")
	namespace := config.Get(ctx, "project:namespace")

	domainHostedZoneId := awsStack.GetOutput(pulumi.String("domainHostedZoneId"))
	fqdn := awsStack.GetOutput(pulumi.String("fqdn"))

	_, err := yamlv2.NewConfigGroup(ctx, "AcmeIssuer", &yamlv2.ConfigGroupArgs{
		Objs: pulumi.Array{
			pulumi.Map{
				"apiVersion": pulumi.String("cert-manager.io/v1"),
				"kind":       pulumi.String("Issuer"),
				"metadata": pulumi.Map{
					"name":      pulumi.String("letsencrypt"),
					"namespace": pulumi.String(namespace),
				},
				"spec": pulumi.Map{
					"acme": pulumi.Map{
						"server": pulumi.String(server),
						"email":  pulumi.String(email),
						"privateKeySecretRef": pulumi.Map{
							"name": pulumi.String("letsencrypt"),
						},
						"solvers": pulumi.Array{
							pulumi.Map{
								"dns01": pulumi.Map{
									"route53": pulumi.Map{
										"hostedZoneID": domainHostedZoneId,
										"region":       pulumi.String(region),
									},
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	// Request certificate in str namespace
	_, err = yamlv2.NewConfigGroup(ctx, "cert", &yamlv2.ConfigGroupArgs{
		Objs: pulumi.Array{
			pulumi.Map{
				"apiVersion": pulumi.String("cert-manager.io/v1"),
				"kind":       pulumi.String("Certificate"),
				"metadata": pulumi.Map{
					"name":      fqdn,
					"namespace": pulumi.String(namespace),
				},
				"spec": pulumi.Map{
					"secretName": fqdn,
					"isCA":       pulumi.Bool(false),
					"usages": pulumi.StringArray{
						pulumi.String("server auth"),
						pulumi.String("client auth"),
					},
					"issuerRef": pulumi.Map{
						"name": pulumi.String("letsencrypt"),
						"kind": pulumi.String("Issuer"),
					},
					"commonName": fqdn,
					"dnsNames": pulumi.Array{
						fqdn,
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func deployEfsChart(ctx *pulumi.Context, awsStack *pulumi.StackReference) error {
	// Create sa
	saName := "efs-csi-controller-sa"
	namespace := "kube-system"

	current, err := aws.GetCallerIdentity(ctx, nil, nil)
	if err != nil {
		return err
	}
	roleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", current.AccountId, getProjectName(ctx)+"-"+fileShareRoleNameSuffix)

	_, err = corev1.NewServiceAccount(ctx, saName, &corev1.ServiceAccountArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(saName),
			Namespace: pulumi.String(namespace),
			Annotations: pulumi.StringMap{
				"eks.amazonaws.com/role-arn": pulumi.String(roleArn),
			},
			Labels: pulumi.StringMap{
				"app.kubernetes.io/name": pulumi.String("aws-efs-csi-driver"),
			},
		},
	})
	if err != nil {
		return err
	}

	// Deploy chart
	_, err = helmv4.NewChart(ctx, "eks-efs", &helmv4.ChartArgs{
		Namespace: pulumi.String(namespace),
		RepositoryOpts: &helmv4.RepositoryOptsArgs{
			Repo: pulumi.String("https://kubernetes-sigs.github.io/aws-efs-csi-driver"),
		},
		Chart: pulumi.String("aws-efs-csi-driver"),
		//Version: pulumi.String("v3.0.5"),
		Values: pulumi.Map{
			"controller": pulumi.Map{
				"serviceAccount": pulumi.Map{
					"create": pulumi.Bool(false),
					"name":   pulumi.String(saName),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func getProjectName(ctx *pulumi.Context) string {
	projectName := config.Get(ctx, "project:name")
	return projectName
}
