# Maple Market API

Mini API for mc-market - 마인크래프트 마켓플레이스 백엔드 서비스

- Experimental Snowflake ID

## 프로젝트 개요

Maple Market API는 Go언어와 Gin 프레임워크로 구현된 마인크래프트 마켓플레이스의 백엔드 서비스입니다. 
RESTful API와 GraphQL을 지원하며, PostgreSQL 데이터베이스를 사용합니다.

## 주요 기능

### 🔐 사용자 관리 (User Management)
- 사용자 등록/인증 (Surge 서비스 연동)
- 세션 관리 및 JWT 토큰 기반 인증
- 사용자 프로필 관리 (닉네임, 아바타)
- 권한 시스템 (관리자, 일반 사용자)
- 사용자 수익 관리 및 캐시 시스템

### 📦 상품 관리 (Products)
- 상품 CRUD 기능 (생성, 조회, 수정, 삭제)
- 상품 카테고리 및 태그 시스템
- 상품 검색 및 필터링 (가격, 카테고리, 키워드)
- 상품 이미지 및 파일 업로드
- 상품 구매 시스템
- 상품 버전 관리 시스템
- 상품 업데이트 로그

### 📝 게시글 시스템 (Articles)
- 게시글 CRUD 기능
- 말머리(Head) 기능으로 게시글 분류
- 조회수 시스템
- 이미지 업로드 지원
- HTML 템플릿 지원
- 웹훅 알림 기능

### 💬 댓글 시스템 (Comments)
- 댓글 CRUD 기능
- 대댓글(답글) 지원
- 페이지네이션 (30개씩)
- 실시간 댓글 수 카운팅

### 👍 추천/비추천 시스템 (Article Likes)
- 게시글 추천/비추천 기능
- 중복 추천 방지
- 추천/비추천 통계

### 💳 결제 시스템 (Payments)
- 결제 생성 및 승인
- 결제 내역 조회
- 토스페이먼츠 연동

### 🎨 배너/광고 관리 (Banner/Adcard)
- 배너 이미지 관리
- 광고 카드 시스템
- 이미지 업로드 및 URL 관리

## 기술 스택

### Backend
- **언어**: Go 1.22.5
- **웹 프레임워크**: Gin
- **데이터베이스**: PostgreSQL
- **ORM**: SQLC (타입 안전한 SQL 코드 생성)
- **API**: REST API, GraphQL (gqlgen)
- **인증**: JWT, JWKS

### 주요 라이브러리
- **HTTP 라우팅**: Gin
- **데이터베이스**: PostgreSQL, pgx/v5 드라이버
- **마이그레이션**: golang-migrate
- **ID 생성**: Snowflake ID
- **이미지 처리**: imaging
- **설정 관리**: envconfig, godotenv
- **로깅**: logrus
- **CLI**: cobra

## 프로젝트 구조

```
├── cmd/                    # CLI 명령어
│   ├── generate.go        # 코드 생성 명령어
│   ├── migrate.go         # DB 마이그레이션 명령어
│   ├── root.go           # 루트 명령어
│   └── serve.go          # 서버 실행 명령어
├── graph/                 # GraphQL 관련
│   ├── generated.go      # GraphQL 생성 코드
│   ├── resolver.go       # GraphQL 리졸버
│   └── schema/           # GraphQL 스키마
├── internal/              # 내부 패키지
│   ├── api/              # API 핸들러
│   │   ├── articles/     # 게시글 관리
│   │   ├── comment/      # 댓글 관리
│   │   ├── user/         # 사용자 관리
│   │   ├── products/     # 상품 관리
│   │   ├── payments/     # 결제 관리
│   │   ├── banner/       # 배너 관리
│   │   └── ...
│   ├── conf/             # 설정 관리
│   ├── middlewares/      # 미들웨어
│   ├── schema/           # 데이터베이스 스키마 (SQLC 생성)
│   ├── storage/          # 스토리지 관리
│   └── surge/            # Surge API 연동
├── pkg/                   # 공용 패키지
│   ├── files/            # 파일 업로드
│   ├── generate/         # 코드 생성 도구
│   └── permissions/      # 권한 관리
├── schema/               # 데이터베이스
│   ├── migrations/       # 마이그레이션 파일
│   └── queries/          # SQL 쿼리 파일
└── .storage/             # 파일 저장소
    ├── contents/         # 컨텐츠 파일
    └── images/           # 이미지 파일
```

