---
layout: default
title: On AWS
parent: Deploy lakeFS
description: 
nav_order: 10
redirect_from:
   - ../deploying-aws/index.html
   - ../deploying-aws/install.html
   - ../deploying-aws/db.html
   - ../deploying-aws/lb_dns.html
---

# Deploy lakeFS on AWS
{: .no_toc }
Expected deployment time: 25 min

{% include toc.html %}

{% include_relative includes/prerequisites.md %}

## Creating the Database on AWS RDS
lakeFS requires a PostgreSQL database to synchronize actions on your repositories.
We will show you how to create a database on AWS RDS but you can use any PostgreSQL database as long as it's accessible by your lakeFS installation.

If you already have a database, take note of the connection string and skip to the [next step](#install-lakefs-on-ec2)

1. Follow the official [AWS documentation](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_GettingStarted.CreatingConnecting.PostgreSQL.html){: target="_blank" } on how to create a PostgreSQL instance and connect to it.
   You may use the default PostgreSQL engine, or [Aurora PostgreSQL](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/Aurora.AuroraPostgreSQL.html){: target="_blank" }. Make sure that you're using PostgreSQL version >= 11.
2. Once your RDS is set up and the server is in `Available` state, take note of the endpoint and port.

   ![RDS Connection String]({{ site.baseurl }}/assets/img/rds_conn.png)

3. Make sure your security group rules allow you to connect to the database instance.

## Installation Options

### On EC2
1. Save the following configuration file as `config.yaml`:

   ```yaml
   ---
   database:
     connection_string: "[DATABASE_CONNECTION_STRING]"
   auth:
     encrypt:
       # replace this with a randomly-generated string:
       secret_key: "[ENCRYPTION_SECRET_KEY]"
   blockstore:
     type: s3
     s3:
       region: us-east-1 # optional, fallback in case discover from bucket is not supported
   ```

1. [Download the binary](../index.md#downloads) to the EC2 instance.
1. Run the `lakefs` binary on the EC2 instance:
   ```bash
   lakefs --config config.yaml run
   ```
   **Note:** It's preferable to run the binary as a service using systemd or your operating system's facilities.

### On ECS
To support container-based environments like AWS ECS, lakeFS can be configured using environment variables. Here is a `docker run` 
command to demonstrate starting lakeFS using Docker:

```sh
docker run \
  --name lakefs \
  -p 8000:8000 \
  -e LAKEFS_DATABASE_CONNECTION_STRING="[DATABASE_CONNECTION_STRING]" \
  -e LAKEFS_AUTH_ENCRYPT_SECRET_KEY="[ENCRYPTION_SECRET_KEY]" \
  -e LAKEFS_BLOCKSTORE_TYPE="s3" \
  treeverse/lakefs:latest run
```

See the [reference](../reference/configuration.md#using-environment-variables) for a complete list of environment variables.

### On EKS
See [Kubernetes Deployment](./k8s.md).

## Load balancing
Depending on how you chose to install lakeFS, you should have a load balancer direct requests to the lakeFS server.  
By default, lakeFS operates on port 8000, and exposes a `/_health` endpoint which you can use for health checks.

### Notes for using an AWS Application Load Balancer
{: .no_toc }

1. Your security groups should allow the load balancer to access the lakeFS server.
1. Create a target group with a listener for port 8000.
1. Setup TLS termination using the domain names you wish to use (e.g., `lakefs.example.com` and potentially `s3.lakefs.example.com`, `*.s3.lakefs.example.com` if using [virtual-host addressing](https://docs.aws.amazon.com/AmazonS3/latest/userguide/VirtualHosting.html)).
1. Configure the health-check to use the exposed `/_health` URL

## Next Steps
Your next step is to [prepare your storage](../setup/storage/index.md). If you already have a storage bucket/container, you're ready to [create your first lakeFS repository](../setup/create-repo.md).
