package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"encoding/json"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/apimachinery/pkg/labels"	
)

// EntityKey defines the structure of an entity's key.
type EntityKey struct {
	Name      string
	Namespace string
	Type      string
}

// Graph represents the structure of a graph.
type Graph struct {
	Nodes         []GraphNode
	Relationships []GraphRelationship
}

// GraphNode represents a node in the graph.
type GraphNode struct {
	Key        EntityKey
	Properties map[string]string
	Revision   int
}

// GraphRelationship represents a relationship between nodes in the graph.
type GraphRelationship struct {
	Source         EntityKey
	Target         EntityKey
	RelationshipType string
	Properties     map[string]string
	Revision       int
}

func listThings() {
	// Find the kubeconfig path (default location)
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")

	// Load kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Error loading kubeconfig: %v", err)
	}

	// Create a clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	// List pods in the default namespace
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error fetching pods: %v", err)
	}

	fmt.Printf("There are %d pods\n", len(pods.Items))

	// Print the pod names
	for _, pod := range pods.Items {
		fmt.Printf(" - %s\n", pod.Name)
	}

    // Get all ReplicaSets in the cluster
    replicasets, err := clientset.AppsV1().ReplicaSets("").List(context.Background(), metav1.ListOptions{})
    if err != nil {
        log.Fatalf("Failed to fetch ReplicaSets: %v", err)
    }

    // Print the names of the ReplicaSets
    fmt.Println("List of ReplicaSets:")
    for _, rs := range replicasets.Items {
		fmt.Printf("Name: %s, Namespace: %s\n", rs.Name, rs.Namespace)
    }

	
	// Get the ReplicaSet you want to find the pods for (by name and namespace)
	namespace := "" // Specify your namespace
	replicaSetName := "hello-minikube-ffcbb5874" // Replace with your ReplicaSet name

	replicaSet, err := clientset.AppsV1().ReplicaSets(namespace).Get(context.Background(), replicaSetName, metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Failed to fetch ReplicaSet: %v", err)
	}

	// Get the label selector from the ReplicaSet
	selector := labels.Set(replicaSet.Spec.Selector.MatchLabels).AsSelector()

	// List the pods in the same namespace that match the ReplicaSet's selector
	pods, err = clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		log.Fatalf("Failed to fetch Pods: %v", err)
	}

	// Print out the pods associated with the ReplicaSet
	fmt.Println("Pods associated with ReplicaSet:", replicaSetName)
	for _, pod := range pods.Items {
		fmt.Printf("Pod Name: %s, Status: %s\n", pod.Name, pod.Status.Phase)
	}
}

func loadClientSet() (*kubernetes.Clientset, error) { 
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Error loading kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	return clientset, nil
}

func getReplicaSets(clientset *kubernetes.Clientset) (*appsv1.ReplicaSetList, error) { 
	// List pods in the default namespace
    replicasets, err := clientset.AppsV1().ReplicaSets("").List(context.Background(), metav1.ListOptions{})
    if err != nil {
        log.Fatalf("Failed to fetch ReplicaSets: %v", err)
    }

    // Print the names of the ReplicaSets
    fmt.Println("List of ReplicaSets:")
    for _, rs := range replicasets.Items {
		fmt.Printf("Name: %s, Namespace: %s\n", rs.Name, rs.Namespace)
    }

	return replicasets, nil	
}

func getPods(clientset *kubernetes.Clientset) (*corev1.PodList, error) { 
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error fetching pods: %v", err)
	}
	
	fmt.Printf("There are %d pods\n", len(pods.Items))

	return pods, nil
}

func getDeployments(clientset *kubernetes.Clientset) (*appsv1.DeploymentList, error) {
	deployments, err := clientset.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error fetching deployments: %v", err)
	}
	
	fmt.Printf("There are %d deployments\n", len(deployments.Items))

	return deployments, nil
}

