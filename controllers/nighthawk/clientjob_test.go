package nighthawk

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"

	ksapi "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("Client Pod", func() {
	Describe("created from CR", func() {
		var cr ksapi.Nighthawk
		var job *batchv1.Job

		BeforeEach(func() {
			cr = ksapi.Nighthawk{
				Spec: ksapi.NighthawkSpec{
					Image: ksapi.ImageSpec{
						Name: "foo",
					},
					ClientConfiguration: ksapi.NighthawkClientConfigurationSpec{
						CmdLineArgs: []string{"--testing", "--things"},
						PodConfigurationSpec: ksapi.PodConfigurationSpec{
							Annotations: map[string]string{
								"annotation_one": "value_one",
							},
						},
					},
				},
			}
			job = NewClientJob(&cr)
		})

		Context("with cmdLineArgs specified", func() {
			It("--testing mode is set", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("--testing"))
			})
		})

		Context("with connectivity to service", func() {
			It("should not match service name", func() {
				service := NewServerService(&cr)
				Expect(job.ObjectMeta.Name).NotTo(
					Equal(service.ObjectMeta.Name))
			})
			It("should target the server service", func() {
				Expect(strings.Join(job.Spec.Template.Spec.Containers[0].Args, " ")).To(
					ContainSubstring(serverServiceNamePort(&cr)))
			})
		})

		Context("by default", func() {
			defaultBackoffLimit := int32(6)
			It("should retry 6 times", func() {
				Expect(job.Spec.BackoffLimit).To(
					Equal(&defaultBackoffLimit))
			})
		})

		Context("with added annotations", func() {
			It("should contain pod annotations", func() {
				Expect(job.ObjectMeta.Annotations).To(HaveKey("annotation_one"))
			})
		})
	})
})
