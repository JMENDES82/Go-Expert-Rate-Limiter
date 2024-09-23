# Go Expert Rate Limiter

Este projeto implementa um **Rate Limiter** configurável em Go, que limita o número de requisições por segundo com base no **IP** ou em um **Token** de acesso. O rate limiter utiliza **Redis** como backend para armazenar as informações sobre as requisições, permitindo que as regras sejam aplicadas de maneira consistente e escalável.

## Funcionalidades

- **Limitação por IP**: Restringe o número de requisições recebidas de um único endereço IP dentro de um intervalo de tempo configurado.
- **Limitação por Token de Acesso**: Permite configurar limites específicos para tokens de acesso enviados no cabeçalho da requisição. As requisições com um token têm prioridade sobre o limite por IP.
- **Redis como backend**: O Redis é utilizado para armazenar as contagens de requisições e o tempo de expiração das chaves, garantindo consistência em um ambiente distribuído.
- **Middleware para fácil integração**: O rate limiter é implementado como middleware, permitindo que seja facilmente integrado a qualquer servidor HTTP em Go.

## Configuração e Execução

### Pré-requisitos

- **Go**: Certifique-se de ter o Go instalado. [Instruções de instalação](https://golang.org/doc/install)
- **Docker**: Certifique-se de ter o Docker e o Docker Compose instalados. [Instruções de instalação](https://docs.docker.com/get-docker/)

### Passos para Executar

1. **Clone o repositório**:

   ```bash
   git clone https://github.com/JMENDES82/Go-Expert-Rate-Limiter.git
   cd Go-Expert-Rate-Limiter
   ```

2. **Crie um arquivo `.env`** na raiz do projeto para definir as configurações de limite e conexão com o Redis. Aqui está um exemplo de um arquivo `.env`:

   ```bash
   REDIS_ADDRESS=redis:6379
   MAX_REQUESTS_IP=10        # Limite de 10 requisições por segundo por IP
   MAX_REQUESTS_TOKEN=50     # Limite de 50 requisições por segundo por Token
   BLOCK_TIME=1m             # Tempo de bloqueio após exceder o limite
   ```


3. **Construa e execute o projeto com Docker Compose**:

   ```bash
   docker-compose up --build
   ```

Isso iniciará o servidor web na porta `8080` e o Redis na porta `6379`. O Redis será usado para armazenar os contadores de requisições e os tempos de expiração.

### Testando o Rate Limiter

Agora você pode testar o rate limiter utilizando o `curl` ou uma ferramenta como o **Postman**.

#### 1. Teste de Limitação por IP:

Envie múltiplas requisições sem um token de acesso para testar o limite por IP:

    ```bash
    for i in {1..12}; do curl -i http://localhost:8080/; done
    ```    

Se o limite estiver configurado para 10 requisições por segundo, a 11ª e 12ª requisições devem ser bloqueadas com o código **429 Too Many Requests**.

#### 2. Teste de Limitação por Token:

Envie requisições com um token de acesso no cabeçalho para testar o limite por token:

    ```bash
    for i in {1..51}; do curl -i -H "API_KEY: abc123" http://localhost:8080/; done
    ```

Aqui, se o limite por token for 50 requisições por segundo, a 51ª requisição será bloqueada.

### Executando os Testes Automatizados

Para garantir que o rate limiter esteja funcionando conforme o esperado, há um conjunto de testes automatizados que você pode executar com o Go.

1. **Instale as dependências**:
   Certifique-se de que as dependências do projeto estão instaladas corretamente. Se necessário, execute:

    ```bash
    go mod tidy
    ```
2. **Execute os testes**:
   Use o seguinte comando para rodar todos os testes automatizados:

   ```bash
   go test ./... -v
   ```

O comando acima vai executar todos os testes da aplicação, incluindo os testes de integração do rate limiter com o Redis.

Os testes verificam se o rate limiter está funcionando corretamente, limitando requisições tanto por **IP** quanto por **Token**, e verificam o comportamento de bloqueio e recuperação após o tempo configurado expirar.

### Comportamento do Redis

- Cada IP ou Token de acesso recebe uma chave no Redis, como `limiter:<IP>` ou `limiter:<TOKEN>`.
- As chaves têm um TTL de 1 segundo, o que significa que o contador é resetado a cada segundo.
- Se o número de requisições exceder o limite configurado, o rate limiter responde com um **HTTP 429 (Too Many Requests)**.

### Estrutura do Projeto



