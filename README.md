# Postfix AWS Cassandra Integration

## Introduction

In our quest to achieve a highly available and scalable mail infrastructure, we needed to connect our distributed Postfix servers to a centralized database. High availability was a primary concern, leading us to choose Cassandra for its robust replication and fault tolerance capabilities. To simplify management and reduce operational overhead, we opted for **AWS Keyspaces**, Amazon's managed Cassandra service.

This package provides a solution to integrate Postfix with AWS Keyspaces, allowing Postfix to query mail routing and policy information from a centralized Cassandra database. By installing this package on your **Amazon Linux EC2 instances**, you can run a service that enables Postfix to fetch up-to-date data from AWS Keyspaces seamlessly.

## How the Package Works

The package installs a Go-based socketmap daemon that runs as a service on your EC2 instance. This daemon listens for queries from Postfix and retrieves the necessary information from AWS Keyspaces.

### Key Features:

- **High Availability**: Leveraging Cassandra's distributed architecture ensures no single point of failure.
- **Managed Service**: AWS Keyspaces eliminates the need to manage Cassandra clusters manually.
- **Scalable**: Easily scale your mail infrastructure without worrying about the underlying database.

### Installation Overview:

1. **Install the RPM Package**: Install the provided RPM package on your Amazon Linux EC2 instance.
2. **Configure Environment Variables**: Set the required environment variables for the daemon.
3. **Start the Service**: Enable and start the _postfix-aws-cassandra_ service.
4. **Configure Postfix**: Update your Postfix configuration to use the socketmap daemon for various lookups.

## Installation and Setup

### Prerequisites

- **Amazon Linux EC2 Instance**: The package is built and tested on Amazon Linux 2.
- **Postfix Installed**: Ensure Postfix is installed and running on your instance.
- **AWS IAM Role**: The EC2 instance should have an IAM role with permissions to access AWS Keyspaces.

### Step 1: Install the RPM Package

Transfer the RPM package to your EC2 instance and install it:

```bash
sudo yum localinstall postfix-aws-cassandra-1.0.0-1.amzn2.x86_64.rpm
```

### Step 2: Configure Environment Variables

The socketmap daemon requires certain environment variables to connect to AWS Keyspaces. You can set these variables in the service's environment file or as part of the systemd service configuration.

#### Required Environment Variables:

- **_AWS_REGION_**: The AWS region where your Keyspaces instance is located.

  - Example: `us-east-1`

- **_KEYSPACE_NAME_**: The name of your Cassandra keyspace.
  - Example: `mailion`

#### Optional Variables

- **_PORT_**: The port on which the socketmap daemon listens.

  - Default: `9999`
  - Example: `PORT=9999`

#### Setting Environment Variables

You can set these variables in the environment file `/etc/default/postfix-aws-cassandra` or `/etc/sysconfig/postfix-aws-cassandra`.

**Example:**

```ini
AWS_REGION=us-east-1
KEYSPACE_NAME=mailion
PORT=9999
```

Alternatively, you can modify the systemd service file at _/usr/lib/systemd/system/postfix-aws-cassandra.service_ and add the environment variables:

```ini
[Service]
Type=simple
ExecStart=/usr/local/bin/postfix-aws-cassandra
User=postfix
Group=postfix
Environment="AWS_REGION=us-east-1"
Environment="KEYSPACE_NAME=mailion"
Environment="PORT=9999"
Restart=on-failure
```

After modifying the service file, reload systemd and restart the service:

```bash
sudo systemctl daemon-reload
sudo systemctl restart postfix-aws-cassandra
```

### Step 3: Enable and Start the Service

Enable the _postfix-aws-cassandra_ service to start on boot and start it immediately:

```bash
sudo systemctl enable postfix-aws-cassandra
sudo systemctl start postfix-aws-cassandra
```

Verify that the service is running:

```bash
sudo systemctl status postfix-aws-cassandra
```

## Configuring Postfix

After installing the package and starting the service, you need to configure Postfix to use the socketmap daemon for various lookups. The socketmap daemon listens on `127.0.0.1:9999` by default.

