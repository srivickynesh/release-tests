package cli

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/tkn"
)

var _ = gauge.Step("Setup CLI paths and assert", func() {
	tkn.SetupCLIPaths()

	tkn.AssertClientVersion("tkn")
	tkn.AssertClientVersion("tkn-pac")
	tkn.AssertClientVersion("opc")
	tkn.AssertServerVersion("opc")

})
