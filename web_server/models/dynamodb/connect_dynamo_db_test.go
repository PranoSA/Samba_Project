package dynamodb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestDynamoDBConnectLogic(t *testing.T) {
	t.Run("IDK", func(t *testing.T) {
		fmt.Print("Reading Credentials")

		c, err := config.LoadDefaultConfig(context.TODO(),
			config.WithSharedConfigProfile("default"),
		)
		if err != nil {
			// handle error
			t.Fatal("IDK")
		}

		svc := dynamodb.NewFromConfig(c)

		out, err := svc.Scan(context.TODO(), &dynamodb.ScanInput{
			TableName: aws.String("SambaShares"),
		})
		if err != nil {
			panic(err)
		}
		t.Log(out.Count)

		t.Log(out.Items[0]["diskid"])

		svc.PutItem(context.Background(), &dynamodb.PutItemInput{
			TableName: aws.String("Samba_Hosts"),
			Item: map[string]types.AttributeValue{
				"hostid": &types.AttributeValueMemberN{Value: "1"},
			},
		})

	})
}
