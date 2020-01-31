package neutronovsagent

import (
	"context"

	neutronv1 "github.com/stuggi/neutron-operator/pkg/apis/neutron/v1"
        appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_neutronovsagent")

const (
        NEUTRON_CONFIGMAP_NAME  string = "neutron-config"
)

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileNeutronOvsAgent{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("neutronovsagent-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource NeutronOvsAgent
	err = c.Watch(&source.Kind{Type: &neutronv1.NeutronOvsAgent{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner NeutronOvsAgent
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &neutronv1.NeutronOvsAgent{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileNeutronOvsAgent implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileNeutronOvsAgent{}

// ReconcileNeutronOvsAgent reconciles a NeutronOvsAgent object
type ReconcileNeutronOvsAgent struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a NeutronOvsAgent object and makes changes based on the state read
// and what is in the NeutronOvsAgent.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileNeutronOvsAgent) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling NeutronOvsAgent")

	// Fetch the NeutronOvsAgent instance
	instance := &neutronv1.NeutronOvsAgent{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

        // Define a new Daemonset object
        ds := newDaemonset(instance)

	// Set NeutronOvsAgent instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, ds, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Daemonset already exists
	found := &appsv1.DaemonSet{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: ds.Name, Namespace: ds.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Daemonset", "ds.Namespace", ds.Namespace, "ds.Name", ds.Name)
		err = r.client.Create(context.TODO(), ds)
		if err != nil {
			return reconcile.Result{}, err
		}

		// ds created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Daemonset already exists - don't requeue
	reqLogger.Info("Skip reconcile: Daemonset already exists", "ds.Namespace", found.Namespace, "ds.Name", found.Name)
	return reconcile.Result{}, nil
}

func newDaemonset(cr *neutronv1.NeutronOvsAgent) *appsv1.DaemonSet {
        var bidirectional corev1.MountPropagationMode = corev1.MountPropagationBidirectional
        var hostToContainer corev1.MountPropagationMode = corev1.MountPropagationHostToContainer
        var trueVar bool = true
        var configVolumeDefaultMode int32 = 0644
        var dirOrCreate corev1.HostPathType = corev1.HostPathDirectoryOrCreate

        daemonSet := appsv1.DaemonSet{
                TypeMeta: metav1.TypeMeta{
                        Kind:       "DaemonSet",
                        APIVersion: "apps/v1",
                },
                ObjectMeta: metav1.ObjectMeta{
                        Name:      cr.Name + "-daemonset",
                        //Name:      fmt.Sprintf("%s-nova-%s",cr.Name, cr.Spec.NodeName),
                        Namespace: cr.Namespace,
                        //OwnerReferences: []metav1.OwnerReference{
                        //      *metav1.NewControllerRef(cr, schema.GroupVersionKind{
                        //              Group:   v1beta1.SchemeGroupVersion.Group,
                        //              Version: v1beta1.SchemeGroupVersion.Version,
                        //              Kind:    "GenericDaemon",
                        //      }),
                        //},
                },
                Spec: appsv1.DaemonSetSpec{
                        Selector: &metav1.LabelSelector{
                                // MatchLabels: map[string]string{"daemonset": cr.Spec.NodeName + cr.Name + "-daemonset"},
                                MatchLabels: map[string]string{"daemonset": cr.Name + "-daemonset"},
                        },
                        Template: corev1.PodTemplateSpec{
                                ObjectMeta: metav1.ObjectMeta{
                                        // Labels: map[string]string{"daemonset": cr.Spec.NodeName + cr.Name + "-daemonset"},
                                        Labels: map[string]string{"daemonset": cr.Name + "-daemonset"},
                                },
                                Spec: corev1.PodSpec{
                                        NodeSelector:   map[string]string{"daemon": cr.Spec.Label},
                                        HostNetwork:    true,
                                        HostPID:        true,
                                        HostAliases:    []corev1.HostAlias{},
                                        InitContainers: []corev1.Container{},
                                        Containers:     []corev1.Container{},
                                },
                        },
                },
        }

        opsHostAliases := []corev1.HostAlias{
                {
                        IP: "172.17.1.83",
                        Hostnames: []string{ "controller-0.internalapi.redhat.local", "controller-0.internalapi"},
                },
                {
                        IP: "172.17.2.16",
                        Hostnames: []string{ "controller-0.tenant.redhat.local", "controller-0.tenant"},
                },
                {
                        IP: "172.17.1.146",
                        Hostnames: []string{ "compute-0.internalapi.redhat.local", "compute-0.internalapi"},
                },
                {
                        IP: "172.17.2.55",
                        Hostnames: []string{ "compute-0.tenant.redhat.local", "compute-0.tenant"},
                },
                {
                        IP: "172.17.1.84",
                        Hostnames: []string{ "compute-1.internalapi.redhat.local", "compute-1.internalapi"},
                },
                {
                        IP: "172.17.2.21",
                        Hostnames: []string{ "compute-1.tenant.redhat.local", "compute-1.tenant"},
                },
                {
                        IP: "172.17.1.29",
                        Hostnames: []string{"overcloud.internalapi.localdomain"},
                },
        }

        for _, opsHostAlias := range opsHostAliases {
                daemonSet.Spec.Template.Spec.HostAliases = append(daemonSet.Spec.Template.Spec.HostAliases, opsHostAlias)
        }

        initContainerSpec := corev1.Container{
                Name:  "ovs-agent-config-init",
                Image: cr.Spec.OpenvswitchImage,
                SecurityContext: &corev1.SecurityContext{
                        Privileged:  &trueVar,
                },
                Command: []string{
                        "/bin/bash", "-c", "export POD_IP_TENANT=$(ip route get 172.17.2.16 | awk '{print $5}') && cp /etc/neutron/plugins/ml2/openvswitch_agent.ini /mnt/openvswitch_agent.ini && crudini --set /mnt/openvswitch_agent.ini ovs local_ip $POD_IP_TENANT",
                },
                Env: []corev1.EnvVar{
                        {
                                Name: "MY_POD_IP",
                                ValueFrom: &corev1.EnvVarSource{
                                        FieldRef: &corev1.ObjectFieldSelector{
                                                FieldPath: "status.podIP",
                                        },
                                },
                        },
                },
                VolumeMounts: []corev1.VolumeMount{
                        {
                                Name:      "neutron-config",
                                ReadOnly:  true,
                                MountPath: "/etc/neutron/neutron.conf",
                                SubPath:   "neutron.conf",
                        },
                        {
                                Name:      "neutron-config",
                                ReadOnly:  true,
                                MountPath: "/etc/neutron/plugins/ml2/openvswitch_agent.ini",
                                SubPath:   "openvswitch_agent.ini",
                        },
                        {
                                Name:      "rendered-config-vol",
                                MountPath: "/mnt",
                                ReadOnly:  false,
                        },
                },
        }
        daemonSet.Spec.Template.Spec.InitContainers = append(daemonSet.Spec.Template.Spec.InitContainers, initContainerSpec)


        neutronOvsAgentContainerSpec := corev1.Container{
                Name:  "neutron-ovs-agent",
                Image: cr.Spec.OpenvswitchImage,
                //ReadinessProbe: &corev1.Probe{
                //        Handler: corev1.Handler{
                //                Exec: &corev1.ExecAction{
                //                        Command: []string{
                //                                "/openstack/healthcheck",
                //                        },
                //                },
                //        },
                //        InitialDelaySeconds: 30,
                //        PeriodSeconds:       30,
                //        TimeoutSeconds:      1,
                //},
                Command: []string{
                        "/usr/bin/neutron-openvswitch-agent", "--config-file", "/usr/share/neutron/neutron-dist.conf", "--config-file", "/etc/neutron/neutron.conf", "--config-file", "/mnt/openvswitch_agent.ini", "--config-dir", "/etc/neutron/conf.d/common", "--log-file=/var/log/neutron/openvswitch-agent.log",
                },
                SecurityContext: &corev1.SecurityContext{
                        Privileged:  &trueVar,
                },
                VolumeMounts: []corev1.VolumeMount{
                        {
                                Name:      "neutron-config",
                                ReadOnly:  true,
                                MountPath: "/etc/neutron/neutron.conf",
                                SubPath:   "neutron.conf",
                        },
                        {
                                Name:      "neutron-config",
                                ReadOnly:  true,
                                MountPath: "/etc/neutron/plugins/ml2/openvswitch_agent.ini",
                                SubPath:   "openvswitch_agent.ini",
                        },
                        {
                                Name:      "lib-modules-volume",
                                MountPath: "/lib/modules",
                                MountPropagation: &hostToContainer,
                        },
                        {
                                Name:      "run-openvswitch-volume",
                                MountPath: "/var/run/openvswitch",
                                MountPropagation: &bidirectional,
                        },
                        {
                                Name:      "neutron-log-volume",
                                MountPath: "/var/log/neutron",
                                MountPropagation: &bidirectional,
                        },
                        {
                                Name:      "rendered-config-vol",
                                MountPath: "/mnt",
                                ReadOnly:  true,
                        },
                },
        }
        daemonSet.Spec.Template.Spec.Containers = append(daemonSet.Spec.Template.Spec.Containers, neutronOvsAgentContainerSpec)

        volConfigs := []corev1.Volume{
                {
                        Name: "hostroot",
                        VolumeSource: corev1.VolumeSource{
                                HostPath: &corev1.HostPathVolumeSource{
                                        Path: "/",
                                },
                        },
                },
                {
                        Name: "boot-volume",
                        VolumeSource: corev1.VolumeSource{
                                HostPath: &corev1.HostPathVolumeSource{
                                        Path: "/boot",
                                },
                        },
                },
                {
                        Name: "run-volume",
                        VolumeSource: corev1.VolumeSource{
                                HostPath: &corev1.HostPathVolumeSource{
                                        Path: "/run",
                                },
                        },
                },
                {
                        Name: "lib-modules-volume",
                        VolumeSource: corev1.VolumeSource{
                                HostPath: &corev1.HostPathVolumeSource{
                                        Path: "/lib/modules",
                                },
                        },
                },
                {
                        Name: "dev-volume",
                        VolumeSource: corev1.VolumeSource{
                                HostPath: &corev1.HostPathVolumeSource{
                                        Path: "/dev",
                                },
                        },
                },
                {
                        Name: "sys-fs-cgroup-volume",
                        VolumeSource: corev1.VolumeSource{
                                HostPath: &corev1.HostPathVolumeSource{
                                        Path: "/sys/fs/cgroup",
                                },
                        },
                },
                {
                        Name: "run-openvswitch-volume",
                        VolumeSource: corev1.VolumeSource{
                                HostPath: &corev1.HostPathVolumeSource{
                                        Path: "/var/run/openvswitch",
                                        Type: &dirOrCreate,
                                },
                        },
                },
                {
                        Name: "neutron-log-volume",
                        VolumeSource: corev1.VolumeSource{
                                HostPath: &corev1.HostPathVolumeSource{
                                        Path: "/var/log/containers/neutron",
                                        Type: &dirOrCreate,
                                },
                        },
                },
                {
                        Name: "neutron-config",
                        VolumeSource: corev1.VolumeSource{
                                ConfigMap: &corev1.ConfigMapVolumeSource{
                                         DefaultMode: &configVolumeDefaultMode,
                                         LocalObjectReference: corev1.LocalObjectReference{
                                                 Name: NEUTRON_CONFIGMAP_NAME,
                                         },
                                },
                        },
                },
                {
                        Name: "rendered-config-vol",
                        VolumeSource: corev1.VolumeSource{
                                EmptyDir: &corev1.EmptyDirVolumeSource{},
                        },
                },
        }
        for _, volConfig := range volConfigs {
                daemonSet.Spec.Template.Spec.Volumes = append(daemonSet.Spec.Template.Spec.Volumes, volConfig)
        }

        return &daemonSet
}
