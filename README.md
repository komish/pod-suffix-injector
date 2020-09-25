# PodSuffixInjector

**PodSuffixInjector** is a controller that injects the generated suffix appended to podnames when `generateName` is used to determine the name of the Pod. The suffix will be stored in a label `pod-suffix` on the pod itself.

Controller logic is functional, in-cluster assets (RBAC, etc) have not been tested.

To use, run the controller and then create pods (with an controller reference) with the label `inject-pod-suffix: "true"`.

Ex:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
        inject-pod-suffix: "true"
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
```

Any pod without this label is considered to have opted out. Any pod with this label having a different value is considered to have opted out. 

The controller makes no other changes to the pod.

Log output for an object that has opted out looks like the following:

```
2020-09-25T10:54:35.003-0500    INFO    controllers.PodSuffixInjector   pod has opted out. doing nothing
2020-09-25T10:54:35.003-0500    DEBUG   controller-runtime.controller   Successfully Reconciled {"controller": "pod", "request": "default/nginx-deployment-6b474476c4-7bj6n"}
2020-09-25T10:54:35.051-0500    INFO    controllers.PodSuffixInjector   pod has opted out. doing nothing
2020-09-25T10:54:35.051-0500    DEBUG   controller-runtime.controller   Successfully Reconciled {"controller": "pod", "request": "default/nginx-deployment-6b474476c4-6gpkb"}
2020-09-25T10:54:35.924-0500    INFO    controllers.PodSuffixInjector   pod has opted out. doing nothing
2020-09-25T10:54:35.924-0500    DEBUG   controller-runtime.controller   Successfully Reconciled {"controller": "pod", "request": "default/nginx-deployment-6b474476c4-6gpkb"}
2020-09-25T10:54:36.318-0500    INFO    controllers.PodSuffixInjector   pod has opted out. doing nothing
2020-09-25T10:54:36.318-0500    DEBUG   controller-runtime.controller   Successfully Reconciled {"controller": "pod", "request": "default/nginx-deployment-6b474476c4-7bj6n"}
```

Log output for an object that has opted in looks like the following (errors here indicate that the resource has changed before the controller was able to work with it so it requeues the work):

```
2020-09-25T10:55:26.519-0500    INFO    controllers.PodSuffixInjector   pod has opted in!
2020-09-25T10:55:26.519-0500    INFO    controllers.PodSuffixInjector   updating object if needed
2020-09-25T10:55:26.577-0500    INFO    controllers.PodSuffixInjector   some error updating the pod
2020-09-25T10:55:26.577-0500    ERROR   controller-runtime.controller   Reconciler error        {"controller": "pod", "request": "default/nginx-deployment-569d6767d5-c82vh", "error": "Operation cannot be fulfilled on pods \"nginx-deployment-569d6767d5-c82vh\": the object has been modified; please apply your changes to the latest version and try again"}
github.com/go-logr/zapr.(*zapLogger).Error
        /Users/dev/.go/pkg/mod/github.com/go-logr/zapr@v0.1.0/zapr.go:128
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).reconcileHandler
        /Users/dev/.go/pkg/mod/sigs.k8s.io/controller-runtime@v0.5.0/pkg/internal/controller/controller.go:258
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).processNextWorkItem
        /Users/dev/.go/pkg/mod/sigs.k8s.io/controller-runtime@v0.5.0/pkg/internal/controller/controller.go:232
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).worker
        /Users/dev/.go/pkg/mod/sigs.k8s.io/controller-runtime@v0.5.0/pkg/internal/controller/controller.go:211
k8s.io/apimachinery/pkg/util/wait.JitterUntil.func1
        /Users/dev/.go/pkg/mod/k8s.io/apimachinery@v0.17.2/pkg/util/wait/wait.go:152
k8s.io/apimachinery/pkg/util/wait.JitterUntil
        /Users/dev/.go/pkg/mod/k8s.io/apimachinery@v0.17.2/pkg/util/wait/wait.go:153
k8s.io/apimachinery/pkg/util/wait.Until
        /Users/dev/.go/pkg/mod/k8s.io/apimachinery@v0.17.2/pkg/util/wait/wait.go:88
2020-09-25T10:55:27.578-0500    INFO    controllers.PodSuffixInjector   pod has opted in!
2020-09-25T10:55:27.578-0500    INFO    controllers.PodSuffixInjector   updating object if needed
2020-09-25T10:55:27.636-0500    DEBUG   controller-runtime.controller   Successfully Reconciled {"controller": "pod", "request": "default/nginx-deployment-569d6767d5-c82vh"}
2020-09-25T10:55:27.678-0500    INFO    controllers.PodSuffixInjector   pod has opted in!
2020-09-25T10:55:27.678-0500    INFO    controllers.PodSuffixInjector   updating object if needed
2020-09-25T10:55:27.734-0500    DEBUG   controller-runtime.controller   Successfully Reconciled {"controller": "pod", "request": "default/nginx-deployment-569d6767d5-c82vh"}
2020-09-25T10:55:28.915-0500    INFO    controllers.PodSuffixInjector   pod has opted in!
2020-09-25T10:55:28.915-0500    INFO    controllers.PodSuffixInjector   updating object if needed
2020-09-25T10:55:28.969-0500    DEBUG   controller-runtime.controller   Successfully Reconciled {"controller": "pod", "request": "default/nginx-deployment-569d6767d5-c82vh"}
```