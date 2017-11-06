package helper

import (
	"time"

	"github.com/cloudfoundry/libbuildpack/cutlass"
	. "github.com/onsi/gomega"
)

func PushApp(app *cutlass.App) {
	Expect(app.Push()).To(Succeed())
	Eventually(app.InstanceStates, 20*time.Second).Should(Equal([]string{"RUNNING"}))
}
