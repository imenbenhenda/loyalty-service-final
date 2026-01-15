# 🚀 Loyalty Points Service

![Go Version](https://img.shields.io/badge/Go-1.21-blue.svg)
![Docker](https://img.shields.io/badge/Docker-Enabled-blue)
![Kubernetes](https://img.shields.io/badge/Kubernetes-Ready-326ce5)
![CI/CD](https://img.shields.io/badge/GitHub%20Actions-CI%2FCD-green)

Un microservice backend robuste écrit en **Go**, conçu pour gérer un système de points de fidélité. Ce projet démontre une approche **DevOps complète** : développement API, conteneurisation, orchestration et pipelines d'intégration continue.

---

## 📋 Fonctionnalités Clés

- **API RESTful** : Gestion des points clients (ajout, consultation).
- **Architecture Cloud-Native** : Conçu pour tourner dans des conteneurs.
- **Observabilité** :
  - Logs structurés en JSON.
  - Endpoint de santé (`/health`) pour les sondes Kubernetes.
  - Métriques Prometheus (`/metrics`).
- **Sécurité** : Scan de vulnérabilités (SAST) intégré dans la CI.
- **CI/CD** : Automatisation complète via GitHub Actions.

---

## 🛠️ Stack Technique

- **Langage** : Go (Golang) 1.21
- **Conteneur** : Docker (Multi-stage build sur Alpine Linux)
- **Orchestration** : Kubernetes (Deployment & Service NodePort)
- **CI/CD** : GitHub Actions (Build, Test, Push to DockerHub)

---

## 📂 Structure du Projet

```
.
├── cmd/api/            # Point d'entrée (main.go)
├── internal/handlers/  # Logique métier (Business Logic)
├── k8s/                # Manifestes d'infrastructure (IaC)
│   ├── deployment.yaml
│   └── service.yaml
├── .github/workflows/  # Pipeline CI/CD
├── Dockerfile          # Image Docker optimisée
├── go.mod              # Dépendances Go
└── README.md           # Documentation
```
---
🚀 Guide d'Installation et Démarrage
📋 Pré-requis
```
1.Go 1.21+ (pour exécution locale)
2.Docker Desktop (avec Kubernetes activé)
3.Git
```
   
---
1️⃣ Exécution Locale (Sans Docker)
