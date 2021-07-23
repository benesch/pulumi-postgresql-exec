# PostgreSQL Exec Pulumi Provider

A [Pulumi](https://pulumi.com) provider that allows running arbitrary SQL
queries against a PostgreSQL database.

**Warning:** It is your responsibility to ensure that you supply both the
forward and reverse SQL statements.

## Usage example

To provision a PostgreSQL table:

```python
import pulumi_postgresql_exec as postgresql_exec

postgresql_exec.Exec(
    "create-table",
    create_sql="CREATE TABLE t (a int)",
    destroy_sql="DROP TABLE t",
)
```