## 환경 설정

### 필수 요구사항
- Go 1.22.5+
- Docker & Docker Compose
- PostgreSQL

### 환경 변수 설정
프로젝트 루트에 `.env` 파일을 생성하고 다음 환경 변수들을 설정하세요:

```bash
# 데이터베이스
MAPLE_DATABASE_URL=postgres://user:password@localhost:5432/maple

# 스토리지
MAPLE_STORAGE_IMAGES_PATH=./.storage/images
MAPLE_STORAGE_CONTENTS_PATH=./.storage/contents

# API
MAPLE_API_HOST=0.0.0.0:3000

# 인증
MAPLE_AUTH_KEYS_ENDPOINT=https://your-auth-server.com/.well-known/jwks.json

# Surge API
MAPLE_SURGE_URL=https://your-surge-api.com
MAPLE_SURGE_SERVICE_KEY=your-service-key

# 로깅
MAPLE_LOGGING_DEBUG=true
MAPLE_LOGGING_REQUESTS=true

# 프로덕션
MAPLE_PRODUCTION=false
```

## 개발 환경 실행

### PostgreSQL 로컬 설정

postgresql 로컬 설정법:

```bash
# PostgreSQL 컨테이너 시작
make dev-postgres-standalone

# 또는 직접 Docker Compose 사용
docker-compose -f docker-compose.postgres.yml up -d
```

### 개발 서버 실행

```bash
# Makefile을 사용한 빌드 및 실행
make dev-run

# 또는 직접 실행
go run main.go serve

# 기본 실행 (마이그레이션 후 서버 시작)
go run main.go
```

### 코드 생성

```bash
# SQLC 및 GraphQL 코드 생성
make dev-generate

# 또는 직접 실행
go run main.go generate
```

### Docker 개발 환경

```bash
# 개발용 Docker 환경 시작
make dev-docker

# 개발 로그 확인
make dev-logs
```

## 프로덕션 배포

### Docker 빌드

```bash
docker build . -t maple-market-api

# 또는 서버 태그로 빌드
docker build . -t server
```

### Docker Compose 배포

```bash
# VPS 환경
docker-compose -f docker-compose.vps.yml up -d

# Dokploy 환경
docker-compose -f docker-compose.dokploy.yml up -d

# 프로덕션 업데이트
git pull && docker-compose down && docker system prune -af && docker-compose up --build -d
```

## API 엔드포인트

### REST API Base URL: `/v1`

#### 사용자 관리 (`/v1/user/`)
- `POST /user/` - 사용자 생성
- `GET /user/:id/` - 사용자 조회
- `GET /user/session/` - 세션 사용자 조회 (인증 필요)
- `POST /user/:id/` - 사용자 정보 수정 (인증 필요)
- `GET /user/revenues/` - 수익 조회 (인증 필요)

#### 상품 관리 (`/v1/products/`)
- `GET /products/` - 상품 목록 조회
- `GET /products/:id/` - 상품 상세 조회
- `POST /products/` - 상품 생성 (인증 필요)
- `POST /products/:id/` - 상품 수정 (인증 필요)
- `DELETE /products/:id/` - 상품 삭제 (인증 필요)
- `POST /products/:id/purchase/` - 상품 구매 (인증 필요)

#### 게시글 관리 (`/v1/articles/`)
- `GET /articles/` - 게시글 목록 조회
- `GET /articles/:id/` - 게시글 상세 조회
- `POST /articles/` - 게시글 생성 (인증 필요)
- `POST /articles/:id/` - 게시글 수정 (인증 필요)
- `DELETE /articles/:id/` - 게시글 삭제 (인증 필요)

#### 댓글 관리 (`/v1/comment/:article_id/`)
- `GET /comment/:article_id/` - 댓글 목록 조회 (페이지네이션: `?page=1`)
- `GET /comment/:article_id/:comment_id/` - 댓글 상세 조회
- `POST /comment/:article_id/` - 댓글 생성 (인증 필요)
- `DELETE /comment/:article_id/:comment_id/` - 댓글 삭제 (인증 필요)

#### 추천/비추천 (`/v1/article_likes/:article_id/`)
- `GET /article_likes/:article_id/` - 추천/비추천 통계 조회
- `POST /article_likes/:article_id/` - 추천/비추천 등록 (인증 필요)

