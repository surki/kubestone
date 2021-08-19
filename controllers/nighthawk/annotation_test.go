package nighthawk

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ksapi "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("Pod Annotations", func() {
	Describe("created from CR", func() {
		var cr ksapi.Nighthawk
		var serverDeployment *appsv1.Deployment
		var clientService *corev1.Service
		var jobSpec *batchv1.Job

		BeforeEach(func() {
			tolerationSeconds := int64(17)
			cr = ksapi.Nighthawk{
				Spec: ksapi.NighthawkSpec{
					Image: ksapi.ImageSpec{
						Name:       "foo",
						PullPolicy: "Always",
						PullSecret: "pull-secret",
					},
					ServerConfiguration: ksapi.NighthawkServerConfigurationSpec{
						CmdLineArgs: []string{"--testing", "--things"},
						PodConfigurationSpec: ksapi.PodConfigurationSpec{
							Annotations: map[string]string{"anno_two": "exists"},
							PodLabels:   map[string]string{"labels": "are", "really": "useful"},
							PodScheduling: ksapi.PodSchedulingSpec{
								Affinity: &corev1.Affinity{
									NodeAffinity: &corev1.NodeAffinity{
										RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
											NodeSelectorTerms: []corev1.NodeSelectorTerm{
												{
													MatchExpressions: []corev1.NodeSelectorRequirement{
														{
															Key:      "mutated",
															Operator: corev1.NodeSelectorOperator(corev1.NodeSelectorOpIn),
															Values:   []string{"nano-virus"},
														},
													},
												},
											},
										},
									},
								},
								Tolerations: []corev1.Toleration{
									{
										Key:               "genetic-code",
										Operator:          corev1.TolerationOperator(corev1.TolerationOpExists),
										Value:             "distressed",
										Effect:            corev1.TaintEffect(corev1.TaintEffectNoExecute),
										TolerationSeconds: &tolerationSeconds,
									},
								},
								NodeSelector: map[string]string{
									"atomized": "spiral",
								},
								NodeName: "energy-spike-07",
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("500m"),
									corev1.ResourceMemory: resource.MustParse("5Gi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("1G"),
									corev1.ResourceMemory: resource.MustParse("10Gi"),
								},
							},
						},
					},
					ClientConfiguration: ksapi.NighthawkClientConfigurationSpec{
						CmdLineArgs: []string{"--testing", "--things"},
						PodConfigurationSpec: ksapi.PodConfigurationSpec{
							Annotations: map[string]string{
								"anno_one": "value_two",
							},
						},
					},
				},
			}

			configMap := corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{Name: "cm"},
			}

			jobSpec = NewClientJob(&cr)
			serverDeployment = NewServerDeployment(&cr, &configMap)
			clientService = NewServerService(&cr)

		})

		Context("job configuration", func() {
			It("has client annotations", func() {
				Expect(jobSpec.Annotations).To(HaveKey("anno_one"))
			})

			It("doesnt have server annotations", func() {
				Expect(jobSpec.Annotations).ToNot(HaveKey("anno_two"))
			})

			Context("template", func() {
				It("has client annotations", func() {
					Expect(jobSpec.Spec.Template.Annotations).To(HaveKey("anno_one"))
				})

				It("doesnt have server annotations", func() {
					Expect(jobSpec.Spec.Template.Annotations).ToNot(HaveKey("anno_two"))
				})
			})
		})

		Context("client configuration", func() {
			It("has server annotations", func() {
				Expect(clientService.Annotations).To(HaveKey("anno_two"))
			})

			It("doesnt have client annotations", func() {
				Expect(clientService.Annotations).ToNot(HaveKey("anno_one"))
			})
		})

		Context("server configuration", func() {
			Context("spec", func() {
				It("has server annotations", func() {
					Expect(serverDeployment.Annotations).To(HaveKey("anno_two"))
				})

				It("doesnt have client annotations", func() {
					Expect(serverDeployment.Annotations).ToNot(HaveKey("anno_one"))
				})
			})
			Context("template", func() {
				It("has server annotations", func() {
					Expect(serverDeployment.Spec.Template.ObjectMeta.Annotations).To(HaveKey("anno_two"))
				})

				It("doesnt have client annotations", func() {
					Expect(serverDeployment.Spec.Template.ObjectMeta.Annotations).ToNot(HaveKey("anno_one"))
				})
			})
		})
	})
})
