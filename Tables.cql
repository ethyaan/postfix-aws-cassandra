CREATE KEYSPACE mailion
WITH REPLICATION = {'class': 'SingleRegionStrategy'};

CREATE TABLE mailion.domains (
    domain text,
    active boolean,
    PRIMARY KEY (domain)
) WITH CUSTOM_PROPERTIES = {'capacity_mode': {'throughput_mode': 'PAY_PER_REQUEST'}};

CREATE TABLE mailion.users (
    email text,
    active boolean,
    PRIMARY KEY (email)
) WITH CUSTOM_PROPERTIES = {'capacity_mode': {'throughput_mode': 'PAY_PER_REQUEST'}};

CREATE TABLE mailion.relay_domains (
    domain text,
    active boolean,
    PRIMARY KEY (domain)
) WITH CUSTOM_PROPERTIES = {'capacity_mode': {'throughput_mode': 'PAY_PER_REQUEST'}};

CREATE TABLE mailion.virtual_aliases (
    alias text,
    destination text,
    PRIMARY KEY (alias)
) WITH CUSTOM_PROPERTIES = {'capacity_mode': {'throughput_mode': 'PAY_PER_REQUEST'}};

CREATE TABLE mailion.transport_maps (
    address text,
    transport text,
    PRIMARY KEY (address)
) WITH CUSTOM_PROPERTIES = {'capacity_mode': {'throughput_mode': 'PAY_PER_REQUEST'}};

CREATE TABLE mailion.access_maps (
    sender text,
    action text,
    PRIMARY KEY (sender)
) WITH CUSTOM_PROPERTIES = {'capacity_mode': {'throughput_mode': 'PAY_PER_REQUEST'}};