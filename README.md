
# AniScale: Distributed Manga Tracking System

AniScale is a cloud-native microservices application designed to track real-time manga updates and deliver instant notifications. Built with **Go** and orchestrated on **Google Kubernetes Engine (GKE)**, the system leverages a decoupled architecture via **Google Pub/Sub** to ensure high availability and scalability.

## Key Features

* **Decoupled Microservices:** Separate services for data scraping and notification delivery to ensure fault isolation.
* **Event-Driven Architecture:** Utilizes **Google Cloud Pub/Sub** as an asynchronous message broker to handle traffic spikes and fan-out notifications.
* **Containerized Orchestration:** Fully containerized with **Docker** and deployed on a **Kubernetes (GKE)** cluster.
* **Secure Config Management:** Implements **12-Factor App** principles using environment variables and **Kubernetes Secrets** for API credentials.
* **Real-time Alerts:** Instant delivery to **Discord** via Webhooks.

---

## Architecture Overview

The system consists of three main pillars:

1. **The Scraper (Producer):** A Go service that polls the MangaDex API, processes chapter metadata, and publishes updates to a Pub/Sub topic.
2. **The Message Broker:** Google Pub/Sub manages the message lifecycle, providing a buffer between ingestion and notification.
3. **The Notifier (Consumer):** A Go service that listens to the Pub/Sub subscription and pushes formatted alerts to Discord.

---

## Tech Stack

* **Language:** Go (Golang)
* **Cloud:** Google Cloud Platform (GKE, Artifact Registry, Pub/Sub)
* **Containers:** Docker, Docker Compose
* **Orchestration:** Kubernetes
* **API:** MangaDex API, Discord Webhooks

---

## Local Development

### Prerequisites

* Go 1.24+
* Docker & Docker Compose
* GCP Service Account Key (`aniscale-key.json`)

### Setup

1. **Clone the repo:**
```bash
git clone https://github.com/yourusername/aniscale.git
cd aniscale

```


2. **Configure environment:**
Create a `.env` file in the root:
```env
GCP_PROJECT_ID=your-project-id
DISCORD_WEBHOOK_URL=your-webhook-url

```


3. **Run with Docker Compose:**
```bash
docker-compose up --build

```



---

## Deployment

The project includes Kubernetes manifests for cloud deployment:

1. **Build and Push images** to Google Artifact Registry.
2. **Create K8s Secrets** for the Service Account key and environment variables.
3. **Apply manifests:**
```bash
kubectl apply -f k8s/

```



---

## Engineering Decisions

* **Why Pub/Sub?** By decoupling the scraper from the notifier, the system can scale "horizontally." We can add more notifiers (Email, SMS, Database logging) without ever touching the scraper code.
* **Why Go?** Used for its small memory footprint and excellent concurrency primitives (Goroutines), making it perfect for lightweight cloud containers.


