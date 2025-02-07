# Go-Based Microservice Coffee Shop Application

This repository contains the code for building a multi-tier microservice architecture, following the "Building Microservices in Go" YouTube series. The application demonstrates how to build various microservices using the Go programming language, focusing on best practices and industry-standard tools.

## NOTE: this is a re-upload of the original project completed in 3/2024

## Table of Contents
- [Overview](#overview)
- [Microservices](#microservices)
  - [Product API](#product-api)
  - [Frontend](#frontend)
  - [Currency](#currency)
  - [Product Images](#product-images)
- [Learning Resources](#learning-resources)
- [Setup Instructions](#setup-instructions)
- [License](#license)

## Overview
This project aims to demonstrate how to design and implement microservices in Go. It covers building RESTful APIs, gRPC services, handling file uploads, validation, and much more. The goal is to create a scalable, modular system that can be easily extended and maintained.

### Key Features
- **Product API**: RESTful Go-based JSON API using the Gorilla framework for CRUD operations.
- **Frontend**: ReactJS frontend that presents data from the Product API.
- **Currency**: gRPC service for live currency exchange rates, supporting unary and bidirectional streaming methods.
- **Product Images**: Go-based image service that supports gzipped content and multipart file uploads.

## Microservices

### Product API
- A **RESTful Go-based JSON API** built using the **Gorilla framework**. It supports CRUD operations on a product list, useful for a coffee shop's product inventory.
  
### Frontend
- A **ReactJS** frontend that consumes the Product API and presents the information in an easy-to-use interface.

### Currency
- A **gRPC service** that supports simple unary and bidirectional streaming methods. This service provides live currency exchange rates used in the coffee shop for international transactions.

### Product Images
- A **Go-based image service** that supports gzipped content, multi-part form uploads, and RESTful methods for uploading and downloading product images.

## Learning Resources
Follow along with the weekly video tutorials in the "Building Microservices in Go" YouTube series to learn how to build and extend this microservice architecture.

- [Playlist Link - Building Microservices in Go](https://www.youtube.com/playlist?list=PLmD8u-IFdreyh6EUfevBcbiuCKzFk0EW_)

## Setup Instructions

### Prerequisites
- Go (version 1.16 or higher)
- Node.js and npm (for the frontend)
- Docker (optional for containerization)


