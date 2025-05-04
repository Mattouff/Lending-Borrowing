# Lending & Borrowing Platform - Backend Documentation

## Overview

This document provides comprehensive technical documentation for the Lending & Borrowing Platform backend. The platform offers a DeFi solution that allows users to deposit collateral, borrow assets, and earn interest on their deposits.

## Table of Contents

- Architecture
- Project Structure
- Getting Started
- Environment Configuration
- API Documentation
- Smart Contract Integration
- Database Schema
- Development Workflow
- Deployment
- Testing
- Troubleshooting

## Architecture

The backend follows a clean architecture pattern with:

- **API Layer**: Handles HTTP requests and responses using Fiber framework
- **Service Layer**: Contains business logic and orchestrates operations
- **Repository Layer**: Data access layer for database operations
- **Domain Layer**: Core business entities and interfaces
- **Infrastructure Layer**: External systems integration including blockchain

### Key Components

- **Fiber**: Fast HTTP framework for REST API
- **GORM**: ORM for PostgreSQL database access
- **JWT**: Used for authentication and authorization
- **Swagger**: API documentation
- **Docker**: Containerization for development and deployment
- **Foundry**: Smart contract deployment and interaction
- **PostgreSQL**: Relational database for persistence

## Project Structure

```
backend/
├── internal/              # Internal packages
│   ├── api/               # API layer (handlers, routes, DTOs)
│   │   ├── dto/           # Data Transfer Objects
│   │   ├── handlers/      # Request handlers
│   │   ├── middleware/    # HTTP middleware
│   │   └── routes/        # API route definitions
│   ├── config/            # Configuration management
│   ├── contracts/         # Generated smart contract bindings
│   ├── domain/            # Core business logic
│   │   ├── models/        # Domain models
│   │   ├── repository/    # Repository interfaces
│   │   └── service/       # Service interfaces
│   ├── infrastructure/    # External systems integration
│   │   ├── blockchain/    # Blockchain client and services
│   │   └── persistence/   # Database implementations
│   └── service/           # Service implementations
├── pkg/                   # Reusable packages
│   ├── blockchain/        # Blockchain utilities
│   └── database/          # Database utilities
├── docs/                  # Swagger documentation
├── scripts/               # Utility scripts
└── main.go                # Application entry point
```

## Getting Started

### Prerequisites

- Go 1.21+
- Docker and Docker Compose
- Foundry (for blockchain development)
- PostgreSQL (or use Docker container)

### Setup Process

1. Clone the repository:

```bash
git clone https://github.com/Mattouff/Lending-Borrowing.git
cd Lending-Borrowing
```

2. Create environment files:

```bash
cp .env.example .env.dev
# Edit .env.dev with your local settings
```

3. **First: Deploy Smart Contracts with Foundry**

   a. Start Anvil (local Ethereum node):

   ```bash
   anvil --mnemonic "test test test test test test test test test test test junk"
   ```

   b. Deploy contracts:

   ```bash
   cd contracts
   forge script script/Deploy.s.sol --fork-url http://localhost:8545 --broadcast
   ```

   c. Copy the contract addresses from the output and update your `.env.dev` file:

   ```
   CONTRACT_ADDRESSES=LendingPool=0x...,Token=0x...,Borrowing=0x...,Collateral=0x...
   ```

4. **Second: Start the Backend Application**

   ```bash
   docker-compose --env-file ./backend/.env.dev -f docker-compose.dev.yml up -d
   ```

5. Access the API:
   - API: http://localhost:8080/api/v1
   - Swagger: http://localhost:8080/swagger/

## Environment Configuration

The application uses different environment files for different deployment contexts:

- `.env.dev` - Development environment
- `.env.prod` - Production environment

### Key Environment Variables

```
# Application settings
APP_NAME=Lending-Borrowing
APP_ENV=[development|production]
APP_DEBUG=[true|false]
LOG_LEVEL=[debug|info|warn|error]

# Server settings
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database settings
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=lending_borrowing
DB_SSLMODE=[disable|require]
DB_TIMEZONE=UTC

# Blockchain settings
BLOCKCHAIN_RPC_URL=http://localhost:8545
BLOCKCHAIN_NETWORK=[local|mainnet|sepolia|etc]
BLOCKCHAIN_CHAIN_ID=31337
BLOCKCHAIN_GAS_LIMIT=3000000
BLOCKCHAIN_GAS_PRICE=20000000000

# Contract addresses
CONTRACT_ADDRESSES=LendingPool=0x...,Token=0x...,Borrowing=0x...,Collateral=0x...

# JWT settings
JWT_SECRET=your_secure_secret
JWT_EXPIRE=24
```

