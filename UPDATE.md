## 추가 정보

### 댓글 기능

게시글에 대한 댓글 crud 기능 추가
게시글 조회 시 게시글에 맞는 댓글 제공
답글인 댓글은 따로 나오지 않고 해당 댓글에 리스트로 같이 나옴.
comments 테이블 정보:

```sql
create table comments
(
    id         bigint                   not null primary key,
    article_id bigint                   not null references articles (id),
    user_id    bigint                   not null references users (id),
    reply_to   bigint                   null references comments (id),
    content    text                     not null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);
```

엔드포인트 시작 : /v1/comment/ (슬래시 꼭 붙여야 함.)
쿼리 스트링으로 댓글 30개씩 페이지네이션 기능 추가 ?page=int
작업 디렉토리: internal/api/comment

### 조회수 기능

테이블 칼럼을 추가하여 특정 게시글을 조회하는 api를 호출할 때마다 조회수를 1씩 늘어 나도록 설정.

```sql
alter table articles
    add column views bigint not null default 0;
```

### 추천 / 비추천 기능

## 업데이트 정보

### 게시글 관련

게시글을 여러 개를 불러올 떄 해당 게시글에 이미지가 포함되어 있는 지의 유무를 전송하도록 함.

### 말머리 기능
