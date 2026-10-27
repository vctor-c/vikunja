# Executar esta versão com Docker

Requisitos: Docker Engine ou Docker Desktop com Docker Compose v2.
O `Dockerfile.local` compila o frontend e o backend deste diretório, incluindo as alterações locais. Go, Node e Mage são instalados durante o build da imagem.

## Iniciar

O arquivo local `.env.docker` já foi criado com um segredo aleatório. Ele não é versionado nem enviado para o contexto de build. Em uma nova cópia do repositório, crie-o antes de iniciar (não substitua um arquivo já existente):

```sh
(umask 077; printf 'VIKUNJA_SERVICE_SECRET=%s\n' "$(openssl rand -hex 32)" > .env.docker)
```

Na raiz do repositório:

```sh
docker compose up -d --build
```

Acesse http://localhost:3456 e crie sua conta. O primeiro build baixa dependências e pode demorar. O banco SQLite e os anexos ficam no volume `vikunja-data` do Compose, montado em `/data` e gravável pelo usuário 1000 do contêiner.

## Verificar e atualizar

```sh
docker compose ps
docker compose logs -f vikunja
```

Depois de alterar o código, execute novamente `docker compose up -d --build`.
Para parar mantendo os dados:

```sh
docker compose down
```

Não use `docker compose down -v` se quiser manter os dados: essa opção remove o volume. Preserve também `.env.docker` para manter o segredo entre reinicializações.

## Acesso pela rede

Por padrão, a porta 3456 é publicada em todas as interfaces do host (`0.0.0.0:3456`). No outro servidor, configure a URL pública com o IP ou domínio usado para acessar a aplicação. Você pode persistir `VIKUNJA_PUBLIC_URL=http://IP_DO_SERVIDOR:3456/` em um arquivo `.env` ao lado do Compose, ou informar no comando:

```sh
VIKUNJA_PUBLIC_URL=http://192.168.1.50:3456/ \
docker compose up -d --build
```

Acesse `http://IP_DO_SERVIDOR:3456`. Se o servidor tiver firewall ou regras de rede, permita conexões TCP à porta 3456. Para restringir a publicação a localhost, configure `VIKUNJA_BIND_ADDRESS=127.0.0.1`.

Ao alterar a porta, ajuste as duas variáveis:

```sh
VIKUNJA_PORT=8080 VIKUNJA_PUBLIC_URL=http://192.168.1.50:8080/ docker compose up -d
```

O Dockerfile original de releases permanece disponível. Este Compose usa explicitamente `Dockerfile.local`.

Referência: [build no Docker Compose](https://docs.docker.com/reference/compose-file/build/).
