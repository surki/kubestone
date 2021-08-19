package nighthawk

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("nighthawk configmap", func() {
	Describe("cr with files for configVolumes", func() {
		var cr perfv1alpha1.Nighthawk
		var configMap *corev1.ConfigMap

		BeforeEach(func() {
			cr = perfv1alpha1.Nighthawk{
				Spec: perfv1alpha1.NighthawkSpec{
					ServerConfiguration: perfv1alpha1.NighthawkServerConfigurationSpec{
						ConfigsVolume: map[string]string{
							"file-1.yml": "content-1",
							"file-2.yml": "content-2",
						},
					},
				},
			}
			configMap = NewConfigMap(&cr)
		})

		Context("with config files specified", func() {
			It("should have them in the configmap", func() {
				Expect(configMap.Data["file-1.yml"]).To(
					Equal(cr.Spec.ServerConfiguration.ConfigsVolume["file-1.yml"]))
				Expect(configMap.Data["file-2.yml"]).To(
					Equal(cr.Spec.ServerConfiguration.ConfigsVolume["file-2.yml"]))
			})
		})
	})
})
