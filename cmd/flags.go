package cmd

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/prdngr/red-sky/internal"
)

type DeploymentType string

const (
	deploymentTypeC2     DeploymentType = "c2"
	deploymentTypeKali   DeploymentType = "kali"
	deploymentTypeNessus DeploymentType = "nessus"
)

func (dt *DeploymentType) String() string {
	return string(*dt)
}

func (dt *DeploymentType) Set(value string) error {
	switch value {
	case string(deploymentTypeNessus), string(deploymentTypeKali), string(deploymentTypeC2):
		*dt = DeploymentType(value)
		return nil
	default:
		return errors.New(`must be one of "nessus", "kali", or "c2"`)
	}
}

func (dt *DeploymentType) Type() string {
	return "deploymentType"
}

type ingressRuleSliceValue struct {
	value *[]internal.IngressRule
}

func newIngressRuleSliceValue(val []internal.IngressRule, p *[]internal.IngressRule) *ingressRuleSliceValue {
	irsv := new(ingressRuleSliceValue)
	irsv.value = p
	*irsv.value = val
	return irsv
}

func (ir *ingressRuleSliceValue) String() string {
	ingressRuleStrings := make([]string, len(*ir.value))
	for index, value := range *ir.value {
		ingressRuleStrings[index] = value.Cidr + fmt.Sprint(value.Port)
	}

	return "[" + strings.Join(ingressRuleStrings, ",") + "]"
}

func (ir *ingressRuleSliceValue) Set(value string) error {
	ingressRules := strings.Split(value, ",")

	out := make([]internal.IngressRule, 0, len(ingressRules))
	for _, ingressRule := range ingressRules {
		parts := strings.Split(ingressRule, ":")

		if len(parts) != 2 {
			return fmt.Errorf("invalid string provided as ingress rule: %s", ingressRule)
		}

		_, cidr, err := net.ParseCIDR(strings.TrimSpace(parts[0]))
		if err != nil {
			return fmt.Errorf("invalid string being converted to CIDR: %s", parts[0])
		}

		port, err := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 32)
		if err != nil {
			return fmt.Errorf("invalid string being converted to port: %s", parts[1])
		}

		out = append(out, internal.IngressRule{
			Cidr: cidr.String(),
			Port: uint(port),
		})
	}

	*ir.value = out

	return nil
}

func (ir *ingressRuleSliceValue) Type() string {
	return "ingressRule"
}
