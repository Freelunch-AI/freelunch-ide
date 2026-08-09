# [dev + ops modes] [declarative] [self-contained] Data Storage Front — lunch-storage: storage runtime + dataops + visual data discovery

A single data storage platform that gives an organization:
- Autoscaled databases (Postgres/MySQL/etc.)
- Autoscaled message brokers (Kafka/NATS/RabbitMQ)
- Unified backup + restore + staging snapshot system
- Automatic PII anonymization + data sampling for makng snaophots for staging environments
- Zero-copy branching for dev/staging workloads to use (like Snowflake/Neon/PlanetScale but engine-agnostic)
- Unified data discovery/visualization/read-only query layer
- Schema migration governance (like Bytebase)
- Cross-engine migration support (e.g., MySQL → PG, PG → ClickHouse)
- Monitoring and Observability on it (with time-travel & replay capabilities)
- Tool-agnostic stable developer SDKs (Developers dont need to change their code when DB engine changes under the hood, unless they need engine-specific features)
- Data Version Control Option (LakeFS for Object Storage & CDC for Databases)

In short:
A Storage-as-a-Platform layer that standardizes storage & management of data. Developers jus use SDK 90% of the time. Data Engineers just change config, observe, approve scaffolded engine migrations by looking at a preview.
