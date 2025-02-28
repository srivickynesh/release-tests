package cli

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/tkn"
)

var _ = gauge.Step("Run OPC tests", func() {
	tkn.SetupCLIPaths()

})
