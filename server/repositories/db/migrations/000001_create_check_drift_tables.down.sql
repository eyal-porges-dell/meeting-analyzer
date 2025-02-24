/*
 * Copyright © 2022 Dell Inc. or its subsidiaries. All Rights Reserved.
 */

-- Drop the tables if they exist
DROP TABLE IF EXISTS node_drifts;

DROP TABLE IF EXISTS deployment_drifts;

DROP TABLE IF EXISTS node_type_category;

-- Drop the ENUM type if it exists and if it is no longer needed
DROP TYPE IF EXISTS status_enum;

DROP TYPE IF EXISTS category_enum;
