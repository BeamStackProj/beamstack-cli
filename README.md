<p align="center">
  <picture>
    <source media="(prefers-color-scheme: light)" srcset="./assets/beamstack-logo.png">
    <img width="300" height="260" src="./assets/beamstack-logo.png">
  </picture>
  <h1 align="center" style="font-size: 24px;">Kubernetes Framework for deploying ML and GenAI Apache Beam workflows</h1>
</p>  

<p align="center">
  <a href="https://beamstack.netlify.app/docs/" rel="nofollow"><strong>Explore Beamstack Documentation »</strong></a>
  <a href="https://beamstack.netlify.app/community/"><strong>Join Beamstack Community »</strong></a>
  <a href="https://beamstack.netlify.app/blog/"><strong>Explore Blogs »</strong></a>
  <a href="https://discord.gg/fYNnNVaEFK"><strong>Join Discord Channel</strong></a>
</p>

</p>
<p align="center">
<a href="https://discord.gg/fYNnNVaEFK"><img src="https://img.shields.io/badge/Join%20us%20on-Discord-e01563.svg" alt="Join Discord"></a>
<a href="http://golang.org"><img src="https://img.shields.io/badge/Made%20with-Go-1f425f.svg" alt="made-with-Go"></a>

## **Beamstack Features**

<details>
  <summary><b>Simplified ML Workflow Deployment:</b></summary>
  <ul>
    <li>Beamstack simplifies the deployment of machine learning workflows on Kubernetes.</li>
  </ul>
</details>

<details>
  <summary><b>Holistic Solution:</b></summary>
  <ul>
    <li>Beamstack offers an all-encompassing solution for managing machine learning pipelines, data processing workflows, and deployment infrastructure.</li>
  </ul>
</details>

<details>
  <summary><b>Abstraction Layers:</b></summary>
  <ul>
    <li>Beamstack introduces abstraction layers that streamline the deployment of various components within ML pipelines.</li>
  </ul>
</details>

<details>
  <summary><b>Leveraged Kubernetes Custom Resource Definitions (CRDs):</b></summary>
  <ul>
    <li>Beamstack uses Kubernetes CRDs to extend the Kubernetes API, allowing smooth integration of machine learning-specific resources.</li>
  </ul>
</details>

<details>
  <summary><b>Seamless Integration with Kubernetes:</b></summary>
  <ul>
    <li>Beamstack empowers users to leverage Kubernetes' features while incorporating machine learning capabilities into the Kubernetes ecosystem.</li>
  </ul>
</details>

<details>
  <summary><b>Easily Monitor and Visualize Deployed Workflows:</b></summary>
  <ul>
    <li>Beamstack seamlessly integrates with Prometheus and Grafana to visualize the states of the deployed workflows in real time.</li>
  </ul>
</details>  
  
---  

## **Architecture of beamstack** 
<p align="center"><img src="./assets/beamstack-arch.png"></p>
  
--- 

## **Components of Beamstack** 
<ul>
  <li>Beamstack CLI</li>
  <li>Beamstack Custom Transforms</li>
  <li>Apache Beam YAML</li>
</ul>  

---

## **Installation**  

### Setup Kubernetes cluster:  
To be able to work with beamstack-cli, an active Kubernetes cluster is required.  

A local Kubernetes cluster can be setup using minikube.  
  
```bash
minikube delete && minikube start --kubernetes-version=v1.23.0 --memory=6g --bootstrapper=kubeadm --extra-config=kubelet.authentication-token-webhook=true --extra-config=kubelet.authorization-mode=Webhook --extra-config=scheduler.bind-address=0.0.0.0 --extra-config=controller-manager.bind-address=0.0.0.0
``` 

### Clone beamstack-cli resources:
   
```bash
git clone https://github.com/BeamStackProj/beamstack-cli.git
```  

### Install beamstack-cli:  
  
```bash
cd beamstack-cli
make install
```
---  

## **Examples of beamstack commands** 

### Initialize a Kubernetes Cluster with Beamstack & Monitoring tools:  

```bash
beamstack init -m
```  
  
### Get the current kubernetes cluster context and profile info:  

```bash
beamstack info
```  

### Display the name, status and age of a cluster:  

```bash
beamstack info cluster
```  

### Create a runner cluster:  

```bash
beamstack create [runner-cluster] [cluster-name]
```  

### Open runner UI:  

```bash
beamstack open [runner] [runner-cluster-name]
```  

### Deploy a pipeline:  

```bash
beamstack deploy pipeline [FILE] [flags]
```  

### Create a vector store:  

```bash
beamstack create vector-store --type=elasticsearch
```  

### Get help:  

```bash
beamstack --help  

beamstack [command] --help
```  

## **Support, Contribution, and Community**
 
### :busts_in_silhouette: Community
 
Get updates on Beamstack's development and chat with project maintainers, contributors, and community members  
- Visit the [Community Page](https://beamstack.netlify.app/community/)
- Raise feature requests, suggest enhancements, and report bugs in our [GitHub Issues](https://github.com/BeamStackProj/beamstack-cli/issues)
- Articles, Howtos, Tutorials - [Beamstack Blogs](https://beamstack.netlify.app/blog/)

### :handshake: Contribute
 
Take a look at our [contributing guidelines](https://beamstack.netlify.app/docs/contribution-guidelines) for information on how to open issues, adhere to coding standards, and understand our development processes. We greatly value your contribution.