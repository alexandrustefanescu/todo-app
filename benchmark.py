#!/usr/bin/env python3
"""
Todo API Benchmark Suite
Comprehensive performance testing for all Todo API implementations
"""

import asyncio
import time
import json
import statistics
from typing import List, Dict, Tuple
from dataclasses import dataclass
from datetime import datetime
import httpx
import sys


@dataclass
class BenchmarkResult:
    """Store benchmark results for an endpoint"""
    name: str
    method: str
    endpoint: str
    latencies: List[float]
    errors: int
    total_requests: int

    @property
    def min_latency(self) -> float:
        return min(self.latencies) if self.latencies else 0

    @property
    def max_latency(self) -> float:
        return max(self.latencies) if self.latencies else 0

    @property
    def avg_latency(self) -> float:
        return statistics.mean(self.latencies) if self.latencies else 0

    @property
    def median_latency(self) -> float:
        return statistics.median(self.latencies) if self.latencies else 0

    @property
    def p95_latency(self) -> float:
        if not self.latencies or len(self.latencies) < 20:
            return 0
        sorted_latencies = sorted(self.latencies)
        index = int(len(sorted_latencies) * 0.95)
        return sorted_latencies[index]

    @property
    def p99_latency(self) -> float:
        if not self.latencies or len(self.latencies) < 100:
            return 0
        sorted_latencies = sorted(self.latencies)
        index = int(len(sorted_latencies) * 0.99)
        return sorted_latencies[index]

    @property
    def success_rate(self) -> float:
        return (self.total_requests - self.errors) / self.total_requests * 100 if self.total_requests > 0 else 0

    @property
    def requests_per_second(self) -> float:
        total_time = sum(self.latencies) / 1000 if self.latencies else 1
        return self.total_requests / (total_time if total_time > 0 else 1)


