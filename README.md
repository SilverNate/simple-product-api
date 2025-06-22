![CI](https://github.com/SilverNate/simple-product-api/actions/workflows/ci.yml/badge.svg)

# 🛍️ Simple Product API

A robust, scalable Golang backend API for product management (Sayuran, Buah, Protein, Snack), built with Fiber, PostgreSQL, Redis, Wire, and Clean Architecture following Solid Pricinple.

---

## 🚀 Features
Create product: ✅ Done

List products: ✅ Done

Search by name and ID: ✅ Done

Filter by type (Sayuran, Buah, Protein, Snack): ✅ Done

Sorting (by name, price, created_at): ✅ Done

Pagination: ✅ Done

Robust, scalable architecture: ✅ Done

SOLID principle, Clean Architecture: ✅ Done

Use Redis cache in usecase layer: ✅ Done

Rate Limiting middleware: ✅ Done

Retry middleware (resilience): ✅ Done

Unit tests (success, failure, edge cases): ✅ Done

Redis simulation tests (miss, error, hit): ✅ Done

Repository layer tests: ✅ Done

Swagger documentation (/docs/index.html): ✅ Done

Docker & Docker Compose (Postgres + Redis + App): ✅ Done

.env config loading (no hardcoded values): ✅ Done

Logging with sirupsen/logrus (injected via Wire): ✅ Done

Dependency Injection with Google Wire: ✅ Done

GitHub Actions CI (go test ./...): ✅ Done

Makefile for automation (make test, make swag, etc.): ✅ Done

---

## 📦 Tech Stack
- Golang + Fiber
- PostgreSQL + Redis
- Google Wire for DI
- Swagger + Docker + GitHub Actions
- SOLID + Clean Architecture
- logrus
- testify + redismock + mockery
- makefile

---

## 📚 API Docs
API: 
- http://localhost:8080/api/v1/products

Visit Swagger after server up:  
- http://localhost:8080/swagger/index.html


---
## 📘 How To Use
Wire dependencies :
-  wire gen ./pkg/di

Start docker server:
- make docker-up

Stop docker server:
- make docker-down

---

## 🧪 Test
Add Mock All interface:
- make mocks-all

Add Mock product interface:
- make mocks-product

Run unit tests:
- make test-product

Run Bencmark test :
- make test-bench-product

---

