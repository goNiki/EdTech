# 🚀 EdTech: Полный гайд по деплою с нуля

> Этот гайд написан для человека **без опыта DevOps**. Каждая команда объяснена. Каждый шаг пронумерован. Ничего не пропускай, иди строго по порядку.

## ⏱ Реалистичная оценка времени

| Этап | Время | Что делаем |
|---|---|---|
| **Этап 0:** Покупка VPS | 10–15 мин | Регистрация, оплата, получение IP |
| **Этап 1:** Настройка сервера | 20–30 мин | SSH, безопасность, Docker |
| **Этап 2:** Доменное имя и DNS | 10–15 мин | Привязка домена к IP |
| **Этап 3:** Подготовка файлов | 15–20 мин | Создание config-файлов для продакшена |
| **Этап 4:** Сборка и запуск | 15–25 мин | Git clone, docker compose up |
| **Этап 5:** SSL-сертификат | 5–10 мин | Let's Encrypt (бесплатный HTTPS) |
| **Этап 6:** Бэкапы | 10 мин | Настройка автоматического резервного копирования |
| **Итого** | **~1.5–2 часа** | Рабочий проект на HTTPS |

> [!TIP]
> Если где-то застрянешь — не паникуй. Первый раз всегда медленнее. Реалистичный worst-case с гуглением ошибок: **3–4 часа**.

---

## Этап 0: Покупка VPS (10–15 минут)

### Что покупать
Рекомендую **Timeweb Cloud** (РФ, оплата российскими картами, быстрый саппорт) или **Hetzner** (дешевле, но зарубежные карты).

**Конфигурация:**
- **ОС:** Ubuntu 24.04 LTS
- **CPU:** 2 vCPU
- **RAM:** 4 GB
- **Диск:** 40–50 GB NVMe SSD
- **Трафик:** Безлимит или 10–20 ТБ

**Цена:** ~600–900 ₽/мес (Timeweb) или ~4.5 €/мес (Hetzner CX22).

### Результат этого этапа
После оплаты хостер пришлёт тебе:
- **IP-адрес** сервера (например, `185.200.100.50`)
- **Пароль root** (либо SSH-ключ, если ты его указал)

---

## Этап 1: Настройка сервера (20–30 минут)

### 1.1. Подключение к серверу по SSH

Открой **PowerShell** (или Windows Terminal) на своём компе и набери:

```powershell
ssh root@185.200.100.50
```

> Замени `185.200.100.50` на **твой реальный IP**, который прислал хостер.

Первый раз система спросит:
```
Are you sure you want to continue connecting (yes/no)?
```
Набери `yes` и нажми Enter. Потом введи пароль (символы **не отображаются** при вводе — это нормально).

### 1.2. Обновление системы

```bash
apt update && apt upgrade -y
```

> **Что это делает:** обновляет список пакетов и устанавливает все обновления безопасности. Без этого система может иметь известные уязвимости.

### 1.3. Создание безопасного пользователя

Работать от `root` — плохая практика. Создаём обычного пользователя с правами sudo:

```bash
# Создать пользователя (введи пароль когда попросят, остальное можно пропустить Enter'ом)
adduser deploy

# Дать права администратора
usermod -aG sudo deploy
```

### 1.4. Установка Docker

Docker — это «контейнеризатор», который упаковывает твоё приложение со всеми зависимостями. Без него придётся вручную ставить Go, Node.js, PostgreSQL и всё настраивать — Docker делает это за тебя.

```bash
# Установка Docker Engine (официальный способ из документации Docker):
curl -fsSL https://get.docker.com | sh

# Даём пользователю deploy право запускать Docker без sudo
usermod -aG docker deploy
```

### 1.5. Настройка файрвола (UFW)

Файрвол — это «стена», которая блокирует все входящие подключения, кроме тех, которые ты явно разрешишь.