## API Documentation

The API is documented using Swagger. Access the documentation at:

```
http://localhost:8080/swagger/
```

### Available Endpoints

#### User Management

- `POST /api/v1/users/register` - Register new user
- `POST /api/v1/users/auth` - Authenticate with signature
- `GET /api/v1/users/nonce/:address` - Get nonce for address
- `GET /api/v1/users/profile` - Get user profile (auth required)
- `PUT /api/v1/users/profile` - Update user profile (auth required)
- `GET /api/v1/users/admin` - List all users (admin only)
- `GET /api/v1/users/admin/:id` - Get user by ID (admin only)
- `GET /api/v1/users/admin/address/:address` - Get user by address (admin only)
- `PUT /api/v1/users/admin/:id/verify` - Verify user (admin only)

#### Lending Operations

- `GET /api/v1/lending/pool-info` - Get lending pool information
- `POST /api/v1/lending/deposit` - Deposit tokens (auth required)
- `POST /api/v1/lending/withdraw` - Withdraw tokens (auth required)
- `GET /api/v1/lending/balance` - Get lending balance (auth required)
- `GET /api/v1/lending/info` - Get lending info (auth required)
- `GET /api/v1/lending/transactions` - Get lending transaction history (auth required)

#### Borrowing Operations

- `GET /api/v1/borrowing/stats` - Get borrowing statistics
- `POST /api/v1/borrowing/borrow` - Borrow tokens (auth required)
- `POST /api/v1/borrowing/repay` - Repay borrowed tokens (auth required)
- `GET /api/v1/borrowing/balance` - Get borrowed amount (auth required)
- `GET /api/v1/borrowing/info` - Get borrowing info (auth required)
- `GET /api/v1/borrowing/transactions` - Get borrowing transaction history (auth required)

#### Collateral Operations

- `POST /api/v1/collateral/deposit` - Deposit collateral (auth required)
- `POST /api/v1/collateral/withdraw` - Withdraw collateral (auth required)
- `GET /api/v1/collateral/balance` - Get collateral balance (auth required)
- `GET /api/v1/collateral/info` - Get collateral info (auth required)

#### Liquidation Operations

- `GET /api/v1/liquidation/positions` - Get liquidatable positions
- `GET /api/v1/liquidation/history` - Get liquidation history
- `GET /api/v1/liquidation/bonus` - Get liquidation bonus
- `POST /api/v1/liquidation/liquidate` - Perform liquidation (auth required)

#### Market Data

- `GET /api/v1/market/overview` - Get market overview
- `GET /api/v1/market/tokens` - Get tokens market data

#### Health Check

- `GET /api/v1/health` - Health check endpoint

### Authentication

The API uses JWT for authentication. Include the token in your request header:

```
Authorization: Bearer <your_jwt_token>
```

### Authentication Workflow with Anvil

For local development, you can use Anvil's private keys to sign authentication messages. Here's a complete workflow:

#### 1. Start Anvil with a Specific Mnemonic

```bash
anvil --mnemonic "test test test test test test test test test test test junk"
```

This ensures consistent addresses and private keys across restarts.

#### 2. Get a Nonce from the API

Request a nonce for the Ethereum address you want to authenticate:

```bash
curl -X 'GET' \
 'http://localhost:8080/api/v1/users/nonce/0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266' \
 -H 'accept: application/json'
```

Response:

```json
{
"message": "Sign this message to verify you are the owner of 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266. Nonce: 123456",
"nonce": "123456"
}
```

#### 3. Sign the Nonce Message

Use the `cast` command-line tool from Foundry to sign the message with the corresponding private key:

```bash
cast sign --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 "Sign this message to verify you are the owner of 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266. Nonce: 123456"
```

This will output a signature string.

#### 4. Authenticate with the Signature

Send the address and signature to the authentication endpoint:

```bash
curl -X 'POST' \
 'http://localhost:8080/api/v1/users/auth' \
 -H 'accept: application/json' \
 -H 'Content-Type: application/json' \
 -d '{
"address": "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
"signature": "YOUR_SIGNATURE_FROM_STEP_3"
}'
```

The response will include your JWT token:

```json
{
"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
"expires_at": "2025-05-05T10:13:45Z"
}
```

#### 5. Use the JWT Token for Authenticated Requests

