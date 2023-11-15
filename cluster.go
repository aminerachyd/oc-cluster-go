package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/tabwriter"
)

type Cluster struct {
	Name     string `yaml:"name"`
	Url      string `yaml:"url"`
	Username string `yaml:"username"`
	Comment  string `yaml:"comment"`
}

func ConnectToCluster(clusterName string, c Config) {
	cluster, found := findCluster(clusterName, c.Clusters)
	if !found {
		log.Fatalf("Cluster %s not found", clusterName)
	}

	cluster.execConnectCommand()
}

func AddCluster(clusterName string, url string, username string, c Config) {
	cluster, exists := findCluster(clusterName, c.Clusters)
	if exists {
		cluster.Url = url
		cluster.Username = username
	} else {
		newCluster := Cluster{
			Name:     clusterName,
			Url:      url,
			Username: username,
			Comment:  "",
		}
		c.Clusters = append(c.Clusters, newCluster)
	}

	_, err := WriteConfig(c)
	if err != nil {
		log.Fatal("Error adding cluster", err)
	}
}

func DeleteCluster(clusterName string, c Config) {
	indexToDelete := -1
	for i, cluster := range c.Clusters {
		if cluster.Name == clusterName {
			indexToDelete = i
		}
	}

	if indexToDelete != -1 {
		c.Clusters = append(c.Clusters[:indexToDelete], c.Clusters[indexToDelete+1:]...)
		WriteConfig(c)
	}
}

func ListClusters(output string, c Config) {
	clusters := c.Clusters

	switch output {
	case "wide":
		printClustersTable(&clusters)
	default:
		printClustersTable(&clusters)
	}
}

func printClustersTable(clusters *[]Cluster) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
		"CLUSTER NAME",
		"USERNAME",
		"URL",
		"COMMENT")

	for _, c := range *clusters {
		fmt.Fprintf(w, "%s", c.toTableCell())
	}
}

func (c *Cluster) toTableCell() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\n", c.Name, c.Username, c.Url, c.Comment)
}

func (c *Cluster) execConnectCommand() {
	cmd := exec.Command("oc",
		"login",
		c.Url,
		"-u",
		c.Username)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	cmd.Run()
}

func findCluster(clusterName string, clusters []Cluster) (*Cluster, bool) {
	var cluster *Cluster
	var found bool

	for _, c := range clusters {
		if c.Name == clusterName {
			cluster = &c
			found = true
		}
	}

	return cluster, found
}
