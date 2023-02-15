/*
Copyright ApeCloud, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package components

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dbaasv1alpha1 "github.com/apecloud/kubeblocks/apis/dbaas/v1alpha1"
	"github.com/apecloud/kubeblocks/controllers/dbaas/components/stateless"
	intctrlutil "github.com/apecloud/kubeblocks/internal/controllerutil"
	testdbaas "github.com/apecloud/kubeblocks/internal/testutil/dbaas"
	testk8s "github.com/apecloud/kubeblocks/internal/testutil/k8s"
)

var _ = Describe("Deployment Controller", func() {
	var (
		randomStr          = testCtx.GetRandomStr()
		clusterDefName     = "stateless-definition1-" + randomStr
		clusterVersionName = "stateless-cluster-version1-" + randomStr
		clusterName        = "stateless1-" + randomStr
	)

	const (
		namespace         = "default"
		statelessCompName = "stateless"
		statelessCompType = "stateless"
	)

	cleanAll := func() {
		// must wait until resources deleted and no longer exist before the testcases start,
		// otherwise if later it needs to create some new resource objects with the same name,
		// in race conditions, it will find the existence of old objects, resulting failure to
		// create the new objects.
		By("clean resources")

		// delete cluster(and all dependent sub-resources), clusterversion and clusterdef
		testdbaas.ClearClusterResources(&testCtx)

		// clear rest resources
		inNS := client.InNamespace(testCtx.DefaultNamespace)
		ml := client.HasLabels{testCtx.TestObjLabelKey}
		// namespaced resources
		testdbaas.ClearResources(&testCtx, intctrlutil.DeploymentSignature, inNS, ml)
		testdbaas.ClearResources(&testCtx, intctrlutil.PodSignature, inNS, ml, client.GracePeriodSeconds(0))
	}

	BeforeEach(cleanAll)

	AfterEach(cleanAll)

	Context("test controller", func() {
		It("", func() {
			testdbaas.NewClusterDefFactory(clusterDefName, testdbaas.MySQLType).
				AddComponent(testdbaas.StatelessNginxComponent, statelessCompType).SetDefaultReplicas(2).
				Create(&testCtx).GetObject()

			cluster := testdbaas.NewClusterFactory(testCtx.DefaultNamespace, clusterName, clusterDefName, clusterVersionName).
				AddComponent(statelessCompName, statelessCompType).Create(&testCtx).GetObject()

			By("patch cluster to Running")
			Expect(testdbaas.ChangeObjStatus(&testCtx, cluster, func() {
				cluster.Status.Phase = dbaasv1alpha1.RunningPhase
			}))

			By("create the deployment of the stateless component")
			deploy := testdbaas.MockStatelessComponentDeploy(testCtx, clusterName, statelessCompName)
			newDeploymentKey := client.ObjectKey{Name: deploy.Name, Namespace: namespace}
			Eventually(testdbaas.CheckObj(&testCtx, newDeploymentKey, func(g Gomega, deploy *appsv1.Deployment) {
				g.Expect(deploy.Generation == 1).Should(BeTrue())
			})).Should(Succeed())

			By("check stateless component phase is Failed")
			Eventually(testdbaas.GetClusterComponentPhase(testCtx, clusterName, statelessCompName)).Should(Equal(dbaasv1alpha1.FailedPhase))

			By("test when a pod of deployment is failed")
			podName := fmt.Sprintf("%s-%s-%s", clusterName, statelessCompName, testCtx.GetRandomStr())
			pod := testdbaas.MockStatelessPod(testCtx, deploy, clusterName, statelessCompName, podName)
			// mock pod container is failed
			errMessage := "Back-off pulling image nginx:latest"
			Expect(testdbaas.ChangeObjStatus(&testCtx, pod, func() {
				pod.Status.ContainerStatuses = []corev1.ContainerStatus{
					{
						State: corev1.ContainerState{
							Waiting: &corev1.ContainerStateWaiting{
								Reason:  "ImagePullBackOff",
								Message: errMessage,
							},
						},
					},
				}
			})).Should(Succeed())
			Eventually(testdbaas.CheckObj(&testCtx, client.ObjectKeyFromObject(pod), func(g Gomega, tmpPod *corev1.Pod) {
				g.Expect(len(tmpPod.Status.ContainerStatuses) == 1).Should(BeTrue())
			})).Should(Succeed())

			// mock failed container timed out
			Expect(testdbaas.ChangeObjStatus(&testCtx, pod, func() {
				pod.Status.Conditions = []corev1.PodCondition{
					{
						Type:               corev1.ContainersReady,
						Status:             corev1.ConditionFalse,
						LastTransitionTime: metav1.NewTime(time.Now().Add(-2 * time.Minute)),
					},
				}
			})).Should(Succeed())
			Eventually(testdbaas.CheckObj(&testCtx, client.ObjectKeyFromObject(pod), func(g Gomega, tmpPod *corev1.Pod) {
				g.Expect(len(tmpPod.Status.Conditions) == 1).Should(BeTrue())
			})).Should(Succeed())

			// wait for component.message contains pod message.
			Eventually(testdbaas.CheckObj(&testCtx, client.ObjectKeyFromObject(cluster), func(g Gomega, tmpCluster *dbaasv1alpha1.Cluster) {
				statusComponent := tmpCluster.Status.Components[statelessCompName]
				g.Expect(statusComponent.Message.GetObjectMessage("Pod", pod.Name)).Should(Equal(errMessage))
			})).Should(Succeed())

			By("mock deployment is ready")
			newDeployment := &appsv1.Deployment{}
			Expect(k8sClient.Get(context.Background(), newDeploymentKey, newDeployment)).Should(Succeed())
			Expect(testdbaas.ChangeObjStatus(&testCtx, newDeployment, func() {
				testk8s.MockDeploymentReady(newDeployment, stateless.NewRSAvailableReason)
			})).Should(Succeed())

			By("test deployment status is ready")
			Eventually(testdbaas.CheckObj(&testCtx, newDeploymentKey, func(g Gomega, deploy *appsv1.Deployment) {
				g.Expect(deploy.Status.AvailableReplicas == newDeployment.Status.AvailableReplicas &&
					deploy.Status.ReadyReplicas == newDeployment.Status.ReadyReplicas &&
					deploy.Status.Replicas == newDeployment.Status.Replicas).Should(BeTrue())
			})).Should(Succeed())

			By("waiting the component is Running")
			Eventually(testdbaas.GetClusterComponentPhase(testCtx, clusterName, statelessCompName)).Should(Equal(dbaasv1alpha1.RunningPhase))
		})
	})
})