func getNodes(clientset *kubernetes.Clientset) (*corev1.NodeList, error) {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error fetching nodes: %v", err)
	}
	
	fmt.Printf("There are %d nodes\n", len(nodes.Items))

	return nodes, nil
}

func getServices (clientset *kubernetes.Clientset) (*corev1.ServiceList, error) {
	services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error fetching services: %v", err)
	}
	
	fmt.Printf("There are %d services\n", len(services.Items))

	return services, nil
}

func getConfigMaps (clientset *kubernetes.Clientset) (*corev1.ConfigMapList, error) {
	configmaps, err := clientset.CoreV1().ConfigMaps("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error fetching services: %v", err)
	}
	
	fmt.Printf("There are %d config maps\n", len(configmaps.Items))

	return configmaps, nil
}

func addReplicaSetsToGraph(graph *Graph, replicasets *appsv1.ReplicaSetList) { 
	for _, rs := range replicasets.Items { 
		revisionnumber, err := strconv.Atoi(rs.Labels["deployment.kubernetes.io/revision"])
		if err != nil { 
			revisionnumber = 1
		}

		rsNode := GraphNode{
			Key: EntityKey{
				Name:      rs.Name,
				Namespace: rs.Namespace,
				Type:      "ReplicaSet",
			},
			Properties: map[string]string{},
			Revision:   revisionnumber,
		}

		// fmt.Printf("%+v\n", rs.Labels)
		// for key, value := range rs.Labels {
		// 	fmt.Printf("Key: %s, Value: %d\n", key, value)
		// }

		graph.Nodes = append(graph.Nodes, rsNode)
	}
}

func addPodsToGraph(graph *Graph, pods *corev1.PodList) { 
	for _, p := range pods.Items { 
		revisionnumber, err := strconv.Atoi(p.Labels["deployment.kubernetes.io/revision"])
		if err != nil { 
			revisionnumber = 1
		}

		pNode := GraphNode{
			Key: EntityKey{
				Name:      p.Name,
				Namespace: p.Namespace,
				Type:      "Pod",
			},
			Properties: map[string]string{},
			Revision:   revisionnumber,
		}

		graph.Nodes = append(graph.Nodes, pNode)
	}
}

func addDeploymentsToGraph(graph *Graph, deployments *appsv1.DeploymentList) { 
	for _, d := range deployments.Items { 
		revisionnumber, err := strconv.Atoi(d.Labels["deployment.kubernetes.io/revision"])
		if err != nil { 
			revisionnumber = 1
		}

		dNode := GraphNode{
			Key: EntityKey{
				Name:      d.Name,
				Namespace: d.Namespace,
				Type:      "Deployment",
			},
			Properties: map[string]string{},
			Revision:   revisionnumber,
		}

		graph.Nodes = append(graph.Nodes, dNode)
	}
}

func addNodesToGraph(graph *Graph, nodes *corev1.NodeList) { 
	for _, n := range nodes.Items { 
		revisionnumber, err := strconv.Atoi(n.Labels["deployment.kubernetes.io/revision"])
		if err != nil { 
			revisionnumber = 1
		}

		nNode := GraphNode{
			Key: EntityKey{
				Name:      n.Name,
				Namespace: n.Namespace,
				Type:      "Node",
			},
			Properties: map[string]string{},
			Revision:   revisionnumber,
		}

		graph.Nodes = append(graph.Nodes, nNode)
	}
}

func addServicesToGraph(graph *Graph, services *corev1.ServiceList) { 
	for _, s := range services.Items { 
		revisionnumber, err := strconv.Atoi(s.Labels["deployment.kubernetes.io/revision"])
		if err != nil { 
			revisionnumber = 1
		}

		sNode := GraphNode{
			Key: EntityKey{
				Name:      s.Name,
				Namespace: s.Namespace,
				Type:      "Service",
			},
			Properties: map[string]string{},
			Revision:   revisionnumber,
		}

		graph.Nodes = append(graph.Nodes, sNode)
	}
}

