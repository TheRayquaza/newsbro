# NewsBro - MLOPS Final Project

Distributed resilient architecture for news recommendation

## Table of Contents
- [App](#App)
- [Project Structure](#project-structure)
- [Documentation](#documentation)
- [Architecture](#architecture)
- [Testing](#testing)
- [Deployment](#deployment)
- [Deployment Specifications](#deployment-specifications)

## App

A small react web app has been designed to interact with articles.

Feel free to have a look at https://app.newsbro.cc website to create your account
and interact with 100k+ articles.

We provide automatic update everyday (9am / 6pm UTC), 10k articles ingested
from various sources, list can be found ./k8s/apps/srvc-scrapping/configmap.yaml

![login page](./assets/login.png)

Feel free to open issues if you see anything relevant that could be added.
App will be down at the end of our project (~ Jan. 2026)

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
│   ├── port-front
│   ├── srvc-scrapping
│   ├── srvc-search
│   └── srvc-inference
├── docs
├── dev                          # dev folder for various dev operations
├── k8s
│   ├── apps                     # deployment of each apps
│   ├── capacitor                # Frontend to manage flux
│   ├── cert-manager             # Manager TLS certificate
│   ├── external-secrets         # ESO to communicate with vault
│   ├── flux/flux-system         # deployment of flux repo (gitops tools)
│   ├── kafka                    # Kafka brokers, controllers, redpanda console and topic definition
│   ├── minio                    # Minio cluster storage (datalake + store for mlflow artifacts)
│   ├── mlflow
│   ├── postgres                 # CNPG CRD, database definition, user definition
│   ├── qdrant                   # Vector database
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
- `srvc-inference`: Inference service to provide use recommendation
- `port-front`: main frontend

![Architecture Diagram](docs/archi/archi_v1.0.png)

## Testing

Individual service tests can be run by navigating to the specific service directory and executing:

```bash
docker compose up -d
cd apps/<service-name> # TODO
cd ../..
docker compose down
```

## Deployment

Our Project is deployed on Kubernetes. The deployment process is handled through our CI/CD pipeline.

To deploy the whole stack, just use:

```bash
kubectl apply -k k8s/ # TODO
```

## Deployment Specifications

TODO

