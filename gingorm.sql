/*
 Navicat Premium Dump SQL

 Source Server         : 192.168.92.133_8.4.3
 Source Server Type    : MySQL
 Source Server Version : 80403 (8.4.3)
 Source Host           : 192.168.92.133:6127
 Source Schema         : gingorm

 Target Server Type    : MySQL
 Target Server Version : 80403 (8.4.3)
 File Encoding         : 65001

 Date: 02/07/2026 18:34:15
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for captcha_code
-- ----------------------------
DROP TABLE IF EXISTS `captcha_code`;
CREATE TABLE `captcha_code`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `captcha_type` tinyint NOT NULL DEFAULT 1 COMMENT '1 mobile 2 email',
  `captcha_account` varchar(55) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `captcha_code` varchar(6) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `is_expired` tinyint UNSIGNED NOT NULL DEFAULT 0 COMMENT '0：未过期  1：已失效',
  `is_used` tinyint NOT NULL DEFAULT 0 COMMENT '1：已使用  0：未使用',
  `expired_time` int NOT NULL DEFAULT 0,
  `created_at` int UNSIGNED NOT NULL DEFAULT 0,
  `updated_at` int UNSIGNED NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of captcha_code
-- ----------------------------
INSERT INTO `captcha_code` VALUES (1, 2, 'user001@user001.com', '111111', 0, 1, 1782982818, 1782982758, 1782983232);
INSERT INTO `captcha_code` VALUES (2, 1, '13800138000', '111111', 0, 1, 1782982896, 1782982836, 1782983179);
INSERT INTO `captcha_code` VALUES (3, 1, '13800138000', '222222', 0, 1, 1782983353, 1782983293, 1782985277);
INSERT INTO `captcha_code` VALUES (4, 1, '13800138000', '141218', 1, 0, 1782986924, 1782986864, 1782987185);
INSERT INTO `captcha_code` VALUES (5, 1, '13800138000', '708613', 0, 1, 1782987202, 1782987142, 1782987200);

-- ----------------------------
-- Table structure for member
-- ----------------------------
DROP TABLE IF EXISTS `member`;
CREATE TABLE `member`  (
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
-- Records of member
-- ----------------------------
INSERT INTO `member` VALUES (1, 'user001', '$2a$10$.GcpveOTBnDGXhWlnYkChuYZ8UbSwzAJAlhRR3RsQt3Gs7er/Vwu.', '13800138000', 'user001@user001.com', 0, 1773824205, '', 1782978436, '127.0.0.1', 0, '', 1, '', '', 1, 1782987606, 1, '');

-- ----------------------------
-- Table structure for member_account
-- ----------------------------
DROP TABLE IF EXISTS `member_account`;
CREATE TABLE `member_account`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `member_id` int NOT NULL DEFAULT 0,
  `amount` decimal(10, 3) NOT NULL,
  `amount_type` tinyint NOT NULL DEFAULT 1,
  `amount_time` int NOT NULL,
  `amount_memos` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of member_account
-- ----------------------------
INSERT INTO `member_account` VALUES (1, 1, 0.000, 1, 0, '');

-- ----------------------------
-- Table structure for member_address
-- ----------------------------
DROP TABLE IF EXISTS `member_address`;
CREATE TABLE `member_address`  (
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
  `updated_at` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 23 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户常用地址表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of member_address
-- ----------------------------
INSERT INTO `member_address` VALUES (1, 1, 'addressaddressaddress11112', '13800138003', 0, 1773909503, 'consignee_name4', '610005', 123546, 1477, 2588, 0, 0);
INSERT INTO `member_address` VALUES (2, 1, 'addressaddressaddress1', '13800138000', 0, 1773909594, 'consignee_name', '610000', 12354, 147, 258, 1778330627, 1778330316);
INSERT INTO `member_address` VALUES (3, 1, 'addressaddressaddress1', '13800138000', 0, 1773909698, 'consignee_name', '610000', 12354, 147, 258, 0, 0);
INSERT INTO `member_address` VALUES (4, 1, 'addressaddressaddress1', '13800138000', 1, 1773909913, 'consignee_name', '610000', 12354, 147, 258, 0, 0);
INSERT INTO `member_address` VALUES (5, 1, 'this is address1', '13800138000', 0, 1782979742, '联系人1', '邮编1', 123456, 147852, 159753, 0, 1782979742);
INSERT INTO `member_address` VALUES (6, 1, 'this is address2', '13800138001', 0, 1782979742, '联系人2', '邮编2', 123456, 147852, 159753, 0, 1782979742);
INSERT INTO `member_address` VALUES (7, 1, 'this is address3', '13800138002', 0, 1782979742, '联系人3', '邮编3', 123456, 147852, 159753, 0, 1782979742);
INSERT INTO `member_address` VALUES (8, 1, 'this is address4', '13800138003', 0, 1782979742, '联系人4', '邮编4', 123456, 147852, 159753, 0, 1782979742);
INSERT INTO `member_address` VALUES (9, 1, 'this is address5', '13800138004', 0, 1782979742, '联系人5', '邮编5', 123456, 147852, 159753, 0, 1782979742);
INSERT INTO `member_address` VALUES (10, 1, 'this is address6', '13800138005', 0, 1782979742, '联系人6', '邮编6', 123456, 147852, 159753, 0, 1782979742);
INSERT INTO `member_address` VALUES (11, 1, 'this is address7', '13800138006', 0, 1782979742, '联系人7', '邮编7', 123456, 147852, 159753, 0, 1782979742);
INSERT INTO `member_address` VALUES (12, 1, 'this is address8', '13800138007', 0, 1782979742, '联系人8', '邮编8', 123456, 147852, 159753, 0, 1782979742);
INSERT INTO `member_address` VALUES (13, 1, 'this is address9', '13800138008', 0, 1782979742, '联系人9', '邮编9', 123456, 147852, 159753, 0, 1782979742);
INSERT INTO `member_address` VALUES (14, 1, 'this is address10', '13800138009', 0, 1782979742, '联系人10', '邮编10', 123456, 147852, 159753, 0, 1782979742);
INSERT INTO `member_address` VALUES (15, 1, 'this is address11', '13800138010', 0, 1782979742, '联系人11', '邮编11', 123456, 147852, 159753, 0, 1782979742);
INSERT INTO `member_address` VALUES (16, 1, 'this is address12', '13800138011', 0, 1782979742, '联系人12', '邮编12', 123456, 147852, 159753, 0, 1782979742);
INSERT INTO `member_address` VALUES (17, 1, 'this is address13', '13800138012', 0, 1782979742, '联系人13', '邮编13', 123456, 147852, 159753, 0, 1782979742);
INSERT INTO `member_address` VALUES (18, 1, 'this is address14', '13800138013', 0, 1782979742, '联系人14', '邮编14', 123456, 147852, 159753, 0, 1782979742);
INSERT INTO `member_address` VALUES (19, 1, 'this is address15', '13800138014', 0, 1782979742, '联系人15', '邮编15', 123456, 147852, 159753, 0, 1782979742);
INSERT INTO `member_address` VALUES (20, 1, 'this is address16', '13800138015', 0, 1782979742, '联系人16', '邮编16', 123456, 147852, 159753, 0, 1782979742);
INSERT INTO `member_address` VALUES (21, 1, 'this is address17', '13800138016', 0, 1782979742, '联系人17', '邮编17', 123456, 147852, 159753, 0, 1782979742);
INSERT INTO `member_address` VALUES (22, 1, 'this is address18', '13800138017', 0, 1782979742, '联系人18', '邮编18', 123456, 147852, 159753, 0, 1782979742);

-- ----------------------------
-- Table structure for member_profile
-- ----------------------------
DROP TABLE IF EXISTS `member_profile`;
CREATE TABLE `member_profile`  (
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
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of member_profile
-- ----------------------------
INSERT INTO `member_profile` VALUES (1, 1, 0.000, 0.000, 0.000, 0.000, 0, 'nick_name', 'true_name', 1, '13800138001', 12345, 12354, '', './storage/public/upload/20260702/20260702182800_eda48d25704708b43e9c58fbd19bb2af.jpg', 'this is address', 14752, 22, 0);

-- ----------------------------
-- Table structure for member_upload
-- ----------------------------
DROP TABLE IF EXISTS `member_upload`;
CREATE TABLE `member_upload`  (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `member_id` int(10) UNSIGNED ZEROFILL NOT NULL,
  `file_name` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `save_file_path` varchar(800) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of member_upload
-- ----------------------------
INSERT INTO `member_upload` VALUES (1, 0000000001, '20260702162955_53ba3b303a21112def8571a550c11554.zip', './storage/public/upload/20260702/20260702162955_53ba3b303a21112def8571a550c11554.zip');

SET FOREIGN_KEY_CHECKS = 1;
