package awsdiscovery

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Filter Output //
func AwsFilter(tagKey, tagValue string) *ec2.DescribeInstancesInput {
	filter1 := &ec2.Filter{
		Name: aws.String("instance-state-name"),
		Values: []*string{
			aws.String("running"),
			aws.String("pending"),
		},
	}
	filter2 := &ec2.Filter{
		Name: aws.String("tag:" + tagKey),
		Values: []*string{
			aws.String(tagValue),
		},
	}
	filter := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{filter1, filter2},
	}
	return filter
}

func AwsSessIon(region string) *ec2.EC2 {
	session := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})
	return session
}

func AwsInstancePrivateIp(session *ec2.EC2, awsFilter *ec2.DescribeInstancesInput) []string {
	var ips []string
	resp, err := session.DescribeInstances(awsFilter)
	if err != nil {
		panic(err)
	}
	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			ips = append(ips, *inst.PrivateIpAddress)
		}
	}
	return ips
}

// func main() {
// 	session := awsSessIon("eu-west-1")
// 	filter := awsFilter("Location", "qa")
// 	ips := awsInstancePrivateIp(session, filter)
//
// 	for _, ip := range ips {
// 		fmt.Println(ip)
// 	}
// }
