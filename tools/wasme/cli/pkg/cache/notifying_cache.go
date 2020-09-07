package cache

import (
	"crypto/md5"
	"fmt"

	"github.com/hashicorp/go-multierror"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// sends Events to kubernetes when images are added to the cache
type Notifier struct {
	kube           kubernetes.Interface
	wasmeNamespace string
	cacheName      string
}

func NewNotifier(kube kubernetes.Interface, wasmeNamespace string, cacheName string) *Notifier {
	return &Notifier{kube: kube, wasmeNamespace: wasmeNamespace, cacheName: cacheName}
}

const (
	// marked as "true" always, for searching
	CacheGlobalLabel = "cache.wasme.io/cache_event"
	// ref to the image
	CacheImageRefLabel = "cache.wasme.io/image_ref"
	Reason_ImageAdded  = "ImageAdded"
	Reason_ImageError  = "ImageError"
)

func (n *Notifier) Notify(err error, image string) error {
	var reason, message string
	if err != nil {
		reason = Reason_ImageError
		message = err.Error()
	} else {
		reason = Reason_ImageAdded
		message = fmt.Sprintf("Image %v added successfully", image)
	}
	_, eventCreateErr := n.kube.CoreV1().Events(n.wasmeNamespace).Create(&v1.Event{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "wasme-cache-event-",
			Namespace:    n.wasmeNamespace,
			Labels:       EventLabels(image),
			Annotations:  EventAnnotations(image),
		},
		InvolvedObject: v1.ObjectReference{
			Kind:       "ConfigMap",
			Namespace:  n.wasmeNamespace,
			Name:       n.cacheName,
			APIVersion: "v1",
		},
		Reason:  reason,
		Message: message,
		Source: v1.EventSource{
			Component: "wasme-cache",
		},
	})

	if eventCreateErr != nil {
		return multierror.Append(err, eventCreateErr)
	}

	return err
}

func EventLabels(image string) map[string]string {
	// take hash for valid label name
	refLabel := fmt.Sprintf("%x", md5.Sum([]byte(image)))
	return map[string]string{
		CacheGlobalLabel:   "true",
		CacheImageRefLabel: refLabel,
	}
}

func EventAnnotations(image string) map[string]string {
	return map[string]string{
		CacheImageRefLabel: image,
	}
}
