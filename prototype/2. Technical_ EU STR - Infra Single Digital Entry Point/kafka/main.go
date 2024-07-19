package main

import (
	"encoding/json"
	"fmt"
	"github.com/pulumi/pulumi-confluentcloud/sdk/go/confluentcloud"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		awsRegion := config.Get(ctx, "aws:region")
		projectName := config.Get(ctx, "project:name")

		// Environment
		confluentEnv, err := confluentcloud.NewEnvironment(ctx, projectName, &confluentcloud.EnvironmentArgs{
			DisplayName: pulumi.String(projectName),
		})
		if err != nil {
			return err
		}

		// Cluster
		cluster, err := confluentcloud.NewKafkaCluster(ctx, projectName, &confluentcloud.KafkaClusterArgs{
			DisplayName:  pulumi.String(projectName),
			Availability: pulumi.String("SINGLE_ZONE"),
			Cloud:        pulumi.String("AWS"),
			Region:       pulumi.String(awsRegion),
			//Basic:        &confluentcloud.KafkaClusterBasicArgs{},
			Standard: &confluentcloud.KafkaClusterStandardArgs{},
			Environment: &confluentcloud.KafkaClusterEnvironmentArgs{
				Id: confluentEnv.ID(),
			},
		})

		// Service account
		serviceAccount, err := confluentcloud.NewServiceAccount(ctx, projectName, &confluentcloud.ServiceAccountArgs{
			DisplayName: pulumi.String(projectName),
			Description: pulumi.String(fmt.Sprintf("Service Account %s", projectName)),
		})
		if err != nil {
			return err
		}

		// Role binding
		roleBinding, err := confluentcloud.NewRoleBinding(ctx, projectName, &confluentcloud.RoleBindingArgs{
			CrnPattern: cluster.RbacCrn,
			Principal:  pulumi.Sprintf("User:%s", serviceAccount.ID()),
			RoleName:   pulumi.String("CloudClusterAdmin"),
		})

		// Cluster API key
		apiKey, err := confluentcloud.NewApiKey(ctx, projectName, &confluentcloud.ApiKeyArgs{
			DisplayName: pulumi.String(projectName),
			Owner: &confluentcloud.ApiKeyOwnerArgs{
				Id:         serviceAccount.ID(),
				Kind:       pulumi.String("ServiceAccount"),
				ApiVersion: serviceAccount.ApiVersion,
			},
			ManagedResource: &confluentcloud.ApiKeyManagedResourceArgs{
				Id:         cluster.ID(),
				Kind:       pulumi.String("Cluster"),
				ApiVersion: cluster.ApiVersion,
				Environment: &confluentcloud.ApiKeyManagedResourceEnvironmentArgs{
					Id: confluentEnv.ID(),
				},
			},
		}, pulumi.DependsOn([]pulumi.Resource{roleBinding}))
		if err != nil {
			return err
		}

		topicsJson := config.Get(ctx, "kafka:topics")
		var topics []string
		err = json.Unmarshal([]byte(topicsJson), &topics)
		if err != nil {
			return err
		}

		for _, topic := range topics {
			fmt.Println("Create Kafka Topic", topic)
			err = createTopic(ctx, topic, cluster, apiKey)
			if err != nil {
				return err
			}
		}

		// Create application service account
		serviceAccountName := fmt.Sprintf("%s-app", projectName)
		serviceAccountApp, err := confluentcloud.NewServiceAccount(ctx, serviceAccountName, &confluentcloud.ServiceAccountArgs{
			DisplayName: pulumi.String(serviceAccountName),
			Description: pulumi.String(fmt.Sprintf("Service Account %s", serviceAccountName)),
		})
		if err != nil {
			return err
		}

		// Create application api key
		appApiKey, err := confluentcloud.NewApiKey(ctx, serviceAccountName, &confluentcloud.ApiKeyArgs{
			DisplayName: pulumi.String(serviceAccountName),
			Owner: &confluentcloud.ApiKeyOwnerArgs{
				Id:         serviceAccountApp.ID(),
				Kind:       pulumi.String("ServiceAccount"),
				ApiVersion: serviceAccountApp.ApiVersion,
			},
			ManagedResource: &confluentcloud.ApiKeyManagedResourceArgs{
				Id:         cluster.ID(),
				Kind:       pulumi.String("Cluster"),
				ApiVersion: cluster.ApiVersion,
				Environment: &confluentcloud.ApiKeyManagedResourceEnvironmentArgs{
					Id: confluentEnv.ID(),
				},
			},
		})
		if err != nil {
			return err
		}

		for _, topic := range topics {
			fmt.Println("Create Kafka Topic acl", topic)
			err = createTopicAcl(ctx, topic, "READ", cluster, apiKey, serviceAccountApp)
			err = createTopicAcl(ctx, topic, "WRITE", cluster, apiKey, serviceAccountApp)
			if err != nil {
				return err
			}
		}

		_, err = confluentcloud.NewKafkaAcl(ctx, "kafkaAclAppConsumerGroup", &confluentcloud.KafkaAclArgs{
			Host:         pulumi.String("*"),
			Operation:    pulumi.String("READ"),
			PatternType:  pulumi.String("LITERAL"),
			Permission:   pulumi.String("ALLOW"),
			Principal:    pulumi.Sprintf("User:%s", serviceAccountApp.ID()),
			ResourceName: pulumi.String("*"),
			ResourceType: pulumi.String("GROUP"),
			RestEndpoint: cluster.RestEndpoint,
			Credentials: &confluentcloud.KafkaAclCredentialsArgs{
				Key:    apiKey.ID(),
				Secret: apiKey.Secret,
			},
			KafkaCluster: &confluentcloud.KafkaAclKafkaClusterArgs{
				Id: cluster.ID(),
			},
		})
		if err != nil {
			return err
		}

		ctx.Export("ClusterBootstrapEndpoint", cluster.BootstrapEndpoint)
		ctx.Export("appApiKey", appApiKey.ID())

		appApiSecret := pulumi.ToSecret(appApiKey.Secret).(pulumi.StringOutput)
		ctx.Export("appApiSecret", appApiSecret)

		return nil
	})
}

