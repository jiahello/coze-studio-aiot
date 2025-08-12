-- IoT hardware devices and TTS settings
-- Align with existing schema conventions: InnoDB + utf8mb4, bigint ms timestamps, no FK constraints

-- 1) hardware_device: AI 硬件设备管理
CREATE TABLE IF NOT EXISTS `hardware_device` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `device_id` varchar(64) NOT NULL COMMENT '设备唯一ID',
  `name` varchar(255) NOT NULL COMMENT '设备名称',
  `description` text NULL COMMENT '设备描述',
  `app_id` bigint unsigned NULL COMMENT '绑定应用ID（可选）',
  `space_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '空间ID',
  `created_user_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建人',
  `updated_user_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新人',
  `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间(毫秒)',
  `updated_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间(毫秒)',
  `firmware_version` varchar(64) NULL COMMENT '固件版本',
  `mac_address` varchar(64) NULL COMMENT 'MAC地址',
  `status` varchar(255) NOT NULL DEFAULT 'offline' COMMENT '设备状态: offline/online/pairing/blocked',
  `last_ping_at` bigint unsigned NULL COMMENT '最后心跳(毫秒)',
  `verify_code` varchar(8) NULL COMMENT '配对验证码',
  `from_end_user_id` varchar(255) NULL COMMENT '来源终端用户ID(可选)',
  `from_account_id` varchar(255) NULL COMMENT '来源账号ID(可选)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_space_device` (`space_id`, `device_id`),
  UNIQUE KEY `uniq_space_verify_code` (`space_id`, `verify_code`),
  KEY `idx_space` (`space_id`),
  KEY `idx_app` (`app_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI 硬件设备管理';


-- 2) tts_voice: TTS 音色库（系统级/空间级）
CREATE TABLE IF NOT EXISTS `tts_voice` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `provider` varchar(255) NOT NULL COMMENT 'TTS 服务商',
  `model` varchar(255) NULL COMMENT '模型/产品线',
  `voice_type` varchar(32) NOT NULL DEFAULT 'system' COMMENT '音色类型: system/custom/uploaded',
  `name` varchar(255) NOT NULL COMMENT '显示名称',
  `voice_code` varchar(255) NOT NULL COMMENT '音色编码',
  `description` text NULL COMMENT '描述',
  `gender` varchar(32) NULL COMMENT '性别',
  `language` varchar(64) NULL COMMENT '语言编码',
  `scenario` varchar(255) NULL COMMENT '场景',
  `sound_quality` varchar(64) NULL COMMENT '音质',
  `sample_rate` varchar(64) NULL COMMENT '采样率',
  `timestamp_support` varchar(16) NULL COMMENT '时间戳支持',
  `erhua_support` varchar(16) NULL COMMENT '儿化音支持',
  `sample_url` varchar(512) NULL COMMENT '样音URL',
  `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间(毫秒)',
  `updated_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间(毫秒)',
  `source_file_id` varchar(255) NULL COMMENT '来源文件ID/URI',
  `sample_file_id` varchar(255) NULL COMMENT '样音文件ID/URI',
  `is_deleted` bool NOT NULL DEFAULT 0 COMMENT '是否删除',
  `deleted_at` datetime NULL COMMENT '删除时间',
  `deleted_by_id` bigint unsigned NULL COMMENT '删除人ID',
  `language_name` varchar(255) NULL COMMENT '语言名',
  `voice_id` varchar(255) NULL COMMENT '外部音色ID',
  `model_type` varchar(64) NULL COMMENT '模型类型',
  `preview_text` text NULL COMMENT '预听文本',
  `emotion_support` bool DEFAULT 0 COMMENT '情感支持',
  `emotions` json NULL COMMENT '情感能力',
  `from_account_id` varchar(255) NULL COMMENT '来源账户',
  `api_type` varchar(64) NULL COMMENT 'API 类型',
  `space_id` bigint unsigned NULL COMMENT '空间ID，NULL表示系统级',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_provider_voice_space` (`provider`, `voice_code`, `space_id`),
  KEY `idx_space` (`space_id`),
  KEY `idx_provider` (`provider`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='TTS 音色库';


-- 3) app_tts_settings: 应用级 TTS 设置（每 app 一条）
CREATE TABLE IF NOT EXISTS `app_tts_settings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `app_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '应用ID',
  `provider` varchar(255) NOT NULL COMMENT 'TTS 服务商',
  `model` varchar(255) NOT NULL COMMENT 'TTS 模型',
  `voice` varchar(255) NOT NULL COMMENT '音色编码(冗余)',
  `voice_ref` bigint unsigned NULL COMMENT '引用 tts_voice.id',
  `created_user_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建人',
  `updated_user_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新人',
  `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间(毫秒)',
  `updated_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间(毫秒)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_app` (`app_id`),
  KEY `idx_voice_ref` (`voice_ref`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='应用级 TTS 设置';


-- 4) hardware_tts_settings: 设备级 TTS 设置（设备优先）
CREATE TABLE IF NOT EXISTS `hardware_tts_settings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `device_id` varchar(255) NOT NULL COMMENT '设备ID',
  `hardware_device_id` bigint unsigned NULL COMMENT '关联 hardware_device.id',
  `provider` varchar(255) NOT NULL COMMENT 'TTS 服务商',
  `model` varchar(255) NOT NULL COMMENT 'TTS 模型',
  `voice` varchar(255) NOT NULL COMMENT '音色编码(冗余)',
  `voice_ref` bigint unsigned NULL COMMENT '引用 tts_voice.id',
  `created_user_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建人',
  `updated_user_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新人',
  `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间(毫秒)',
  `updated_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间(毫秒)',
  `is_deleted` bool NOT NULL DEFAULT 0 COMMENT '是否删除',
  `deleted_at` datetime NULL COMMENT '删除时间',
  `deleted_by_id` bigint unsigned NULL COMMENT '删除人ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_device_id` (`device_id`),
  KEY `idx_hardware_device_id` (`hardware_device_id`),
  KEY `idx_voice_ref` (`voice_ref`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设备级 TTS 设置';