package aws

import (
	"log"

	"github.com/tidwall/gjson"
)

var whitelist = map[string]bool{"t2.nano": true, "t2.micro": true}

// EvaluateInstance determines what actions should be triggered
func EvaluateInstance(object string) ([]string, error) {
	log.Printf("evaluating aws instance")

	// Ignore everything without a "mimosa" tag
	if gjson.Get(object, "Tags.#(Key==\"mimosa\").Value").String() != "true" {
		return nil, nil
	}

	// Evaluate the rules
	var topics []string
	if shouldReap(object) {
		topics = append(topics, "actions-reaper")
		log.Printf("instance should be reaped")
	} else {
		log.Printf("instance should NOT be reaped")
	}
	if shouldRunTask(object) {
		topics = append(topics, "actions-runtask")
	}

	return topics, nil
}

func shouldReap(object string) bool {
	instanceType := gjson.Get(object, "InstanceType").String()
	if instanceType == "" {
		log.Printf("could not determine instance type")
		return false
	}
	_, present := whitelist[instanceType]
	return !present
}

func shouldRunTask(object string) bool {

	// Check criteria here

	return false

}