### Main Configuration (_/etc/postfix/main.cf_)

Add or update the following configurations in your _main.cf_ file:

#### 1. Virtual Alias Maps

```ini
virtual_alias_maps = socketmap:tcp:127.0.0.1:9999:virtual_aliases
```

#### 2. Virtual Mailbox Domains

```ini
virtual_mailbox_domains = socketmap:tcp:127.0.0.1:9999:domains
```

#### 3. Relay Domains

```ini
relay_domains = socketmap:tcp:127.0.0.1:9999:relay_domains
```

#### 4. Transport Maps

```ini
transport_maps = socketmap:tcp:127.0.0.1:9999:transport_maps
```

#### 5. Access Maps

```ini
smtpd_sender_restrictions =
check_sender_access socketmap:tcp:127.0.0.1:9999:access_maps,
permit
```

### Explanation of Configuration Options

- **virtual_alias_maps**: Redirects emails from one address to another.
- **virtual_mailbox_domains**: Specifies domains for which Postfix will accept mail and store it locally.
- **relay_domains**: Specifies domains for which Postfix will relay mail.
- **transport_maps**: Defines custom transport methods for specific destinations.
- **access_maps**: Controls access based on the sender's email address.

### Reload Postfix

After updating the configuration, reload Postfix to apply the changes:

```bash
sudo systemctl reload postfix
```

## AWS Keyspaces Table Structures

Below are the Cassandra Query Language (CQL) statements to create the necessary tables in AWS Keyspaces. Replace _your_keyspace_ with the name of your keyspace if different.

### 1. Domains Table

```cql
CREATE TABLE your_keyspace.domains (
domain text,
active boolean,
PRIMARY KEY (domain)
) WITH CUSTOM_PROPERTIES = {'capacity_mode': {'throughput_mode': 'PAY_PER_REQUEST'}};
```

### 2. Users Table

```cql
CREATE TABLE your_keyspace.users (
email text,
active boolean,
PRIMARY KEY (email)
) WITH CUSTOM_PROPERTIES = {'capacity_mode': {'throughput_mode': 'PAY_PER_REQUEST'}};
```

### 3. Relay Domains Table

```cql
CREATE TABLE your_keyspace.relay_domains (
domain text,
active boolean,
PRIMARY KEY (domain)
) WITH CUSTOM_PROPERTIES = {'capacity_mode': {'throughput_mode': 'PAY_PER_REQUEST'}};
```

### 4. Virtual Aliases Table

```cql
CREATE TABLE your_keyspace.virtual_aliases (
alias text,
destination text,
PRIMARY KEY (alias)
) WITH CUSTOM_PROPERTIES = {'capacity_mode': {'throughput_mode': 'PAY_PER_REQUEST'}};
```

### 5. Transport Maps Table

```cql
CREATE TABLE your_keyspace.transport_maps (
address text,
transport text,
PRIMARY KEY (address)
) WITH CUSTOM_PROPERTIES = {'capacity_mode': {'throughput_mode': 'PAY_PER_REQUEST'}};
```

### 6. Access Maps Table

```cql
CREATE TABLE your_keyspace.access_maps (
sender text,
action text,
PRIMARY KEY (sender)
) WITH CUSTOM_PROPERTIES = {'capacity_mode': {'throughput_mode': 'PAY_PER_REQUEST'}};
```

### Notes:

- **Capacity Mode**: The tables are set to use `PAY_PER_REQUEST` mode for billing. Adjust as needed.
- **Keyspace Replication**: AWS Keyspaces uses `SingleRegionStrategy` for replication.

## Contributing

We welcome contributions to this project! If you have any questions, suggestions, or issues, please feel free to [create an issue](https://github.com/your-repo/issues) in the repository. Your feedback helps us improve and adapt the solution to better meet the community's needs.

## License

This project is licensed under the MIT License. See the [Apache Lisence](https://www.apache.org/licenses/LICENSE-2.0) file for details.

---

**Note**: Replace _your-repo_ and _your_keyspace_ with your actual GitHub repository URL and keyspace name, respectively.
