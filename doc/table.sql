-- public.aggregated_rules definition

-- Drop table

-- DROP TABLE public.aggregated_rules;

CREATE TABLE public.aggregated_rules (
	id serial4 NOT NULL,
	"name" varchar(255) NOT NULL,
	services _varchar DEFAULT '{}'::character varying[] NOT NULL,
	operation_filters jsonb DEFAULT '{}'::jsonb NOT NULL,
	configurations jsonb DEFAULT '{}'::jsonb NOT NULL,
	priority int4 DEFAULT 0 NOT NULL,
	disabled bool DEFAULT false NOT NULL,
	created_by varchar(255) NOT NULL,
	updated_by varchar(255) NOT NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	CONSTRAINT aggregated_rules_pkey PRIMARY KEY (id)
);


-- public.alerts definition

-- Drop table

-- DROP TABLE public.alerts;

CREATE TABLE public.alerts (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	"name" varchar(255) NOT NULL,
	"type" varchar(255) NOT NULL,
	description text NULL,
	notification_configs _uuid DEFAULT '{}'::uuid[] NOT NULL,
	notification_format varchar(255) DEFAULT ''::character varying NOT NULL,
	disabled bool DEFAULT false NULL,
	created_by varchar(255) NOT NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	organization_id varchar NOT NULL,
	workspace_id varchar NOT NULL,
	CONSTRAINT alerts_pkey PRIMARY KEY (id)
);


-- public.apdex_settings definition

-- Drop table

-- DROP TABLE public.apdex_settings;

CREATE TABLE public.apdex_settings (
	service_name text NOT NULL,
	threshold float8 NOT NULL,
	exclude_status_codes text NOT NULL,
	CONSTRAINT apdex_settings_pkey PRIMARY KEY (service_name)
);


-- public.auth_apis definition

-- Drop table

-- DROP TABLE public.auth_apis;

CREATE TABLE public.auth_apis (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	created_by varchar(50) NOT NULL,
	created_at timestamptz(6) DEFAULT CURRENT_TIMESTAMP NULL,
	updated_by varchar(50) NULL,
	updated_at timestamptz(6) DEFAULT CURRENT_TIMESTAMP NULL,
	permission_id uuid NULL,
	"path" varchar(255) NOT NULL,
	"method" varchar(10) NOT NULL,
	model_name _text DEFAULT '{}'::text[] NULL,
	is_universal bool DEFAULT false NULL,
	license_model_name varchar(255) NULL,
	description varchar(255) NULL,
	license_enable bool DEFAULT false NULL,
	CONSTRAINT auth_apis_pkey PRIMARY KEY (id),
	CONSTRAINT uk_auth_apis_path_method UNIQUE (path, method)
);
CREATE INDEX idx_auth_apis_path_method ON public.auth_apis USING btree (path, method);
CREATE INDEX idx_auth_apis_permission_id ON public.auth_apis USING btree (permission_id);


-- public.auth_invite definition

-- Drop table

-- DROP TABLE public.auth_invite;

