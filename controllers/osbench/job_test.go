package osbench

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("osbench job", func() {
	Describe("cr with cmd args", func() {
		var cr perfv1alpha1.Osbench
		var job *batchv1.Job

		BeforeEach(func() {
			cr = perfv1alpha1.Osbench{
				Spec: perfv1alpha1.OsbenchSpec{
					Image: perfv1alpha1.ImageSpec{
						Name: "surki/osbench:latest",
					},
					Options:  "/tmp",
					TestName: "create_files",
				},
			}
			job = NewJob(&cr)
		})

		Context("with command line arguments", func() {
			It("should have the same options", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("/tmp"))
			})
			It("should have the same testName", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Command).To(
					ContainElement("create_files"))
			})
		})
	})
})
