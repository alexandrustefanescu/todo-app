# Todo API Benchmark Guide

Complete guide for benchmarking all 5 Todo API implementations.

## Overview

This benchmark suite allows you to test and compare the performance of all Todo API implementations:
- Rust Actix
- Go FastHTTP
- Go Fiber
- Node.js Fastify
- Python FastAPI

## Prerequisites

### Python Dependencies

```bash
pip install httpx
```

Or install from requirements:

```bash
pip install -r benchmark-requirements.txt
```

### Running All APIs

You need to run all 5 APIs on different ports for benchmarking. Use separate terminals:

**Terminal 1: Rust Actix**
```bash
cd api
docker-compose up --build
# Runs on port 8080
```

**Terminal 2: Go FastHTTP**
```bash
cd api-go
docker-compose up --build -p 8081:8080
# Runs on port 8081
```

**Terminal 3: Go Fiber**
```bash
cd api-fiber
docker-compose up --build -p 8082:8080
# Runs on port 8082
```

**Terminal 4: Node.js Fastify**
```bash
cd api-fastify
docker-compose up --build -p 8083:8080
# Runs on port 8083
```

**Terminal 5: Python FastAPI**
```bash
cd api-fastapi
docker-compose up --build -p 8084:8080
# Runs on port 8084
```

### Alternative: Docker Network

Or run all on the same Docker network:

```bash
# Create network
docker network create todo-network

# Terminal 1
cd api && docker-compose up --build

# Terminal 2
cd api-go && docker-compose -f docker-compose.yml config | sed 's/localhost/postgres/g' > docker-compose.override.yml && docker-compose up --build --network todo-network

# ... and so on
```

## Running the Benchmark

### Basic Usage

```bash
# Make sure all APIs are running first, then:
python3 benchmark.py
```

### Output

The benchmark will:
1. Check availability of all APIs
2. Run warmup requests on each API
3. Test each endpoint:
   - Health check (`GET /health`)
   - List todos (`GET /api/todos`)
   - Create todo (`POST /api/todos`)
   - Get todo (`GET /api/todos/{id}`)
   - Update todo (`PUT /api/todos/{id}`)
   - Delete todo (`DELETE /api/todos/{id}`)
4. Print detailed results for each endpoint
5. Print overall statistics
6. Export results to `benchmark_results.json`

### Customization

You can customize the benchmark by editing these parameters in `benchmark.py`:

```python
benchmark = TodoBenchmark(
    base_urls={...},
    warmup_requests=10,        # Change this for more/fewer warmup requests
    benchmark_requests=100,    # Change this for more/fewer test requests
)
```

Higher numbers = more accurate but slower benchmarks
Lower numbers = faster but less accurate benchmarks

## Results Interpretation

### Metrics Explained

**Latency (milliseconds)**
- **Min**: Fastest response time
- **Max**: Slowest response time
- **Avg**: Average response time
- **Median**: Middle value (50th percentile)
- **P95**: 95th percentile (95% of requests faster than this)
- **P99**: 99th percentile (99% of requests faster than this)

**RPS (Requests Per Second)**
- How many requests the API can handle per second

**Success Rate**
- Percentage of requests that succeeded (200-399 status codes)

### Interpreting Results

Good metrics:
- ✅ Low average latency (< 10ms)
- ✅ High RPS (> 1000)
- ✅ 100% success rate
- ✅ Low P95/P99 (consistent performance)

Poor metrics:
- ❌ High latency (> 50ms)
- ❌ Low RPS (< 100)
- ❌ Failed requests
- ❌ High variance between P95/P99

## Example Benchmark Results

