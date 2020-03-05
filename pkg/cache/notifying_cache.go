package cache

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/opencontainers/go-digest"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// sends Events to kubernetes when images are added
type NotifyingCache struct {
	kube kubernetes.Interface
	Cache
	wasmeNamespace string
	cacheName      string
}

func NewNotifyingCache(kube kubernetes.Interface, cache Cache, wasmeNamespace string, cacheName string) *NotifyingCache {
	return &NotifyingCache{kube: kube, Cache: cache, wasmeNamespace: wasmeNamespace, cacheName: cacheName}
}

const (
	// marked as "true" always, for searching
	CacheGlobalLabel = "cache.wasme.io/cache_event"
	// ref to the image
	CacheImageRefLabel = "cache.wasme.io/image_ref"
	Reason_ImageAdded  = "ImageAdded"
	Reason_ImageError  = "ImageError"
)

func (n *NotifyingCache) Add(ctx context.Context, image string) (digest.Digest, error) {
	var reason, message string
	dgst, err := n.Cache.Add(ctx, image)
	if err != nil {
		reason = Reason_ImageError
		message = err.Error()
	} else {
		reason = Reason_ImageAdded
		message = fmt.Sprintf("Image %v with digest %+v added successfully", image, dgst)
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

	return dgst, multierror.Append(err, eventCreateErr)
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
