/*
 Navicat Premium Data Transfer

 Source Server         : 192.168.92.133_8.0.42
 Source Server Type    : MySQL
 Source Server Version : 80042
 Source Host           : 192.168.92.133:35136
 Source Schema         : gingorm

 Target Server Type    : MySQL
 Target Server Version : 80042
 File Encoding         : 65001

 Date: 23/03/2026 12:08:27
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for ps_member
-- ----------------------------
DROP TABLE IF EXISTS `ps_member`;
CREATE TABLE `ps_member`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `member_name` varchar(35) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名字',
  `member_pass` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码',
  `mobile` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号码',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '邮箱',
  `parent_id` int NULL DEFAULT 0,
  `created_at` int NOT NULL DEFAULT 0 COMMENT '创建时间',
  `created_ip` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '创建IP',
  `login_time` int NULL DEFAULT 0 COMMENT '上次登陆时间',
  `login_ip` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '上次登陆IP',
  `login_times` int NULL DEFAULT 0 COMMENT '总登陆次数',
  `salt` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '用户密码salt',
  `status` int UNSIGNED NOT NULL DEFAULT 1 COMMENT '会员状态 1 正常 0 已禁用',
  `open_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '',
  `union_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '',
  `reg_type` tinyint UNSIGNED NULL DEFAULT 1 COMMENT '注册类型 1 普通 2 手机号码密码注册 3 微信小程序注册 4 手机号码验证码注册',
  `updated_at` int NOT NULL DEFAULT 0,
  `member_type_id` int NULL DEFAULT 1 COMMENT '用户等级',
  `invite_code` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `mobile`(`mobile` ASC) USING BTREE,
  INDEX `email`(`email` ASC) USING BTREE,
  INDEX `open_id`(`open_id` ASC) USING BTREE,
  INDEX `member_name`(`member_name` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for ps_member_account
-- ----------------------------
DROP TABLE IF EXISTS `ps_member_account`;
CREATE TABLE `ps_member_account`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `member_id` int NOT NULL DEFAULT 0,
  `amount` decimal(10, 3) NOT NULL,
  `amount_type` tinyint NOT NULL DEFAULT 1,
  `amount_time` int NOT NULL,
  `amount_memos` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for ps_member_address
-- ----------------------------
DROP TABLE IF EXISTS `ps_member_address`;
CREATE TABLE `ps_member_address`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `member_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT ' 用户id唯一标识',
  `address` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '收货地址',
  `tel` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '收货人手机号码',
  `is_default` tinyint NOT NULL DEFAULT 0 COMMENT '是否默认',
  `created_at` int NOT NULL DEFAULT 0 COMMENT '创建时间',
  `consignee_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '收货人姓名',
  `post` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '610000' COMMENT '邮编',
  `province_id` int NOT NULL DEFAULT 0 COMMENT '省份名称',
  `city_id` int NOT NULL DEFAULT 0 COMMENT '城市名称',
  `area_id` int NOT NULL DEFAULT 0 COMMENT '区县名称',
  `deleted_at` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户常用地址表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for ps_member_profile
-- ----------------------------
DROP TABLE IF EXISTS `ps_member_profile`;
CREATE TABLE `ps_member_profile`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `member_id` int UNSIGNED NOT NULL DEFAULT 0,
  `account_balance` decimal(10, 3) NULL DEFAULT 0.000,
  `available_balance` decimal(10, 3) NULL DEFAULT 0.000,
  `non_available_balance` decimal(10, 3) NULL DEFAULT 0.000,
  `total_spend_balance` decimal(10, 3) NULL DEFAULT 0.000,
  `total_consumpe_times` int NULL DEFAULT 0,
  `nick_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '',
  `true_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '',
  `sex` tinyint NULL DEFAULT 0,
  `mobile` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '',
  `province_id` int UNSIGNED NULL DEFAULT 0,
  `city_id` int UNSIGNED NULL DEFAULT 0,
  `open_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '',
  `head_img` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '',
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '',
  `area_id` int UNSIGNED NULL DEFAULT 0,
  `age` int NULL DEFAULT 1,
  `scores` int NULL DEFAULT 0,
  `is_remind_check_in` tinyint NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

SET FOREIGN_KEY_CHECKS = 1;