func addConfigMapsToGraph(graph *Graph, services *corev1.ConfigMapList) { 
	for _, c := range services.Items { 
		revisionnumber, err := strconv.Atoi(c.Labels["deployment.kubernetes.io/revision"])
		if err != nil { 
			revisionnumber = 1
		}

		cNode := GraphNode{
			Key: EntityKey{
				Name:      c.Name,
				Namespace: c.Namespace,
				Type:      "ConfigMap",
			},
			Properties: map[string]string{},
			Revision:   revisionnumber,
		}

		graph.Nodes = append(graph.Nodes, cNode)
	}
}

func associatePodsToReplicaSets(clientset *kubernetes.Clientset, graph *Graph, replicasets *appsv1.ReplicaSetList) { 
	for _, rs := range replicasets.Items { 
		selector := labels.Set(rs.Spec.Selector.MatchLabels).AsSelector()
		pods, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{
			LabelSelector: selector.String(),
		})
		if err != nil {
			log.Fatalf("Failed to fetch Pods: %v", err)
		}
			
		rsentitykey := EntityKey{Name: rs.Name, Namespace: rs.Namespace, Type: "ReplicaSet"}
		for _, p := range pods.Items { 
			pentitykey := EntityKey{Name: p.Name, Namespace: p.Namespace, Type: "Pod"}

			relationship := GraphRelationship{
				Source: rsentitykey,
				Target: pentitykey, 
				RelationshipType: "Manages",
				Properties: map[string]string{},
				Revision: 1,
			}
			graph.Relationships = append(graph.Relationships, relationship)
		}

	}
}

func associateReplicaSetsToDeployments(clientset *kubernetes.Clientset, graph *Graph, deployments *appsv1.DeploymentList) { 
	for _, d := range deployments.Items { 
 		selector := labels.Set(d.Spec.Selector.MatchLabels).AsSelector()
		replicasets, err := clientset.AppsV1().ReplicaSets("").List(context.Background(), metav1.ListOptions{
			LabelSelector: selector.String(),
		})
		if err != nil {
			log.Fatalf("Failed to fetch ReplicaSets: %v", err)
		}
			
		dentitykey := EntityKey{Name: d.Name, Namespace: d.Namespace, Type: "Deployment"}
		for _, rs := range replicasets.Items { 
			rsentitykey := EntityKey{Name: rs.Name, Namespace: rs.Namespace, Type: "ReplicaSet"}

			relationship := GraphRelationship{
				Source: dentitykey,
				Target: rsentitykey, 
				RelationshipType: "Controls",
				Properties: map[string]string{},
				Revision: 1,
			}
			graph.Relationships = append(graph.Relationships, relationship)
		}
	}
}

func associatePodsToNodes(clientset *kubernetes.Clientset, graph *Graph, nodes *corev1.NodeList, pods *corev1.PodList) { 
	for _, n := range nodes.Items { 
		nentitykey := EntityKey{Name: n.Name, Namespace: n.Namespace, Type: "Node"}
		for _, p := range pods.Items { 
			if p.Spec.NodeName == n.Name { 
				pentitykey := EntityKey{Name: p.Name, Namespace: p.Namespace, Type: "Pod"}
	
				relationship := GraphRelationship{
					Source: nentitykey,
					Target: pentitykey, 
					RelationshipType: "Runs",
					Properties: map[string]string{},
					Revision: 1,
				}
				graph.Relationships = append(graph.Relationships, relationship)
			}
		}
	}
}