```bash
# Разрешить SSH (чтобы не потерять доступ к серверу!)
ufw allow OpenSSH

# Разрешить веб-трафик (HTTP и HTTPS)
ufw allow 80/tcp
ufw allow 443/tcp

# Включить файрвол
ufw enable
```

> [!CAUTION]
> **Обязательно** сначала разреши SSH (`ufw allow OpenSSH`), а потом включай файрвол! Если забудешь — файрвол заблокирует SSH и ты потеряешь доступ к серверу.

### 1.6. Переключение на пользователя deploy

```bash
# Выйти из root
exit

# Подключиться заново, но уже как deploy
ssh deploy@185.200.100.50
```

Теперь ты работаешь от безопасного пользователя. Для команд, требующих права администратора, используй `sudo` перед командой.

---

## Этап 2: Доменное имя и DNS (10–15 минут)

### 2.1. Зачем нужен домен?

Без домена сайт работает только по IP (например `http://185.200.100.50`). А для SSL-сертификата (HTTPS) **обязательно нужен домен**.

### 2.2. Где купить?

- **reg.ru** — ~200–500 ₽/год за домен `.ru`
- **Cloudflare Registrar** — дёшево, без наценок

Купи домен (например, `edtech-platform.ru`) и в DNS-панели регистратора создай **2 записи**:

| Тип | Имя | Значение | TTL |
|---|---|---|---|
| **A** | `@` | `185.200.100.50` | 300 |
| **A** | `www` | `185.200.100.50` | 300 |

> DNS-записи могут «проснуться» через 5–30 минут. Можешь проверить:
> ```bash
> ping edtech-platform.ru
> ```
> Если отвечает твоим IP — DNS заработал.

---

## Этап 3: Подготовка конфигурационных файлов (15–20 минут)

> [!IMPORTANT]
> Все файлы из этого этапа нужно **добавить в твой Git-репозиторий** на локальной машине, закоммитить и запушить. На сервере ты потом просто сделаешь `git clone`.

### 3.1. Список файлов, которые нужно создать

В корне проекта нужно создать / обновить:

```
EdTech/
├── Dockerfile                  # ← ОБНОВИТЬ (безопасная версия)
├── docker-compose.prod.yml     # ← СОЗДАТЬ
├── nginx.conf                  # ← СОЗДАТЬ
├── .env.prod.example           # ← СОЗДАТЬ (шаблон для продовых переменных)
├── deploy.sh                   # ← СОЗДАТЬ (скрипт быстрого деплоя)
├── backup.sh                   # ← СОЗДАТЬ (скрипт бэкапов)
└── frontend/
    ├── Dockerfile              # ← СОЗДАТЬ
    ├── .dockerignore           # ← СОЗДАТЬ
    └── next.config.ts          # ← ОБНОВИТЬ (добавить output: 'standalone')
```

> [!NOTE]
> Все эти файлы я описал ниже по тексту. Также они будут созданы в проекте отдельно.

### 3.2. `frontend/next.config.ts` — обновить

```typescript
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  allowedDevOrigins: ['172.19.0.1'],
  async rewrites() {
    return [
      {
        source: '/static/:path*',
        destination: 'http://localhost:8082/static/:path*',
      },
    ];
  },
};

export default nextConfig;
```

> Ключевое изменение: `output: "standalone"`. Без этого Next.js не создаст автономный серверный билд, и Docker-образ будет весить 1+ ГБ вместо 150 МБ.

### 3.3. `frontend/Dockerfile` — создать

```dockerfile
# 1. Зависимости
FROM node:20-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --ignore-scripts

# 2. Сборка
FROM node:20-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
ENV NEXT_TELEMETRY_DISABLED=1
RUN npm run build

# 3. Минимальный рантайм
FROM node:20-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs
EXPOSE 3000
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

CMD ["node", "server.js"]
```

### 3.4. `frontend/.dockerignore` — создать

```
node_modules
.next
.env*
*.tsbuildinfo
next-env.d.ts
```

