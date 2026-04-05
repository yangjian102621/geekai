# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

### Go Backend (api/)
- **Development**: `cd api && go run main.go` (uses config.toml)
- **Build**: `cd api && make` (builds both amd64 and arm64 binaries)
- **Individual builds**: `make amd64` or `make arm64`
- **Clean**: `make clean`
- **Config**: Copy `config.sample.toml` to `config.toml` and configure

### Web Frontend (web/)
- **Development**: `cd web && npm run dev` (runs on Vite dev server with --host)
- **Build**: `cd web && npm run build`
- **Lint**: `cd web && npm run lint` (ESLint with auto-fix)

### Testing
- Backend tests: `cd api/test && bash run_crawler_test.sh`
- No specific frontend test configuration found

## Project Architecture

### Backend (Go)
- **Framework**: Gin web framework with dependency injection via uber-go/fx
- **Database**: GORM with MySQL, Redis for caching, LevelDB for local storage
- **Authentication**: JWT tokens with Redis session storage
- **Middleware**: CORS, authorization, parameter handling, static resource serving
- **Structure**:
  - `handler/`: HTTP request handlers (REST API endpoints)
  - `service/`: Business logic services (AI integrations, payments, etc.)
  - `store/`: Database models and data access layer
  - `core/`: Application server and middleware configuration
  - `utils/`: Utility functions and helpers

### Frontend (Vue.js)
- **Framework**: Vue 3 with Composition API
- **UI Components**: Element Plus + Vant (mobile components)
- **State Management**: Pinia
- **Routing**: Vue Router with nested routes
- **Build Tool**: Vite
- **CSS**: Stylus preprocessor with Tailwind CSS utilities
- **Features**: Responsive design (desktop/mobile views), theme switching (dark/light)

### Key Features
- **AI Chat**: Multiple chat models and conversation management
- **Image Generation**: MidJourney, Stable Diffusion, DALL-E integration
- **Audio/Video**: Suno music creation, Luma/KeLing video generation
- **User Management**: Authentication, payments, power logs, invitations
- **Admin Panel**: Comprehensive management interface

### Database Models
Key entities: User, ChatItem, ChatMessage, ChatRole, ChatModel, Order, Product, AdminUser, and various job types for AI services.

### API Structure
- User APIs: `/api/user/*` (auth, profile, settings)
- Chat APIs: `/api/chat/*` (conversations, messages)
- AI Service APIs: `/api/mj/*`, `/api/sd/*`, `/api/dall/*`, `/api/suno/*`, `/api/video/*`
- Admin APIs: `/api/admin/*` (management functions)

### Configuration
- Backend: TOML configuration file (`config.toml`)
- Database: MySQL with automatic migrations
- Services: Redis, various AI API integrations
- File Storage: Local, Aliyun OSS, MinIO, Qiniu options