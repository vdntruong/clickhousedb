# ClickHouse DB Go Application

This project demonstrates efficient data operations with ClickHouse using Go, focusing on high-performance bulk insert operations and database migrations.

## Overview

This application connects to a ClickHouse database, performs migrations, and demonstrates best practices for bulk data ingestion. It uses the official ClickHouse Go driver and implements optimized batch insert patterns.

## Getting Started

### Prerequisites

- Go 1.24+
- ClickHouse server (local or remote)
- Docker (optional, for containerized ClickHouse)

### Setup

1. Clone the repository
   ```bash
   git clone https://github.com/yourusername/clickhousedb
   cd clickhousedb
   ```

2. Set up your environment
   ```bash
   cp .env.example .env
   # Edit .env with your ClickHouse connection parameters
   ```

3. Run the application
   ```bash
   go run .
   ```

## Configuration

The application uses environment variables for configuration, which can be set in the `.env` file:

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | ClickHouse host | localhost |
| `DB_PORT` | ClickHouse port | 9000 |
| `DB_USER` | Database username | default |
| `DB_PASS` | Database password | (empty) |
| `DB_NAME` | Database name | default |
| `MIGRATIONS_PATH` | Path to migrations | migrations |

## ClickHouse Bulk Insert Best Practices

### Optimal Batch Size

For efficient data insertion into ClickHouse, the following batch sizes are recommended:

- **Minimum**: 1,000 rows per batch
- **Optimal**: 10,000-100,000 rows per batch

Large batch sizes amortize network overhead and significantly improve throughput.

### Key Considerations

1. **Create New Batch After Sending**
    - Always create a new batch with `PrepareBatch()` after sending a batch
    - Reusing a sent batch will result in "batch has already been sent" errors

2. **Memory Management**
    - Monitor memory usage with large batches
    - Use appropriately sized batches for your available memory

3. **Pre-sorting Data**
    - When possible, pre-sort data by primary key before insertion
    - This improves insertion efficiency, especially for large batches

4. **Context and Timeouts**
    - Use context with appropriate timeouts for batch operations
    - Long-running batch operations should have longer timeouts

5. **Error Handling**
    - Implement proper error handling for batch operations
    - Abort batches when errors occur to prevent partial inserts

## Database Migrations

This project uses the `golang-migrate` library to manage database schema changes. Migrations are stored in the `migrations` directory and are automatically applied when the application starts.

## Project Structure

- `main.go` - Application entry point
- `repository.go` - Data access layer with bulk insert implementation
- `db.go` - Database connection and migration management
- `env.go` - Environment configuration
- `migrations/` - SQL migration files

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [ClickHouse](https://clickhouse.com/) - High-performance columnar database
- [ClickHouse Go Driver](https://github.com/ClickHouse/clickhouse-go) - Official Go driver for ClickHouse