### 3.5. `Dockerfile` (бэкенд, корень проекта) — обновить

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
ENV CGO_ENABLED=0 GOOS=linux

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-w -s" -o /bin/edtech ./cmd/edtech/main.go
RUN go build -ldflags="-w -s" -o /bin/migration ./cmd/migration/migration.go

FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata curl && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app
COPY --from=builder /bin/edtech /app/edtech
COPY --from=builder /bin/migration /app/migration
COPY --from=builder /app/migrators /app/migrators

RUN mkdir -p /app/uploads && chown -R appuser:appgroup /app

USER appuser
EXPOSE 8082

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8082/api/v1/categories || exit 1

CMD ["/app/edtech"]
```

### 3.6. `docker-compose.prod.yml` — создать

```yaml
services:
  db:
    image: postgres:15-alpine
    restart: unless-stopped
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER} -d ${DB_NAME}"]
      interval: 5s
      timeout: 5s
      retries: 5
    networks:
      - internal

  migrator:
    build:
      context: .
      dockerfile: Dockerfile
    env_file:
      - .env.prod
    environment:
      - DB_HOST=db
      - DB_PORT=5432
    command: ["/app/migration", "up"]
    depends_on:
      db:
        condition: service_healthy
    networks:
      - internal

  backend:
    build:
      context: .
      dockerfile: Dockerfile
    restart: unless-stopped
    env_file:
      - .env.prod
    environment:
      - DB_HOST=db
      - DB_PORT=5432
      - SERVER_HOST=0.0.0.0
      - SERVER_PORT=8082
    volumes:
      - uploads_data:/app/uploads
    depends_on:
      migrator:
        condition: service_completed_successfully
      db:
        condition: service_healthy
    networks:
      - internal

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
      args:
        - NEXT_PUBLIC_API_URL=https://${DOMAIN}/api/v1
    restart: unless-stopped
    environment:
      - NEXT_PUBLIC_API_URL=https://${DOMAIN}/api/v1
    depends_on:
      - backend
    networks:
      - internal

  nginx:
    image: nginx:alpine
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - uploads_data:/var/www/uploads:ro
      - certbot_www:/var/www/certbot:ro
      - certbot_certs:/etc/letsencrypt:ro
    depends_on:
      - backend
      - frontend
    networks:
      - internal

  # Контейнер для автоматического обновления SSL-сертификатов
  certbot:
    image: certbot/certbot
    volumes:
      - certbot_www:/var/www/certbot
      - certbot_certs:/etc/letsencrypt

volumes:
  postgres_data:
  uploads_data:
  certbot_www:
  certbot_certs:

networks:
  internal:
    driver: bridge
```

### 3.7. `nginx.conf` — создать

```nginx
worker_processes auto;

