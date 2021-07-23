import pulumi
import pulumi_postgresql_exec as postgresql_exec

provider = postgresql_exec.Provider(
    "provider",
    host="localhost",
    port=5432,
)

postgresql_exec.Exec(
    "create-table",
    create_sql="CREATE TABLE t (a int)",
    destroy_sql="DROP TABLE t",
    opts=pulumi.ResourceOptions(provider=provider)
)