CREATE TABLE public.auth_invite (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	created_by varchar(50) NOT NULL,
	created_at timestamptz(6) DEFAULT CURRENT_TIMESTAMP NULL,
	updated_by varchar(50) NULL,
	updated_at timestamptz(6) DEFAULT CURRENT_TIMESTAMP NULL,
	email varchar(100) NOT NULL,
	inviter varchar(50) NOT NULL,
	invitation_time timestamptz(6) NULL,
	join_time timestamptz(6) NULL,
	status varchar(20) NOT NULL,
	"type" varchar(20) NOT NULL,
	organization_id uuid NULL,
	workspace_id uuid NULL,
	target_organization_id uuid NULL,
	target_workspace_id uuid NULL,
	target_organization_role_ids _text DEFAULT '{}'::text[] NULL,
	target_workspace_role_ids _text DEFAULT '{}'::text[] NULL,
	CONSTRAINT auth_invite_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_auth_invite_email ON public.auth_invite USING btree (email);
CREATE INDEX idx_auth_invite_status ON public.auth_invite USING btree (status);


-- public.auth_organizations definition

-- Drop table

-- DROP TABLE public.auth_organizations;

CREATE TABLE public.auth_organizations (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	created_by varchar(50) NOT NULL,
	created_at timestamptz(6) DEFAULT CURRENT_TIMESTAMP NULL,
	updated_by varchar(50) NULL,
	updated_at timestamptz(6) DEFAULT CURRENT_TIMESTAMP NULL,
	"name" varchar(100) NOT NULL,
	owner_id uuid NOT NULL,
	description varchar(255) NULL,
	CONSTRAINT auth_organizations_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_auth_organizations_owner_id ON public.auth_organizations USING btree (owner_id);


-- public.auth_permissions definition

-- Drop table

-- DROP TABLE public.auth_permissions;

CREATE TABLE public.auth_permissions (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	created_by varchar(50) NOT NULL,
	created_at timestamptz(6) DEFAULT CURRENT_TIMESTAMP NULL,
	updated_by varchar(50) NULL,
	updated_at timestamptz(6) DEFAULT CURRENT_TIMESTAMP NULL,
	parent_id uuid NULL,
	"name" varchar(100) NOT NULL,
	permission_type varchar(20) NOT NULL,
	description varchar(255) NULL,
	parent_name varchar(255) NULL,
	CONSTRAINT auth_permissions_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_auth_permissions_parent_id ON public.auth_permissions USING btree (parent_id);


-- public.auth_role_permission definition

-- Drop table

-- DROP TABLE public.auth_role_permission;

CREATE TABLE public.auth_role_permission (
	role_id uuid NOT NULL,
	permission_id uuid NOT NULL,
	created_by varchar(50) NOT NULL,
	created_at timestamptz(6) DEFAULT CURRENT_TIMESTAMP NULL,
	updated_by varchar(50) NULL,
	updated_at timestamptz(6) DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT auth_role_permission_pkey PRIMARY KEY (role_id, permission_id)
);
CREATE INDEX idx_auth_role_permission_permission_id ON public.auth_role_permission USING btree (permission_id);
CREATE INDEX idx_auth_role_permission_role_id ON public.auth_role_permission USING btree (role_id);


-- public.auth_roles definition

-- Drop table

-- DROP TABLE public.auth_roles;

CREATE TABLE public.auth_roles (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	created_by varchar(50) NOT NULL,
	created_at timestamptz(6) DEFAULT CURRENT_TIMESTAMP NULL,
	updated_by varchar(50) NULL,
	updated_at timestamptz(6) DEFAULT CURRENT_TIMESTAMP NULL,
	"name" varchar(100) NOT NULL,
	description text NULL,
	organization_id uuid NULL,
	workspace_id uuid NULL,
	role_type varchar(20) NOT NULL,
	status varchar(20) NULL,
	"source" varchar(20) NULL,
	CONSTRAINT auth_roles_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_auth_roles_organization_id ON public.auth_roles USING btree (organization_id);
CREATE INDEX idx_auth_roles_workspace_id ON public.auth_roles USING btree (workspace_id);


-- public.auth_user_organization definition

-- Drop table

-- DROP TABLE public.auth_user_organization;

CREATE TABLE public.auth_user_organization (
	user_id uuid NOT NULL,
	organization_id uuid NOT NULL,
	created_by varchar(50) NOT NULL,
	created_at timestamptz(6) DEFAULT CURRENT_TIMESTAMP NULL,
	updated_by varchar(50) NULL,
	updated_at timestamptz(6) DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT auth_user_organization_pkey PRIMARY KEY (user_id, organization_id)
);
CREATE INDEX idx_auth_user_organization_org_id ON public.auth_user_organization USING btree (organization_id);
CREATE INDEX idx_auth_user_organization_user_id ON public.auth_user_organization USING btree (user_id);


-- public.auth_user_role definition

-- Drop table

-- DROP TABLE public.auth_user_role;

CREATE TABLE public.auth_user_role (
	user_id uuid NOT NULL,
	role_id uuid NOT NULL,
	created_by varchar(50) NOT NULL,
	created_at timestamptz(6) DEFAULT CURRENT_TIMESTAMP NULL,
	updated_by varchar(50) NULL,
	updated_at timestamptz(6) DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT auth_user_role_pkey PRIMARY KEY (user_id, role_id)
);
CREATE INDEX idx_auth_user_role_role_id ON public.auth_user_role USING btree (role_id);
CREATE INDEX idx_auth_user_role_user_id ON public.auth_user_role USING btree (user_id);


-- public.auth_user_workspace definition

-- Drop table

-- DROP TABLE public.auth_user_workspace;

CREATE TABLE public.auth_user_workspace (
	user_id uuid NOT NULL,
	workspace_id uuid NOT NULL,
	created_by varchar(50) NOT NULL,
	created_at timestamptz(6) DEFAULT CURRENT_TIMESTAMP NULL,
	updated_by varchar(50) NULL,
	updated_at timestamptz(6) DEFAULT CURRENT_TIMESTAMP NULL,
	user_status varchar(20) DEFAULT 'active'::character varying NOT NULL,
	CONSTRAINT auth_user_workspace_pkey PRIMARY KEY (user_id, workspace_id)
);
CREATE INDEX idx_auth_user_workspace_user_id ON public.auth_user_workspace USING btree (user_id);
CREATE INDEX idx_auth_user_workspace_workspace_id ON public.auth_user_workspace USING btree (workspace_id);


-- public.auth_users definition

-- Drop table

-- DROP TABLE public.auth_users;

CREATE TABLE public.auth_users (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	created_by varchar(50) NOT NULL,
	created_at timestamptz(6) DEFAULT CURRENT_TIMESTAMP NULL,
	updated_by varchar(50) NULL,
	updated_at timestamptz(6) DEFAULT CURRENT_TIMESTAMP NULL,
	username varchar(50) NOT NULL,
	password_hash varchar(255) NOT NULL,
	nickname varchar(50) NULL,
	email varchar(255) NULL,
	clerk_user_id varchar(255) DEFAULT ''::character varying NOT NULL,
	oauth_provider varchar(50) DEFAULT ''::character varying NOT NULL,
	avatar_url text DEFAULT ''::text NOT NULL,
	email_verified bool DEFAULT false NULL,
	CONSTRAINT auth_users_pkey PRIMARY KEY (id),
	CONSTRAINT auth_users_username_key UNIQUE (username)
);
CREATE INDEX idx_auth_users_username ON public.auth_users USING btree (username);


-- public.auth_workspaces definition

-- Drop table

-- DROP TABLE public.auth_workspaces;

CREATE TABLE public.auth_workspaces (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	created_by varchar(50) NOT NULL,
	created_at timestamptz(6) DEFAULT CURRENT_TIMESTAMP NULL,
	updated_by varchar(50) NULL,
	updated_at timestamptz(6) DEFAULT CURRENT_TIMESTAMP NULL,
	"name" varchar(100) NOT NULL,
	organization_id uuid NOT NULL,
	owner_id uuid NOT NULL,
	status varchar(20) NOT NULL,
	description varchar(255) NULL,
	feature_menu jsonb NULL,
	CONSTRAINT auth_workspaces_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_auth_workspaces_organization_id ON public.auth_workspaces USING btree (organization_id);
CREATE INDEX idx_auth_workspaces_owner_id ON public.auth_workspaces USING btree (owner_id);


-- public.casbin_rule definition

-- Drop table

-- DROP TABLE public.casbin_rule;

CREATE TABLE public.casbin_rule (
	id text NOT NULL,
	ptype text NULL,
	v0 text NULL,
	v1 text NULL,
	v2 text NULL,
	v3 text NULL,
	v4 text NULL,
	v5 text NULL,
	CONSTRAINT casbin_rule_pkey PRIMARY KEY (id),
	CONSTRAINT idx_casbin_rule_unique UNIQUE (ptype, v0, v1, v2, v3, v4, v5)
);
CREATE INDEX idx_casbin_rule_ptype ON public.casbin_rule USING btree (ptype);
CREATE INDEX idx_casbin_rule_v0 ON public.casbin_rule USING btree (v0);
CREATE INDEX idx_casbin_rule_v1 ON public.casbin_rule USING btree (v1);
CREATE INDEX idx_casbin_rule_v2 ON public.casbin_rule USING btree (v2);


-- public.casbin_rules definition

-- Drop table

-- DROP TABLE public.casbin_rules;

CREATE TABLE public.casbin_rules (
	id int8 GENERATED BY DEFAULT AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1 NO CYCLE) NOT NULL,
	ptype varchar DEFAULT ''::character varying NOT NULL,
	v0 varchar DEFAULT ''::character varying NOT NULL,
	v1 varchar DEFAULT ''::character varying NOT NULL,
	v2 varchar DEFAULT ''::character varying NOT NULL,
	v3 varchar DEFAULT ''::character varying NOT NULL,
	v4 varchar DEFAULT ''::character varying NOT NULL,
	v5 varchar DEFAULT ''::character varying NOT NULL,
	CONSTRAINT casbin_rules_pkey PRIMARY KEY (id)
);


-- public.container_metrics definition

-- Drop table

-- DROP TABLE public.container_metrics;

CREATE TABLE public.container_metrics (
	id uuid NOT NULL,
	created_at timestamptz(6) NOT NULL,
	updated_at timestamptz(6) NOT NULL,
	container_name varchar NOT NULL,
	status int8 NULL,
	image varchar DEFAULT ''::character varying NOT NULL,
	ip varchar DEFAULT ''::character varying NOT NULL,
	host_name varchar NOT NULL,
	cpu_usage float8 NULL,
	memory_usage float8 NULL,
	container_id varchar NULL,
	CONSTRAINT container_metrics_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_container_metrics_main ON public.container_metrics USING btree (updated_at, container_id, container_name, status);


-- public.dashboards definition

-- Drop table

-- DROP TABLE public.dashboards;

CREATE TABLE public.dashboards (
	id serial4 NOT NULL,
	"uuid" text NOT NULL,
	title text NOT NULL,
	description text NULL,
	created_at int8 NOT NULL,
	updated_at int8 NOT NULL,
	"data" text NOT NULL,
	created_by text NULL,
	updated_by text NULL,
	"locked" int4 DEFAULT 0 NULL,
	organization_id varchar(64) DEFAULT '1'::character varying NOT NULL,
	workspace_id varchar(64) DEFAULT '1'::character varying NOT NULL,
	CONSTRAINT dashboards_pkey PRIMARY KEY (id),
	CONSTRAINT dashboards_uuid_key UNIQUE (uuid)
);
CREATE INDEX idx_dashboards_org_workspace ON public.dashboards USING btree (organization_id, workspace_id);
COMMENT ON INDEX public.idx_dashboards_org_workspace IS 'Composite index for organization and workspace filtering';


-- public.dashboards_variable definition

-- Drop table

-- DROP TABLE public.dashboards_variable;

CREATE TABLE public.dashboards_variable (
	id varchar NOT NULL,
	dashboard_uuid varchar DEFAULT ''::character varying NULL,
	code varchar(255) DEFAULT ''::character varying NOT NULL,
	"name" varchar(255) DEFAULT ''::character varying NULL,
	data_source varchar(255) DEFAULT ''::character varying NULL,
	value_sort varchar(255) DEFAULT ''::character varying NULL,
	hide bool DEFAULT false NULL,
	include_start bool DEFAULT true NULL,
	multiple bool DEFAULT true NULL,
	definition varchar DEFAULT ''::character varying NULL,
	extend varchar DEFAULT ''::character varying NULL,
	"type" varchar(255) DEFAULT ''::character varying NULL,
	created_at timestamptz(6) NOT NULL,
	updated_at timestamptz(6) NOT NULL,
	deleted bool DEFAULT false NULL,
	CONSTRAINT dashboard_variable_pkey PRIMARY KEY (id)
);


-- public.host_metrics definition

-- Drop table

-- DROP TABLE public.host_metrics;

CREATE TABLE public.host_metrics (
	id uuid NOT NULL,
	created_at timestamptz(6) NOT NULL,
	updated_at timestamptz(6) NOT NULL,
	host_name varchar NOT NULL,
	operating_system varchar(64) NULL,
	host_status int8 NULL,
	cpu_usage float8 NULL,
	memory_usage float8 NULL,
	cpu_load float8 NULL,
	system_kernel varchar(255) NULL,
	system_architecture varchar(255) NULL,
	processor_vendor varchar(255) NULL,
	processor_model varchar(255) NULL,
	processor_cores int2 NULL,
	processor_frequency_mhz float8 NULL,
	ipv6_address text NULL,
	mac_address text NULL,
	total_memory_gb float8 NULL,
	total_swap_gb float8 NULL,
	disk_partition_name text NULL,
	disk_partition_total_bytes text NULL,
	ipv4_address text NULL,
	network_card_name text NULL,
	visibale bool DEFAULT false NOT NULL,
	platform varchar(255) DEFAULT ''::character varying NULL,
	CONSTRAINT host_metrics_pkey PRIMARY KEY (id)
);


-- public.license definition

-- Drop table

-- DROP TABLE public.license;

CREATE TABLE public.license (
	license_id uuid NOT NULL,
	customer_id text NOT NULL,
	customer_name text NOT NULL,
	lighthouse_id text NOT NULL,
	model_name text NOT NULL,
	issue_time int8 NOT NULL,
	expire_time int8 NOT NULL,
	resources_count int8 NOT NULL,
	resources_use_count int8 NOT NULL,
	license_status text NOT NULL,
	enabled bool DEFAULT false NOT NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	organization_id varchar NOT NULL,
	workspace_id varchar NULL,
	CONSTRAINT sys_license_pkey PRIMARY KEY (license_id)
);


-- public.license_resource_stats definition

-- Drop table

-- DROP TABLE public.license_resource_stats;

CREATE TABLE public.license_resource_stats (
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	lighthouse_id varchar NOT NULL,
	model_name varchar NOT NULL,
	resource_use_count int8 NOT NULL,
	organization_id varchar NOT NULL,
	workspace_id varchar NOT NULL,
	CONSTRAINT uk_license_resource_stats_unique UNIQUE (organization_id, workspace_id, model_name, lighthouse_id)
);
CREATE INDEX idx_license_resource_stats_license_id ON public.license_resource_stats USING btree (lighthouse_id);
CREATE INDEX idx_license_resource_stats_organization_id ON public.license_resource_stats USING btree (organization_id);
CREATE INDEX idx_license_resource_stats_workspace_id ON public.license_resource_stats USING btree (workspace_id);


-- public.metric definition

-- Drop table

-- DROP TABLE public.metric;

CREATE TABLE public.metric (
	id varchar NOT NULL,
	created_at timestamptz(6) NOT NULL,
	updated_at timestamptz(6) NOT NULL,
	"name" varchar NOT NULL,
	metrictype varchar NOT NULL,
	fieldtype varchar NOT NULL,
	metricunit varchar NULL,
	memoen varchar NULL,
	metricsetid varchar NOT NULL,
	temporality varchar(255) NULL,
	monotonic bool NULL,
	enabled bool DEFAULT false NULL,
	memocn varchar(255) NULL,
	metricunitshort varchar(255) DEFAULT ''::character varying NULL,
	CONSTRAINT metric_pkey PRIMARY KEY (id)
);


-- public.metric_set definition

-- Drop table

-- DROP TABLE public.metric_set;

CREATE TABLE public.metric_set (
	id varchar NOT NULL,
	created_at timestamptz(6) NOT NULL,
	updated_at timestamptz(6) NOT NULL,
	"name" varchar NOT NULL,
	CONSTRAINT metric_set_pkey PRIMARY KEY (id)
);


-- public.metric_time_series definition

-- Drop table

-- DROP TABLE public.metric_time_series;

CREATE TABLE public.metric_time_series (
	created_at timestamptz(6) NOT NULL,
	updated_at timestamptz(6) NOT NULL,
	metric_id varchar NOT NULL,
	day_series int8 DEFAULT 20250314 NOT NULL,
	num_series int8 DEFAULT 0 NULL,
	CONSTRAINT metric_set_copy1_pkey PRIMARY KEY (metric_id, day_series)
);


-- public.notification_configs definition

-- Drop table

-- DROP TABLE public.notification_configs;

CREATE TABLE public.notification_configs (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	"name" varchar(255) NOT NULL,
	"type" varchar(255) NOT NULL,
	payload json NULL,
	disabled bool DEFAULT false NOT NULL,
	created_by varchar(255) NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	organization_id varchar NOT NULL,
	workspace_id varchar NOT NULL,
	CONSTRAINT notification_configs_pkey PRIMARY KEY (id)
);


-- public.org_billing definition

-- Drop table

-- DROP TABLE public.org_billing;

CREATE TABLE public.org_billing (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	org_id uuid NULL,
	"month" varchar(7) NULL,
	usage_count int8 NULL,
	free_quota int8 NULL,
	overage int8 NULL,
	amount numeric(10, 2) NULL,
	status varchar(20) NULL,
	CONSTRAINT org_billing_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_org_billing_month ON public.org_billing USING btree (month);
CREATE INDEX idx_org_billing_org_id ON public.org_billing USING btree (org_id);


-- public.org_usage definition

-- Drop table

-- DROP TABLE public.org_usage;

CREATE TABLE public.org_usage (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	organization_id varchar NOT NULL,
	workspace_id varchar NOT NULL,
	"month" varchar(7) NOT NULL,
	"usage" jsonb NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT org_usage_pkey PRIMARY KEY (id)
);


-- public.payments definition

-- Drop table

-- DROP TABLE public.payments;

CREATE TABLE public.payments (
	id varchar(255) NOT NULL,
	customer_id varchar(255) NOT NULL,
	subscription_id varchar(255) NULL,
	amount int8 NOT NULL,
	currency varchar(10) NOT NULL,
	status varchar(50) NOT NULL,
	stripe_event_id varchar(255) NULL,
	metadata text NULL,
	created_at timestamptz DEFAULT now() NULL,
	updated_at timestamptz DEFAULT now() NULL,
	CONSTRAINT payments_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_payments_customer_id ON public.payments USING btree (customer_id);
CREATE INDEX idx_payments_subscription_id ON public.payments USING btree (subscription_id);

-- Table Triggers

create trigger update_payments_updated_at before
update
    on
    public.payments for each row execute function update_updated_at_column();


-- public.process_metrics definition

-- Drop table

-- DROP TABLE public.process_metrics;

CREATE TABLE public.process_metrics (
	id uuid NOT NULL,
	created_at timestamptz(6) NOT NULL,
	updated_at timestamptz(6) NOT NULL,
	process_name varchar NOT NULL,
	status int8 NULL,
	host_name varchar NOT NULL,
	cpu_usage float8 NULL,
	memory_usage float8 NULL,
	load_time varchar(64) NOT NULL,
	process_user varchar NULL,
	pid int8 DEFAULT 0 NOT NULL,
	cmdline varchar NULL,
	load_time_original int8 DEFAULT 0 NULL,
	CONSTRAINT process_metrics_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_process_metrics_main ON public.process_metrics USING btree (updated_at, pid, process_name, process_user, host_name, status);


-- public.processed_events definition

-- Drop table

-- DROP TABLE public.processed_events;

CREATE TABLE public.processed_events (
	event_id varchar(255) NOT NULL,
	created_at timestamptz DEFAULT now() NULL,
	CONSTRAINT processed_events_pkey PRIMARY KEY (event_id)
);


-- public.project definition

-- Drop table

-- DROP TABLE public.project;

CREATE TABLE public.project (
	id text NOT NULL,
	"name" text NOT NULL,
	prometheus text NULL,
	settings text NULL,
	CONSTRAINT project_name_key UNIQUE (name),
	CONSTRAINT project_pkey PRIMARY KEY (id)
);


-- public.query_config definition

-- Drop table

-- DROP TABLE public.query_config;

CREATE TABLE public.query_config (
	id uuid NOT NULL,
	db_name text NOT NULL,
	table_name text NOT NULL,
	column_name text NOT NULL,
	data_type text NOT NULL,
	desc_en text NOT NULL,
	desc_zh text NOT NULL,
	capability _text NULL,
	enabled bool DEFAULT false NOT NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	unit text DEFAULT ''::text NULL,
	is_column bool DEFAULT true NOT NULL,
	map_key_name text DEFAULT ''::text NULL,
	CONSTRAINT query_config_pkey PRIMARY KEY (id)
);


-- public.resource_usage definition

-- Drop table

-- DROP TABLE public.resource_usage;

CREATE TABLE public.resource_usage (
	organization_id varchar NOT NULL,
	workspace_id varchar NOT NULL,
	"usage" jsonb NOT NULL,
	CONSTRAINT resource_usage_pkey PRIMARY KEY (organization_id, workspace_id)
);


-- public.rum_applications definition

-- Drop table

-- DROP TABLE public.rum_applications;

CREATE TABLE public.rum_applications (
	id serial4 NOT NULL,
	"name" varchar(255) NOT NULL,
	code varchar(255) NOT NULL,
	"type" varchar(50) DEFAULT 'web'::character varying NOT NULL,
	message jsonb NULL,
	create_time timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	update_time timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"disable" bool DEFAULT false NOT NULL,
	CONSTRAINT rum_applications_pkey PRIMARY KEY (id),
	CONSTRAINT rum_applications_type_check CHECK (((type)::text = ANY (ARRAY[('web'::character varying)::text, ('ios'::character varying)::text, ('android'::character varying)::text, ('miniprogram'::character varying)::text])))
);


-- public.schema_migrations definition

-- Drop table

-- DROP TABLE public.schema_migrations;

CREATE TABLE public.schema_migrations (
	"version" int8 NOT NULL,
	dirty bool NOT NULL,
	CONSTRAINT schema_migrations_pkey PRIMARY KEY (version)
);


-- public.settings definition

-- Drop table

-- DROP TABLE public.settings;

CREATE TABLE public.settings (
	"name" text NOT NULL,
	value text NOT NULL,
	organization_id varchar NOT NULL,
	workspace_id varchar NOT NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	CONSTRAINT settings_pkey PRIMARY KEY (name, organization_id, workspace_id)
);


-- public.subscription_plans definition

-- Drop table

-- DROP TABLE public.subscription_plans;

CREATE TABLE public.subscription_plans (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	tier_name varchar(50) NOT NULL,
	pricing_monthly numeric(10, 2) DEFAULT 0.00 NOT NULL,
	pricing_quarterly numeric(10, 2) DEFAULT 0.00 NOT NULL,
	pricing_yearly numeric(10, 2) DEFAULT 0.00 NOT NULL,
	limits jsonb DEFAULT '{}'::jsonb NOT NULL,
	features jsonb NULL,
	target_users text NULL,
	upgrade_path text NULL,
	is_custom bool DEFAULT false NULL,
	default_flow_package public.flow_package_type NULL,
	is_active bool DEFAULT true NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	stripe_price_id_monthly varchar(255) NULL,
	stripe_price_id_quarterly varchar(255) NULL,
	stripe_price_id_yearly varchar(255) NULL,
	CONSTRAINT subscription_plans_pkey PRIMARY KEY (id),
	CONSTRAINT subscription_plans_tier_name_key UNIQUE (tier_name)
);
CREATE INDEX idx_subscription_plans_is_active ON public.subscription_plans USING btree (is_active);
CREATE INDEX idx_subscription_plans_price_monthly ON public.subscription_plans USING btree (stripe_price_id_monthly) WHERE (stripe_price_id_monthly IS NOT NULL);
CREATE INDEX idx_subscription_plans_price_quarterly ON public.subscription_plans USING btree (stripe_price_id_quarterly) WHERE (stripe_price_id_quarterly IS NOT NULL);
CREATE INDEX idx_subscription_plans_price_yearly ON public.subscription_plans USING btree (stripe_price_id_yearly) WHERE (stripe_price_id_yearly IS NOT NULL);
CREATE INDEX idx_subscription_plans_tier_name ON public.subscription_plans USING btree (tier_name);


-- public.users definition

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE public.users (
	id serial4 NOT NULL,
	email text NOT NULL,
	"name" text NOT NULL,
	"password" text NOT NULL,
	roles text NOT NULL,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (id)
);


-- public.webhook_events definition

-- Drop table

-- DROP TABLE public.webhook_events;

CREATE TABLE public.webhook_events (
	pk_id bigserial NOT NULL,
	stripe_event_id varchar(255) NOT NULL,
	event_type varchar(100) NOT NULL,
	event_data jsonb NOT NULL,
	api_version varchar(50) NULL,
	object_id varchar(255) NULL,
	processed bool DEFAULT false NULL,
	processed_at timestamp NULL,
	processing_result text NULL,
	error_message text NULL,
	retry_count int4 DEFAULT 0 NULL,
	request_id varchar(255) NULL,
	source_ip varchar(45) NULL,
	user_agent text NULL,
	received_at timestamp DEFAULT now() NOT NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	CONSTRAINT webhook_events_pkey PRIMARY KEY (pk_id)
);
CREATE INDEX idx_webhook_events_created_at ON public.webhook_events USING btree (created_at DESC);
CREATE INDEX idx_webhook_events_event_type ON public.webhook_events USING btree (event_type);
CREATE INDEX idx_webhook_events_object_id ON public.webhook_events USING btree (object_id);
CREATE INDEX idx_webhook_events_processed ON public.webhook_events USING btree (processed);
CREATE INDEX idx_webhook_events_received_at ON public.webhook_events USING btree (received_at DESC);
CREATE INDEX idx_webhook_events_stripe_event_id ON public.webhook_events USING btree (stripe_event_id);


-- public.alert_rules definition

-- Drop table

-- DROP TABLE public.alert_rules;

CREATE TABLE public.alert_rules (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	alert_id uuid NOT NULL,
	"type" varchar(255) NOT NULL,
	composite_query json NOT NULL,
	conditions json NOT NULL,
	duration int4 NOT NULL,
	"period" int4 NOT NULL,
	"options" json NULL,
	created_by varchar(255) DEFAULT ''::character varying NOT NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	organization_id varchar NOT NULL,
	workspace_id varchar NOT NULL,
	CONSTRAINT alert_rules_pkey PRIMARY KEY (id),
	CONSTRAINT fk_alert_id FOREIGN KEY (alert_id) REFERENCES public.alerts(id)
);


-- public.alert_schedules definition

-- Drop table

-- DROP TABLE public.alert_schedules;

CREATE TABLE public.alert_schedules (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	alert_id uuid NOT NULL,
	next_run timestamptz NOT NULL,
	payload json NULL,
	active bool DEFAULT true NULL,
	"content" text NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	delay int4 DEFAULT 60 NOT NULL,
	organization_id varchar NOT NULL,
	workspace_id varchar NOT NULL,
	CONSTRAINT alert_schedules_pkey PRIMARY KEY (id),
	CONSTRAINT fk_alert_id FOREIGN KEY (alert_id) REFERENCES public.alerts(id)
);


-- public.application_deployment definition

-- Drop table

-- DROP TABLE public.application_deployment;

CREATE TABLE public.application_deployment (
	project_id text NOT NULL,
	application_id text NOT NULL,
	"name" text NOT NULL,
	started_at int4 NOT NULL,
	finished_at int4 DEFAULT 0 NOT NULL,
	details text NULL,
	metrics_snapshot text NULL,
	notifications text NULL,
	CONSTRAINT application_deployment_pkey PRIMARY KEY (project_id, application_id, started_at),
	CONSTRAINT application_deployment_project_id_fkey FOREIGN KEY (project_id) REFERENCES public.project(id)
);


-- public.application_settings definition

-- Drop table

-- DROP TABLE public.application_settings;

CREATE TABLE public.application_settings (
	project_id text NOT NULL,
	application_id text NOT NULL,
	settings text NOT NULL,
	CONSTRAINT application_settings_pkey PRIMARY KEY (project_id, application_id),
	CONSTRAINT application_settings_project_id_fkey FOREIGN KEY (project_id) REFERENCES public.project(id)
);


-- public.check_configs definition

-- Drop table

-- DROP TABLE public.check_configs;

CREATE TABLE public.check_configs (
	project_id text NOT NULL,
	application_id text NOT NULL,
	configs text NULL,
	CONSTRAINT check_configs_pkey PRIMARY KEY (project_id, application_id),
	CONSTRAINT check_configs_project_id_fkey FOREIGN KEY (project_id) REFERENCES public.project(id)
);


-- public.incident definition

-- Drop table

-- DROP TABLE public.incident;

CREATE TABLE public.incident (
	project_id text NOT NULL,
	application_id text NOT NULL,
	"key" text NOT NULL,
	opened_at int4 NOT NULL,
	resolved_at int4 DEFAULT 0 NOT NULL,
	severity int4 NOT NULL,
	CONSTRAINT incident_pkey PRIMARY KEY (project_id, application_id, opened_at),
	CONSTRAINT incident_project_id_fkey FOREIGN KEY (project_id) REFERENCES public.project(id)
);
CREATE UNIQUE INDEX incident_key ON public.incident USING btree (project_id, key);


-- public.incident_notification definition

-- Drop table

-- DROP TABLE public.incident_notification;

CREATE TABLE public.incident_notification (
	project_id text NOT NULL,
	application_id text NOT NULL,
	incident_key text NOT NULL,
	status int4 NOT NULL,
	destination text NOT NULL,
	"timestamp" int4 NOT NULL,
	sent_at int4 DEFAULT 0 NOT NULL,
	external_key text DEFAULT ''::text NOT NULL,
	details text NULL,
	CONSTRAINT incident_notification_project_id_fkey FOREIGN KEY (project_id) REFERENCES public.project(id)
);


-- public.subscription_users definition

-- Drop table

-- DROP TABLE public.subscription_users;

CREATE TABLE public.subscription_users (
	id varchar DEFAULT gen_random_uuid() NOT NULL,
	user_id uuid NOT NULL,
	plan_id uuid NOT NULL,
	status varchar(20) DEFAULT 'active'::character varying NOT NULL,
	billing_cycle varchar(20) DEFAULT 'monthly'::character varying NOT NULL,
	start_date timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	end_date timestamp NULL,
	payment_method varchar(50) NULL,
	last_billed_at timestamp NULL,
	trial_days_used int4 DEFAULT 0 NULL,
	organization_id varchar NOT NULL,
	notes text NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	stripe_session_id varchar(255) NULL,
	stripe_session_data jsonb NULL,
	CONSTRAINT subscription_users_billing_cycle_check CHECK (((billing_cycle)::text = ANY ((ARRAY['monthly'::character varying, 'quarterly'::character varying, 'yearly'::character varying])::text[]))),
	CONSTRAINT subscription_users_pkey PRIMARY KEY (id),
	CONSTRAINT subscription_users_status_check CHECK (((status)::text = ANY ((ARRAY['active'::character varying, 'trial'::character varying, 'canceled'::character varying, 'expired'::character varying, 'suspended'::character varying])::text[]))),
	CONSTRAINT subscription_users_plan_id_fkey FOREIGN KEY (plan_id) REFERENCES public.subscription_plans(id) ON DELETE RESTRICT
);
CREATE INDEX idx_subscription_users_plan_id ON public.subscription_users USING btree (plan_id);
CREATE INDEX idx_subscription_users_status ON public.subscription_users USING btree (status);
CREATE INDEX idx_subscription_users_stripe_session_id ON public.subscription_users USING btree (stripe_session_id) WHERE (stripe_session_id IS NOT NULL);
CREATE INDEX idx_subscription_users_user_id ON public.subscription_users USING btree (user_id);


-- public.alert_history definition

-- Drop table

-- DROP TABLE public.alert_history;

CREATE TABLE public.alert_history (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
	rule_id uuid NOT NULL,
	severity varchar(255) NOT NULL,
	message text NOT NULL,
	"start" int8 NOT NULL,
	"end" int8 NOT NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"oid" varchar(255) NOT NULL,
	status varchar(10) DEFAULT 'completed'::character varying NOT NULL,
	reason text NULL,
	alert_name varchar(255) DEFAULT ''::character varying NOT NULL,
	organization_id varchar NOT NULL,
	workspace_id varchar NOT NULL,
	CONSTRAINT alert_history_pkey PRIMARY KEY (id),
	CONSTRAINT fk_rule_id FOREIGN KEY (rule_id) REFERENCES public.alert_rules(id)
);
