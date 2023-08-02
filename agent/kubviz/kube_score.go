package main

import (
	"context"
	"encoding/json"
	"log"
	exec "os/exec"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/intelops/kubviz/constants"
	"github.com/intelops/kubviz/model"
	"github.com/nats-io/nats.go"
	"github.com/zegl/kube-score/renderer/json_v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func RunKubeScore(clientset *kubernetes.Clientset, js nats.JetStreamContext, wg *sync.WaitGroup, errCh chan error) {
	defer wg.Done()
	nsList, err := clientset.CoreV1().
		Namespaces().
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Println("Error occurred while getting client set for kube-score: ", err)
		return
	}

	log.Printf("Namespace size: %d", len(nsList.Items))
	for _, ns := range nsList.Items {
		var report json_v2.ScoredObject
		out, err := executeCommand("kubectl api-resources --verbs=list --namespaced -o name | xargs -n1 -I{} sh -c \"kubectl get {} -n " + ns.Name + " -oyaml && echo ---\" | kube-score score - -o json ")
		if err != nil {
			log.Printf("Error scanning image %s: %v", ns, err)
			continue // Move on to the next image in case of an error
		}

		parts := strings.SplitN(out, "{", 2)
		if len(parts) <= 1 {
			log.Println("No output from command", err)
			continue // Move on to the next image if there's no output
		}

		log.Println("Command logs", parts[0])
		jsonPart := "{" + parts[1]
		log.Println("First 200 lines output", jsonPart[:200])
		log.Println("Last 200 lines output", jsonPart[len(jsonPart)-200:])

		err = json.Unmarshal([]byte(jsonPart), &report)
		if err != nil {
			log.Printf("Error occurred while Unmarshalling json: %v", err)
			continue // Move on to the next image in case of an error
		}
		publishKubescoreMetrics(report, js, errCh)
		// If you want to publish the report or perform any other action with it, you can do it here

	}
}

func publishKubescoreMetrics(report json_v2.ScoredObject, js nats.JetStreamContext, errCh chan error) {
	metrics := model.KubeScoreRecommendations{
		ID:          uuid.New().String(),
		ClusterName: ClusterName,
		Report:      report,
	}
	metricsJson, _ := json.Marshal(metrics)
	_, err := js.Publish(constants.KUBESCORE_SUBJECT, metricsJson)
	if err != nil {
		errCh <- err
		return
	}
	log.Printf("KubeScore report with ID:%s has been published\n", metrics.ID)
	errCh <- nil
}

func executeCommand(command string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", command)
	stdout, err := cmd.Output()

	if err != nil {
		log.Println("Execute Command Error", err.Error())
		return "", err
	}

	// Print the output
	log.Println(string(stdout))
	return string(stdout), nil
}
