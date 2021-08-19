package nighthawk

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ksapi "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("Server Service", func() {
	Describe("created from CR", func() {
		var cr ksapi.Nighthawk
		var service *corev1.Service

		BeforeEach(func() {
			cr = ksapi.Nighthawk{
				Spec: ksapi.NighthawkSpec{
					Image: ksapi.ImageSpec{
						Name: "foo",
					},
					ServerConfiguration: ksapi.NighthawkServerConfigurationSpec{
						Port: 1234,
					},
				},
			}
			service = NewServerService(&cr)
		})

		Context("with default settings", func() {
			It("should use TCP protocol", func() {
				Expect(service.Spec.Ports[0].Protocol).To(
					Equal(corev1.ProtocolTCP))
			})
		})

		Context("crosschecked with server deployment", func() {
			configMap := corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{Name: "cm"},
			}

			service := NewServerService(&cr)
			deployment := NewServerDeployment(&cr, &configMap)
			It("should match on port", func() {
				Expect(service.Spec.Ports[0].Protocol).To(
					Equal(deployment.Spec.Template.Spec.Containers[0].Ports[0].Protocol))
			})
			It("should match on selectors", func() {
				Expect(service.Spec.Selector).To(
					Equal(deployment.Spec.Template.ObjectMeta.Labels))
			})
			It("should match on namespace", func() {
				Expect(service.ObjectMeta.Namespace).To(
					Equal(deployment.ObjectMeta.Namespace))
			})
		})
	})
})