class TodoBenchmark:
    """Benchmark suite for Todo APIs"""

    def __init__(self, base_urls: Dict[str, str], warmup_requests: int = 10, benchmark_requests: int = 100):
        """
        Initialize benchmark suite

        Args:
            base_urls: Dictionary mapping implementation name to base URL
            warmup_requests: Number of warmup requests per endpoint
            benchmark_requests: Number of benchmark requests per endpoint
        """
        self.base_urls = base_urls
        self.warmup_requests = warmup_requests
        self.benchmark_requests = benchmark_requests
        self.results: Dict[str, List[BenchmarkResult]] = {name: [] for name in base_urls.keys()}
        self.test_todo_ids: Dict[str, str] = {}

    async def test_endpoint(
        self,
        client: httpx.AsyncClient,
        method: str,
        url: str,
        name: str,
        data: Dict = None,
        num_requests: int = None,
    ) -> BenchmarkResult:
        """
        Test a single endpoint

        Args:
            client: HTTPX async client
            method: HTTP method (GET, POST, etc.)
            url: Full URL to test
            name: Name of the test
            data: Request body data
            num_requests: Number of requests to make

        Returns:
            BenchmarkResult with latency statistics
        """
        if num_requests is None:
            num_requests = self.benchmark_requests

        latencies = []
        errors = 0

        for _ in range(num_requests):
            try:
                start = time.perf_counter()

                if method == "GET":
                    response = await client.get(url, timeout=30.0)
                elif method == "POST":
                    response = await client.post(url, json=data, timeout=30.0)
                elif method == "PUT":
                    response = await client.put(url, json=data, timeout=30.0)
                elif method == "DELETE":
                    response = await client.delete(url, timeout=30.0)

                elapsed = (time.perf_counter() - start) * 1000  # Convert to milliseconds

                if response.status_code >= 400:
                    errors += 1
                else:
                    latencies.append(elapsed)

            except Exception as e:
                errors += 1
                print(f"  ❌ Error in {name}: {e}")

        result = BenchmarkResult(
            name=name,
            method=method,
            endpoint=url,
            latencies=latencies,
            errors=errors,
            total_requests=num_requests,
        )

        return result

    async def warmup(self, client: httpx.AsyncClient, base_url: str) -> bool:
        """
        Warmup requests to stabilize the API

        Args:
            client: HTTPX async client
            base_url: Base URL of the API

        Returns:
            True if warmup succeeded, False otherwise
        """
        try:
            print(f"  🔥 Warming up with {self.warmup_requests} requests...")
            for i in range(self.warmup_requests):
                try:
                    response = await client.get(f"{base_url}/health", timeout=30.0)
                    if response.status_code == 200:
                        continue
                except Exception:
                    pass

            print("  ✅ Warmup complete")
            return True
        except Exception as e:
            print(f"  ❌ Warmup failed: {e}")
            return False

    async def benchmark_implementation(self, name: str, base_url: str) -> bool:
        """
        Benchmark a single implementation

        Args:
            name: Implementation name
            base_url: Base URL of the API

        Returns:
            True if benchmarking succeeded, False otherwise
        """
        print(f"\n📊 Benchmarking {name.upper()}")
        print(f"   URL: {base_url}")
        print("   " + "=" * 60)

        try:
            async with httpx.AsyncClient() as client:
                # Warmup
                if not await self.warmup(client, base_url):
                    print(f"  ❌ Failed to warmup {name}")
                    return False

                # Test health check
                print(f"  📌 Testing HEALTH endpoint...")
                result = await self.test_endpoint(
                    client,
                    "GET",
                    f"{base_url}/health",
                    "Health Check",
                    num_requests=self.benchmark_requests,
                )
                self.results[name].append(result)
                self._print_result(result)

                # Test LIST endpoint
                print(f"  📌 Testing LIST endpoint...")
                result = await self.test_endpoint(
                    client,
                    "GET",
                    f"{base_url}/api/todos",
                    "List Todos",
                    num_requests=self.benchmark_requests,
                )
                self.results[name].append(result)
                self._print_result(result)

                # Test CREATE endpoint
                print(f"  📌 Testing CREATE endpoint...")
                todo_data = {
                    "title": "Benchmark Test Todo",
                    "description": "Testing performance",
                }
                result = await self.test_endpoint(
                    client,
                    "POST",
                    f"{base_url}/api/todos",
                    "Create Todo",
                    data=todo_data,
                    num_requests=self.benchmark_requests // 2,  # Fewer creates
                )
                self.results[name].append(result)
                self._print_result(result)

                # Get a todo ID for GET/UPDATE/DELETE tests
                try:
                    response = await client.get(f"{base_url}/api/todos", timeout=30.0)
                    if response.status_code == 200:
                        todos = response.json()
                        if todos:
                            self.test_todo_ids[name] = todos[0]["id"]
                except Exception:
                    pass

                # Test GET endpoint
                if name in self.test_todo_ids:
                    print(f"  📌 Testing GET endpoint...")
                    todo_id = self.test_todo_ids[name]
                    result = await self.test_endpoint(
                        client,
                        "GET",
                        f"{base_url}/api/todos/{todo_id}",
                        "Get Todo",
                        num_requests=self.benchmark_requests,
                    )
                    self.results[name].append(result)
                    self._print_result(result)

                    # Test UPDATE endpoint
                    print(f"  📌 Testing UPDATE endpoint...")
                    update_data = {"completed": True, "title": "Updated Test Todo"}
                    result = await self.test_endpoint(
                        client,
                        "PUT",
                        f"{base_url}/api/todos/{todo_id}",
                        "Update Todo",
                        data=update_data,
                        num_requests=self.benchmark_requests // 2,
                    )
                    self.results[name].append(result)
                    self._print_result(result)

                    # Test DELETE endpoint
                    print(f"  📌 Testing DELETE endpoint...")
                    result = await self.test_endpoint(
                        client,
                        "DELETE",
                        f"{base_url}/api/todos/{todo_id}",
                        "Delete Todo",
                        num_requests=10,  # Just a few deletes
                    )
                    self.results[name].append(result)
                    self._print_result(result)

                print("  ✅ Benchmarking complete")
                return True

        except Exception as e:
            print(f"  ❌ Benchmarking failed: {e}")
            return False

    def _print_result(self, result: BenchmarkResult):
        """Print a single result"""
        if result.latencies:
            print(f"     {result.name}:")
            print(f"       Min:     {result.min_latency:.2f}ms")
            print(f"       Max:     {result.max_latency:.2f}ms")
            print(f"       Avg:     {result.avg_latency:.2f}ms")
            print(f"       Median:  {result.median_latency:.2f}ms")
            if result.p95_latency:
                print(f"       P95:     {result.p95_latency:.2f}ms")
            if result.p99_latency:
                print(f"       P99:     {result.p99_latency:.2f}ms")
            print(f"       RPS:     {result.requests_per_second:.0f}")
            print(f"       Success: {result.success_rate:.1f}%")
        else:
            print(f"     {result.name}: ❌ No successful requests")

    async def run_all_benchmarks(self):
        """Run benchmarks for all implementations"""
        print("\n" + "=" * 70)
        print("🚀 TODO API BENCHMARK SUITE")
        print("=" * 70)
        print(f"\nStarting benchmarks at {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        print(f"  Warmup requests per endpoint: {self.warmup_requests}")
        print(f"  Benchmark requests per endpoint: {self.benchmark_requests}")

        for name, base_url in self.base_urls.items():
            await self.benchmark_implementation(name, base_url)

    def print_summary(self):
        """Print benchmark summary and comparison"""
        print("\n" + "=" * 70)
        print("📊 BENCHMARK SUMMARY")
        print("=" * 70)

        # Group results by endpoint
        endpoints = {}
        for impl_name, results in self.results.items():
            for result in results:
                endpoint_key = result.name
                if endpoint_key not in endpoints:
                    endpoints[endpoint_key] = {}
                endpoints[endpoint_key][impl_name] = result

        # Print comparison for each endpoint
        for endpoint_name, implementations in endpoints.items():
            print(f"\n📌 {endpoint_name}")
            print("   " + "-" * 60)
            print(f"   {'Implementation':<15} {'Avg (ms)':<12} {'Min (ms)':<12} {'Max (ms)':<12} {'RPS':<10}")
            print("   " + "-" * 60)

            # Sort by average latency
            sorted_impls = sorted(
                implementations.items(),
                key=lambda x: x[1].avg_latency if x[1].latencies else float('inf')
            )

            for impl_name, result in sorted_impls:
                if result.latencies:
                    print(
                        f"   {impl_name:<15} {result.avg_latency:<12.2f} "
                        f"{result.min_latency:<12.2f} {result.max_latency:<12.2f} "
                        f"{result.requests_per_second:<10.0f}"
                    )
                else:
                    print(f"   {impl_name:<15} ❌ FAILED")

        # Overall statistics
        print("\n" + "=" * 70)
        print("📈 OVERALL STATISTICS")
        print("=" * 70)

        for impl_name, results in self.results.items():
            total_latencies = []
            total_errors = 0
            total_requests = 0

            for result in results:
                total_latencies.extend(result.latencies)
                total_errors += result.errors
                total_requests += result.total_requests

            if total_latencies:
                avg_latency = statistics.mean(total_latencies)
                min_latency = min(total_latencies)
                max_latency = max(total_latencies)
                success_rate = (total_requests - total_errors) / total_requests * 100 if total_requests > 0 else 0

                print(f"\n{impl_name.upper()}")
                print(f"  Total Requests:  {total_requests}")
                print(f"  Successful:      {total_requests - total_errors}")
                print(f"  Failed:          {total_errors}")
                print(f"  Success Rate:    {success_rate:.1f}%")
                print(f"  Avg Latency:     {avg_latency:.2f}ms")
                print(f"  Min Latency:     {min_latency:.2f}ms")
                print(f"  Max Latency:     {max_latency:.2f}ms")
            else:
                print(f"\n{impl_name.upper()}: ❌ NO DATA")

    def export_results(self, filename: str = "benchmark_results.json"):
        """Export results to JSON file"""
        export_data = {}

        for impl_name, results in self.results.items():
            export_data[impl_name] = []
            for result in results:
                export_data[impl_name].append({
                    "name": result.name,
                    "method": result.method,
                    "endpoint": result.endpoint,
                    "min_latency_ms": result.min_latency,
                    "max_latency_ms": result.max_latency,
                    "avg_latency_ms": result.avg_latency,
                    "median_latency_ms": result.median_latency,
                    "p95_latency_ms": result.p95_latency,
                    "p99_latency_ms": result.p99_latency,
                    "success_rate": result.success_rate,
                    "requests_per_second": result.requests_per_second,
                    "total_requests": result.total_requests,
                    "errors": result.errors,
                })

        with open(filename, 'w') as f:
            json.dump(export_data, f, indent=2)

        print(f"\n✅ Results exported to {filename}")


async def main():
    """Main benchmark runner"""
    # Define API implementations to benchmark
    base_urls = {
        "rust_actix": "http://localhost:8080",
        "go_fasthttp": "http://localhost:8081",
        "go_fiber": "http://localhost:8082",
        "node_fastify": "http://localhost:8083",
        "python_fastapi": "http://localhost:8084",
    }

    # Filter out APIs that aren't running
    print("\n🔍 Checking API availability...")
    available_apis = {}

    async with httpx.AsyncClient() as client:
        for name, url in base_urls.items():
            try:
                response = await client.get(f"{url}/health", timeout=5.0)
                if response.status_code == 200:
                    available_apis[name] = url
                    print(f"  ✅ {name.upper()} at {url}")
                else:
                    print(f"  ⚠️  {name.upper()} not responding correctly")
            except Exception as e:
                print(f"  ❌ {name.upper()} not available: {e}")

    if not available_apis:
        print("\n❌ No APIs available for benchmarking!")
        print("\nMake sure to start the APIs first:")
        print("  cd api && docker-compose up --build")
        print("  cd api-go && docker-compose up --build")
        print("  cd api-fiber && docker-compose up --build")
        print("  cd api-fastify && docker-compose up --build")
        print("  cd api-fastapi && docker-compose up --build")
        sys.exit(1)

    # Create benchmark suite
    benchmark = TodoBenchmark(
        base_urls=available_apis,
        warmup_requests=10,
        benchmark_requests=100,
    )

    # Run benchmarks
    await benchmark.run_all_benchmarks()

    # Print summary
    benchmark.print_summary()

    # Export results
    benchmark.export_results()


if __name__ == "__main__":
    asyncio.run(main())