events {
    worker_connections 512;
}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;
    sendfile      on;
    keepalive_timeout 65;

    # Ограничение размера загружаемых файлов (50 МБ — под презентации)
    client_max_body_size 55M;

    # Gzip-сжатие (ускоряет загрузку страниц)
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml;
    gzip_min_length 256;

    # Rate Limiting: защита от брутфорса
    limit_req_zone $binary_remote_addr zone=api_limit:10m rate=20r/s;
    limit_req_zone $binary_remote_addr zone=login_limit:10m rate=5r/m;

    # Блок для проверки сертификатов Let's Encrypt
    server {
        listen 80;
        server_name _;

        # Certbot проверяет владение доменом через эту папку
        location /.well-known/acme-challenge/ {
            root /var/www/certbot;
        }

        # Всё остальное перенаправляем на HTTPS
        location / {
            return 301 https://$host$request_uri;
        }
    }

    # Основной HTTPS-сервер
    server {
        listen 443 ssl;
        server_name YOUR_DOMAIN_HERE;

        # SSL-сертификаты (появятся после запуска certbot)
        ssl_certificate     /etc/letsencrypt/live/YOUR_DOMAIN_HERE/fullchain.pem;
        ssl_certificate_key /etc/letsencrypt/live/YOUR_DOMAIN_HERE/privkey.pem;

        # Современные настройки SSL
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_prefer_server_ciphers off;

        # Статические файлы курсов: Nginx отдаёт напрямую с диска
        location /static/uploads/ {
            alias /var/www/uploads/;
            expires 7d;
            add_header Cache-Control "public, immutable";
            try_files $uri =404;
        }

        # API: проксируем на Go-бэкенд
        location /api/ {
            limit_req zone=api_limit burst=40 nodelay;
            proxy_pass http://backend:8082;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # Эндпоинт логина: более строгий rate-limit
        location /api/v1/auth/login {
            limit_req zone=login_limit burst=3 nodelay;
            proxy_pass http://backend:8082;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # Фронтенд (Next.js): всё остальное
        location / {
            proxy_pass http://frontend:3000;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
    }
}
```

### 3.8. `.env.prod.example` — создать

Это **шаблон**. На сервере ты скопируешь его в `.env.prod` и заполнишь реальными значениями.

```env
# === Домен ===
DOMAIN=edtech-platform.ru

# === Сервер ===
SERVER_HOST=0.0.0.0
SERVER_PORT=8082
SERVER_TIMEOUT=4s
SERVER_IDLETIMEOUT=60s
CORS_ALLOWED_ORIGINS=https://edtech-platform.ru,https://www.edtech-platform.ru

# === База данных ===
DB_HOST=db
DB_PORT=5432
DB_USER=edtech_user
DB_PASSWORD=CHANGE_ME_STRONG_PASSWORD_HERE
DB_NAME=edtech_db
DB_SSLMODE=disable
DB_MAXCONNS=10
DB_MINCONNS=2
DB_MAXCONNLIFETIME=1h
DB_MAXCONNIDLETIME=30m
DB_HEALTHCHECKPERIOD=1m

# === JWT (обязательно замени секрет!) ===
JWT_SECRET=CHANGE_ME_RANDOM_STRING_64_CHARS
JWT_ACCESSEXP=1h
JWT_REFRESHEXP=24h

# === Логгер ===
ENV=production
LOGGER_LEVEL=info
LOGGER_CONSOLE_ENABLED=true
LOGGER_CONSOLE_FORMAT=json
LOGGER_CONSOLE_OUTPUT=stdout
LOGGER_FILE_ENABLED=false
```

### 3.9. Коммит и push

На своей **локальной машине** (Windows), после создания всех файлов:

```powershell
cd C:\Users\gogol\OneDrive\Desktop\EdTech
git add -A
git commit -m "feat: add production deployment configs"
git push origin main
```

---

## Этап 4: Сборка и запуск на сервере (15–25 минут)

Теперь переключаемся на **сервер** (SSH от пользователя `deploy`).

### 4.1. Клонирование репозитория

```bash
cd ~
git clone https://github.com/goNiki/EdTech.git
cd EdTech
```

### 4.2. Создание файла с продовыми переменными

```bash
# Скопировать шаблон
cp .env.prod.example .env.prod

# Отредактировать (откроется простой текстовый редактор в терминале)
nano .env.prod
```

> **Внутри `nano`:**
> - Стрелками перемещайся по файлу
> - Замени `CHANGE_ME_STRONG_PASSWORD_HERE` на реальный пароль для БД
> - Замени `CHANGE_ME_RANDOM_STRING_64_CHARS` на длинную случайную строку
> - Замени `edtech-platform.ru` на твой реальный домен
> - Для сохранения: `Ctrl+O`, затем `Enter`
> - Для выхода: `Ctrl+X`

**Как сгенерировать надёжный пароль и JWT-секрет прямо на сервере:**
```bash
# Пароль для БД (24 символа):
openssl rand -base64 24

# JWT-секрет (64 символа):
openssl rand -hex 32
```

### 4.3. Замена домена в nginx.conf

```bash
# Заменить плейсхолдер ВАША_ДОМЕН на реальный домен (одной командой):
sed -i 's/YOUR_DOMAIN_HERE/edtech-platform.ru/g' nginx.conf
```

### 4.4. Первый запуск (без SSL — чтобы certbot мог проверить домен)

Перед получением SSL-сертификата нужно временно поднять Nginx только на порту 80.

Создай временный Nginx-конфиг:
```bash
cat > nginx-init.conf << 'EOF'
worker_processes auto;
events { worker_connections 512; }
http {
    server {
        listen 80;
        server_name _;
        location /.well-known/acme-challenge/ {
            root /var/www/certbot;
        }
        location / {
            return 200 'EdTech server is alive';
            add_header Content-Type text/plain;
        }
    }
}
EOF
```

Запусти Nginx с временным конфигом:
```bash
# Подставляем временный конфиг
cp nginx.conf nginx.conf.real
cp nginx-init.conf nginx.conf

# Запускаем только nginx и certbot volumes
docker compose -f docker-compose.prod.yml up -d nginx
```

---

## Этап 5: SSL-сертификат (5–10 минут)

### 5.1. Получение сертификата

```bash
docker compose -f docker-compose.prod.yml run --rm certbot certonly \
  --webroot \
  --webroot-path /var/www/certbot \
  -d edtech-platform.ru \
  -d www.edtech-platform.ru \
  --email your-email@gmail.com \
  --agree-tos \
  --no-eff-email
```

> Замени `edtech-platform.ru` на свой домен и `your-email@gmail.com` на свою почту.

Если видишь `Congratulations!` — сертификат получен!

### 5.2. Переключение на полный конфиг и запуск всего проекта

```bash
# Остановить временный nginx
docker compose -f docker-compose.prod.yml down

# Вернуть боевой конфиг
cp nginx.conf.real nginx.conf

# Запустить ВСЁ
docker compose -f docker-compose.prod.yml up -d --build
```

Первая сборка Docker-образов займёт **5–10 минут** (скачивание базовых образов + компиляция Go + сборка Next.js). Все последующие сборки будут быстрее благодаря кэшированию.

### 5.3. Проверка

```bash
# Посмотреть, все ли контейнеры запущены:
docker compose -f docker-compose.prod.yml ps

# Посмотреть логи бэкенда:
docker compose -f docker-compose.prod.yml logs backend --tail 50

# Посмотреть логи фронтенда:
docker compose -f docker-compose.prod.yml logs frontend --tail 50
```

Ожидаемый результат `ps`:
```
NAME         STATUS
db           Up (healthy)
backend      Up (healthy)
frontend     Up
nginx        Up
```

Теперь открой `https://edtech-platform.ru` в браузере — должен загрузиться сайт! 🎉

### 5.4. Автообновление SSL-сертификата (каждые 60 дней)

Let's Encrypt выдаёт сертификаты на 90 дней. Настрой автоматическое продление:

```bash
# Открыть планировщик задач
sudo crontab -e

# Добавить строку (обновление каждый день в 3:00 ночи):
0 3 * * * cd /home/deploy/EdTech && docker compose -f docker-compose.prod.yml run --rm certbot renew --quiet && docker compose -f docker-compose.prod.yml exec nginx nginx -s reload
```

---

## Этап 6: Настройка бэкапов (10 минут)

### 6.1. Создать скрипт бэкапа

```bash
sudo mkdir -p /opt/scripts
sudo nano /opt/scripts/backup_edtech.sh
```

Вставь:
```bash
#!/usr/bin/env bash
set -eo pipefail

TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_DIR="/var/backups/edtech"
COMPOSE_DIR="/home/deploy/EdTech"

mkdir -p "$BACKUP_DIR"

echo "[$TIMESTAMP] Starting backup..."

# Дамп базы данных
docker compose -f "$COMPOSE_DIR/docker-compose.prod.yml" exec -T db \
  pg_dump -U edtech_user -d edtech_db | gzip -9 > "$BACKUP_DIR/db_${TIMESTAMP}.sql.gz"

# Проверка что файл не пустой
if [ ! -s "$BACKUP_DIR/db_${TIMESTAMP}.sql.gz" ]; then
    echo "ERROR: Database backup is empty!" >&2
    exit 1
fi

# Бэкап загруженных файлов (uploads)
UPLOADS_VOL=$(docker volume inspect edtech_uploads_data -f '{{.Mountpoint}}' 2>/dev/null || true)
if [ -n "$UPLOADS_VOL" ] && [ -d "$UPLOADS_VOL" ]; then
    tar -czf "$BACKUP_DIR/uploads_${TIMESTAMP}.tar.gz" -C "$UPLOADS_VOL" .
fi

# Ротация: удалить бэкапы старше 7 дней
find "$BACKUP_DIR" -type f -name "*.gz" -mtime +7 -delete

echo "[$TIMESTAMP] Backup completed. Files in $BACKUP_DIR:"
ls -lh "$BACKUP_DIR"
```

```bash
sudo chmod +x /opt/scripts/backup_edtech.sh
```

### 6.2. Добавить в cron (автозапуск каждые 6 часов)

```bash
sudo crontab -e

# Добавить строку:
0 */6 * * * /opt/scripts/backup_edtech.sh >> /var/log/edtech_backup.log 2>&1
```

### 6.3. Проверка бэкапа (запустить вручную)

```bash
sudo /opt/scripts/backup_edtech.sh
```

Должен появиться файл в `/var/backups/edtech/`.

---

## 🧯 Шпаргалка по частым проблемам

### «Сайт не открывается после docker compose up»

```bash
# 1. Проверь статус контейнеров
docker compose -f docker-compose.prod.yml ps

# 2. Если backend упал — смотри логи:
docker compose -f docker-compose.prod.yml logs backend

# 3. Если фронтенд не поднялся:
docker compose -f docker-compose.prod.yml logs frontend

# 4. Если nginx показывает 502 Bad Gateway:
docker compose -f docker-compose.prod.yml logs nginx
```

### «Certbot говорит: Domain not reachable»

Значит DNS ещё не обновился. Подожди 10–15 минут и проверь:
```bash
dig +short edtech-platform.ru
```
Должен показать твой IP.

### «Как обновить код после изменений?»

```bash
cd ~/EdTech
git pull origin main
docker compose -f docker-compose.prod.yml up -d --build
```

### «Как восстановить базу из бэкапа?»

```bash
# Распаковать и залить дамп
gunzip -c /var/backups/edtech/db_20261004_120000.sql.gz | \
  docker compose -f docker-compose.prod.yml exec -T db \
  psql -U edtech_user -d edtech_db
```

### «Как посмотреть, сколько места занимает сервер?»

```bash
df -h           # Диск
free -h          # RAM
docker system df # Docker (образы, контейнеры, volumes)
```

---

## 📋 Чеклист перед запуском

- [ ] VPS куплен, SSH-доступ работает
- [ ] Пользователь `deploy` создан, файрвол настроен
- [ ] Домен куплен, DNS A-записи указывают на IP сервера
- [ ] Все файлы (Dockerfile, docker-compose.prod.yml, nginx.conf) закоммичены в Git
- [ ] `.env.prod` заполнен **реальными** паролями и доменом
- [ ] `YOUR_DOMAIN_HERE` в `nginx.conf` заменён на реальный домен
- [ ] SSL-сертификат получен через Certbot
- [ ] `docker compose up -d --build` запущен без ошибок
- [ ] Сайт открывается по `https://твой-домен.ru`
- [ ] Бэкап работает (`sudo /opt/scripts/backup_edtech.sh`)
