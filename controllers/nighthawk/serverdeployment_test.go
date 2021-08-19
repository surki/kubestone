package nighthawk

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	ksapi "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("Server Deployment", func() {
	Describe("created from CR", func() {
		var cr ksapi.Nighthawk
		var deployment *appsv1.Deployment

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
						ConfigsVolume: map[string]string{
							"the-config.yml":    "config content",
							"included-file.yml": "included content",
						},
						ConfigFile: "the-config.yml",
						Port:       1234,
						PodConfigurationSpec: ksapi.PodConfigurationSpec{
							PodLabels: map[string]string{"labels": "are", "really": "useful"},
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
				},
			}

			configMap := corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{Name: "cm"},
			}

			deployment = NewServerDeployment(&cr, &configMap)
		})

		Context("with Image details specified", func() {
			It("should match on Image.Name", func() {
				Expect(deployment.Spec.Template.Spec.Containers[0].Image).To(
					Equal(cr.Spec.Image.Name))
			})
			It("should match on Image.PullPolicy", func() {
				Expect(deployment.Spec.Template.Spec.Containers[0].ImagePullPolicy).To(
					Equal(corev1.PullPolicy(cr.Spec.Image.PullPolicy)))
			})
			It("should match on Image.PullSecret", func() {
				Expect(deployment.Spec.Template.Spec.ImagePullSecrets[0].Name).To(
					Equal(cr.Spec.Image.PullSecret))
			})
		})

		Context("when existent configfile is referred", func() {
			It("CR Validation should succeed", func() {
				valid, err := IsCrValid(&cr)
				Expect(valid).To(BeTrue())
				Expect(err).To(BeNil())
			})
		})

		Context("when non-existent configFile is referred", func() {
			invalidCr := ksapi.Nighthawk{
				Spec: ksapi.NighthawkSpec{
					Image: ksapi.ImageSpec{
						Name: "foo/foo:test",
					},
					ServerConfiguration: ksapi.NighthawkServerConfigurationSpec{
						ConfigsVolume: map[string]string{
							"the-config.yml": "config content",
						},
						ConfigFile: "non-existent-config.yml",
					},
				},
			}

			It("CR Validation should fail", func() {
				valid, err := IsCrValid(&invalidCr)
				Expect(valid).To(BeFalse())
				Expect(err).NotTo(BeNil())
			})
		})

		Context("with cmdLineArgs specified", func() {
			It("--testing mode is set", func() {
				Expect(deployment.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("--testing"))
			})
		})

		Context("with podLabels specified", func() {
			It("should contain all podLabels", func() {
				for k, v := range cr.Spec.ServerConfiguration.PodLabels {
					Expect(deployment.Spec.Template.ObjectMeta.Labels).To(
						HaveKeyWithValue(k, v))
				}
			})
		})

		Context("with podAffinity specified", func() {
			It("should match with Affinity", func() {
				Expect(deployment.Spec.Template.Spec.Affinity).To(
					Equal(cr.Spec.ServerConfiguration.PodScheduling.Affinity))
			})
			It("should match with Tolerations", func() {
				Expect(deployment.Spec.Template.Spec.Tolerations).To(
					Equal(cr.Spec.ServerConfiguration.PodScheduling.Tolerations))
			})
			It("should match with NodeSelector", func() {
				Expect(deployment.Spec.Template.Spec.NodeSelector).To(
					Equal(cr.Spec.ServerConfiguration.PodScheduling.NodeSelector))
			})
			It("should match with NodeName", func() {
				Expect(deployment.Spec.Template.Spec.NodeName).To(
					Equal(cr.Spec.ServerConfiguration.PodScheduling.NodeName))
			})
		})

		Context("with readiness probe", func() {
			It("should have protocol port", func() {
				Expect(deployment.Spec.Template.Spec.Containers[0].ReadinessProbe.Handler).To(
					Equal(corev1.Handler{
						TCPSocket: &corev1.TCPSocketAction{
							Port: intstr.IntOrString{Type: intstr.Int, IntVal: 1234},
						}}))
			})
		})

		Context("with resources specified", func() {
			It("should request the given CPU", func() {
				Expect(deployment.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu()).To(
					BeEquivalentTo(cr.Spec.ServerConfiguration.Resources.Requests.Cpu()))
			})
			It("should request the given memory", func() {
				Expect(deployment.Spec.Template.Spec.Containers[0].Resources.Requests.Memory()).To(
					BeEquivalentTo(cr.Spec.ServerConfiguration.Resources.Requests.Memory()))
			})
			It("should limit to the given CPU", func() {
				Expect(deployment.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu()).To(
					BeEquivalentTo(cr.Spec.ServerConfiguration.Resources.Limits.Cpu()))
			})
			It("should limit to the given memory", func() {
				Expect(deployment.Spec.Template.Spec.Containers[0].Resources.Limits.Memory()).To(
					BeEquivalentTo(cr.Spec.ServerConfiguration.Resources.Limits.Memory()))
			})
		})
	})
})