func createTopic(ctx *pulumi.Context, topicName string, cluster *confluentcloud.KafkaCluster, apiKey *confluentcloud.ApiKey) error {
	_, err := confluentcloud.NewKafkaTopic(ctx, topicName, &confluentcloud.KafkaTopicArgs{
		TopicName: pulumi.String(topicName),

		KafkaCluster: &confluentcloud.KafkaTopicKafkaClusterArgs{
			Id: cluster.ID(),
		},
		PartitionsCount: pulumi.Int(2),
		RestEndpoint:    cluster.RestEndpoint,
		Credentials: &confluentcloud.KafkaTopicCredentialsArgs{
			Key:    apiKey.ID(),
			Secret: apiKey.Secret,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func createTopicAcl(ctx *pulumi.Context,
	topic string,
	operation string,
	cluster *confluentcloud.KafkaCluster,
	apiKey *confluentcloud.ApiKey,
	serviceAccountApp *confluentcloud.ServiceAccount,
) error {
	// Create application acl
	name := fmt.Sprintf("kafkaAcl%s%s", topic, operation)
	_, err := confluentcloud.NewKafkaAcl(ctx, name, &confluentcloud.KafkaAclArgs{
		Host:         pulumi.String("*"),
		Operation:    pulumi.String(operation),
		PatternType:  pulumi.String("LITERAL"),
		Permission:   pulumi.String("ALLOW"),
		Principal:    pulumi.Sprintf("User:%s", serviceAccountApp.ID()),
		ResourceName: pulumi.String(topic),
		ResourceType: pulumi.String("TOPIC"),
		RestEndpoint: cluster.RestEndpoint,
		Credentials: &confluentcloud.KafkaAclCredentialsArgs{
			Key:    apiKey.ID(),
			Secret: apiKey.Secret,
		},
		KafkaCluster: &confluentcloud.KafkaAclKafkaClusterArgs{
			Id: cluster.ID(),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
