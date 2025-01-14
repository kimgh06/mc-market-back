# Maple

Mini API for mc-market

- Experimental Snowflake ID

## 환경 설정

postgresql 로컬 설정법:

```bash
docker-compose -f docker-compose.postgres.yml up -d
```

## 웹 서버 실행:

이 프로젝트에서는 golang을 사용해서 웹 서버를 구성합니다.

```bash
go run main.go
```

## 빌드

```bash
docker build . -t server
```
