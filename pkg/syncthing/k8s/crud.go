package k8s

import (
	"fmt"
	"strings"

	"github.com/okteto/okteto/pkg/log"
	"github.com/okteto/okteto/pkg/model"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var devTerminationGracePeriodSeconds int64

//Deploy creates or updates a syncthing stateful set
func Deploy(dev *model.Dev, c *kubernetes.Clientset) error {
	ss := translate(dev)

	if exists(ss, c) {
		if err := update(ss, c); err != nil {
			return err
		}
	} else {
		if err := create(ss, c); err != nil {
			return err
		}
	}
	return nil
}

func exists(ss *appsv1.StatefulSet, c *kubernetes.Clientset) bool {
	ss, err := c.AppsV1().StatefulSets(ss.Namespace).Get(ss.Name, metav1.GetOptions{})
	if err != nil {
		return false
	}
	return ss.Name != ""
}

func create(ss *appsv1.StatefulSet, c *kubernetes.Clientset) error {
	log.Infof("creating syncthing statefulset '%s", ss.Name)
	_, err := c.AppsV1().StatefulSets(ss.Namespace).Create(ss)
	if err != nil {
		return fmt.Errorf("error creating kubernetes syncthing statefulset: %s", err)
	}
	log.Infof("syncthing statefulset '%s' created", ss.Name)
	return nil
}

func update(ss *appsv1.StatefulSet, c *kubernetes.Clientset) error {
	log.Infof("updating syncthing statefulset '%s'", ss.Name)
	if _, err := c.AppsV1().StatefulSets(ss.Namespace).Update(ss); err != nil {
		return fmt.Errorf("error updating kubernetes syncthing statefulset: %s", err)
	}
	log.Infof("syncthing statefulset '%s' updated", ss.Name)
	return nil
}

// Destroy destroys a database
func Destroy(dev *model.Dev, c *kubernetes.Clientset) error {
	log.Infof("destroying syncthing statefulset '%s' ...", dev.Name)
	sfsClient := c.AppsV1().StatefulSets(dev.Namespace)
	if err := sfsClient.Delete(dev.GetStatefulSetName(), &metav1.DeleteOptions{GracePeriodSeconds: &devTerminationGracePeriodSeconds}); err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("couldn't destroy syncthing statefulset: %s", err)
		}
	}
	log.Infof("syncthing statefulset '%s' destroyed", dev.Name)
	return nil
}
