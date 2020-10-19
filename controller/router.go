package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	schedulerapi "k8s.io/kube-scheduler/extender/v1"
	"log"
	"net/http"
)

type toworker struct {
	name string
	nums int
}

var changeNode = make(map[string]toworker, 0)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome to sample-scheduler-extender!\n")
}

func Prioritize(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var buf bytes.Buffer
	body := io.TeeReader(r.Body, &buf)

	var extenderArgs schedulerapi.ExtenderArgs
	var hostPriorityList *schedulerapi.HostPriorityList
	if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
		hostPriorityList = &schedulerapi.HostPriorityList{}
	} else {
		hostPriorityList = prioritize(extenderArgs)
	}

	if response, err := json.Marshal(hostPriorityList); err != nil {
		log.Fatalln(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func prioritize(args schedulerapi.ExtenderArgs) *schedulerapi.HostPriorityList {
	pod := args.Pod

	//_,ok:=pod.ObjectMeta.Labels["just-types"]

	cpusNeed := pod.Spec.Containers[0].Resources.Requests["cpu"]
	memNeed := pod.Spec.Containers[0].Resources.Requests["memory"]
	fmt.Println("Pod name:", pod.Name, "Resource need : cpu:", cpusNeed, "mem:", memNeed)
	nodes := args.Nodes.Items

	hostPriorityList := make(schedulerapi.HostPriorityList, len(nodes))
	for i, node := range nodes {
		var score int64 = 0

		/*
			if ok && node.ObjectMeta.Name=="worker01"{
				score = 1000000
			}
		*/
		//log.Printf("pod %v/%v is lucky to get score %v\\n", pod.Name, pod.Namespace, score)

		hostPriorityList[i] = schedulerapi.HostPriority{
			Host:  node.Name,
			Score: score,
		}
	}

	return &hostPriorityList
}