func associatePodsToServices(clientset *kubernetes.Clientset, graph *Graph, services *corev1.ServiceList) { 
	for _, s := range services.Items { 
		selector := s.Spec.Selector
		if len(selector) == 0 {
			continue
		}
	
		var labelSelector string
		for key, value := range selector {
			if labelSelector != "" {
				labelSelector += ","
			}
			labelSelector += fmt.Sprintf("%s=%s", key, value)
		}
	
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		if err != nil {
			log.Fatalf("Failed to fetch Pods: %v", err)
		}
			
		sentitykey := EntityKey{Name: s.Name, Namespace: s.Namespace, Type: "Service"}
		for _, p := range pods.Items { 
			pentitykey := EntityKey{Name: p.Name, Namespace: p.Namespace, Type: "Pod"}

			relationship := GraphRelationship{
				Source: sentitykey,
				Target: pentitykey, 
				RelationshipType: "Exposes",
				Properties: map[string]string{},
				Revision: 1,
			}
			graph.Relationships = append(graph.Relationships, relationship)
		}
	}
}

func associatePodsToConfigMaps(clientset *kubernetes.Clientset, graph *Graph, configmaps *corev1.ConfigMapList, pods *corev1.PodList) { 
	for _, c := range configmaps.Items { 
		centitykey := EntityKey{Name: c.Name, Namespace: c.Namespace, Type: "ConfigMap"}
		for _, p := range pods.Items {
			var uses bool = false

			for _, volume := range p.Spec.Volumes {
				if (volume.ConfigMap != nil && volume.ConfigMap.Name == c.Name) {
					pentitykey := EntityKey{Name: p.Name, Namespace: p.Namespace, Type: "Pod"}
		
					relationship := GraphRelationship{
						Source: pentitykey,
						Target: centitykey, 
						RelationshipType: "Uses",
						Properties: map[string]string{},
						Revision: 1,
					}
					graph.Relationships = append(graph.Relationships, relationship)

					uses = true
					break
				}
			}

			if (uses) { 
				continue
			}
	
			for _, container := range p.Spec.Containers {
				for _, envFrom := range container.EnvFrom {
					if envFrom.ConfigMapRef != nil && envFrom.ConfigMapRef.Name == c.Name {
						pentitykey := EntityKey{Name: p.Name, Namespace: p.Namespace, Type: "Pod"}
		
						relationship := GraphRelationship{
							Source: pentitykey,
							Target: centitykey, 
							RelationshipType: "Uses",
							Properties: map[string]string{},
							Revision: 1,
						}
						graph.Relationships = append(graph.Relationships, relationship)
					}
				}
			}
		}
	}
}



func main() {
	graph := Graph{
		Nodes:         []GraphNode{},
		Relationships: []GraphRelationship{},
	}
	clientset, _ := loadClientSet()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Start an infinite loop
	run_iteration := 0

	for {
		select {
		case <-ticker.C:
			// Execute the function when the ticker ticks
			replicasets, _ := getReplicaSets(clientset)
			addReplicaSetsToGraph(&graph, replicasets)
		
			pods, _ := getPods(clientset)
			addPodsToGraph(&graph, pods)
		
			deployments, _ := getDeployments(clientset)
			addDeploymentsToGraph(&graph, deployments)
		
			nodes, _ := getNodes(clientset)
			addNodesToGraph(&graph, nodes)
		
			services, _ := getServices(clientset)
			addServicesToGraph(&graph, services)
			
			configmaps, _ := getConfigMaps(clientset)
			addConfigMapsToGraph(&graph, configmaps)
		
			associatePodsToReplicaSets(clientset, &graph, replicasets)
			associateReplicaSetsToDeployments(clientset, &graph, deployments)
			associatePodsToNodes(clientset, &graph, nodes, pods)
			associatePodsToServices(clientset, &graph, services)
			associatePodsToConfigMaps(clientset, &graph, configmaps, pods)
		
			// fmt.Printf("Graph with nodes: %+v\n", graph)
		
			jsonDataIndented, err := json.MarshalIndent(graph, "", "  ")
			if err != nil {
				log.Fatalf("Error marshaling struct: %v", err)
			}
		
			fmt.Println("\nJSON in iteration #%d:", run_iteration)
			fmt.Println(string(jsonDataIndented))	

			run_iteration += 1
		}
	}	
	
}

