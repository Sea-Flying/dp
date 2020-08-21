CREATE KEYSPACE IF NOT EXISTS artifact
WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

CREATE TABLE IF NOT EXISTS artifact.class (
  name ASCII,
  ci_url TEXT,
  created_time TIMESTAMP,
  default_nomad_template ASCII,
  default_template_params MAP<ASCII, TEXT>,
  unit_timeout_seconds INT,
  description TEXT,
  git_url TEXT,
  kind ASCII,
  repo_name ASCII,
  PRIMARY KEY((name))
);

CREATE INDEX IF NOT EXISTS idx_created_time ON artifact.class (created_time);
CREATE INDEX IF NOT EXISTS idx_kind ON artifact.class (kind);
CREATE INDEX IF NOT EXISTS idx_repo_name ON artifact.class (repo_name);

CREATE TABLE IF NOT EXISTS artifact.entity (
  class_name ASCII,
  generated_time TIMESTAMP,
  checksum ASCII,
  class_kind ASCII,
  repo_name ASCII,
  uploader ASCII,
  url TEXT,
  version ASCII,
  valid BOOLEAN,
  PRIMARY KEY ((class_name), generated_time)
)
WITH CLUSTERING ORDER BY (generated_time DESC);

CREATE INDEX IF NOT EXISTS idx_class_kind ON artifact.entity (class_kind);
CREATE INDEX IF NOT EXISTS idx_repo_name ON artifact.entity (repo_name);

CREATE MATERIALIZED VIEW IF NOT EXISTS artifact.mv_entity_by_version
  AS SELECT version, class_name, generated_time, url
  FROM artifact.entity
  WHERE version IS NOT NULL AND class_name IS NOT NULL AND generated_time IS NOT NULL
  PRIMARY KEY (version, class_name, generated_time);

CREATE TABLE IF NOT EXISTS artifact.kind (
  name ASCII,
  created_time TIMESTAMP,
  default_repo ASCII,
  default_template ASCII,
  default_unit_timeout_seconds int,
  PRIMARY KEY (name)
);

CREATE TABLE IF NOT EXISTS artifact.repo (
  name ASCII,
  base_url TEXT,
  created_time TIMESTAMP,
  description TEXT,
  kind ASCII,
  web_url TEXT,
  PRIMARY KEY (name)
);

CREATE INDEX IF NOT EXISTS idx_kind ON artifact.repo (kind);
CREATE INDEX IF NOT EXISTS idx_created_time ON artifact.repo (created_time);


CREATE KEYSPACE IF NOT EXISTS schedule
WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

CREATE TABLE IF NOT EXISTS schedule.job
(
  class_name ASCII,
  created_time TIMESTAMP,
  entity_generated_time TIMESTAMP,
  entity_version ASCII,
  executed_time TIMESTAMP,
  kind ASCII,
  nomad_template_name ASCII,
  nomad_template_params MAP<ASCII, TEXT>,
  unit_timeout_seconds INT,
  status ASCII,
  PRIMARY KEY ((class_name), created_time)
);

CREATE MATERIALIZED VIEW IF NOT EXISTS schedule.mv_job_by_entity_version
  AS SELECT class_name, entity_generated_time, created_time, entity_version, executed_time, status
  FROM schedule.job
  WHERE class_name IS NOT NULL AND entity_generated_time IS NOT NULL AND created_time IS NOT NULL
    AND entity_version IS NOT NULL AND executed_time IS NOT NULL AND status IS NOT NULL
  primary key (entity_version, class_name, created_time);

CREATE MATERIALIZED VIEW IF NOT EXISTS schedule.mv_job_by_entity_generated_time
  AS SELECT class_name, entity_generated_time, created_time, entity_version, executed_time, status
  FROM schedule.job
  WHERE class_name IS NOT NULL AND entity_generated_time IS NOT NULL AND created_time IS NOT NULL
    AND entity_version IS NOT NULL AND executed_time IS NOT NULL AND status IS NOT NULL
primary key (entity_generated_time, class_name, created_time);

CREATE MATERIALIZED VIEW IF NOT EXISTS schedule.mv_job_by_status
  AS SELECT status, class_name, entity_generated_time, created_time, entity_version, executed_time
  FROM schedule.job
  WHERE class_name IS NOT NULL AND entity_generated_time IS NOT NULL AND created_time IS NOT NULL
    AND entity_version IS NOT NULL AND executed_time IS NOT NULL AND status IS NOT NULL
primary key (status, class_name, created_time);

CREATE TABLE IF NOT EXISTS schedule.nomad_template
(
  name ASCII PRIMARY KEY,
  git_url TEXT,
  params MAP<ASCII, TEXT>,
  params_description MAP<ASCII, TEXT>,
  tags SET<ASCII>
);

CREATE KEYSPACE IF NOT EXISTS app
    WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

CREATE TABLE IF NOT EXISTS app.status_history(
    app_name ascii,
    time timestamp,
    status_changed_to ascii,
    PRIMARY KEY ( (app_name), time )
) WITH CLUSTERING ORDER BY (time DESC);

CREATE TABLE IF NOT EXISTS app.stopped_app_cache(
    app_name ascii PRIMARY KEY ,
    nomad_job_json TEXT
);