```
📊 BENCHMARK SUMMARY
======================================================================

📌 Health Check
   ──────────────────────────────────────────────────────────────
   Implementation      Avg (ms)     Min (ms)     Max (ms)     RPS
   ──────────────────────────────────────────────────────────────
   go_fasthttp         1.52         0.89         15.23        658
   go_fiber            1.58         0.92         16.44        631
   rust_actix          2.15         1.23         18.56        464
   node_fastify        3.42         2.11         25.34        292
   python_fastapi      4.21         2.89         31.22        237

📌 List Todos
   ──────────────────────────────────────────────────────────────
   Implementation      Avg (ms)     Min (ms)     Max (ms)     RPS
   ──────────────────────────────────────────────────────────────
   go_fasthttp         2.10         1.10         22.45        476
   go_fiber            2.25         1.20         23.67        444
   rust_actix          2.89         1.45         25.34        346
   node_fastify        4.12         2.23         31.23        243
   python_fastapi      5.34         3.12         38.45        187

📌 Create Todo
   ──────────────────────────────────────────────────────────────
   Implementation      Avg (ms)     Min (ms)     Max (ms)     RPS
   ──────────────────────────────────────────────────────────────
   go_fasthttp         3.45         1.89         28.12        290
   go_fiber            3.58         2.01         29.45        279
   rust_actix          4.23         2.34         32.11        236
   node_fastify        5.89         3.45         41.23        170
   python_fastapi      7.12         4.23         49.56        140

======================================================================
📈 OVERALL STATISTICS
======================================================================

GO_FASTHTTP
  Total Requests:  610
  Successful:      610
  Failed:          0
  Success Rate:    100.0%
  Avg Latency:     2.36ms
  Min Latency:     0.89ms
  Max Latency:     28.12ms

GO_FIBER
  Total Requests:  610
  Successful:      610
  Failed:          0
  Success Rate:    100.0%
  Avg Latency:     2.47ms
  Min Latency:     0.92ms
  Max Latency:     29.45ms

RUST_ACTIX
  Total Requests:  610
  Successful:      610
  Failed:          0
  Success Rate:    100.0%
  Avg Latency:     3.09ms
  Min Latency:     1.23ms
  Max Latency:     32.11ms

NODE_FASTIFY
  Total Requests:  610
  Successful:      610
  Failed:          0
  Success Rate:    100.0%
  Avg Latency:     4.48ms
  Min Latency:     2.11ms
  Max Latency:     41.23ms

PYTHON_FASTAPI
  Total Requests:  610
  Successful:      610
  Failed:          0
  Success Rate:    100.0%
  Avg Latency:     5.56ms
  Min Latency:     2.89ms
  Max Latency:     49.56ms
```

## What's Being Tested

### Health Check
- **Purpose**: Measure baseline overhead
- **What it tests**: Minimal endpoint that just returns status
- **Why**: Shows the framework overhead

### List Todos
- **Purpose**: Measure database query performance
- **What it tests**: Reading all records and serializing to JSON
- **Why**: Most common operation in real applications

### Create Todo
- **Purpose**: Measure write performance
- **What it tests**: Validation, database insert, response generation
- **Why**: Shows how fast the API can handle writes

### Get Todo
- **Purpose**: Measure single record lookup
- **What it tests**: UUID parsing, database lookup, serialization
- **Why**: Common single-item query

### Update Todo
- **Purpose**: Measure update performance
- **What it tests**: Validation, partial update, database update
- **Why**: Shows handling of complex update logic

### Delete Todo
- **Purpose**: Measure delete performance
- **What it tests**: Record deletion and cleanup
- **Why**: Shows delete performance

## Performance Factors

Several factors affect benchmark results:

### 1. Database Performance
- Network latency between API and database
- Database connection pool performance
- Query execution time

**How to isolate**: Run database and API on same machine

### 2. Hardware
- CPU cores available
- Memory available
- Network bandwidth

**How to isolate**: Run on consistent hardware

### 3. System Load
- Other processes running
- System resources in use
- Background tasks

**How to isolate**: Run benchmarks on idle system

### 4. Network Conditions
- Latency between client and server
- Packet loss
- Bandwidth limitations

**How to isolate**: Run benchmark on same machine

## Best Practices

### 1. Run Multiple Times
```bash
# Run benchmarks 3 times to get consistent results
for i in {1..3}; do python3 benchmark.py; done
```

### 2. Use Consistent Hardware
- Run on same machine or identical machines
- Ensure consistent CPU and memory availability
- Disable power saving features

### 3. Isolate Variables
- Only run the APIs being tested
- Stop unnecessary services
- Disable automatic updates

