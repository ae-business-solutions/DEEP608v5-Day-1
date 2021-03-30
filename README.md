# DEEP608v5 - Day 1
### Objectives for Today
* __Review Lab Architecture__
* __Setup required accounts__
	* GitHub (if you don't already have one)
	* GCP - Provided for you
	* Terraform Cloud - Provided for you
* __Install required tools__
	* Git
	* Visual Studio Code ("VS Code")
		* Plugin: GitHub
		* Plugin: Terraform
* __Clone *DEEP608v5-Day-1* Repo Template__
* __Setup *GitHub Workflows* CI/CD & integrate with *Terraform Cloud* and *GCP*__
	* To align with Git best practices, move to a *pull request* based workflow.
* __Build Containers for our Services & Push to Google Container Registry (GCR)__
	* __Block Page__ - *Static web page simulating URL Filtering block page*
	* __Request Unblock__ - *Rust app serving a web form used to request a URL be unblocked*
	* __EDL Admin__ - *Golang app serving an admin page to approve/deny URL unblock requests*
	* __EDL__ - *Python app that serves out a plaintext list of approved URLs; these can be consumed by NGFWs or other URL filters*
* __Deploy resources on GCP (Google Cloud Platform)__
	* Kubernetes Cluster (GKE)
	* Redis Instance (Cloud Memorystore)

### Lab Guide
* [Chapter 1](guide/chapter1.md) - **Lab Architecture**
* [Chapter 2](guide/chapter2.md) - **Accounts & Tools**
* [Chapter 3](guide/chapter3.md) - **GitHub Workflows (CI/CD)**
* [Chapter 4](guide/chapter4.md) - **Terraform: GCP**