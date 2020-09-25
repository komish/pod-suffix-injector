/*


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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var (
	optInLabel  = "inject-pod-suffix"
	suffixLabel = "pod-suffix"
)

// PodReconciler reconciles a Pod object
type PodReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods/status,verbs=get;update;patch

func (r *PodReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("pod", req.NamespacedName)

	// your logic here

	p := corev1.Pod{}
	err := r.Get(context.TODO(), req.NamespacedName, &p)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// object was deleted before we could do anything with it
			return ctrl.Result{}, nil
		}

		// we got some other error
		return ctrl.Result{}, err
	}

	// get labels from the pod
	labels := p.GetLabels()
	// check that the opt-in label exists
	optedIn, keyFound := labels[optInLabel]

	if !keyFound || optedIn != "true" {
		// the object didn't have the opt-in key or it was not set to true
		r.Log.Info("pod has opted out. doing nothing")
		return ctrl.Result{}, nil
	}
	r.Log.Info("pod has opted in!")

	// the pod has opted in, add the label to the pod.
	baseName := p.GetGenerateName()
	if baseName == "" {
		// the pod isn't using generateName so we can't reliably determine
		// the suffix
		r.Log.Info("unable to determine suffix because the pod is not controlled by an object using generateName")
		return ctrl.Result{}, err
	}

	r.Log.Info("updating object if needed")
	suffix := getSuffix([]byte(baseName), []byte(p.GetName()))

	// with suffix, now add the label
	new := p.DeepCopy()
	new.ObjectMeta.Labels[suffixLabel] = string(suffix)
	// ideally this should be a patch?
	err = r.Update(context.TODO(), new, &client.UpdateOptions{})
	if err != nil {
		r.Log.Info("some error updating the pod")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *PodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// we don't care about deletes
	predicates := predicate.Funcs{
		DeleteFunc: func(e event.DeleteEvent) bool { return false },
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		WithEventFilter(predicates).
		Complete(r)
}

func getSuffix(base, suffixed []byte) []byte {
	return suffixed[len(base):]
}