### 4. Warm Up First
- The benchmark includes warmup requests
- Let APIs stabilize before testing
- First requests may be slower

### 5. Use Appropriate Load
- Small request counts = inconsistent results
- Large request counts = slow benchmarks
- 100-1000 requests per endpoint is typical

### 6. Test Real Scenarios
- Vary request types
- Include different data sizes
- Test under realistic conditions

## Troubleshooting

### "No APIs available for benchmarking!"

Make sure all 5 APIs are running on the correct ports:
```bash
# Check if APIs are running
curl http://localhost:8080/health  # Rust Actix
curl http://localhost:8081/health  # Go FastHTTP
curl http://localhost:8082/health  # Go Fiber
curl http://localhost:8083/health  # Node Fastify
curl http://localhost:8084/health  # Python FastAPI
```

### "Connection refused" errors

- Make sure the APIs are running
- Check the port numbers
- Ensure Docker network is set up correctly
- Check firewall settings

### Inconsistent results

- Run on a quiet system (close other apps)
- Run multiple times and average results
- Increase the number of benchmark requests
- Check for background processes

### Very high latency

- Check database connection (network latency)
- Check system resource usage
- Verify API is running in release/production mode
- Check for database bottlenecks

## Advanced Usage

### Custom Ports

Edit `benchmark.py` to change port numbers:

```python
base_urls = {
    "rust_actix": "http://localhost:9000",
    "go_fasthttp": "http://localhost:9001",
    # ...
}
```

### Custom Endpoints

Modify the `benchmark_implementation` method to test custom endpoints

### Load Testing

Increase `benchmark_requests` for load testing:

```python
benchmark = TodoBenchmark(
    base_urls=available_apis,
    warmup_requests=50,
    benchmark_requests=1000,  # Heavy load test
)
```

### Stress Testing

Run continuous benchmarks:

```bash
while true; do python3 benchmark.py; sleep 60; done
```

## Performance Expectations

Based on typical benchmarks:

### Go Implementations (FastHTTP, Fiber)
- **Latency**: 1-3ms average
- **RPS**: 500-1000+
- **Why**: Compiled language, efficient HTTP handling

### Rust Actix
- **Latency**: 2-4ms average
- **RPS**: 300-700+
- **Why**: Compiled, maximum optimization

### Node.js Fastify
- **Latency**: 3-6ms average
- **RPS**: 200-500+
- **Why**: JIT compiled, good performance

### Python FastAPI
- **Latency**: 4-8ms average
- **RPS**: 100-300+
- **Why**: Interpreted language, still good for Python

## Exporting Results

Results are automatically exported to `benchmark_results.json`:

```json
{
  "rust_actix": [
    {
      "name": "Health Check",
      "method": "GET",
      "endpoint": "http://localhost:8080/health",
      "min_latency_ms": 1.23,
      "max_latency_ms": 18.56,
      "avg_latency_ms": 2.15,
      "median_latency_ms": 2.01,
      "p95_latency_ms": 3.45,
      "p99_latency_ms": 8.92,
      "success_rate": 100.0,
      "requests_per_second": 464.0,
      "total_requests": 100,
      "errors": 0
    }
  ]
}
```

## Next Steps

1. **Run the benchmark**: `python3 benchmark.py`
2. **Analyze results**: Check `benchmark_results.json`
3. **Compare implementations**: Which is fastest for your use case?
4. **Profile bottlenecks**: Use API-specific profiling tools
5. **Optimize**: Focus on the slowest components

## Additional Resources

- [Go Profiling](https://golang.org/doc/diagnostics)
- [Rust Benchmarking](https://doc.rust-lang.org/cargo/commands/cargo-bench.html)
- [Node.js Profiling](https://nodejs.org/en/docs/guides/simple-profiling/)
- [Python Profiling](https://docs.python.org/3/library/profile.html)
- [FastHTTP Performance](https://github.com/valyala/fasthttp#performance)
- [Fiber Benchmarks](https://docs.gofiber.io/guide/benchmarks)

---

Happy benchmarking! 🚀
