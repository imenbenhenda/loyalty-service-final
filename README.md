# Loyalty Points Service 🚀

Un microservice backend écrit en Go pour gérer un système de points de fidélité.

## 📋 Fonctionnalités
- **API REST** : Gestion des points.
- **Observabilité** : Logs JSON et Métriques Prometheus (`/api/v1/metrics`).
- **Sécurité** : Scans SAST (Gosec) et DAST (OWASP ZAP).
- **CI/CD** : Pipeline GitHub Actions automatisée.

## 🛠️ Installation & Démarrage

### Pré-requis
- Docker Desktop avec Kubernetes activé.

### Déploiement Rapide
```bash
# 1. Déployer
kubectl apply -f kubernetes/deployment.yaml
kubectl apply -f kubernetes/service.yaml

# 2. Accéder (Port-Forwarding)
kubectl port-forward service/loyalty-service 8081:80