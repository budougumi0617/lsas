package lsas

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
)

// Tag is AMI tag name and value.
type Tag struct {
	Key, Value string
}

// LoadConfig loads AWS setting with option.
func LoadConfig(region string) (aws.Config, error) {
	// FIXME Need to set flexible amount of option
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return aws.Config{}, err
	}
	if len(region) != 0 {
		cfg.Region = region
	}
	return cfg, nil
}

// Execute is main logic.
func Execute(region string, showHeader bool) error {
	cfg, err := LoadConfig(region)
	if err != nil {
		return err
	}

	svc := autoscaling.New(cfg)
	input := &autoscaling.DescribeAutoScalingGroupsInput{
		MaxRecords: aws.Int64(100),
	}
	req := svc.DescribeAutoScalingGroupsRequest(input)
	result, err := req.Send()
	if err != nil {
		return err
	}

	var tags []Tag
	for _, arg := range flag.Args() {
		if strings.Contains(arg, "=") {
			t := strings.Split(arg, "=")
			tags = append(tags, Tag{Key: t[0], Value: t[1]})
		}
	}
	results := []autoscaling.Group{}
	results = append(results, result.AutoScalingGroups...)

	var nt *string
	nt = result.NextToken

	for {
		if nt != nil {
			nt, err = func(nt *string) (*string, error) {
				// fmt.Println("exist more group yet" + aws.StringValue(nt))
				ri := &autoscaling.DescribeAutoScalingGroupsInput{
					MaxRecords: aws.Int64(100),
					NextToken:  nt,
				}
				rreq := svc.DescribeAutoScalingGroupsRequest(ri)
				rresult, err := rreq.Send()
				if err != nil {
					return nil, err
				}
				results = append(results, rresult.AutoScalingGroups...)
				return rresult.NextToken, nil
			}(nt)
			if err != nil {
				return err
			}
		} else {
			break
		}
	}

	filterd := []autoscaling.Group{}
	for _, asg := range results {
		// FIXME Need async output
		var matched int
		for _, astag := range asg.Tags {
			for _, t := range tags {
				if aws.StringValue(astag.Key) == t.Key && aws.StringValue(astag.Value) == t.Value {
					// fmt.Printf("%s = %+v\n", tag2[0], aws.StringValue(tag.Value))
					matched++
				}
			}
		}
		if matched == len(tags) {
			filterd = append(filterd, asg)
		}
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
	if showHeader {
		w.Write([]byte(fmt.Sprintf("autoscaling-group-name\tinstances\tdesired\tmin\tmax\tLaunch-configuration-name\n")))
	}
	// autoscaling group name
	// instances
	// desired capacity
	// min
	// max
	// launch configuration名
	for _, asg := range filterd {
		w.Write([]byte(fmt.Sprintf("%s\t%d\t%d\t%d\t%d\t%s\n",
			aws.StringValue(asg.AutoScalingGroupName),
			len(asg.Instances),
			aws.Int64Value(asg.DesiredCapacity),
			aws.Int64Value(asg.MinSize),
			aws.Int64Value(asg.MaxSize),
			aws.StringValue(asg.LaunchConfigurationName))),
		)
	}
	if err := w.Flush(); err != nil {
		return err
	}
	return nil
}
