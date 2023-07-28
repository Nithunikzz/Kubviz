package model

import "github.com/aquasecurity/trivy/pkg/types"

type Trivy struct {
	ID          string
	ClusterName string
	Report      types.Report
}
