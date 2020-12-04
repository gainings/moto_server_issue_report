package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

func main() {
	endpoint := "http://moto:5000"
	session := session.Must(session.NewSession(
		&aws.Config{
			Region:      aws.String("ap-northeast-1"),
			Endpoint:    aws.String(endpoint),
			Credentials: credentials.NewStaticCredentials("dummy", "dummy", ""),
		}))
	elbv2Client := elbv2.New(session)
	ec2Client := ec2.New(session)
	ctx := context.Background()

	subnets := CreateVPC(ec2Client, elbv2Client)

	dlbo, err := elbv2Client.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{})
	if err != nil {
		fmt.Println("------Second Run: Time parse error-------")
		log.Fatal(err)
	}
	fmt.Println(dlbo)

	clbi := &elbv2.CreateLoadBalancerInput{
		Name:    aws.String("test-lb"),
		Subnets: aws.StringSlice(subnets),
		Type:    aws.String("application"),
	}
	_, err = elbv2Client.CreateLoadBalancer(clbi)
	if err != nil {
		fmt.Println("------First Run: Got 400 without error message, but success create LoadBalancer-----")
		log.Fatal(err)
	}
	_, err = elbv2Client.DescribeLoadBalancersWithContext(ctx, &elbv2.DescribeLoadBalancersInput{
		Names: aws.StringSlice([]string{"test-lb"}),
	})
	if err != nil {
		log.Fatal(err)
	}
}

func CreateVPC(ec2Client *ec2.EC2, elbv2Client *elbv2.ELBV2) []string {
	ctx := context.Background()
	cvo, err := ec2Client.CreateVpcWithContext(ctx, &ec2.CreateVpcInput{
		CidrBlock: aws.String("10.0.0.0/16"),
	})
	if err != nil {
		log.Fatal(err)
	}

	cigo, err := ec2Client.CreateInternetGatewayWithContext(ctx, &ec2.CreateInternetGatewayInput{})
	if err != nil {
		log.Fatal(err)
	}
	ec2Client.AttachInternetGatewayWithContext(ctx, &ec2.AttachInternetGatewayInput{
		VpcId:             cvo.Vpc.VpcId,
		InternetGatewayId: cigo.InternetGateway.InternetGatewayId,
	})

	az1CSO, err := ec2Client.CreateSubnetWithContext(ctx, &ec2.CreateSubnetInput{
		CidrBlock:        aws.String("10.0.0.0/24"),
		VpcId:            cvo.Vpc.VpcId,
		AvailabilityZone: aws.String("ap-northeast-1a"),
	})
	if err != nil {
		log.Fatal(err)
	}
	rt1CRTO, err := ec2Client.CreateRouteTableWithContext(ctx, &ec2.CreateRouteTableInput{
		VpcId: cvo.Vpc.VpcId,
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = ec2Client.AssociateRouteTableWithContext(ctx, &ec2.AssociateRouteTableInput{
		SubnetId:     az1CSO.Subnet.SubnetId,
		RouteTableId: rt1CRTO.RouteTable.RouteTableId,
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = ec2Client.CreateRouteWithContext(ctx, &ec2.CreateRouteInput{
		RouteTableId:         rt1CRTO.RouteTable.RouteTableId,
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		GatewayId:            cigo.InternetGateway.InternetGatewayId,
	})
	if err != nil {
		log.Fatal(err)
	}

	az2CSO, err := ec2Client.CreateSubnetWithContext(ctx, &ec2.CreateSubnetInput{
		CidrBlock:        aws.String("10.0.1.0/24"),
		VpcId:            cvo.Vpc.VpcId,
		AvailabilityZone: aws.String("ap-northeast-1c"),
	})
	if err != nil {
		log.Fatal(err)
	}
	rt2CRTO, err := ec2Client.CreateRouteTableWithContext(ctx, &ec2.CreateRouteTableInput{
		VpcId: cvo.Vpc.VpcId,
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = ec2Client.AssociateRouteTableWithContext(ctx, &ec2.AssociateRouteTableInput{
		SubnetId:     az2CSO.Subnet.SubnetId,
		RouteTableId: rt2CRTO.RouteTable.RouteTableId,
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = ec2Client.CreateRouteWithContext(ctx, &ec2.CreateRouteInput{
		RouteTableId:         rt2CRTO.RouteTable.RouteTableId,
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		GatewayId:            cigo.InternetGateway.InternetGatewayId,
	})
	if err != nil {
		log.Fatal(err)
	}
	return []string{
		*az1CSO.Subnet.SubnetId,
		*az2CSO.Subnet.SubnetId,
	}
}