#### 결제 관리 (`/v1/payments/`)
- `GET /payments/` - 결제 목록 조회 (인증 필요)
- `POST /payments/` - 결제 생성 (인증 필요)
- `POST /payments/:orderId/approve/` - 결제 승인 (인증 필요)

### GraphQL
- **GraphQL Endpoint**: `/graphql`
- **GraphQL Playground**: 개발 환경에서 사용 가능

### 상태 확인
- `GET /status/` - 서버 상태 및 버전 정보

## 데이터베이스 스키마

### 주요 테이블

#### Users (사용자)
```sql
create table users (
    id         bigint primary key,
    nickname   varchar(255),
    permissions integer default 0,
    cash       integer default 0,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);
```

#### Products (상품)
```sql
create table products (
    id             bigint primary key,
    creator        bigint not null references users(id),
    name           varchar(255) not null,
    description    text,
    category       varchar(100),
    price          integer not null,
    price_discount integer,
    created_at     timestamp with time zone default now(),
    updated_at     timestamp with time zone default now()
);
```

#### Articles (게시글)
```sql
create table articles (
    id      bigint primary key,
    title   varchar(255) not null,
    content text not null,
    author  bigint not null references users(id),
    head    varchar(50),
    views   bigint default 0,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);
```

#### Comments (댓글)
```sql
create table comments (
    id         bigint primary key,
    article_id bigint not null references articles(id),
    user_id    bigint not null references users(id),
    reply_to   bigint references comments(id),
    content    text not null,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);
```

## 개발 도구

### Makefile 명령어
- `make dev-run` - 로컬 서버 실행
- `make dev-generate` - 코드 생성 (SQLC, GraphQL)
- `make dev-docker` - Docker 개발 환경 시작
- `make dev-postgres-standalone` - PostgreSQL 컨테이너 시작
- `make dev-postgres-standalone-stop` - PostgreSQL 컨테이너 중지
- `make dev-postgres-standalone-reset` - PostgreSQL 컨테이너 재생성
- `make dev-logs` - 개발 로그 확인

### CLI 명령어
```bash
# 서버 실행
go run main.go serve

# 데이터베이스 마이그레이션
go run main.go migrate

# 코드 생성
go run main.go generate

# 새 마이그레이션 생성
make migrate-new name=migration_name
```

### 테스트
HTTP 요청 테스트는 `test.http` 파일을 사용하세요.

## 설정 파일

- `gqlgen.yml` - GraphQL 코드 생성 설정
- `sqlc.yml` - SQLC 설정 (SQL을 Go 코드로 생성)
- `go.mod` - Go 모듈 의존성
- `Makefile` - 개발 명령어 모음

## 특별한 기능

### Snowflake ID
- 분산 시스템에서 고유 ID 생성
- 시간 기반 정렬 가능
- 64비트 정수형 ID

### Surge API 연동
- 외부 인증 서비스 연동
- 사용자명 해석 및 검증
- JWT 토큰 기반 인증

### 권한 시스템
- 비트 플래그 기반 권한 관리
- 관리자 권한 (2147483647)
- 세분화된 권한 제어

### 파일 업로드
- 이미지 및 컨텐츠 파일 업로드
- 안전한 파일 저장
- 파일 접근 권한 제어

## 문제 해결

### 로그 확인
```bash
# 개발 로그
make dev-logs

# 프로덕션 로그 (nohup 사용시)
tail -f nohup.out
```

### 데이터베이스 문제
```bash
# PostgreSQL 컨테이너 재시작
make dev-postgres-standalone-reset

# 마이그레이션 재실행
go run main.go migrate
```

### 코드 생성 문제
```bash
# 코드 재생성
make dev-generate
```

## 버전 정보

현재 버전: `development44.1`

## 라이센스

이 프로젝트는 실험적인 Snowflake ID 구현을 포함하고 있습니다.

---

## 개발 팀을 위한 추가 정보

### 코드 컨벤션
- Go 표준 포맷팅 사용
- SQLC를 통한 타입 안전한 데이터베이스 접근
- 구조화된 에러 핸들링 (`perrors` 패키지)
- 미들웨어 기반 인증 및 권한 관리

### 배포 전 체크리스트
1. 환경 변수 설정 확인
2. 데이터베이스 마이그레이션 실행
3. 코드 생성 및 빌드 테스트
4. API 엔드포인트 테스트
5. 로그 레벨 및 프로덕션 모드 설정
