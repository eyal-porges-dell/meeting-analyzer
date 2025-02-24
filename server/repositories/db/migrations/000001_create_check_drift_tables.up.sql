/*
 * Copyright © 2022 Dell Inc. or its subsidiaries. All Rights Reserved.
 */

-- Drop the tables if they exist since we need to reload them each time
DROP TABLE IF EXISTS node_type_category;

DROP TYPE IF EXISTS category_enum;

-- Create ENUM types for status and category
CREATE TYPE deployment_status_enum AS ENUM ('PENDING','IN_PROGRESS',  'DONE', 'FAILED', 'CANCELLED',  'DELETED_ON_UPDATE');
CREATE TYPE node_status_enum AS ENUM ('IN_PROGRESS' ,'DONE', 'FAILED');
CREATE TYPE category_enum AS ENUM ('CONF', 'INFRA', 'SEC', 'SOFT');

-- Create the deployment_drifts table
CREATE TABLE IF NOT EXISTS deployment_drifts (
    deployment_id VARCHAR(256) PRIMARY KEY NOT NULL,
    execution_id VARCHAR(256),
    "timestamp" TIMESTAMPTZ NOT NULL,
    total_nodes INTEGER,
    status deployment_status_enum NOT NULL,
    error TEXT
);

-- Create the node_drifts table
CREATE TABLE IF NOT EXISTS node_drifts (
    node_id VARCHAR(256) NOT NULL,
    deployment_id VARCHAR(256) NOT NULL,
    node_type VARCHAR(256) NOT NULL,
    diff_count INTEGER,
    diff TEXT,
    status node_status_enum,
    "timestamp" TIMESTAMPTZ,
    category category_enum NOT NULL DEFAULT 'CONF',    -- Include 'DEFAULT' for the category column
    error VARCHAR(256),
    PRIMARY KEY (node_id, deployment_id),
    FOREIGN KEY (deployment_id) REFERENCES deployment_drifts(deployment_id)
    ON DELETE CASCADE
);

-- Create the node_type_category table
CREATE TABLE node_type_category (
    node_type VARCHAR(255) PRIMARY KEY, -- Column for node type, set as the primary key
    category category_enum NOT NULL     -- Column for the category, must be one of the enum values
);


-- Categories Insert statements
INSERT INTO node_type_category (node_type, category) VALUES
      ('nativeedge.nodes.vsphere.IPPool', 'INFRA'),
      ('nativeedge.nodes.vsphere.NIC', 'INFRA'),
      ('nativeedge.nodes.vsphere.Server', 'INFRA'),
      ('nativeedge.nodes.vsphere.Storage', 'INFRA'),
      ('nativeedge.nodes.template.NativeEdgeCompose', 'CONF'),
      ('nativeedge.nodes.template.NativeEdgeVM', 'INFRA'),
      ('nativeedge.nodes.template.ApplicationVM', 'INFRA'),
      ('nativeedge.nodes.template.Config', 'INFRA'),
      ('nativeedge.nodes.template.BinaryImage', 'SOFT'),
      ('nativeedge.nodes.template.RegistryLogin', 'SEC');
