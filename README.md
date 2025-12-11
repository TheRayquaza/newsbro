# NewsBro - MLOPS Final Project

Distributed resilient architecture for news recommendation

Feel free to have a look at our [website](https://app.newsbro.cc) to create your account
and interact with 100k+ articles.

![login page](./assets/login.png) ![home page](./assets/home.png)

We provide automatic update everyday (9am / 6pm UTC), 1k articles ingested
from various sources. We also provide new recommendation based on your liked articles.

Feel free to open issues if you see anything relevant that could be added.

![presentation.pdf](./assets/presentation.pdf))

⚠ App will be down at the end of our project (~ January 2026)

## Table of Contents
- [Project Structure](#project-structure)
- [Architecture](#architecture)
    - [Light Speed Recommendation](#light-speed-recommendation)
    - [North-South Traffic](#nort-south-traffic)
    - [Secrets Rollout](#secrets-rollout)
- Models
    - [TF-IDF / SVD](docs/model/tf_idf.md#tf-idf)
    - [SBERT](docs/model/sbert.md#sbert)
- [Development](#dev)
    - [CI/CD](#cicd)
    - [Commits](#commits)
- [Services](./apps/README.md)
    - [port-front](./apps/port-front/README.md)
    - [repo-account](./apps/repo-account/README.md)
    - [repo-article](./apps/repo-article/README.md)
    - [repo-feed](./apps/repo-feed/README.md)
    - [srvc-inference](./apps/srvc-inference/README.md)
    - [srvc-scrapping](./apps/srvc-scrapping/README.md)
    - [srvc-drift](./apps/srvc-drift/README.md)
- [Deployment](#deployment)
    - [Deployment Specifications](#deployment-specifications)

### Privacy Note

We only use your feedback to improve our system recommenders.
- Passwords are never stored, authentication is handled through secure and short-lived JWT.
- Internal credentials are randomly generated for each session and securely stored in Vault, never embedded in code, only injected as env variables.
- All feedback data is anonymized before being used for training or evaluation (only `user_id` is used), ensuring that no personally identifiable information is retained.
- Data transmission is encrypted (TLS) between client & dns provider, and between dns provider and server.

All internal data will be securely erased once the project concludes.

## Project Structure

```
.
├── apps
│   ├── repo-account
│   ├── repo-article
│   ├── repo-feed
│   ├── port-front
│   ├── srvc-scrapping
│   ├── srvc-search
│   ├── srvc-drift
│   └── srvc-inference
├── docs
├── dev                          # dev folder for various dev operations
├── k8s
│   ├── apps                     # deployment of each apps
│   ├── capacitor                # Frontend to manage flux
│   ├── cert-manager             # Manager TLS certificate
│   ├── external-secrets         # ESO to communicate with vault
│   ├── flux/flux-system         # deployment of flux repo (gitops tools)
│   ├── ingress-nginx            # Ingress nginx
│   ├── kafka                    # Kafka brokers, controllers, redpanda console and topic definition
│   ├── metallb                  # Metal LB
│   ├── minio                    # Minio cluster storage (datalake + store for mlflow artifacts)
│   ├── mlflow
│   ├── postgres                 # CNPG CRD, database definition, user definition
│   ├── qdrant                   # Vector database
│   ├── redis                    # Redis Operator, Insight and Sentinels
│   └── vault                    # Vault server to store secrets, act as authenticator for SA in kubernetes
├── .github/workflows/
└── README.md
```

## Architecture

The project is built with a microservices architecture including the following components :

- `repo-account`: User account and login management
- `repo-article`: Article & feedbacks management
- `srvc-scrapping`: Scrapping for articles, ingest to kafka
- `srvc-search`: Search service
- `srvc-inference`: Inference service to provide update user feeds
- `srvc-drift`: Drift service to provide insights on potential data drift
- `repo-feed`: Service providing the user feed
- `port-front`: main frontend

![Architecture Diagram](docs/archi/archi_v1.2.png)

### Light Speed Recommendation

We want light speed recommendation and scalable infrastructure. Behind the scene:
- Kafka consumers consuming articles, processing articles by batch and creating new recommendations.
- Qdrant storing vector embeddings to pre compute recommendation as they come. These representation are then used in various way.
- Redis Sentinels with [zsets](https://redis.io/docs/latest/develop/data-types/sorted-sets/) to create efficient feed queue.
- Repo Feed, a dedicated service managing feed, storing feedbacks / articles for short periods?

![User Feed](docs/misc/user_feed_2.png)

### Nort-South Traffic

North-South traffic is handled by two separate ingress controllers.
- Public Ingress: on a dedicated ip behind a NAT, exposing public services including frontend and backend services.
- Private Ingress: using VPN ip handled by [wireguard](https://www.wireguard.com/), exposing critical and back office services.

We use also a dedicated load balancer called [Metal LB](https://metallb.io/) to expose ip pools for the two ingresses.

![North-South](docs/misc/north-south.png)

### Secrets Rollout

We use [external-secret](https://external-secrets.io/latest/) to create our secrets

Secrets are centralized within a [vault](https://developer.hashicorp.com/vault), provisioned by Terraform.

![External Secrets](https://external-secrets.io/v0.5.8/pictures/diagrams-high-level-ns-detail.png)

## Dev

We provide a generic docker compose to run a minimal stack at `apps/docker-compose.yml`

### CI/CD

- Our CI/CD will lint, build and test backend codes.
- To deploy a new version just tag `<srvc-name>-vx.x.x`.

### Commits

Example of commit we use:

```
k8s: redis: increase redis sentinel for minimal config
apps: port-front: now using get profile to retrieve user info as cookies could not be read by frontend
apps: review swagger, specify dto for errors
```

## Deployment

Our Project is deployed with Kubernetes. We [FluxCD](https://fluxcd.io/) as GitOps tool to automate the deployment of our new manifests.

To deploy the whole stack, just use:

```bash
kubectl apply -k k8s/flux/flux-system
```