Include the token in the Authorization header for protected endpoints:

```bash
curl -X 'GET' \
 'http://localhost:8080/api/v1/borrowing/balance' \
 -H 'accept: application/json' \
 -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...'
```

#### Default Anvil Accounts

For convenience, here are the first few accounts from Anvil's default mnemonic:

| Account    | Address                                    | Private Key                                                        |
| ---------- | ------------------------------------------ | ------------------------------------------------------------------ |
| Account #0 | 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 | 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 |
| Account #1 | 0x70997970C51812dc3A010C7d01b50e0d17dc79C8 | 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d |
| Account #2 | 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC | 0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a |
| Account #3 | 0x90F79bf6EB2c4f870365E785982E1f101E93b906 | 0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6 |
| Account #4 | 0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65 | 0x47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a |

## Smart Contract Integration

The backend interacts with several smart contracts:

1. **Token**: ERC-20 token contract for the platform's native token
2. **LendingPool**: Manages deposits and withdrawals
3. **Borrowing**: Manages borrowing operations
4. **Collateral**: Manages collateral deposits and withdrawals

### Contract ABIs

Contract ABIs are stored in `pkg/blockchain/contracts/` and are used to generate Go bindings via abigen.

### Generating Contract Bindings

After updating contract ABIs, regenerate the Go bindings:

```bash
cd backend
./scripts/abigen/generate.sh
```

## Database Schema

The application uses PostgreSQL with the following key models:

1. **User**: Represents users of the platform
2. **Transaction**: Records all financial transactions
3. **Position**: Represents lending, borrowing, and collateral positions

### Migrations

Database migrations are handled automatically on startup using GORM's auto-migration feature. The migration code is in `pkg/database/migrations.go`.

## Development Workflow

### Hot Reload Development

The development setup includes hot reload via Air:

1. First deploy smart contracts with Foundry:

```bash
# Start local blockchain
anvil --mnemonic "test test test test test test test test test test test junk"

# In another terminal, deploy contracts
cd contracts
forge script script/Deploy.s.sol --fork-url http://localhost:8545 --broadcast
```

2. Update `.env.dev` with contract addresses from Foundry output

3. Run development environment:

```bash
docker-compose --env-file ./backend/.env.dev -f docker-compose.dev.yml up -d
```

4. Code changes will automatically trigger recompilation and server restart

### Adding New API Endpoints

1. Create DTO in `internal/api/dto/`
2. Add handler function in `internal/api/handlers/`
3. Register route in `internal/api/routes/`
4. Update Swagger annotations
5. Run `swag init` to update documentation

### Middleware

The application includes several middleware components:

- **Authentication**: JWT-based auth (`middleware/auth.go`)
- **CORS**: Cross-Origin Resource Sharing (`middleware/cors.go`)
- **Error Handler**: Centralized error handling (`middleware/error_handler.go`)
- **Logger**: Request logging (`middleware/logger.go`)

## Deployment

### Production Deployment

1. Create a production environment file:

```bash
cp backend/.env.example backend/.env.prod
# Edit .env.prod with production settings
```

2. Deploy using Docker Compose:

```bash
docker-compose --env-file ./backend/.env.prod -f docker-compose.yml up -d
```

### Container Security Considerations

- Use non-root users in containers
- Implement proper secret management (e.g., Docker secrets)
- Keep dependencies updated
- Scan images for vulnerabilities

## Testing

### Running Tests

Run unit tests:

```bash
cd backend
go test ./...
```

Run specific test:

```bash
go test github.com/Mattouff/Lending-Borrowing/internal/api/handlers -v
```

## Troubleshooting

### Common Issues

1. **Database Connection Errors**

   - Check PostgreSQL container is running: `docker ps`
   - Verify connection settings in `.env` file
   - Check network configuration between containers

2. **Blockchain Connection Issues**

   - Ensure Anvil is running with the correct mnemonic
   - Verify contract addresses in `.env` file
   - Check RPC URL and port forwarding

3. **CORS Issues**

   - Update CORS configuration in `middleware/cors.go`
   - For authentication with credentials, specific origins must be set

4. **Docker Issues**
   - Check logs: `docker logs lending-borrowing-api-dev`
   - Rebuild containers: `docker-compose -f docker-compose.dev.yml up --build -d`

### Logs

Access container logs:

```bash
docker logs lending-borrowing-api-dev
```

---

## Maintenance and Support

For questions or issues, please contact the development team or create an issue in the repository.

---

_This documentation is for internal use only._
