<div align="center">

<h1 style="font-size: 3em; text-decoration: none; border-bottom: none;">🦟 모기</h1>

[![Status](https://img.shields.io/badge/Status-Development-orange?style=flat-square)](https://github.com/revil/mogi)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=flat-square)](https://github.com/revil/mogi)
[![Tests](https://img.shields.io/badge/Tests-Passing-brightgreen?style=flat-square)](https://github.com/revil/mogi)
[![Version](https://img.shields.io/badge/Version-1.0.0-blue?style=flat-square)](https://github.com/revil/mogi)
[![License](https://img.shields.io/badge/License-MIT-yellow?style=flat-square)](LICENSE)
[![Discord](https://img.shields.io/badge/Discord-7289DA?style=flat-square&logo=discord&logoColor=white)](https://discord.gg/YOUR_INVITE_CODE)

</div>

## ⚠️ 꼭 읽어주세요
그것을 위한 서드파티 서비스입니다.
<br>
언제 사라질지 모르는 토이 프로젝트이며 라이센스에 맞게 이용하시면 됩니다.

문의/건의사항은 [디스코드](https://discord.gg/YOUR_INVITE_CODE)에 남겨주세요.

>서비스 이용을 위해 오신 분들은 [가이드](https://example.com/guide)를 읽고 사용하시면 됩니다.
<br>
사용 시 불이익에 대해선 아무 책임을 지지 않습니다.

⭐ **코드 및 서비스에 대한 피드백은 언제나 환영입니다.**


## 🚀 로컬 개발

### 환경 설정

```bash
# 초기 세팅
./scripts/setup.sh

# 의존성 설치
bun install
go mod tidy

# 컨테이너 시작
docker-compose up -d
```

### 실행

```bash
bun run dev:api-server
bun run dev:scraping-server
bun run dev:website

make watch-suction-server
make watch-suction-client
```

> 실시간 패킷을 가져오기 위해서는 suction-client는 Network Bridge 모드가 활성화된 VMWare에서 실행해야합니다.
로컬 환경에서는 미리 캡처한 samples/*.pcap 파일로 개발합니다.


## 🛠️ 기술 스택

### Bun Frontend & Backend
[![Bun](https://img.shields.io/badge/Bun-1.2.0+-000000?style=flat-square&logo=bun)](https://bun.sh/)
[![Next.js](https://img.shields.io/badge/Next.js-15.4.0+-000000?style=flat-square&logo=next.js)](https://nextjs.org/)
[![NestJS](https://img.shields.io/badge/NestJS-11.1.0+-E0234E?style=flat-square&logo=nestjs)](https://nestjs.com/)
[![Socket.IO](https://img.shields.io/badge/Socket.IO-4.0+-010101?style=flat-square&logo=socket.io)](https://socket.io/)

### Golang Ecosystem
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![Go Fiber](https://img.shields.io/badge/Go%20Fiber-2.0+-00ADD8?style=flat-square&logo=go)](https://gofiber.io/)
[![Uber FX](https://img.shields.io/badge/Uber%20FX-1.0+-000000?style=flat-square&logo=go)](https://github.com/uber-go/fx)
[![Uber Zap](https://img.shields.io/badge/Uber%20Zap-1.0+-000000?style=flat-square&logo=go)](https://github.com/uber-go/zap)
[![gopacket](https://img.shields.io/badge/gopacket-1.0+-00ADD8?style=flat-square&logo=go)](https://github.com/google/gopacket)

### Protocols & Data
[![QUIC](https://img.shields.io/badge/QUIC%20Protocol-1.0+-000000?style=flat-square)](https://quicwg.org/)
[![Protocol Buffers](https://img.shields.io/badge/Protocol%20Buffers-3.0+-000000?style=flat-square&logo=protobuf)](https://developers.google.com/protocol-buffers)
[![Snappy Compression](https://img.shields.io/badge/Snappy%20Compression-1.0+-000000?style=flat-square)](https://github.com/golang/snappy)

### Infrastructure
#### Common
[![MongoDB](https://img.shields.io/badge/MongoDB-7.0+-47A248?style=flat-square&logo=mongodb)](https://www.mongodb.com/)
[![Redis Stack](https://img.shields.io/badge/Redis%20Stack-7.0+-DC382D?style=flat-square&logo=redis&logoColor=white)](https://redis.io/docs/stack/)

#### Local
[![Docker](https://img.shields.io/badge/Docker-20.10+-2496ED?style=flat-square&logo=docker)](https://www.docker.com/)

#### Production
[![Kubernetes](https://img.shields.io/badge/Kubernetes-326CE5?style=flat-square&logo=kubernetes&logoColor=white)](https://kubernetes.io/)

## 🎯 애플리케이션

### Website
수집 정보 및 사용자 정보 제공 사이트

### API Server
웹사이트를 위한 Backend for Frontend (BFF) REST API 및 WebSocket 서버

### Scraping Server
그 곳의 정보를 수집하는 서버

### Suction Server
클라이언트에서 수집한 데이터를 받아 API Server에서 사용할 수 있도록 재가공하는 서버

### Suction Client
그 것의 패킷을 캡쳐 후 분석하고 서버로 전달하는 클라이언트

### Architecture

```mermaid
graph LR
    subgraph "Frontend"
        Website[🌐 Website] 
    end
    
    subgraph "Backend Services"
        APIServer[🍞 API Server]
        ScrapingServer[🍞 Scraping Server]
        SuctionServer[🐹 Suction Server]
    end

    subgraph "External Sources"
        That[🎮 그 것]
        There[🌐 그 곳]
    end

    subgraph "Data Collection"
        SuctionClient[🐹 Suction Client]
    end
    
    subgraph "Infrastructure"
        MongoDB[🗄️ MongoDB]
        RedisStack[🔴 RedisStack]
    end

    That -->|패킷 캡처| SuctionClient
    SuctionClient -->|가공한 데이터 전송| SuctionServer
    SuctionServer -->|임시 데이터 저장 및<br>완료된 데이터 수신| RedisStack
    SuctionServer -->|완료된 데이터 저장| MongoDB
    There -->|스크래핑| ScrapingServer
    ScrapingServer -->|스크랩한 데이터 가공 후 저장| MongoDB
    Website <-->|HTTP, WebSocket| APIServer
    APIServer -->|정보 업데이트 요청, 수신| ScrapingServer
    APIServer -->|실시간 데이터| RedisStack
    APIServer -->|영구 데이터| MongoDB
```

## 🤝 같이 코딩 하실래요?

1. 레파지토리를 포크해주세요.
2. 목적에 맞는 이름으로 브랜치를 생성해주세요.(`git checkout -b {feature,bugfix,hotfix...}/amazing-feature`)
3. 작업 후 내역을 알 수 있는 메시지를 작성하여 커밋해주세요. (`git commit -m 'Add some amazing feature'`)
4. 작업 한 브랜치를 푸시한 후 (`git push origin feature/amazing-feature`)
5. PR을 생성해주세요.

## 📝 라이센스

이 프로젝트는 MIT 라이센스 하에 배포됩니다. 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.
