-- 记录所有网络请求，为防止存储空间爆炸，默认清理7天前的数据。可修改prune.sql文件以覆盖默认配置。
CREATE TABLE IF NOT EXISTS netguard_requests (
        id SERIAL PRIMARY KEY,
		request_id int8,
		client_ip VARCHAR(45),
		x_forwarded_for VARCHAR(255),
		user_agent VARCHAR(500),
		http_referer VARCHAR(255),
		request_url varchar(1000) NOT NULL,
		request_headers json NOT NULL,
		raw_url varchar(1000) NOT NULL,
		deleted_at timestamp NULL,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp DEFAULT CURRENT_TIMESTAMP
    );
CREATE INDEX IF NOT EXISTS "IDX_client_ip" ON netguard_requests USING btree (client_ip);
CREATE INDEX IF NOT EXISTS "IDX_created_at_client_ip" ON netguard_requests USING btree (created_at, client_ip);
CREATE INDEX IF NOT EXISTS "IDX_http_referer" ON netguard_requests USING btree (http_referer);

-- 记录被拦截的网络请求，为防止存储空间爆炸，默认清理30天前的数据。可修改prune.sql文件以覆盖默认配置。
CREATE TABLE IF NOT EXISTS netguard_block_requests (
        id SERIAL PRIMARY KEY,
		request_id int8,
		client_ip VARCHAR(45),
		x_forwarded_for VARCHAR(255),
		user_agent VARCHAR(500),
		http_referer VARCHAR(255),
		request_url varchar(1000) NOT NULL,
		request_headers json NOT NULL,
		raw_url varchar(1000) NOT NULL,
		block_type SMALLINT NOT NULL DEFAULT 0,
		deleted_at timestamp NULL,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON COLUMN netguard_block_requests.block_type IS '拦截阻断的理由类别。0=IP黑名单拦截 1=漏洞扫描拦截 2=网络爬虫 3=异常的UserAgent';
CREATE INDEX IF NOT EXISTS "IDX_client_ip_block_requests" ON netguard_block_requests USING btree (client_ip);
CREATE INDEX IF NOT EXISTS "IDX_block_requests_block_type" ON netguard_block_requests USING btree (block_type);

-- 记录归档的网络请求。不会被自动删除。
CREATE TABLE IF NOT EXISTS netguard_archived_requests (
        id SERIAL PRIMARY KEY,
		request_id int8,
		client_ip VARCHAR(45),
		x_forwarded_for VARCHAR(255),
		user_agent VARCHAR(500),
		http_referer VARCHAR(255),
		request_url varchar(1000) NOT NULL,
		request_headers json NOT NULL,
		raw_url varchar(1000) NOT NULL,
		archived_type SMALLINT NOT NULL DEFAULT 0,
		remark VARCHAR(64),
		created_at timestamp DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS "IDX_client_ip_archived_requests" ON netguard_archived_requests USING btree (client_ip);
-- 添加字段注释
COMMENT ON COLUMN netguard_archived_requests.archived_type IS '归档类型：0默认，1漏洞扫描，2网络爬虫';

--IP白名单
CREATE TABLE IF NOT EXISTS netguard_ip_white_list (
        id SERIAL PRIMARY KEY,
		ip VARCHAR(45) NOT NULL,
		title VARCHAR(64) DEFAULT NULL,
		deleted_at timestamp NULL,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);
-- 为ip字段添加唯一约束（此操作会自动创建唯一索引）
CREATE UNIQUE INDEX IF NOT EXISTS "UQE_ip_white_list_ip" ON netguard_ip_white_list USING btree (ip);

-- IP黑名单
CREATE TABLE IF NOT EXISTS netguard_ip_black_list (
        id SERIAL PRIMARY KEY,
		ip VARCHAR(45) NOT NULL,
		title VARCHAR(64) DEFAULT NULL,
		black_type SMALLINT NOT NULL DEFAULT 0,
		deleted_at timestamp NULL,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);
-- 为ip字段添加唯一约束（此操作会自动创建唯一索引）
CREATE UNIQUE INDEX IF NOT EXISTS "UQE_ip_black_list_ip" ON netguard_ip_black_list USING btree (ip);
-- 为black_type字段创建普通索引（非唯一）
CREATE INDEX IF NOT EXISTS "IDX_ip_black_list_black_type" ON netguard_ip_black_list (black_type);


-- statis 数据统计
CREATE TABLE IF NOT EXISTS netguard_statis (
	id SERIAL PRIMARY KEY,
	bucket_id SMALLINT NOT NULL DEFAULT 0,
	statis_date DATE NOT NULL DEFAULT CURRENT_DATE,
	request_count int8 NOT NULL DEFAULT 0,
	blocked_count int8 NOT NULL DEFAULT 0,
	blocked_black_count int8 NOT NULL DEFAULT 0,
	blocked_scanvul_count int8 NOT NULL DEFAULT 0,
	blocked_webspider_count int8 NOT NULL DEFAULT 0,
	blocked_useragent_count int8 NOT NULL DEFAULT 0,
	request_size int8 NOT NULL DEFAULT 0,
	blocked_size int8 NOT NULL DEFAULT 0,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON COLUMN netguard_statis.bucket_id IS '存储空间ID';
COMMENT ON COLUMN netguard_statis.statis_date IS '统计日期';
COMMENT ON COLUMN netguard_statis.request_count IS '请求次数';
COMMENT ON COLUMN netguard_statis.blocked_count IS '拦截请求次数';
COMMENT ON COLUMN netguard_statis.blocked_black_count IS '黑名单拦截请求次数';
COMMENT ON COLUMN netguard_statis.blocked_scanvul_count IS '漏洞扫描拦截请求次数';
COMMENT ON COLUMN netguard_statis.blocked_webspider_count IS '网络爬虫拦截请求次数';
COMMENT ON COLUMN netguard_statis.blocked_useragent_count IS '异常用户代理拦截请求次数';
COMMENT ON COLUMN netguard_statis.request_size IS '请求消耗流量大小';
COMMENT ON COLUMN netguard_statis.blocked_size IS '拦截流量大小';
CREATE INDEX IF NOT EXISTS "IDX_statis_bucket_id" ON netguard_statis (bucket_id);
CREATE INDEX IF NOT EXISTS "IDX_statis_statis_date" ON netguard_statis (statis_date);
