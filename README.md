# beamstack
<img src="https://github.com/BeamStackProj/beamstack-cli/blob/main/logo/logo.png" width="300">
----

beamstack is a revolutionary tool that has been meticulously crafted to simplify and revolutionize 
the deployment of machine learning (ML) workflows on [Kubernetes](https://kubernetes.io/docs/concepts/overview/). It offers a holistic solution by 
introducing abstraction layers that streamline the deployment of diverse components 
of ML pipelines, data processing workflows, and deployment infrastructure.


At the heart of beamstack's capabilities lie Kubernetes Custom Resource Definitions (CRDs). These CRDs 
serve as a powerful mechanism for extending the Kubernetes API, enabling the seamless integration 
of ML-specific resources into the Kubernetes ecosystem. Through this innovative approach, 
beamstack empowers users to leverage the robust features and functionalities of 
Kubernetes while unlocking the immense potential of ML.

----

## minikube
```sh
minikube delete && minikube start --kubernetes-version=v1.23.0 --memory=6g --bootstrapper=kubeadm --extra-config=kubelet.authentication-token-webhook=true --extra-config=kubelet.authorization-mode=Webhook --extra-config=scheduler.bind-address=0.0.0.0 --extra-config=controller-manager.bind-address=0.0.0.0
```
## To start using beamstack

Basic to advanced knowledge on kubernetes is advised. 

## Installing beamstack

```
git clone https://github.com/BeamStackProj/beamstack-cli.git
cd beamstack-cli
make install
```
