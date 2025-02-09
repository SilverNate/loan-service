# Loan Service

## Overview

This Loan Service is a robust, scalable system designed to handle the entire lifecycle of a loan from creation (proposed) through approval, investments, and disbursement. The system calculates key metrics (such as total interest and investor ROI) and integrates with SendGrid to send dummy agreement letter PDFs once a loan is fully funded. The solution follows SOLID principles and employs a layered architecture to ensure maintainability, testability, and scalability.

## Architecture and Design

### Key Design Principles

- **Separation of Concerns / SOLID Principles:**
    - **Models:** Define domain entities (Loan, Borrower, Investor, Investment).
    - **Repositories:** Abstract data access (using PostgreSQL) with the Repository pattern.
    - **Service Layer:** Contains business logic:
        - Calculates `total_interest` (as `principal_amount * rate / 100`).
        - Validates and tracks investments.
        - Updates loan status when fully funded.
        - Integrates with SendGrid for sending agreement letters (dummy function :-) ).
    - **Handlers:** Provide HTTP endpoints using Go-Chi and offer a universal JSON response format.

- **RESTful API:**
    - Built with the lightweight [Go-Chi](https://github.com/go-chi/chi) router.
    - Consistent API responses via a shared JSON response helper.

- **Event Stream Processing :**
    - Kafka integration for decoupling state change events from processing (e.g., consumer services can later trigger notifications).

- **Containerization:**
    - Docker and Docker Compose are used to containerize the API service, database (PostgreSQL), Kafka, and Zookeeper, ensuring consistent deployments across environments.

- **Test-Driven Development (TDD):**
    - Unit tests are written using [Testify](https://github.com/stretchr/testify) and mocks are generated using tools like `mockgen` to ensure that services, repositories, and handlers work as expected.

- **File Uploads & SendGrid Integration:**
    - The disbursement handler accepts file uploads (via `multipart/form-data`) for uploading agreement letter proofs.
    - Uploaded files are saved locally (or can be integrated with cloud storage) and the service sends agreement letters via SendGrid when a loan is fully funded.

## Setup and Installation

### Prerequisites

- **Go 1.20 or later**
- **Docker and Docker Compose** (for containerized deployment)
- **PostgreSQL** (if running locally without Docker)
- **Kafka and Zookeeper** (if using the Kafka event streaming functionality)
- **SendGrid API Key** (for sending emails)

### Database Setup

- Use the SQL migration script (found in the `migrations/` folder) to set up your PostgreSQL database. The migration creates four tables:
    - **borrowers**
    - **investors**
    - **loans** (with a generated `total_interest` column)
    - **investments**

### Running with Docker

1. **Build and Start Containers:**

   ```bash
   docker-compose up --build
