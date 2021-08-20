package perfbench

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("perfbench job", func() {
	Describe("cr with cmd args", func() {
		var cr perfv1alpha1.Perfbench
		var job *batchv1.Job

		BeforeEach(func() {
			cr = perfv1alpha1.Perfbench{
				Spec: perfv1alpha1.PerfbenchSpec{
					Image: perfv1alpha1.ImageSpec{
						Name:       "xridge/perfbench:test",
						PullPolicy: "Always",
						PullSecret: "a-pull-secret",
					},
					CmdLineArgs: []string{"--testing", "--something"},
				},
			}
			job = NewJob(&cr)
		})

		Context("with default settings", func() {
			It("should have perf bench command", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Command).To(
					ContainElement("perf"))
				Expect(job.Spec.Template.Spec.Containers[0].Command).To(
					ContainElement("bench"))
				Expect(*job.Spec.Template.Spec.Containers[0].SecurityContext.Privileged).To(
					BeTrue())
			})
		})

		Context("with command line args specified", func() {
			It("should have the same args", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("--testing"))
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("--something"))
			})
		})
	})
})
